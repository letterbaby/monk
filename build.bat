@echo off

set GO111MODULE=on
set GOPROXY=https://goproxy.cn,direct
set GOPATH=%cd%
set GOPRIVATE=github.com

echo %GOPATH%

cd src
@rem make.file
set servers=pingsvr;pongsvr

for %%I in (%servers%) do (
	echo build %%I
	@echo on
	go build -o ../bin -v ./%%I
	@echo off
)

echo build all done

pause