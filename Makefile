build:
	@CGO_ENABLED=0 go build -ldflags="-s -w" .

run: build
	@./kmscan $(ARGS)

test:
	BASE_PATH="${shell pwd}" go test -cover ./...

clean:
	rm -f kmscan

install: build
	mkdir -p /usr/local/bin
	cp -f kmscan /usr/local/bin
	chmod 755 /usr/local/bin/kmscan

uninstall:
	rm -f /usr/local/bin/kmscan
