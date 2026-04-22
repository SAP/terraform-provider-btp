# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the official Terraform Provider for SAP Business Technology Platform (SAP BTP). It uses the HashiCorp Terraform Plugin Framework to expose SAP BTP resources, data sources, list resources, provider functions, and actions.

## Essential Commands

Build and development:
- `make fmt` - Format code with gofmt
- `make fix` - Run go fix to update code to newer Go versions
- `make lint` - Run golangci-lint (must pass before commits)
- `make build` - Compile the provider
- `make install` - Build and install to `$GOBIN` for local Terraform dev override
- `make generate` - Generate documentation from code annotations and templates

**CRITICAL: After every code change, always run in order:**
1. `make lint` - Check for linting issues
2. `make fix` - Apply automatic fixes
3. `make build` - Verify compilation

Testing:
- `make test` - Run unit tests with coverage (tagged tests included)
- `make testacc` - Run acceptance tests (requires `TF_ACC=1`, long-running, needs live BTP credentials)
- `go test -v -run TestResourceSubaccountServiceInstance ./btp/provider/` - Run specific test

Development setup:
- Configure Terraform CLI dev override in `~/.terraformrc` (Mac/Linux) or `%APPDATA%/terraform.rc` (Windows):
  ```hcl
  provider_installation {
    dev_overrides {
      "sap/btp" = "/path/to/go/bin"
    }
    direct {}
  }
  ```
- Do NOT run `terraform init` when using dev overrides
- Verify setup: `cd examples/provider/ && terraform validate`

Pre-commit hooks (via Lefthook):
- `make lefthook` - Install Lefthook and register the pre-commit hooks
- Hooks run automatically on commit: `go fmt`, `golangci-lint --fix`, `terraform fmt`
- Install once after cloning: `make lefthook`

## Architecture

### Code Organization

```
├── main.go                          # Provider server entrypoint
├── btp/provider/                    # Core provider implementation
│   ├── provider.go                  # Provider schema and configuration
│   ├── resource_<scope>_<entity>.go # Resource implementations (CRUD)
│   ├── datasource_<scope>_<entity>.go # Data source implementations (read-only)
│   ├── list_resource_<scope>_<entity>.go # List resources (Terraform 1.14+, read-only lists)
│   ├── function_<name>.go           # Provider functions (pure transformation helpers)
│   ├── action_<name>.go             # Terraform actions (imperative operations)
│   ├── type_<scope>_<entity>.go    # Type models for schema translation
│   ├── *_test.go                   # Tests (paired with implementations)
│   └── helper*.go                  # Shared utilities
├── internal/
│   ├── btpcli/                     # SAP BTP API client facade (HTTP abstraction)
│   ├── tfutils/                    # Terraform + Go helper utilities
│   ├── validation/                 # Centralized schema validators
│   └── version/                    # Provider version utilities
├── docs/                           # Generated documentation (DO NOT EDIT MANUALLY)
├── templates/                      # Doc generation templates
├── examples/                       # Terraform configuration examples
├── tests/
│   ├── integration-test/           # Integration test scenarios
│   └── regression-test/            # Regression test scenarios
└── regression-test/                # Regression test scenarios (legacy location)
```

### Naming Conventions

Files:
- Resources: `resource_<scope>_<entity>.go` (e.g., `resource_subaccount_service_instance.go`)
- Data sources: `datasource_<scope>_<entity>.go` (e.g., `datasource_globalaccount_entitlements.go`)
- List resources: `list_resource_<scope>_<entity>.go` (e.g., `list_resource_subaccount.go`)
- Provider functions: `function_<name>.go` (e.g., `function_extract_cf_api_url.go`)
- Actions: `action_<name>.go` (e.g., `action_restore_subaccount.go`)
- Types: `type_<scope>_<entity>.go` or camelCase for hierarchical types (e.g., `type_directoryHierarchy.go`)
- Tests: Add `_test.go` suffix to match implementation file

Scopes: `globalaccount`, `directory`, `subaccount`

### Key Architectural Patterns

**API Interaction:**
- All BTP API calls go through `internal/btpcli` - avoid raw HTTP in resource/datasource files
- Keep resource files declarative - business logic belongs in btpcli client layer

**Schema Design:**
- Use framework types explicitly: `types.String`, `types.Int64`, etc.
- Validate early with framework validators and `internal/validation`
- Mark sensitive fields with `Sensitive: true`
- Use nested attributes for structured data (avoid JSON strings)

**Resource Implementation:**
- **Resources** (CRUD): Copy an analogous `resource_*.go` file, implement Create/Read/Update/Delete
- **List Resources** (Terraform 1.14+): Copy an analogous `list_resource_*.go` file, implement `list.ListResource` interface
  - List resources are read-only, optimize for listing/filtering entities
  - Must match managed resource TypeName (e.g., `btp_subaccount`)
  - Implement `ListResourceConfigSchema()` and `List()` methods
  - Support filtering via schema attributes
- **Provider Functions**: Copy an analogous `function_*.go` file, implement `function.Function` interface
  - Pure transformation helpers (e.g., parse Cloud Foundry or Kyma environment labels)
  - No side effects; implement `Metadata()`, `Definition()`, and `Run()` methods
- **Actions**: Copy an analogous `action_*.go` file, implement `action.Action` interface
  - Imperative operations that don't map to standard CRUD (e.g., restore a subaccount)
  - Implement `Metadata()`, `Schema()`, `Run()`, and `ConfigValidators()` methods
- Include comprehensive schema with validators, plan modifiers, and timeouts
- Add corresponding `_test.go` with happy path + error cases + import test

**Testing:**
- Uses `terraform-plugin-testing` framework
- VCR (go-vcr) recordings in `fixtures/` reduce live API dependency
- Test naming: `TestResource<Name>` or `TestDataSource<Name>`
- Include import state verification in tests

## Documentation Generation

- **NEVER** manually edit files in `docs/` - they are generated
- Modify code comments (especially schema MarkdownDescription fields) and `templates/` instead
- Run `make generate` to regenerate docs
- Generated docs power the Terraform Registry documentation

## Development Workflow

1. Start with similar existing resource/datasource/list_resource/function/action as template
2. Implement schema with proper types, validators, descriptions
3. Add CRUD/List logic delegating to `internal/btpcli`
4. Write tests in `*_test.go` with VCR fixtures
5. **MANDATORY after every change:**
   - `make lint` - Fix any linting issues
   - `make fix` - Apply automatic fixes
   - `make build` - Verify compilation succeeds
6. Test: `make test`
7. Generate docs: `make generate`
8. Install locally: `make install`
9. Verify with example: `cd examples/provider/ && terraform validate`

## Commit Conventions

Follow [Conventional Commits](https://www.conventionalcommits.org/):
- `feat: add resource for subaccount subscription`
- `fix: handle nil pointer in service instance read`
- `docs: update examples for trust configuration`
- `refactor!: breaking change to schema`
- `feat(btp_subaccount): scoped feature addition`

## Common Pitfalls

1. **Package declarations**: Each Go file has exactly ONE `package` declaration. When editing existing files, preserve the existing package line - never duplicate it.

2. **Dev overrides**: When using local dev overrides, do NOT run `terraform init` - it's unnecessary and will error.

3. **Test failures**: If acceptance tests fail, ensure:
   - `BTP_USERNAME` and `BTP_PASSWORD` env vars are set
   - VCR fixtures exist or test is marked for live API calls
   - Timeout is sufficient for long-running operations

4. **Generated docs**: Changes to `docs/*.md` will be overwritten. Update code comments and run `make generate`.

5. **Error handling**: Always return diagnostics via `resp.Diagnostics.Append()` - never panic in provider code.

6. **Schema stability**: Keep attribute names stable across versions. Use deprecation warnings for schema changes.

## Testing Strategy

- Unit tests: Fast, use VCR recordings where possible
- Acceptance tests: Slower, may require live BTP account
- Integration tests: In `tests/integration-test/` folder, full Terraform scenarios
- Regression tests: In `tests/regression-test/` and `regression-test/` folders

VCR setup in tests:
```go
rec, user := setupVCR(t, "fixtures/resource_subaccount_service_instance.wo_parameters")
defer stopQuietly(rec)
```

## Security Considerations

- No hardcoded credentials - use environment variables
- Mark sensitive attributes with `Sensitive: true` in schema
- Redact sensitive data in logs and VCR recordings
- Keep dependencies updated for security patches
