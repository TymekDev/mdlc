PREFIX = /usr/local

mdlsc:
	 go build -ldflags "-X main.version=$$(git describe --always --dirty)" .

clean: 
	rm -f mdlsc

install:
	install -d \
		$(PREFIX)/bin

	install -pm 0755 mdlsc $(PREFIX)/bin/mdlsc

uninstall:
	rm -f \
		$(PREFIX)/bin/mdlsc

.PHONY: mdlsc clean install uninstall
