.PHONY: run
run:
	docker stop lokle_api || true
	docker rm lokle_api || true
	docker build . -t lokle_api
	docker run --restart always --add-host=host.docker.internal:host-gateway --name lokle_api -p 3001:3001 -d lokle_api
