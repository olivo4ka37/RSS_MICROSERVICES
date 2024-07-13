# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY newsApp /app

CMD [ "/app/newsApp" ]