.PHONY: helm-install
helm-install:
	helm install -f ./helm/values.yaml block-user-service ./helm --namespace=softcorp-blocks --create-namespace

.PHONY: helm-delete
helm-delete:
	helm delete block-user-service --namespace=softcorp-blocks

helm-package:
	rm -rf * ./helm/charts/ &&
	rm ./helm/index.yaml || true && touch ./helm/index.yaml && \
	helm package ./helm -d ./helm/charts --version=$(tag) && \
	helm repo index ./helm