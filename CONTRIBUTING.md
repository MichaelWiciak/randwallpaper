# Contributing

Thanks for your interest in contributing to `randwallpaper`! This project is small and friendly — bug reports, feature requests, and pull requests are all welcome.

## Getting started

```bash
# Run the test suite
go test ./...

# Build the CLI binary
go build ./cmd/randwallpaper

# Static analysis
go vet ./...

# Optional: stronger linter (also run in CI)
go run honnef.co/go/tools/cmd/staticcheck@latest ./...
```

## Development workflow

1. Create a feature branch from `main`.
2. Make your changes and add tests where appropriate.
3. Run `go test ./...` and `go vet ./...` locally to make sure everything passes.
4. Open a pull request against `main`.
5. CI runs `go vet`, `staticcheck`, and `go test -race` on every pull request.

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

## Reporting issues

Please include the Go version (`go version`) and, if useful, the command or
code you ran and the expected vs. actual output.
