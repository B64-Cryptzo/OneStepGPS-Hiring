@echo off

:: Run the backend server
cd Backend
start /B GolangBackend.exe

:: Run the build serve command
cd ../Frontend
start /B serve -s dist

:: Wait for the server to start, then copy the server address to the clipboard (if required)
timeout /t 5 /nobreak

:: Clear the screen
cls

:: Get clipboard content (URL)
for /f "delims=" %%i in ('powershell Get-Clipboard') do set clipboard_content=%%i

:: Check if clipboard contains 'localhost'
echo %clipboard_content% | findstr /i "localhost" >nul
if %errorlevel%==0 (
    :: Open the address in the default browser
    start %clipboard_content%
) else (
    echo No localhost URL found in clipboard.
)

:: Keep the terminal open
pause > nul
