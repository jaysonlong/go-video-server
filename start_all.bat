@echo off

set WORK_DIR=%GOPATH%\bin\video-server
set CURR_DIR=%cd%

cd %WORK_DIR%
start /min .\api
start /min .\streamserver
start /min .\scheduler
start /min .\web

echo Press any key to exit all programs...
pause > nul
taskkill /F /IM web.exe
taskkill /F /IM scheduler.exe
taskkill /F /IM streamserver.exe
taskkill /F /IM api.exe

cd %CURR_DIR%