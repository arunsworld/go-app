FROM alpine:latest

COPY go-app /bin/

EXPOSE 80

CMD ["go-app"]