all:
	go build -o rockman.exe
	rockman.exe
dev:
	go build -o rockman.exe
	rockman.exe --config data/config_debug.yaml
dev-mac:
	GOOS=windows go build -o rockman.exe
	wine rockman.exe --config data/config_debug.yaml
res:
	.\tmp\DXArchive\DxaEncode.exe .\data\private\images
	move .\data\private\images.dxa ".\data\"
	.\tmp\DXArchive\DxaEncode.exe .\data\private\sounds
	move .\data\private\sounds.dxa ".\data\"
	git add .\data\images.dxa .\data\sounds.dxa
protoc:
	cd pkg/net/netconnpb && \
	protoc --go_out=plugins=grpc:. *.proto && \
	move .\github.com\sh-miyoshi\go-rockmanexe\pkg\net\netconnpb\netconn.pb.go . && \
	rd /s /q github.com
protoc-linux:
	protoc --go_out=. --go-grpc_out=. --go-grpc_opt require_unimplemented_servers=false ./pkg/net/netconnpb/netconn.proto && \
	cp github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb/* ./pkg/net/netconnpb/ && \
	rm -rf github.com
docker:
	docker build -t docker.io/smiyoshi/rockmanexe-router -f build/Dockerfile.server .
	docker tag docker.io/smiyoshi/rockmanexe-router docker.io/smiyoshi/rockmanexe-router:v0.13
localnet:
	cd tools/localnet-boot && \
	go build -o localnet-boot.exe && \
	localnet-boot.exe
localnet-mac:
	cd tools/localnet-boot-mac && \
	go build -o localnet-boot.out && \
	./localnet-boot.out
skill:
	cd cmd/skill-simulator && \
	go build -o app.exe && \
	app.exe
skill-mac:
	cd cmd/skill-simulator && \
	GOOS=windows go build -o app.out && \
	wine64 app.out
skill-gen:
	go run tools/skill-code-generator/local/main.go
skill-gen-net:
	go run tools/skill-code-generator/network/main.go
