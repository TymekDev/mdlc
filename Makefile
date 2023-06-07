PREFIX = /usr/local

mdlc:
	 go build -ldflags "-X main.version=$$(git describe --always --dirty)" .

clean: 
	rm -f mdlc

install:
	install -d \
		$(PREFIX)/bin

	install -pm 0755 mdlc $(PREFIX)/bin/mdlc

uninstall:
	rm -f \
		$(PREFIX)/bin/mdlc

.PHONY: mdlc clean install uninstall
