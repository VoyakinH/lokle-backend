.PHONY: docker-run
docker-run:
	docker build . -t lokle_api
	docker run --name lokle_api -p 3001:3001 -d lokle_api

.PHONY: docker-stop
docker-stop:
	docker stop lokle_api

.PHONY: docker-prune
docker-prune:
	docker system prune -a --volumes
