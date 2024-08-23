.build-tailwind:
	npx --yes tailwindcss -i ./css/input.css -o ./css/style.css

.build-application:
	go build

build: .build-tailwind .build-application

deploy: build
	fly deploy

dev:
	air

