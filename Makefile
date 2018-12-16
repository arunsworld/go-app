
kick: clean
	statics -i=assets/static -o=assets/static.go -pkg=assets -group=Assets  -prefix=assets/static
	statics -i=assets/templates -o=assets/templates.go -pkg=assets -group=Templates  -prefix=assets/templates
	cd client; gopherjs build
	cd assets/static; ln -s ../../client/client.js* .
	export PORT=9095; export DEV=$(PWD); export ENV=DEV; go-kick -appPath=$(PWD) -mainSourceFile=main.go -gopherjsAppPath=client

build: buildprep
	go build

dockerize: buildprep
	export GOOS=linux; export GOARCH=amd64; \
		go build
	docker build -t arunsworld/go-app:1.0.0 .

buildprep: clean
	cd client; gopherjs build -m
	cd assets/static; cp ../../client/client.js* .
	statics -i=assets/static -o=assets/static.go -pkg=assets -group=Assets  -prefix=assets/static
	statics -i=assets/templates -o=assets/templates.go -pkg=assets -group=Templates  -prefix=assets/templates

clean: 
	rm -f assets/static/client.js
	rm -f assets/static/client.js.map
	rm -f client/client.js
	rm -f client/client.js.map
	rm -f go-app
	statics -o=assets/static.go -pkg=assets -group=Assets -init=true
	statics -o=assets/templates.go -pkg=assets -group=Templates -init=true

push:
	docker push arunsworld/go-app:1.0.0
	# docker login -u arunsworld

run:
	docker run -d --name go-app -p 80:80 arunsworld/go-app:1.0.0

rungcp:
	docker run -d --name go-app --network apps-net arunsworld/go-app:1.0.0
