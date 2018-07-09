@echo off
setlocal
set GOARCH=386
set "EXE=%~dp0seek.exe"
call :"%1"
endlocal
exit /b

:""
    go build
    exit /b

:"upgrade"
    for /F %%I in ('where seek.exe') do if not "%%I" == "%EXE%" copy /-y /v "%EXE%" "%%I"
    exit /b
