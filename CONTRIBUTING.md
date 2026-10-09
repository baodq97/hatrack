# Contributing

Issues and pull requests are welcome.

## Before you open a pull request

```sh
gofmt -l .        # prints nothing
go vet ./...
go test ./...
GOOS=windows go vet ./...   # the tray app only builds for Windows
GOOS=darwin go vet ./cmd/hat ./internal/claude ./internal/gui   # macOS keychain code; the menu bar app needs a Mac
```

CI runs the same on Linux, Windows and macOS. On a Mac, `HATRACK_KEYCHAIN_TEST=1 go test ./...` also
round-trips a throwaway item through the real login keychain.

## Ground rules

- hatrack handles login tokens. Changes must not send anything over the network, log tokens, or
  write them anywhere but Claude's login files and `~/.hatrack`.
- The user stays in control: no hidden steps. If hatrack does something on the user's behalf, it says so.
- Logic that touches login files goes in `internal/claude` and comes with a test.
- Keep changes small and dependencies few.

## Releasing

Maintainers tag a version and push it; the release workflow tests, builds, attests and publishes:

```sh
git tag v0.1.0 && git push origin v0.1.0
```
