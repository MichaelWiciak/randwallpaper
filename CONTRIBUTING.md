# Contributing

Thanks for your interest in contributing to `randwallpaper`! This project is small and developing. Bug reports, feature requests, and pull requests are all very welcome.

## Getting started

```bash
# Run the test suite
go test ./...

# Build the CLI binary
go build ./cmd/randwallpaper

# Static analysis
gofmt -l .
go vet ./...
go run honnef.co/go/tools/cmd/staticcheck@v0.8.2 ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
```

## Development workflow

1. Create a feature branch from `main`.
2. Make your changes and add tests where appropriate.
3. Run `gofmt`, `go vet ./...`, and `go test -race ./...` locally to make sure everything passes.
4. Open a pull request against `main`.
5. CI runs `gofmt`, `go vet`, `staticcheck`, `govulncheck`, and `go test -race` on Linux, macOS, and Windows, plus the latest stable Go release.

## Commit messages

Releases are **automated by [Release Please](https://github.com/googleapis/release-please)**.
It reads your PR title/commits to decide the next version, so please use
[Conventional Commits](https://www.conventionalcommits.org/):

| Prefix       | Version effect | Examples                      |
| ------------ | -------------- | ----------------------------- |
| `feat!:`     | Major          | `feat!: remove WithSeed`      |
| `feat:`      | Minor          | `feat: add -version flag`     |
| `fix:` / `perf:` | Patch      | `fix: clamp colour values`    |
| `chore:` / `docs:` / `style:` / `test:` / `ci:` | No release | `chore: tidy go.mod` |

A `!` after the type (e.g. `feat!:` or `refactor!:`) or a
`BREAKING CHANGE:` footer marks a breaking change and triggers a major release.

## Style

- Keep the public API small: prefer unexported implementation types behind
  the single `Generate` entry point.
- Run `gofmt` on your changes.
- Follow the existing comment style: exported identifiers get doc comments
  starting with the identifier name.
