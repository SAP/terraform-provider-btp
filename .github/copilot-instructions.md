# Project Overview

The repository implements the official (open source) Terraform Provider for SAP Business Technology Platform (SAP BTP). It enables infrastructure-as-code management of SAP BTP.

Primary goals:
- Provide comprehensive, consistent CRUD / data source coverage of SAP BTP platform entities.
- Model SAP BTP primitives as Terraform resources and data sources with clear schemas, validation, and lifecycle handling.
- Offer predictable diffs and idempotent behavior; minimize drift and surface configuration errors early.
- Maintain high test coverage (unit + acceptance; VCR where feasible) and documentation parity (generated docs + examples).
- Adhere to HashiCorp Terraform Plugin Framework best practices (framework-based types, plan modifiers, validators, timeouts, proper diagnostics).

## Folder Structure

- `main.go`: Provider entrypoint (serves the Terraform provider server executable).
- `go.mod` / `go.sum`: Module definition; Go 1.25; declares HashiCorp plugin framework + supporting libs.
- `Makefile`: Standard developer commands (build, install, lint, generate, fmt, test, testacc).
- `README.md`: User-facing overview and badges.
- `DEVELOPER.md`: Local / container / Codespaces setup and dev overrides, committing guidance.
- `CONTRIBUTING.md`: Contribution process, DCO, AI usage guidelines.
- `sonar-project.properties`: SonarCloud static analysis configuration.
- `terraform-registry-manifest.json`: Registry publishing metadata.
- `examples/`: Practical Terraform configurations showcasing provider usage (provider, data-sources, resources subfolders).
- `docs/`: Generated provider documentation (data sources and resources). Do not hand-edit generated docs (use code + templates then `make generate`).
- `guides/`: Manually maintained topic guides (drift detection, import, sensitive data, etc.).
- `templates/`: Go text/template sources for generating documentation pages.
- `integration/`: Integration test/example Terraform configs (manual / scripted validation across resources).
- `regression-test/`: Additional Terraform scenarios for regression detection.
- `internal/`: Internal (non-exported) Go packages:
  - `btpcli/`: API client / facade for SAP BTP REST endpoints (abstraction layer over HTTP).
  - `tfutils/`: Shared Terraform + Go helper utilities (conversions, plan helpers, etc.).
  - `validation/`: Centralized schema validation helpers.
  - `version/`: Provider versioning utilities.
- `btp/provider/`: Core provider implementation code:
  - Resource files: `resource_<scope>_<entity>.go`
  - Data source files: `datasource_<scope>_<entity>.go`
  - Type models: `type_<scope>_<entity>.go` / variations for hierarchy models.
  - Tests: Paired `_test.go` files next to the implementation (unit and acceptance style). Data sources and resources each have a corresponding test file.
  - Helpers: `helper.go`, `helper_doc_values_format.go` for documentation value formatting and shared logic.
  - `provider.go` / `provider_test.go`: Provider schema and configuration implementation and tests.
  - `fixtures/`: (If present) serialized recordings / static payloads aiding tests (e.g., with go-vcr).

Naming scheme highlights inside `btp/provider/`:
- `datasource_<account|directory|subaccount|...>_<entity>.go`: Data source implementation.
- `resource_<account|directory|subaccount|...>_<entity>.go`: Resource implementation.
- `type_<scope><Entity>.go` or `type_<scope>_<entity>.go`: Strongly typed internal model or structure representing API data (used in schema translation).

## Libraries and Frameworks

Runtime and Framework:
- HashiCorp Terraform Plugin Framework (`github.com/hashicorp/terraform-plugin-framework`): Core abstraction for provider schema and CRUD.
- HashiCorp Terraform Plugin Go (`github.com/hashicorp/terraform-plugin-go`): Underlying protocol and plumbing.
- Framework Add-ons:
  - `terraform-plugin-framework-validators`: Attribute validators.
  - `terraform-plugin-framework-timeouts`: Declarative operation timeouts.

Testing and Quality:
- terraform-plugin-testing: Acceptance and unit test helpers.
- stretchr/testify: Assertions.
- go-vcr (`gopkg.in/dnaeon/go-vcr.v3`): Recording/replaying HTTP interactions (when applicable) to reduce live API dependency.
- golangci-lint: Aggregated linting.

Auxiliary:
- internal/btpcli: SAP BTP API client abstraction (respect API versioning and error mapping here first before resource layer changes).
- internal/validation: Central location for cross-schema validators.
- internal/version: Provider version injection (used in User-Agent strings, etc.).

## Coding Standards

### General Instructions

- Write simple, clear, and idiomatic Go code
- Favor clarity and simplicity over cleverness
- Follow the principle of least surprise
- Keep the happy path left-aligned (minimize indentation)
- Return early to reduce nesting
- Prefer early return over if-else chains; use if condition { return } pattern to avoid else blocks
- Make the zero value useful
- Write self-documenting code with clear, descriptive names
- Document exported types, functions, methods, and packages
- Use Go modules for dependency management
- Leverage the Go standard library instead of reinventing the wheel (e.g., use strings.Builder for string concatenation, filepath.Join for path construction)
- Prefer standard library solutions over custom implementations when functionality exists
- Write comments in English by default; translate only upon user request
- Avoid using emoji in code and comments
- Module path: `github.com/SAP/terraform-provider-btp`.
- Use modern Go language features where it is possible
- Adhere to best practises when using the Terraform Plugin Framework

### Naming Conventions

#### Packages

- Use lowercase, single-word package names
- Avoid underscores, hyphens, or mixedCaps
- Choose names that describe what the package provides, not what it contains
- Avoid generic names like `util`, `common`, or `base`
- Package names should be singular, not plural

##### Package Declaration Rules (CRITICAL):
- **NEVER duplicate `package` declarations** - each Go file must have exactly ONE `package` line
- When editing an existing `.go` file:
  - **PRESERVE** the existing `package` declaration - do not add another one
  - If you need to replace the entire file content, start with the existing package name
- When creating a new `.go` file:
  - **BEFORE writing any code**, check what package name other `.go` files in the same directory use
  - Use the SAME package name as existing files in that directory
  - If it's a new directory, use the directory name as the package name
  - Write **exactly one** `package <name>` line at the very top of the file
- When using file creation or replacement tools:
  - **ALWAYS verify** the target file doesn't already have a `package` declaration before adding one
  - If replacing file content, include only ONE `package` declaration in the new content
  - **NEVER** create files with multiple `package` lines or duplicate declarations

#### Filenames

File names must follow the conventions below to ensure consistency and clarity:

- Resource filenames: `resource_<scope>_<entity>.go` (e.g., `resource_subaccount_service_instance.go`).
- Data source filenames: `datasource_<scope>_<entity>.go` (e.g., `datasource_globalaccount_entitlements.go`).
- Test files mirror implementation: append `_test.go` (e.g., `datasource_subaccount_service_instance_test.go`).
- Types: `type_<scope>_<entity>.go` or camel-case variant where hierarchical (e.g., `type_directoryHierarchy.go`). Keep internal naming consistent; prefer lower-case file-local helpers.

#### Variables and Functions

- Use mixedCaps or MixedCaps (camelCase) rather than underscores
- Keep names short but descriptive
- Use single-letter variables only for very short scopes (like loop indices)
- Exported names start with a capital letter
- Unexported names start with a lowercase letter
- Avoid stuttering (e.g., avoid `http.HTTPServer`, prefer `http.Server`)

#### Interfaces

- Name interfaces with -er suffix when possible (e.g., `Reader`, `Writer`, `Formatter`)
- Single-method interfaces should be named after the method (e.g., `Read` → `Reader`)
- Keep interfaces small and focused

#### Constants

- Use MixedCaps for exported constants
- Use mixedCaps for unexported constants
- Group related constants using `const` blocks
- Consider using typed constants for better type safety

### Error Handling Patterns

#### Creating Errors

- Use `errors.New` for simple static errors
- Use `fmt.Errorf` for dynamic errors
- Create custom error types for domain-specific errors
- Export error variables for sentinel errors
- Use `errors.Is` and `errors.As` for error checking

#### Error Propagation

- Add context when propagating errors up the stack
- Don't log and return errors (choose one)
- Handle errors at the appropriate level
- Consider using structured errors for better debugging

### Code Style and Formatting

#### Formatting

- Always use `gofmt -s` or `make fmt` to format the code.
- Run `golangci-lint run` locally or via pre-commit (hook) before pushing.
- Avoid unnecessary exported identifiers—keep scope minimal (internal packages or unexported symbols where possible).
- Use `goimports` to manage imports automatically
- Keep line length reasonable (no hard limit, but consider readability)
- Add blank lines to separate logical groups of code

#### Comments

- Strive for self-documenting code; prefer clear variable names, function names, and code structure over comments
- Write comments only when necessary to explain complex logic, business rules, or non-obvious behavior
- Write comments in complete sentences in English by default
- Translate comments to other languages only upon specific user request
- Start sentences with the name of the thing being described
- Package comments should start with "Package [name]"
- Use line comments (`//`) for most comments
- Use block comments (`/* */`) sparingly, mainly for package documentation
- Document why, not what, unless the what is complex
- Avoid emoji in comments and code

#### Error Handling

- Check errors immediately after the function call
- Don't ignore errors using `_` unless you have a good reason (document why)
- Wrap errors with context using `fmt.Errorf` with `%w` verb
- Create custom error types when you need to check for specific errors
- Place error returns as the last return value
- Name error variables `err`
- Keep error messages lowercase and don't end with punctuation
- Always return diagnostics via the Terraform Plugin Framework patterns (`resp.Diagnostics.Append(...)`). Avoid panics in provider logic.
- Wrap lower-level errors with context (use `%w` to preserve chain). Avoid leaking secrets in error strings.
- Be concise, actionable, and avoid exposure of raw internal response bodies unless needed. Prefer: `failed to create subaccount service instance: <reason>`.

### Architecture and Project Structure

#### Package Organization

- Follow standard Go project layout conventions
- Keep `main` packages in `cmd/` directory
- Put reusable packages in `pkg/` or `internal/`
- Use `internal/` for packages that shouldn't be imported by external projects
- Group related functionality into packages
- Avoid circular dependencies

#### Dependency Management

- Use Go modules (`go.mod` and `go.sum`)
- Keep dependencies minimal
- Regularly update dependencies for security patches
- Use `go mod tidy` to clean up unused dependencies
- Vendor dependencies only when necessary

### Type Safety and Language Features

#### Type Definitions

- Define types to add meaning and type safety
- Use struct tags for JSON, XML, database mappings
- Prefer explicit type conversions
- Use type assertions carefully and check the second return value
- Prefer generics over unconstrained types; when an unconstrained type is truly needed, use the predeclared alias `any` instead of `interface{}` (Go 1.18+)

#### Pointers vs Values

- Use pointer receivers for large structs or when you need to modify the receiver
- Use value receivers for small structs and when immutability is desired
- Use pointer parameters when you need to modify the argument or for large structs
- Use value parameters for small structs and when you want to prevent modification
- Be consistent within a type's method set
- Consider the zero value when choosing pointer vs value receivers

#### Interfaces and Composition

- Accept interfaces, return concrete types
- Keep interfaces small (1-3 methods is ideal)
- Use embedding for composition
- Define interfaces close to where they're used, not where they're implemented
- Don't export interfaces unless necessary

### Terraform Schema and Resource Modeling

#### Schema Design

- Validate user inputs early with framework validators (`framework-validators` module) and custom functions in `internal/validation`.
- Use explicit attribute types (framework `types.String`, `types.Int64`, etc.) and plan modifiers for computed or optional attributes.
- Mark sensitive attributes with `Sensitive: true` where they must not appear in plans or state outputs.
- Prefer nested attribute blocks for structured data rather than encoding JSON strings.

#### Resource Design

- Start new resource/data source by copying an analogous existing file (closest scope and complexity) and adjust names to maintain pattern consistency.
- Include comprehensive schema: types, required/optional/computed flags, validators, plan modifiers, timeouts (if long operations), and description comments (used for docs generation).
- For CRUD functions (Create / Read / Update / Delete): centralize API logic through `internal/btpcli` to keep resource files declarative; avoid raw `http` usage directly in resource files.
- Always add or update a corresponding `_test.go` file; include at least a happy path and one error or edge case (invalid input / missing attribute / permission denial simulation if feasible) as well as an import.
- Keep attribute names stable

## Security and Secrets:

- No hardcoded credentials. Read from environment variables or configuration.
- Redact sensitive fields in logs and diagnostics and test recordings.
- Keep dependency updates reasonable; prefer minimal surface area.

## Documentation Generation

- Do not manually edit generated `docs/` resource or data source pages—modify code comments and templates under `templates/` then run `make generate`.

## Development and Build Workflow

We usually do not use plain Go commands directly; instead, we leverage the provided Makefile targets for consistency. The main commands are:

- `make fmt`: Format code.
- `make lint`: Run linters.
- `make build`: Compile provider.
- `make install`: Install binary into `$GOBIN` for Terraform dev override usage.

The commands to generate the documentation artifacts:

- `make generate`: Generate documentation artifacts from code annotations/templates.

The commands for testing that are usually nit run manually:

- `make test`: Run unit + (tagged) tests with coverage.
- `make testacc`: Run extended acceptance tests (long-running).
