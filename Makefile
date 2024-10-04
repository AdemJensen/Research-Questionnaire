all: clean fmt resources build

fmt:
	goimports -w .

test:
	go test -v ./...

build:
	go build -o output/darwin
	GOOS=linux GOARCH=amd64 go build -o output/linux-amd64

resources:
	mkdir output
	cp -r static output/static

clean:
	rm -rf output