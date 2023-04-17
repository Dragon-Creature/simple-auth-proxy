install:
	docker build -t docker.ssns.se/frozendragon498/simple-auth-proxy ./

push:
	docker push docker.ssns.se/frozendragon498/simple-auth-proxy

build: install push