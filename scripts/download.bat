@echo off

REM GitHub repository and program name
set GITHUB_USER=MultiAdaptive
set REPO_NAME=multiAdaptive-cli
set PROGRAM_NAME=multiAdaptive-cli

REM Get system and architecture information
for /f "tokens=2 delims==" %%A in ('wmic os get osarchitecture /format:list') do set "ARCH=%%A"

REM Map system and architecture to Go architecture names
if "%ARCH%"=="64-bit" (
    set ARCH=amd64
) else if "%ARCH%"=="ARM64-based PC" (
    set ARCH=arm64
) else if "%ARCH%"=="32-bit" (
    set ARCH=386
) else (
    echo Unsupported architecture: %ARCH%
    exit /b 1
)

REM Set OS to windows
set OS=windows
set ARCH=%ARCH%.exe

REM Construct the srs file download URL
set SRS_DOWNLOAD_URL=https://github.com/%GITHUB_USER%/%REPO_NAME%/releases/latest/download/srs

REM Download the srs file
echo Downloading srs file from %SRS_DOWNLOAD_URL%...
curl -L -o srs %SRS_DOWNLOAD_URL%

REM Construct the download URL
set DOWNLOAD_URL=https://github.com/%GITHUB_USER%/%REPO_NAME%/releases/latest/download/%PROGRAM_NAME%-%OS%-%ARCH%

REM Download the file
echo Downloading %PROGRAM_NAME% for %OS%/%ARCH% from %DOWNLOAD_URL%...
curl -L -o %PROGRAM_NAME% %DOWNLOAD_URL%

REM Set execute permission
echo Download complete! You can now run %PROGRAM_NAME%.
