FROM fedora:30

ENV DEBIAN_FRONTEND noninteractive
RUN dnf install -y make mock wget

RUN useradd build
RUN gpasswd -a build mock

VOLUME /defender-release
WORKDIR /defender-release

ENTRYPOINT [ "make", "build" ]

