FROM alpine:latest

RUN mkdir /app

COPY .env /app

COPY ./src/libraryService /app

CMD [ "/app/libraryService"]