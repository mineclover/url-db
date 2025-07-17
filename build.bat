@echo off
echo Building URL-DB Server...

REM Set Go environment
set GO111MODULE=on
set GOOS=windows
set GOARCH=amd64

REM Clean previous builds
if exist "bin" rmdir /s /q bin
mkdir bin

REM Build the application
echo Building for Windows...
go build -o bin\url-db.exe cmd\server\main.go
if errorlevel 1 (
    echo Build failed!
    exit /b 1
)

echo Build completed successfully!
echo Executable created: bin\url-db.exe

REM Run tests
echo Running tests...
go test -v ./...
if errorlevel 1 (
    echo Tests failed!
    exit /b 1
)

echo All tests passed!
echo.
echo To run the server:
echo   bin\url-db.exe
echo.
echo Default configuration:
echo   Port: 8080
echo   Database: file:./url-db.sqlite
echo   Tool Name: url-db