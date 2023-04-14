FROM golang:1.20-alpine as build

COPY . /build
WORKDIR /build

EXPOSE 3000

RUN CGO_ENABLED=0 go build ./app.go

FROM gcr.io/distroless/static

COPY --from=build /build/app /app

CMD ["/app"]
