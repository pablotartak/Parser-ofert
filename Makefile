hello:
	echo "Hello"

build:
	mkdir -p bin
	go build -o bin/parser parser.go

run:
	go run parser.go

clean:
	rm -f bin/parser
