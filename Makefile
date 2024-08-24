.build-tailwind:
	npx --yes tailwindcss -i ./input.css -o ./assets/css/style.css

.build-application:
	go build

.build-templ:
	templ generate

build: .build-tailwind .build-templ .build-application

deploy: 
	fly deploy

build-dev: .build-tailwind .build-templ
	go build -o ./tmp/main .

dev:
	air serve

