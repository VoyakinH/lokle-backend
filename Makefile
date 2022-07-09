.PHONY: run
run:
	docker stop lokle_api || true
	docker rm lokle_api || true
	docker build . -t lokle_api
	docker run --restart always --network host -h 127.0.0.1 --name lokle_api -p 3001:3001 -d lokle_api
