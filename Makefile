APP=p2p
CC=go
PACK=goupx
VERSION=$(shell git describe)
OS=$(shell uname -s)
ARCH=$(shell uname -m)
sinclude config.make

build: $(APP)
ifdef UPX_BIN
	release: pack
endif

all: release

$(APP): help.go instance.go main.go
	$(CC) build -ldflags="-w -s -X main.AppVersion=$(VERSION)" -o $@ -v $^

pack: $(APP)
	$(PACK) $(APP)

clean:
	-rm -f $(APP)
	-rm -f $(APP)-$(OS)*

mrproper: clean
mrproper:
	-rm -f config.make

test:  $(APP)
	go test ./...

release: build
release:
	-mv $(APP) $(APP)-$(OS)-$(ARCH)-$(VERSION)

install: 
	@mkdir -p $(DESTDIR)/bin
	@cp $(APP) $(DESTDIR)/bin
