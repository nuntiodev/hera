.PHONY: helm-install
helm-install:
	helm install -f ./helm/values.yaml nuntio-user-block ./helm --namespace=nuntio-blocks --create-namespace

.PHONY: helm-delete
helm-delete:
	helm delete nuntio-user-block --namespace=nuntio-blocks

.PHONY: helm-package
helm-package:
	rm -rf ./helm/charts/* && \
	rm ./helm/index.yaml || true && touch ./helm/index.yaml && \
	helm package ./helm -d ./helm/charts --version=$(tag) && \
	helm repo index ./helm