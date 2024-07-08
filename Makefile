build:
	@CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=0.1" .

run: build
	@./kmscan $(ARGS)

test:
	BASE_PATH="${shell pwd}" go test -cover ./...

clean:
	rm -f kmscan

install:
	mkdir -p /usr/local/bin
	cp -f kmscan /usr/local/bin
	chmod 755 /usr/local/bin/kmscan

uninstall:
	rm -f /usr/local/bin/kmscan
