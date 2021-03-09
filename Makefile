all:
	go build -o rockman.exe
	rockman.exe
dev:
	go build -o rockman.exe
	rockman.exe --debug
release:
	.\tmp\DXArchive\DxaEncode.exe .\data\private\images
	move .\data\private\images.dxa ".\data\"
