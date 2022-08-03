all: rpm-build
build: source
	mock -r fedora-30-x86_64 \
		--source=unix-defender-source.tar.gz \
		--spec=unix-defender.spec \
		--rebuild \
		--isolation=simple \
		--resultdir=output 

bump:
	./bump-version.sh

source: 
		tar cvzf unix-defender-source.tar.gz unix-defender-*/*


clean:
		rm -rf output unix-defender/usr/bin/unix-defender
		docker rmi rpm-builder

docker-build:
	docker build -t rpm-builder .

rpm-build: docker-build
	docker run --rm -i -e TARGET_UID=$(shell id -u) -e TARGET_GID=$(shell id -g) -v $(shell pwd):/defender-release --security-opt apparmor:unconfined --cap-add=SYS_ADMIN rpm-builder


