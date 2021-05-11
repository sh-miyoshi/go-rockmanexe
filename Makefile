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
