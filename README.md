# tf-depricated-provider
This is a Go program that can be run in the terraform directory to check if the provider is deprecated.
This program needs to be ran in in the terraform work directory. This go binary uses `terraform providers` command and checks if there is "This provider is deprecated" on the official terraform docs sites.

To run `terraform providers` command, first the `terraform init` command needs to be run. 

To build a binary
```
go build -o terraform-deprecation-checker
```

To  build a binary for amd64
```
GOOS=linux GOARCH=amd64 go build -o terraform-deprecation-checker
```
