@echo off

set SRC_DIR=E:\Workspace\Go\src\github.com\midmis\go-video-server
set BIN_DIR=E:\Workspace\Go\bin\video-server
set GOBIN=%BIN_DIR%

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
cd /d %BIN_DIR%
mkdir templates videos
xcopy %SRC_DIR%\templates .\templates /e /y

cd %~dp0