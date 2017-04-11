all: windows linux

install-dependencies:
	glide restore

windows: windows-amd64 windows-386

windows-amd64:
	mkdir -p bin/windows/amd64
	GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/haproxy-dynagent.exe

windows-386:
	mkdir -p bin/windows/386
	GOOS=windows GOARCH=386 go build -o bin/windows/386/haproxy-dynagent.exe

linux: linux-386 linux-amd64 linux-arm linux-arm64

linux-386:
	mkdir -p bin/linux/386
	GOOS=linux GOARCH=386 go build -o bin/linux/386/haproxy-dynagent

linux-amd64:
	mkdir -p bin/linux/amd64
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/haproxy-dynagent

linux-arm:
	mkdir -p bin/linux/arm
	GOOS=linux GOARCH=arm go build -o bin/linux/arm/haproxy-dynagent

linux-arm64:
	mkdir -p bin/linux/arm64
	GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/haproxy-dynagent

clean:
	rm -rf bin/
