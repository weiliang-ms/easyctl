@echo off
cd ../
go test ./... -v -cover > coverage.txt