FROM node:latest as tailwind

WORKDIR /app

COPY . /app
RUN npx --yes tailwindcss -i ./css/input.css -o ./css/style.css

FROM ghcr.io/a-h/templ:latest as templ

WORKDIR /app
COPY --chown=65532:65532 . /app
RUN ["templ", "generate"]

FROM golang:1.22 as build

WORKDIR /app

RUN apt-get install -y ca-certificates
RUN update-ca-certificates

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY --from=templ /app /app
COPY --from=tailwind /app/css /app/css

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o bahngleise

FROM scratch

WORKDIR /app
COPY --from=build /app/bahngleise /app/bahngleise
COPY --from=build /app/config.fly.yml /app/config.yml 
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

# Run
CMD ["/app/bahngleise"]
