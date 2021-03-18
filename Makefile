install-kind:
	-brew install kind
cluster:
	kind create cluster --name image-cloner-test --config kind.yaml
clean:
	kind delete cluster --name image-cloner-test
run:
	go run main.go
