# sshm

Blazingly fast SSH manager written on [Bubble Tea](https://github.com/charmbracelet/bubbletea)

## Installation

macOS/Linux:
```sh
curl -Lo sshm "https://github.com/krabiworld/sshm/releases/latest/download/sshm-$(uname -s | tr A-Z a-z)-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')" && chmod +x sshm
```

Windows:
```powershell
iwr "https://github.com/krabiworld/sshm/releases/latest/download/sshm-windows-$($env:PROCESSOR_ARCHITECTURE.ToLower())" -OutFile sshm.exe
```

## Building from source

```sh
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o sshm ./cmd/sshm
```
