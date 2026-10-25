# AGENTS.md

## Purpose and scope

Use this file as the working agreement for AI agents contributing to this repository. These rules are intended to apply across Go projects, libraries, services, and applications; do not assume the repository is a particular package or module. Follow the actual codebase and its established conventions where they provide more specific, verified requirements.

Treat implementation, tests, and executable configuration as the source of truth for behavior and tooling. Inspect `go.mod`, repository configuration, CI workflows, tests, and existing documentation before making assumptions. Never copy project-specific invariants, module paths, test names, tool versions, coverage thresholds, or commands from another project without verifying them here.

## Primary responsibility: documentation quality

Clear, accurate documentation is a first-class deliverable. When adding or changing behavior, review all documentation surfaces that may be affected: declaration comments, package documentation in `doc.go`, `README.md`, and runnable examples. Update the surfaces whose content or promises changed; do not mechanically modify every file for an implementation-only change.

Write source comments, Go documentation, examples, README content, and other project documentation in clear, idiomatic English unless the repository explicitly establishes another language.

Documentation must describe actual behavior, not intended or assumed behavior. Never invent guarantees, edge cases, complexity claims, thread-safety properties, compatibility promises, or error behavior. If the implementation and tests do not establish a claim, verify it or phrase the uncertainty accurately.

## Documentation at declaration level versus inside function bodies

These two kinds of comments serve different purposes and must not be conflated.

### Declaration comments (the header above a function, method, type, or field)

Explain **what the API is for, why a caller would use it, and which meaningful use cases or input conditions it supports**. Describe the caller-visible contract rather than narrating the implementation.

For exported functions and methods, as applicable:

- Start the comment with the exact declared name, following GoDoc conventions.
- State the purpose and intended use case directly.
- Describe supported modes and meaningful input cases, including the semantics of zero, negative, empty, nil, omitted, or optional values when relevant.
- Explain return values, error conditions, panic conditions, side effects, ownership, mutation, resource lifecycle, and concurrency behavior when these are part of the contract.
- Distinguish absent data from zero-valued or nil data where the API supports that distinction.
- Explain important precedence, fallback, inheritance, override, ordering, or cancellation semantics when callers need to understand them.
- Link to related Go identifiers with GoDoc links such as `[Context]` or `[NewClient]` when those identifiers exist and the links improve navigation.
- Use separate paragraphs for materially different aspects of the contract. Be complete about subtle behavior without padding the comment with generic wording.

Do not explain every internal step in the declaration comment. Callers should be able to understand when and how to use the API without having to read the implementation.

### Comments inside function and method bodies

An inline comment should explain **why the implementation does something this way**, especially when the reason is not obvious from the code. Capture intent, invariants, trade-offs, constraints, or a subtle language/library behavior that a maintainer could otherwise accidentally break.

Good reasons for a body comment include:

- preserving a non-obvious invariant or observable behavior;
- choosing one implementation over a plausible alternative for a concrete reason;
- handling a subtle edge case imposed by the language, runtime, dependency, protocol, or platform;
- ordering operations to preserve correctness, avoid a race, release a resource, or maintain error precedence;
- retaining code that appears redundant but is required for compatibility or correctness.

Do not add comments that merely translate code into prose, such as `// Increment i` above `i++`, or repeat an obvious condition. Prefer clear names and structure over explanatory comments. If a design decision is important to callers, document its observable effect in the declaration comment as well; explain the implementation rationale inside the body only when it adds value.

Private types, helpers, and fields should receive comments when their purpose or invariants are not obvious, or when the project convention requires them. Public comments should not expose private implementation details unless doing so helps explain a caller-visible guarantee.

## Local Tooling

`TOOLS.md` is a concise, practical guide to the development tools used by the project. It helps developers and AI agents install the required tools locally and understand their basic usage.

### Requirements

- Inspect the project configuration and `.github/workflows/` to identify the development tools actually used by the project.
- Document the tools needed to develop, test, lint, format, build, or maintain the project.
- Include tools installed or invoked by CI workflows when they are relevant to local development.
- Do not document every GitHub Action, CI service, upload action, or GitHub-hosted automation. Include only actual tools that developers may need to install or use.
- Keep the document concise. Do not include a complete CI inventory, workflow analysis, coverage reports, version audits, or lengthy explanations of CI internals.
- Do not duplicate instructions already covered by `AGENTS.md`.

### Information for Each Tool

For each relevant tool, document:

- **Purpose:** What the tool is used for.
- **Version:** The required or pinned version, when applicable.
- **Installation:** The command needed to install it. Project-managed CLI tools must be installed into `./bin/`.
- **Basic usage:** The most useful commands developers can run in this repository.
- **Verification:** How to check the installed version, when applicable.

Include the language toolchain when it is a prerequisite for working with the project. Explain how to obtain it and identify the required version.

Use simple, practical examples that can be copied and executed from the repository root.

### Maintenance

- Keep `TOOLS.md` consistent with the actual project configuration and relevant CI workflows.
- When adding or removing a development tool, update the document accordingly.
- Use exact version pins when the project specifies them. Do not invent versions or installation commands.
- If a tool is provided only by CI and has no relevant local usage, omit it.
- Do not add separate sections listing CI-only actions, version inconsistencies, or CI implementation details unless they directly affect a developer's ability to use the project.
- Prefer a short, useful guide over exhaustive documentation.

## Package documentation (`doc.go`)

For packages with meaningful public APIs, maintain package-level documentation in `doc.go`. It must explain the package as a coherent whole rather than duplicate every declaration comment. The package comment should start with `Package <name> ...`, in accordance with GoDoc conventions.

Include the following sections or clearly identifiable topics. Headings are encouraged for substantial packages; use Go documentation heading syntax such as `// # Overview`. Keep each section concise but substantive. If a topic genuinely does not apply to a package, do not invent content just to fill it; cover the applicable topics and make important limitations explicit.

1. **Overview and purpose** — What problem the package solves, what it provides, and when it is appropriate to use it.
2. **Core concepts and usage model** — The principal types or abstractions and how they fit together. Explain the user's mental model, not private implementation mechanics.
3. **Basic usage** — A small, idiomatic usage path or code example when it helps readers understand the API. Keep larger or scenario-specific examples in `example*_test.go` and refer readers to them where useful.
4. **Behavior and guarantees** — Package-wide semantics that affect multiple APIs, such as lifecycle, inheritance, precedence, ownership, mutation, concurrency, error handling, cancellation, or resource management, when applicable.
5. **Limitations and important caveats** — Constraints, unsupported scenarios, security considerations, assumptions, or cases in which another package or approach is more appropriate, when applicable.
6. **API and further documentation** — Point readers to the main entry-point identifiers, relevant examples, or related packages using valid GoDoc links where useful.

A suggested structure for a substantial package is:

```go
// Package example provides the functionality described below.
//
// # Overview
// Explain the problem and intended use.
//
// # Core concepts
// Describe the main abstractions and their relationship.
//
// # Basic usage
// Show a compact usage path or direct readers to an executable example.
//
// # Behavior and guarantees
// Describe package-wide behavior callers can rely on.
//
// # Limitations
// Describe relevant caveats and unsupported cases.
//
// # See also
// Link to key APIs and related documentation.
package example
```

Treat this as a content checklist, not a requirement to create empty headings in a tiny package. The first package comment and all package-wide claims must remain accurate. Avoid duplicating long function comments, embedding a full README into `doc.go`, or documenting behavior that is only an incidental implementation detail. Inspect the rendered result with `go doc` after material changes.

## README (`README.md`)

Write the README for developers evaluating, adopting, installing, or operating the project. Adapt the structure to the type and maturity of the repository. For a reusable library, it will normally cover:

- a clear summary of what the project does and the problem it addresses;
- the main use cases and, where useful, when not to use it;
- installation/import instructions using the real module path from `go.mod`;
- a quick-start example that uses the actual public API;
- important behavior, guarantees, constraints, and caveats a user needs to know;
- links to detailed API docs, package docs, examples, configuration, or further guides;
- project-specific setup, test, build, or run instructions when relevant.

For applications and services, prioritize prerequisites, configuration, local development, operation, deployment, and troubleshooting as appropriate. Do not force library-specific sections onto other project types.

Keep the README user-oriented. It complements `doc.go` rather than copying all of it or becoming a dump of declaration comments. Ensure code snippets compile where practical and do not drift from the real API. Prefer linking to executable examples for details that can be tested automatically.

## Runnable examples and example test files

Examples are part of the public documentation and, where configured, executable tests. They must be idiomatic, realistic, focused, and consistent with the public API.

- Use Go's standard example test convention and place examples in `_test.go` files. A package-level example can use `func Example()`; use named examples such as `func ExampleClient_Do()` or `func ExampleParse()` when they demonstrate a particular identifier and follow Go's example naming rules.
- Multiple example files are allowed and encouraged when they improve organization. Use names such as `example_test.go`, `example_usage_test.go`, `example_errors_test.go`, `example_client_test.go`, or other descriptive `example_**_test.go` names. Choose names by topic, feature, workflow, or API area rather than putting every example in one large file.
- Prefer examples in an external test package (for example, `package example_test`) when the goal is to demonstrate the public API as a user would consume it. Use the internal package only when access to unexported details is itself necessary for a legitimate test purpose; do not expose private implementation details in user-facing examples.
- Include exact `// Output:` comments for examples whose output is deterministic and should be checked. Use `// Unordered output:` only when ordering is genuinely irrelevant and the Go testing convention supports it. Avoid nondeterministic output, network calls, wall-clock dependence, and environmental assumptions unless demonstrating those behaviors is the point and the example is made reliable.
- Keep examples small enough to understand, but complete enough to show setup, the important call, and the resulting behavior. Use realistic names and error handling; do not hide relevant failures merely to shorten a snippet.
- Organize the package-level overview example separately from topic-specific examples when that makes the documentation easier to navigate. Package-level and API-specific examples can live across multiple `example_**_test.go` files.
- Run `go test ./...` and, where useful, `go test -run Example ./...`. Confirm that example output comments match actual output. Do not claim examples ran unless the commands were executed successfully.

## Consistency across documentation surfaces

Each documentation surface has a distinct job:

- Declaration comments explain why an API exists, how it should be used, and which cases its contract covers.
- Comments inside function bodies explain why a non-obvious implementation choice exists and what invariant it preserves.
- `doc.go` explains the package model, common semantics, usage path, guarantees, and limitations.
- `README.md` helps a developer evaluate, adopt, install, configure, or operate the project.
- `example_**_test.go` files demonstrate concrete workflows and keep examples executable where possible.

These surfaces should reinforce one another without repeating large blocks of text. When updating a project:

1. Inspect the relevant implementation, tests, public API, and existing docs first.
2. Identify the actual contract and the questions a user or maintainer needs answered.
3. Update each affected documentation surface and add or adjust examples when they clarify important use cases.
4. Run formatting and relevant tests; inspect `go doc` output after material GoDoc changes.
5. Review all claims for accuracy, consistency, and useful detail.

Do not change implementation behavior only to make a draft sound compelling. If intended behavior is unclear or source and tests conflict, investigate and report the ambiguity instead of making up a contract.

## Documentation review checklist

Before reporting documentation work as complete, verify all four points:

1. **Structure.** `doc.go` covers the applicable sections from the package documentation section of this file and uses their headings; `Basic usage` is present whenever the package has a concrete usage path.
2. **Contract versus rationale.** Declaration comments state purpose, caller-visible behavior, and edge cases rather than narrating implementation. Body comments explain only non-obvious implementation decisions, invariants, and the problems that alternatives would cause.
3. **Verified guarantees.** Every claim about thread safety, copying, ownership, concurrency, error handling, or performance is checked against the implementation and tests, not inferred from names or intent. If the evidence is partial, narrow the claim instead of presenting it as a guarantee.
4. **Runnable examples.** Any statement that examples are runnable or checked is verified against actual `Example...` functions with `// Output:` comments, and those examples are executed successfully. A code block in `doc.go` or `README.md` is not a runnable example by itself.

## Go implementation and API design

Write idiomatic, formatted Go. Prefer explicit contracts, clear control flow, composition, and maintainable abstractions over clever shortcuts.

### Design for intentional flexibility

Design functions and methods to support meaningful, foreseeable usage modes instead of unnecessarily restricting callers to one narrow path. Identify valid states and define their behavior explicitly. Flexibility must be intentional: do not add speculative parameters, options, abstractions, or branches without a concrete use case.

A useful design pattern, when it fits the API, is to validate required dependencies first, validate input ranges explicitly, distinguish meaningful zero/negative/positive values, derive and cancel a child context when a timeout is requested, otherwise respect the caller's context when that is part of the contract, and propagate the underlying operation's error. This is a pattern, not a universal rule; define the semantics appropriate to each API and document them.

### Correctness and maintainability

- Validate inputs early and consistently. Keep validation order deterministic when it affects which error is returned.
- Prefer explicit dependency injection and caller-controlled contexts over package globals and hidden environmental behavior.
- For request-scoped or cancellable I/O, propagate the caller's context. Avoid replacing it with `context.Background()` or `context.TODO()` inside ordinary helper functions unless detaching from cancellation is deliberate and documented.
- Handle errors explicitly. Wrap errors with `%w` when callers should be able to inspect the underlying error. Follow established project conventions for sentinel or typed errors.
- Do not use `panic` for routine runtime failures. Reserve it for deliberately violated programmer contracts where appropriate and document such behavior.
- Release resources reliably, including cancellation functions, response bodies, locks, and temporary resources, following the relevant API conventions.
- Avoid hidden mutation, unnecessary global state, and needless coupling. Document ownership, mutation, and concurrency constraints when they are part of the contract.
- Preserve compatibility with the Go version and dependencies established by the repository.
- Test meaningful boundary conditions, error paths, and alternate supported modes. Tests must assert behavior, not merely execute lines to increase coverage.
- Do not claim performance, allocation, thread-safety, or complexity properties without evidence.

## Repository inspection, commands, and verification

Before changing code, inspect the repository structure, current branch, working tree, `go.mod`, relevant source files, tests, project configuration, and CI workflows. Preserve unrelated user changes.

Use the commands established by the repository. Typical Go checks may include:

```sh
gofmt -w <changed-go-files>
go test ./...
go test -race ./...
go vet ./...
golangci-lint run
go mod tidy
```

These are examples, not commands that must blindly be run in every repository. Check whether the required tools are installed, follow pinned versions and project scripts, and consult CI for authoritative gates. Run `go mod tidy` when relevant and check whether it changes `go.mod` or `go.sum`; do not discard changes without review. Apply coverage thresholds only when the repository actually enforces them. If a repository requires exact coverage, test each meaningful branch and assert its behavior rather than writing superficial tests.

Before reporting completion, list the checks that actually ran and whether they passed. Never claim a command passed if it was not executed successfully. If a check could not run, state which check and why.

## Local tools and installation commands

Keep repository-local downloaded tools and their binaries in the project-level `./bin/` directory. Do not install project tools globally by default. After a tool is installed for the repository, invoke it through its local path (for example, `./bin/golangci-lint run`), not through a global executable or a bare command that might resolve to another version. Use `./bin/` for all repository-managed tools and commands whenever a local tool is available.

Keep every tool installation command in one file named exactly `tmp/commands`, with no file extension. Do not create separate per-tool installation scripts or a `tmp/commands/` directory. Add each installation command to this single file, using `./bin/` as the installation destination. Keep the commands readable and reproducible; use shell comments to identify a tool or explain a non-obvious argument where needed.

- Inspect CI and repository configuration first, and install the version pinned or expected by the project rather than an unpinned latest release.
- Use the installer option that writes to `./bin/` whenever available.
- After installation, run tools using paths such as `./bin/<tool>` and verify the installed version when practical.
- Preserve existing ignore rules and repository conventions. Do not modify `.gitignore` merely to make a generated binary trackable; downloaded local binaries should normally remain ignored unless the project explicitly requires otherwise.
- If a tool cannot reasonably be installed or run from `./bin/`, explain the constraint and follow the repository's existing convention instead of silently falling back to a different global version.

For example, the single file `tmp/commands` may contain this command when the repository requires `golangci-lint` v2.15.0:

```sh
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ./bin v2.15.0
```

Then run the installed tool as `./bin/golangci-lint run`. The version is an example only; use the version required by the current repository.

## Branch rules

Create a separate branch for every new feature. Do not combine independent features into one branch. Use a short, descriptive, lowercase kebab-case name with an appropriate prefix, for example:

- `feat/add-request-retries`
- `fix/handle-empty-configuration`
- `docs/improve-package-guide`
- `refactor/simplify-client-lifecycle`
- `test/cover-timeout-edge-cases`
- `chore/update-tooling`

Use `feat/` for features, `fix/` for bug fixes, and matching prefixes for documentation, refactoring, tests, or maintenance. Follow any more specific branch naming policy already established by the repository. Before switching branches, inspect the working tree and preserve existing changes. Do not delete branches or discard work unless explicitly authorized.

## Commit rules

**Never create a commit until the user has reviewed the changes and explicitly confirmed that all changes are approved.** A request to implement, test, document, or review code is not permission to commit. Do not push unless the user explicitly authorizes pushing.

When commits are authorized:

- Each commit must represent one complete, coherent, logically related unit of work.
- Follow Conventional Commits.
- Write the commit message entirely in English.
- Use a single-line header only: no body, footer, scope in parentheses, blank continuation lines, or additional metadata.
- Do not add a period at the end of the description.

Correct examples:

```text
feat: add request retry support
fix: handle empty configuration
docs: improve package documentation
test: cover timeout edge cases
refactor: simplify client lifecycle
```


## Expected workflow

1. Inspect the repository, current branch and working tree, source, tests, documentation, and relevant configuration.
2. For a new feature, create or use a dedicated branch with an appropriate name, preserving existing changes.
3. Implement the requested work with maintainability and the actual API contract in mind.
4. Update declaration comments, `doc.go`, `README.md`, and relevant `example_**_test.go` files whenever their content is affected.
5. Format changed Go files and run the relevant tests, examples, linters, and project-specific checks.
6. Summarize the changes and report actual verification results. Wait for the user's explicit approval of all changes before creating any commit.
