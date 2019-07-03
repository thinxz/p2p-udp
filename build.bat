del /F/Q cmd/conner.exe
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -v -o conn.exe cmd/conner.go

set GOOS=linux
go build -v -o conn cmd/conner.go