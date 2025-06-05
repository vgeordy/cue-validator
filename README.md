# MrNesBits Validator

## Demo

There are two ways to validate:

# 1. Programmatically using Go
navigate to root directory and run 
```bash
go run validator.go
```
validator.go has the function that needs to be integrated into the main application. It demos a sample scenario

## 2. Using CUE command
1. install 
```
cuelang.org/go/cmd/cue@latest
```
for more information: https://cuelang.org/docs/introduction/installation/

2. 
```bash
cue vet -d "#CpInitRoot" schema.cue test.yaml
```