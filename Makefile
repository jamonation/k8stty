.PHONY: terminal websocket namespace networkpolicy pod service kind-images
CGO_ENABLED := 0
GOOS := linux
GOARCH := amd64

export CGO_ENABLED GOOS GOARCH

terminal websocket namespace networkpolicy pod service:
	cd src && go build -o ../builds/$@/server ./cmd/$@/
	docker build -q --platform linux/amd64 builds/$@/ -t localhost:5001/$@:latest -f docker/Dockerfile && docker push localhost:5001/$@:latest

kind-images:
	cd src/cmd/terminal; ko build -B .
	cd src/cmd/websocket; ko build -B .
	cd src/cmd/namespace; ko build -B .
	cd src/cmd/networkpolicy; ko build -B .
	cd src/cmd/pod; ko build -B .
	cd src/cmd/service; ko build -B . 

