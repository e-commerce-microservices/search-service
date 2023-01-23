.PHONY: rebuild
rebuild:
	docker build -t ngoctd/ecommerce-search:latest . && \
	docker push ngoctd/ecommerce-search

.PHONY: redeploy
redeploy:
	kubectl rollout restart deployment depl-search