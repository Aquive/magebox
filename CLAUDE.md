# MageBox

MageBox is a Magento 2 development environment using native PHP-FPM/Nginx with Docker for stateful services only. Built in Go with Cobra CLI framework.

## Build & Test

```bash
make build          # Build binary for current platform
make lint           # Run golangci-lint (required before committing)
make test           # Run tests (go test -v ./...)
make fmt            # Format code (go fmt + goimports)
make build-all      # Cross-compile for darwin/linux amd64/arm64
```

Single package test: `go test ./internal/config/... -v`

## Pre-commit requirements

Always run `make lint` and `make test` before committing Go code. Both must pass.

## Project structure

- `cmd/magebox/` - CLI commands (one file per command, uses Cobra)
- `internal/` - Private packages (config, platform, php, nginx, ssl, docker, dns, varnish, project, etc.)
- `lib/templates/` - Embedded config/Docker templates
- `vitepress/` - Documentation site (magebox.dev)
- `docs/` - Internal design docs
- `VERSION` - Single-line semver, drives auto-tag CI workflow
- `CHANGELOG.md` - Keep a Changelog format

## Release process

1. Update `VERSION` file with new semver
2. Add release entry to `CHANGELOG.md` (Keep a Changelog format, date as YYYY-MM-DD)
3. Commit and push to main
4. CI auto-tags from VERSION file, then builds/releases binaries and updates Homebrew tap

## Code conventions

- Go 1.24, CGO_ENABLED=0
- Return errors, don't panic
- Table-driven tests preferred
- Linters: errcheck, gosimple, govet, ineffassign, staticcheck, unused, gofmt, goimports, misspell
- Commit messages: imperative mood, reference PR numbers with `(#XX)`
