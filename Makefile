all:
	go build -o rockman.exe
	rockman.exe
dev:
	go build -o rockman.exe
	rockman.exe --config data/config_debug.yaml
res:
	.\tmp\DXArchive\DxaEncode.exe .\data\private\images
	move .\data\private\images.dxa ".\data\"
	.\tmp\DXArchive\DxaEncode.exe .\data\private\sounds
	move .\data\private\sounds.dxa ".\data\"
	git add .\data\images.dxa .\data\sounds.dxa
protoc:
	cd pkg/oldnet/netconnpb && \
	protoc --go_out=plugins=grpc:. *.proto && \
	move .\github.com\sh-miyoshi\go-rockmanexe\pkg\net\netconnpb\netconn.pb.go . && \
	rd /s /q github.com
protoc-linux:
	protoc --go_out=. --go-grpc_out=. --go-grpc_opt require_unimplemented_servers=false ./pkg/net/netconnpb/netconn.proto && \
	cp github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb/* ./pkg/net/netconnpb/ && \
	rm -rf github.com
docker:
	docker build -t smiyoshi/rockmanexe-matcher -f build/Dockerfile.matcher .
	docker build -t asia-northeast1-docker.pkg.dev/rockmanexe/router/router -f build/Dockerfile.router .
localnet:
	cd test/helpers/localnet-boot && \
	go build -o localnet-boot.exe && \
	localnet-boot.exe
