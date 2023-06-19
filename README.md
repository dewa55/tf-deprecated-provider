# tf-depricated-provider
This is a Go program which can be run in the terraform directory to check if provider is depricated

To build a binary
```
go build -o terraform-deprecation-checker
```

To  build binary for amd64
```
GOOS=linux GOARCH=amd64 go build -o terraform-deprecation-checker
```
