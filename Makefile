.PHONY: helm-install
helm-install:
	helm install -f ./helm/values.yaml softcorp-user-service ./helm --namespace=softcorp-user --create-namespace

.PHONY: helm-delete
helm-delete:
	helm delete softcorp-user-service --namespace=softcorp-user

.PHONY: helm-package
helm-package:
	helm package ./helm -d ./helm/packages --version=$(tag) && \
	helm repo index ./helm