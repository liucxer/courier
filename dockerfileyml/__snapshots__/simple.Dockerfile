FROM busybox:latest

WORKDIR /todo

ENV key=hello

COPY x ./

ENTRYPOINT ["sh"]

CMD ["-c","echo","${key}"]

