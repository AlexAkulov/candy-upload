NAME := candy-upload
VERSION := $(shell git describe --always --tags --abbrev=0 | tail -c +2)
RELEASE := $(shell git describe --always --tags | awk -F- '{ if ($$2) dot="."} END { printf "1%s%s%s%s\n",dot,$$2,dot,$$3}')
GO_VERSION := $(shell go version | cut -d' ' -f3)
BUILD_DATE := $(shell date --iso-8601=second)
LDFLAGS := -ldflags "-X main.version=${VERSION}-${RELEASE} -X main.goVersion=${GO_VERSION} -X main.buildDate=${BUILD_DATE}"

default: build

test: prepare_test
	go test -v

prepare_test:
	go get "github.com/smartystreets/goconvey"

prepare:
	echo "All ready"

build: clean prepare
	mkdir -p build/root/usr/bin/
	go build  ${LDFLAGS} -o build/root/usr/bin/${NAME}

tar: build
	mkdir -p build/root/etc/${NAME}
	cp config.yml build/root/etc/${NAME}/config.yml
	tar -czvPf build/${NAME}-${VERSION}-${RELEASE}.tar.gz -C build/root .

rpm:
	fpm -t rpm \
		-s "tar" \
		--description "Very simple backend for upload and processing files" \
		--vendor "Alex Akulov" \
		--url "https://github.com/AlexAkulov/statsd-ha-proxy" \
		--license "GPLv3" \
		--name "${NAME}" \
		--version "${VERSION}" \
		--iteration "${RELEASE}" \
		--depends logrotate \
		--config-files "/etc/${NAME}/config.yml" \
		-p build \
		build/${NAME}-${VERSION}-${RELEASE}.tar.gz

clean:
	rm -rf build

.PHONY: test
