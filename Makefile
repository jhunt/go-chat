build:
	go build .

examples: build
	go build ./example/...
