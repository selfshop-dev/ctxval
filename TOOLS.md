# Project Tools

Tools used to develop, test, lint, and check this module locally. All commands
are run from the repository root. Project-managed CLI tools are installed into
`./bin/`, which is not tracked by Git.

## Go toolchain

The compiler and the standard tooling (`go test`, `go vet`, `go mod`) are
prerequisites for every other tool below.

**Version**

`1.26.10`, as declared by the `go` directive in `go.mod`. Newer Go releases
also satisfy this requirement.

**Installation**

Install Go 1.26.10 from <https://go.dev/dl/> or with a Go version manager, and
make sure `go` is on `PATH`. The toolchain is not installed into `./bin/`.

**Verification**

```sh
go version
```

The reported version must satisfy the requirement in `go.mod`.

**Usage**

```sh
go vet ./...
go test ./...
go test -race ./...
```

Run `go mod tidy` when dependencies or module metadata change. Formatting is
handled by `./bin/golangci-lint fmt`, which uses the project configuration.

The race detector requires cgo (`CGO_ENABLED=1`) and a C compiler such as gcc.
Without a C compiler, `go test -race ./...` cannot run.

## golangci-lint

Runs the configured linters and formatters over the Go sources. Configuration
lives in `.golangci.yml`.

**Version**

`v2.15.0`.

**Installation**

```sh
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ./bin v2.15.0
```

**Verification**

```sh
./bin/golangci-lint --version
./bin/golangci-lint config verify
```

**Usage**

```sh
./bin/golangci-lint fmt
./bin/golangci-lint run ./...
```

`fmt` rewrites files in place using the formatters configured in `.golangci.yml`
and is the project's formatting command; run it before committing Go changes.
`run` reports lint issues without changing files.

## govulncheck

Reports known vulnerabilities that affect this module and its dependencies.

**Version**

`v1.1.4`.

**Installation**

`GOBIN` must be an absolute path:

```sh
GOBIN="$PWD/bin" go install golang.org/x/vuln/cmd/govulncheck@v1.1.4
```

**Verification**

```sh
./bin/govulncheck -version
```

The output includes `Scanner: govulncheck@v1.1.4`.

**Usage**

```sh
./bin/govulncheck ./...
```

## trivy

Scans the module dependencies for known vulnerabilities. CI runs the same
scanner through `aquasecurity/trivy-action`.

**Version**

`v0.76.1`.

**Installation**

```sh
curl -sSfL https://github.com/aquasecurity/trivy/releases/download/v0.76.1/trivy_0.76.1_Linux-64bit.tar.gz | tar xz -C ./bin trivy
```

**Verification**

```sh
./bin/trivy --version
```

**Usage**

```sh
./bin/trivy fs --scanners vuln --severity HIGH,CRITICAL --vuln-type library --exit-code 1 .
```

## actionlint

Checks the GitHub Actions workflow files in `.github/workflows/` for syntax and
expression errors. It is a local check for workflow changes and is not part of
CI.

**Version**

`v1.7.12`.

**Installation**

```sh
GOBIN="$PWD/bin" go install github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
```

**Verification**

```sh
./bin/actionlint -version
```

**Usage**

```sh
./bin/actionlint
```
