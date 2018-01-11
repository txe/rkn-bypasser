name = rkn-bypasser

all: build start

build:
	docker build -t $(name) .

start:
	docker run -d --restart=unless-stopped -p 127.0.1.1:8000:8000 --name $(name) $(name)

kill:
	docker kill $(name)

remove:
	docker rm $(name)

stop: kill remove

restart: stop start

rebuild: stop build start