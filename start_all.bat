@echo off

set BIN_DIR=E:\Workspace\Go\bin\video-server

cd %BIN_DIR%
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

cd %~dp0