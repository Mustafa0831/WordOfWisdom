FROM alpine:latest

ADD out/bin/server /usr/local/bin/server

RUN chmod u+x /usr/local/bin/server

ENTRYPOINT ["server"]
