
kick: clean
	cd client; gopherjs build
	cd assets/static; ln -s ../../client/client.js* .
	export PORT=9095; kick -appPath=$(PWD) -mainSourceFile=main.go -gopherjsAppPath=client

build: buildprep
	packr2 build; packr2 clean

dockerize: buildprep
	export GOOS=linux; export GOARCH=amd64; \
		packr2 build; packr2 clean
	docker build -t arunsworld/go-app:1.0.0 .

buildprep: clean
	cd client; gopherjs build
	cd assets/static; cp ../../client/client.js* .

clean: 
	rm -f assets/static/client.js
	rm -f assets/static/client.js.map
	rm -f client/client.js
	rm -f client/client.js.map
	rm -f go-app

push:
	docker push arunsworld/go-app:1.0.0
	# docker login -u arunsworld

run:
	docker run -d --name go-app -p 80:80 arunsworld/go-app:1.0.0

rungcp:
	docker run -d --name go-app --network apps-net arunsworld/go-app:1.0.0
