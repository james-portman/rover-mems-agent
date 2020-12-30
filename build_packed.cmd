del *%GOARCH%.exe
go build -ldflags="-s -w" -o rover-mems_%GOARCH%.exe
if %errorlevel% neq 0 exit /b %errorlevel%
upx -9 rover-mems_%GOARCH%.exe


REM make a zip file:
REM tar.exe -a -c -f rover-mems-%GOARCH%.zip *.exe web-static

REM make a winrar sfx:
"%ProgramFiles%\WinRAR\WinRAR.exe" a -afzip -cfg- -ed -ep1 -k -m5 -r -tl "-sfx%ProgramFiles%\WinRAR\Zip.sfx" "-zsfxoptions_%GOARCH%.txt" "rover-mems-agent_%GOARCH%.exe" "web-static" "rover-mems_%GOARCH%.exe"
