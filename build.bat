SET CGO_ENABLED=0
SET GOARCH=amd64
set GOPATH=
set GOROOT=
SET GOOS=linux
go-bindata -o=./asset/script.go -pkg=asset static/script/... static/tmpl/... static/conf/...
go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s"