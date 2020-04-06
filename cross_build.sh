mkdir -p output/bin
GOOS=windows GOARCH=amd64 go build -o output/bin/poseidon.exe