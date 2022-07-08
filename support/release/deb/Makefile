all: deb-build

clean:
	rm -rf output unix-defender/usr/local/bin/unix-defender
	docker rmi -f deb-builder



docker-build:
	docker build -t deb-builder .

deb-build: docker-build
	docker run --rm -i -e TARGET_UID=$(shell id -u) -e TARGET_GID=$(shell id -g) -v $(shell pwd):/defender-release deb-builder make -C /defender-release build

build:
	mkdir -p output
	sh make-dists.sh

bump:
	./bump-version.sh