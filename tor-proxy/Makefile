name = tor-proxy

install: build start

build:
	docker build -t $(name) .

start:
	docker run -d --restart=unless-stopped --name $(name) $(name)

stop:
	docker rm -f $(name)

rebuild: stop build start

restart: stop start

log:
	docker logs -f $(name)