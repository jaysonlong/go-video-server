@echo off

set SRC_DIR=%~dp0
set WORK_DIR=%GOPATH%\bin\video-server
set GOBIN=%WORK_DIR%
set CURR_DIR=%cd%

REM Build and install project
cd /d %SRC_DIR%
cd api
go install
cd ../streamserver
go install
cd ../scheduler
go install
cd ../web
go install

REM Move files into work directory
cd /d %WORK_DIR%
mkdir templates videos
xcopy %SRC_DIR%\templates .\templates /e /y

cd /d %CURR_DIR%