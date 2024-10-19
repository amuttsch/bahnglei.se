.build-tailwind:
	npx --yes tailwindcss -i ./input.css -o ./assets/css/style.css

.build-application:
	go build

.build-templ:
	templ generate

.build-sqlc:
	sqlc generate

build: .build-tailwind .build-templ .build-sqlc .build-application

deploy: 
	fly deploy

build-dev: .build-tailwind .build-templ .build-sqlc
	go build -o ./tmp/main .

dev:
	air serve

