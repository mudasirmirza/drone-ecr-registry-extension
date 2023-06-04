FROM golang:1.20-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

# Disable CGO and set GOOS and GOARCH for the build.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o drone-ecr-registry-extension

FROM alpine:3.18
RUN apk add -U --no-cache ca-certificates
COPY --from=builder /app/drone-ecr-registry-extension /bin/drone-ecr-registry-extension

EXPOSE 3000
ENTRYPOINT ["/bin/drone-ecr-registry-extension"]
