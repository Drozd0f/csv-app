FROM golang:1.19.2 as base
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /csv-app ./cmd


FROM alpine:3.16.2 as prod

COPY --from=base /csv-app /usr/bin/csv-app
CMD ["sh", "-c", "csv-app run"]
