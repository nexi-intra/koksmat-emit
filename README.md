# Emitter

This application reads events from a queue, processes them, and triggers external services (like GitHub Actions) as needed.

## Structure

```
tree | set-clipboard

.
├── Dockerfile
├── LICENSE
├── README.md
├── api
│   ├── register.go
│   ├── webhook-github.go
│   └── webhook-microsoftgraph.go
├── cmd
│   ├── root.go
│   └── serve.go
├── config
│   └── setup.go
├── config.go
├── dependencies
│   ├── github.go
│   ├── magicmix.go
│   ├── nats
│   │   └── main.go
│   └── testmain_test.go
├── go.mod
├── go.sum
├── internal
│   ├── emitter
│   │   └── app.go
│   └── observability
│       └── observability.go
├── main.go
├── scripts
│   └── release.ps1
└── services
    └── webhooks.go

11 directories, 21 files
```
