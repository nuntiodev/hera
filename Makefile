.PHONY: helm-install
helm-install:
	helm install -f ./helm/values.yaml block-user-service ./helm --namespace=softcorp-blocks --create-namespace

.PHONY: helm-delete
helm-delete:
	helm delete block-user-service --namespace=softcorp-blocks

.PHONY: helm-package
helm-package:
	helm package ./helm -d ./helm/packages --version=$(tag) && \
	helm repo index ./helm