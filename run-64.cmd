@echo off
set GOARCH=amd64
del *.exe
go build -ldflags="-s -w" -o rover-mems_%GOARCH%.exe
if %errorlevel% neq 0 exit /b %errorlevel%
rover-mems_%GOARCH%.exe
