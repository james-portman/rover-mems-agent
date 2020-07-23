@echo off
del *.exe
go build
if %errorlevel% neq 0 exit /b %errorlevel%
rover-mems-agent.exe


REM go build -ldflags="-s -w"
