FROM pmdcosta/golang:1.13 AS builder
WORKDIR /code

# Add go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Add code and compile it
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app ./cmd/app

# Final image
FROM gcr.io/distroless/base
COPY --from=builder /app ./
ENTRYPOINT ["./app"]