mkdir -p ./builds

# Build binary for:

# Linux
GOOS=linux GOARCH=386 go build -o ./builds/app_linux_386
GOOS=linux GOARCH=amd64 go build -o ./builds/app_linux_x64
# MacOS X
GOOS=darwin GOARCH=amd64 go build -o ./builds/app_osx_x64
GOOS=darwin GOARCH=arm64 go build -o ./builds/app_osx_arm
# Windows
GOOS=windows GOARCH=386 go build -o ./builds/app_win_386.exe
GOOS=windows GOARCH=amd64 go build -o ./builds/app_win_x64.exe
# Raspberry Pi
GOOS=linux GOARCH=arm GOARM=7 go build -o ./builds/app_linux_rpi

# Save commit info
git log -1 > ./builds/commit.info

echo "Building complete"
