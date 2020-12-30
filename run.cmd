@echo off
del *.exe
go build -ldflags="-s -w" -o rover-mems.exe
if %errorlevel% neq 0 exit /b %errorlevel%
rover-mems.exe
