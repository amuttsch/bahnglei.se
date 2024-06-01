FROM golang:1.22 as build

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o bahngleise

FROM scratch

WORKDIR /app
COPY --from=build /app/bahngleise /app/config.yml /app

EXPOSE 8080

# Run
CMD ["/app/bahngleise"]
