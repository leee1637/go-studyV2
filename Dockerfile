FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app/cmd/app
COPY app_binary ./app
COPY .env /app/.env
EXPOSE 8081
CMD ["./app"]