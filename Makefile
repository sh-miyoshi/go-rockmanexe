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
	git commit -m "update resources"
router:
	go build ./cmd/router
	router.exe --config cmd/router/config.yaml
protoc:
	cd pkg/net/routerpb && \
	protoc --go_out=plugins=grpc:. *.proto && \
	move .\github.com\sh-miyoshi\go-rockmanexe\pkg\net\routerpb\router.pb.go . && \
	rd /s /q github.com
docker:
	docker build -t smiyoshi/rockmanexe-matcher -f build/Dockerfile.matcher .
	docker build -t smiyoshi/rockmanexe-router -f build/Dockerfile.router .
