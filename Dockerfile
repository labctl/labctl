FROM alpine:3.18

LABEL maintainer="Johann Kellerman <kellerza@gmail.com>"
LABEL documentation="https://labctl.net"
LABEL repo="https://github.com/labctl/labctl"

RUN apk add --no-cache bash \
	curl \
	docker-cli \
	git \
	openssh

COPY labctl /usr/bin/labctl

CMD ["/usr/bin/labctl"]
