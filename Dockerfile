# builder image
FROM golang:1.18 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o aboba


# generate clean, final image for end users
FROM alpine:latest
COPY --from=builder /build/aboba .
COPY migrations .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# executable
ENTRYPOINT [ "./aboba" ]
EXPOSE 80
