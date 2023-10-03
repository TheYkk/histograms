# builder image
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /bin/apikey

FROM gcr.io/distroless/static:nonroot
WORKDIR /app/
COPY --from=builder /bin/apikey /app/apikey
ENTRYPOINT ["/app/apikey"]
