@echo off

tbltool.exe -f=roletbl.md -po=./roletbl.proto -do=../src/common/roletble/roletble.go -so=./roletble.tsql

protoc.exe --go_out=../src/ss_proto ss_proto.proto
protoc.exe --go_out=../src/common/roletble roletbl.proto

pause


