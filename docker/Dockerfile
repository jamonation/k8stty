FROM busybox:latest

COPY server /usr/local/bin/server

RUN adduser -D -H server

USER server

CMD /usr/local/bin/server
