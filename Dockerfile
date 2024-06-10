FROM golang:1.22 as build

WORKDIR /app

RUN apt-get install -y ca-certificates
RUN update-ca-certificates

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o bahngleise

FROM scratch

WORKDIR /app
COPY --from=build /app/bahngleise /app/bahngleise
COPY --from=build /app/config.fly.yml /app/config.yml 
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

# Run
CMD ["/app/bahngleise"]
