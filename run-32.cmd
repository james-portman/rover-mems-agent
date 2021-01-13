set GOARCH=386
del *.exe
go build -ldflags="-s -w" -o rover-mems_%GOARCH%-TEST.exe
if %errorlevel% neq 0 exit /b %errorlevel%
rover-mems_%GOARCH%-TEST.exe
