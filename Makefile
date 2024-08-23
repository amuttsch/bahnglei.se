.build-tailwind:
	npx --yes tailwindcss -i ./css/input.css -o ./css/style.css

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

