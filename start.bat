@echo off

cd bin

rem 清理日志
del /Q /F /S .\log\*.*

@rem start.file
set servers=pingsvr;pongsvr

for %%I in (%servers%) do (
	echo start %%I
	start %%I.exe
)

echo start all done

cd ../

rem pause