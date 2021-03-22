IMG ?= ishankhare07/image-cloner

install-kind:
	-brew install kind
cluster:
	kind create cluster --name image-cloner-test --config kind.yaml
clean:
	kind delete cluster --name image-cloner-test
build:
	docker build . -t ${IMG}
push:
	docker push ${IMG}
deploy:
	cd config/manifests && kustomize edit set image controller=${IMG}
	kustomize build config/manifests | kubectl apply -f -
run:
	DEV_MODE=ON go run -race main.go
