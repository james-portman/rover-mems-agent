del *.exe
go build -ldflags="-s -w" -o rover-mems.exe
if %errorlevel% neq 0 exit /b %errorlevel%
upx -9 *.exe


REM make a zip file:
REM tar.exe -a -c -f rover-mems.zip *.exe web-static

REM make a winrar sfx:
"%ProgramFiles%\WinRAR\WinRAR.exe" a -afzip -cfg- -ed -ep1 -k -m5 -r -tl "-sfx%ProgramFiles%\WinRAR\Zip.sfx" "-zsfxoptions.txt" "rover-mems-agent.exe" "web-static" "rover-mems.exe"
