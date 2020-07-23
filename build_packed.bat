del *.exe
go build -ldflags="-s -w"
if %errorlevel% neq 0 exit /b %errorlevel%
upx -9 rover-mems-agent.exe
