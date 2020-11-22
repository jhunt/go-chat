build:
	go build .

examples: build
	for x in ./example/*; do go build $$x; done
