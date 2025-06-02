@echo off
REM Check if the bin directory exists, if not create it
if not exist bin (
    mkdir bin
)

REM Build the Go program
go build -tags windows -o bin/main.exe ./src

REM Check if the build was successful
if %ERRORLEVEL% NEQ 0 (
    echo Build failed.
    exit /b %ERRORLEVEL%
)

echo Build succeeded. Output is in the bin folder.
