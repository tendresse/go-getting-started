FROM alpine:latest

MAINTAINER Tendresse App <tendresseapp@gmail.com>

WORKDIR "/opt"

ADD .docker_build/go-getting-started /opt/bin/go-getting-started

CMD ["/opt/bin/go-getting-started"]

