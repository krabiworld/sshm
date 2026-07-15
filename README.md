# sshm

## How to build

```sh
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o sshm ./cmd/sshm
```
