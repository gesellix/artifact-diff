FROM golang:1.20.6-alpine AS builder
LABEL builder=true

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV APPPATH /app

#RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc
RUN apk add --update -t build-deps git

COPY . $APPPATH
RUN cd $APPPATH && go get -d \
 && go test -short ./... \
 && go build \
    -a \
    -ldflags '-s -w -extldflags "-static"' \
    -o /bin/main \
    github.com/gesellix/artifact-diff/cmd/artifact-diff

FROM alpine:3.18.2
LABEL maintainer="Tobias Gesellchen <tobias@gesellix.de> (@gesellix)"

ENTRYPOINT [ "/main" ]
CMD [ ]

RUN apk --no-cache add ca-certificates \
 && adduser -DH user
USER user

COPY --from=builder /bin/main /main
