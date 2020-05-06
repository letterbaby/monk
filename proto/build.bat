@echo off

protoc.exe --go_out=../src/ss_proto ss_proto.proto
