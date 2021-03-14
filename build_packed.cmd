set GOARCH=386
del *%GOARCH%.exe

REM -ldflags="-s -w" - upsets windows defender (strips exe)
go build -o rover-mems_%GOARCH%.exe
if %errorlevel% neq 0 exit /b %errorlevel%
upx -9 rover-mems_%GOARCH%.exe


REM make a zip file:
tar.exe -a -c -f rover-mems-%GOARCH%.zip rover-mems_%GOARCH%.exe web-static

REM the winrar SFX archives were getting flagged as a trojan.. I was only trying to save space
REM make a winrar sfx:
REM "%ProgramFiles%\WinRAR\WinRAR.exe" a -afzip -cfg- -ed -ep1 -k -m5 -r -tl "-sfx%ProgramFiles%\WinRAR\Zip.sfx" "-zsfxoptions_%GOARCH%.txt" "rover-mems-agent_%GOARCH%.exe" "web-static" "rover-mems_%GOARCH%.exe"
