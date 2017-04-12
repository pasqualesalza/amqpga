FROM alpine
MAINTAINER Pasquale Salza <pasquale.salza@gmail.com>

# Sets environment variables.
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Copies files into container.
ADD . /go/src/github.com/pasqualesalza/amqpga

# Installs AMQPGA.
RUN \
    apk add --no-cache --virtual .build-dependencies bash gcc musl-dev openssl git go godep && \
    # Creates the directories.
    mkdir -p "$GOPATH/src" "$GOPATH/bin" && \
    chmod -R 777 "$GOPATH" && \
    # Installs project dependencies.
    cd $GOPATH/src/github.com/pasqualesalza/amqpga && \
    godep restore && \
    # Installs the app.
    cd $GOPATH && \
    go install github.com/pasqualesalza/amqpga && \
    # Cleaning stuff.
    apk del .build-dependencies && \
    rm -rf $GOPATH/pkg && \
    rm -rf $GOPATH/src

# Sets the command.
ENTRYPOINT ["amqpga"]
