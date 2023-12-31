SHELL = /bin/sh

PREFIX = /usr/local


.PHONY: mdlc
mdlc:
	 go build -ldflags "-X main.version=$$(git describe --always --dirty)" .

.PHONY: install
install:
	install -d \
		$(PREFIX)/bin

	install -pm 0755 mdlc $(PREFIX)/bin/mdlc

.PHONY: uninstall
uninstall:
	rm -f \
		$(PREFIX)/bin/mdlc

.PHONY: clean
clean: 
	rm -f mdlc
