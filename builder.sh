# Build binary for:

# Linux
GOOS=linux GOARCH=386 go build -o ./builds/app_linux_386
GOOS=linux GOARCH=amd64 go build -o ./builds/app_linux_amd64
# MacOS X
GOOS=darwin GOARCH=amd64 go build -o ./builds/app_osx
# Windows
GOOS=windows GOARCH=386 go build -o ./builds/app_win_386.exe
GOOS=windows GOARCH=amd64 go build -o ./builds/app_win_x64.exe

# Save commit info
git log -1 > ./builds/commit.info

echo "Building complete"
