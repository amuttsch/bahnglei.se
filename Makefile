.build-tailwind:
	npx --yes tailwindcss-cli -i ./input.css -o ./assets/css/style.css

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

i18n-extract:
	goi18n extract -outdir translations
	
i18n-merge:
	cd translations; goi18n merge active.*.toml

i18n-finish:
	cd translations; goi18n merge active.*.toml translate.*.toml

