#!/bin/bash

showHelp() {
	cat << EOF

Usage:  ./fixtureDriftDetect.sh (file|dir|prefix) [--revision <rev>] [--regex <expr>]
    -h, -help, --help       Display help.
    -revision, --revision   Specify the revision to compare with, default: 'origin/main'
    -regex,    --regex      Specify regular expression for lines to compare, default: '(id: |url: |status: )'

Dependencies:
    - diff
    - git
    - grep

EOF
}

revision='origin/main'
regex='(id: |url: |status: )'

while [[ $# -gt 0 ]]; do
	case $1 in
		-h|--help)
			showHelp
			exit 0
			;;
		-revision|--revision)
			revision="$2"
			shift
			shift
			;;
		-regex|--regex)
			regex="$2"
			shift
			shift
			;;
		--*|-*)
			>&2 echo "unknown option $1"
			exit 1
			;;
		*)
      if [[ -n "$pathspec" ]]; then
  			>&2 echo "unknown positional argument $1"
	  		exit 1
      fi
      pathspec=$1
      shift
			;;
	esac
done

if [[ -z "$pathspec" ]]; then
  >&2 echo "missing positional argument (file|dir|prefix)"
	exit 1
fi

for file in $pathspec*; do
  if [[ -f "$file" ]]; then
	  found="true"
	  git cat-file -e $revision:$file 2> /dev/null
	  if [[ $? -gt 0 ]]; then
		  printf "ignoring file '$file' since it does not exist in revision '$revision'\n"
		  continue
		fi
    diff=$(diff <(git cat-file blob $revision:$file | grep -E "$regex") <(cat $file | grep -E "$regex"))
		if [[ -n "$diff" ]]; then
		  printf "\n\nChanges in file '$file' to revision '$revision' with regards to regex '$regex':\n\n$diff\n\n"
			failed="true"
		fi
  fi 
done

if [[ -z "$found" ]]; then
  printf "\nWARNING: no files processed, check pathspec '$pathspec' (e.g. add trailing slash for directories)\n\n"
  exit 1
fi

if [[ -n "$failed" ]]; then
  printf "\nWARNING: Changes to revision '$revision' detected with regards to regex '$regex' (see above)\n\n"
  exit 1
fi

printf "\nSUCCESS: No changes to revision '$revision' detected with regards to regex '$regex':\n\n"
exit 0
