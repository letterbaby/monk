@echo off


@rem kill.file
set servers=pingsvr;pongsvr

for %%I in (%servers%) do (
	echo stop %%I
	taskkill /F /IM %%I.exe
	ping -n 2 127.1 >nul
)

echo kill all done

rem pause