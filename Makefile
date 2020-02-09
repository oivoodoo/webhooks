DOCKER_COMPOSE=pipenv run docker-compose

init:
	(docker network create webhooks-network || true)
	pip install --user pipenv
	pipenv install
.PHONY: init

server:
	$(DOCKER_COMPOSE) up
.PHONY: master

stop:
	$(DOCKER_COMPOSE) stop || true
.PHONY: stop

test:
	go test -short -v gitlab.com/oivoodoo/webhooks/pkg/db/
	go test -short -v gitlab.com/oivoodoo/webhooks/pkg/master/v1/
	go test -short -v gitlab.com/oivoodoo/webhooks/pkg/slave/v1/
.PHONY: test
