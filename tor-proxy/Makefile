name = tor-proxy

all: build start

build:
	docker build -t $(name) .

start:
	docker run -d --restart=unless-stopped --network rkn-bypasser --name $(name) $(name)

kill:
	docker kill $(name)

remove:
	docker rm $(name)

stop: kill remove

restart: stop start

rebuild: stop build start