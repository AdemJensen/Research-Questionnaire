all: clean fmt resources build

fmt:
	goimports -w .

test:
	go test -v ./...

build:
	go build -o output/

resources:
	mkdir output
	cp -r static output/static

clean:
	rm -rf output