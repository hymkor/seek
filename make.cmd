@echo off
setlocal
set GOARCH=386
call :"%1"
endlocal
exit /b

:""
    go build
    exit /b

:"update"
    for /F "skip=1" %%I in ('where seek.exe') do copy /-Y seek.exe "%%I"
    exit /b
