# TODO

## Releases / distribution

- [ ] **Homebrew tap** (once `pkg.go.dev` is indexing the module)
  - Create a separate `homebrew-tap` repo (e.g. `MichaelWiciak/homebrew-tap`)
  - Add a formula (`randwallpaper.rb`) that installs from source via `go install`
  - Publish the formula with `brew tap MichaelWiciak/tap && brew install randwallpaper`
  - Revisit binary assets once releases are flowing; source-install is the simplest for now
- [ ] Confirm the Release Please workflow produces the first release and that
      `go install github.com/MichaelWiciak/randwallpaper/cmd/randwallpaper@latest` resolves it

## Code / product

- [ ] Add a `-version` flag to the CLI that reports the release version
- [ ] `test_images/` currently only holds a `.gitkeep` — decide whether to add a
      sample or remove the directory