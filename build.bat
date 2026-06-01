@echo off
echo Building terminal-agent...
go build -o bin\terminal-agent.exe .
if %errorlevel% neq 0 (
    echo Build failed.
    timeout /t 5 >nul
    exit /b %errorlevel%
)
echo Build successful: bin\terminal-agent.exe
timeout /t 2 >nul
