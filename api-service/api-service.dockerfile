# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY apiApp /app

CMD [ "/app/apiApp" ]