# sshm

Blazingly fast SSH manager written on [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Installation

You can download the binary from the [GitHub Releases](https://github.com/krabiworld/sshm/releases).

### Homebrew (macOS/Linux)

```sh
brew install --cask krabiworld/tap/sshm
```

## Building from source

```sh
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o sshm ./cmd/sshm
```
