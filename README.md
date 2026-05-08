# i5s

`i5s` is a k9s-inspired terminal UI for Incus. It uses your existing Incus client configuration and lets you browse instances, open shells, view logs, edit config, run lifecycle actions, and switch remotes/projects without changing Incus defaults.

## Requirements

- Go
- A configured Incus client

## Run

```sh
go run ./cmd/i5s
```

Build a local binary:

```sh
go build ./cmd/i5s
```

Common flags:

```text
--remote string     Incus remote for this session
--project string    Incus project for this session
--refresh duration  Auto-refresh interval (default 5s)
--debug             Enable debug logging
--version           Print version and commit, then exit
```

Version output:

```text
Version: 0.1.0
Commit: abc1234
```

## Keys

```text
j/k, arrows  move selection
Enter        shell into running instance
e            edit instance config
l            view logs
c            view console logs
s            stop running instance
S            start stopped instance
d            delete stopped instance
R            switch remote
p            switch project
/            filter instances
r            refresh
?            help
q            quit
```

Remote and project changes are session-local. `i5s` does not mutate Incus CLI defaults.

## Releases

Releases are built with GoReleaser from semantic git tags.

Create and push a version tag:

```sh
git tag v0.1.0
git push origin v0.1.0
```

Run a release:

```sh
goreleaser release --clean
```

GoReleaser injects the tag version and short commit SHA into the binary, so a release built from `v0.1.0` prints:

```text
Version: 0.1.0
Commit: abc1234
```

Check the release config or build a local snapshot without publishing:

```sh
goreleaser check
goreleaser release --snapshot --clean
```
