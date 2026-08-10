#!/usr/bin/env bash
# Maps changed btp/provider/*.go files to a go test -run regex.
# Reads changed file paths from stdin, one per line.
# Outputs a single regex string or "." (meaning: run everything).
#
# Usage:
#   git diff --name-only origin/main...HEAD -- 'btp/provider/*.go' | bash .github/scripts/get_test_filter.sh
#
# Examples:
#   resource_subaccount_entitlement.go   → TestResourceSubaccountEntitlement
#   datasource_directory_role.go         → TestDataSourceDirectoryRole
#   list_resource_subaccount.go          → TestSubaccountListResource
#   action_restore_subaccount.go         → TestActionRestoreSubaccount
#   helper_*.go / provider.go            → . (full suite, shared files affect everything)

# Converts snake_case to PascalCase: sub_account_role → SubAccountRole
to_pascal() {
  echo "$1" | awk 'BEGIN{FS="_"; OFS=""} {for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2); print}'
}

patterns=()

while IFS= read -r file; do
  [[ -z "$file" ]] && continue

  # Get the base filename without .go and strip _test suffix
  base=$(basename "$file" .go)
  base="${base%_test}"

  case "$base" in
    resource_*)    patterns+=("TestResource$(to_pascal "${base#resource_}")")    ;;
    datasource_*)  patterns+=("TestDataSource$(to_pascal "${base#datasource_}")") ;;
    list_resource*)
      # list_resource_subaccount → SubaccountListResource (naming convention in this repo)
      name=$(echo "$base" | sed 's/^list_resource[-_]//')
      patterns+=("Test$(to_pascal "$name")ListResource")
      ;;
    action_*)      patterns+=("TestAction$(to_pascal "${base#action_}")")         ;;
    *)
      # helper_*, provider*, function_*, type_*, or anything else shared
      # → can't isolate, run everything
      echo "."
      exit 0
      ;;
  esac
done

if [[ ${#patterns[@]} -eq 0 ]]; then
  echo "."
  exit 0
fi

# Deduplicate and join with | for use as a go test -run regex
printf '%s\n' "${patterns[@]}" | sort -u | tr '\n' '|' | sed 's/|$//'
