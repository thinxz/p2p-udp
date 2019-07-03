rm -rf cmd/conner
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o conn cmd/conner.go