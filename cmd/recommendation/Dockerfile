FROM golang:1.19.3-alpine3.16 as build
RUN apk add --no-cache gcc musl-dev
RUN mkdir /app
WORKDIR /app
COPY . .
RUN go buid -o bin/recommendation -v ./cmd/recommendation/


FROM alpine:3.16
EXPOSE 8886
COPY --from=build /app/bin/recommendation .
CMD ["/recommendation"]
