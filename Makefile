compile: compile-arm64 compile-amd64
	lipo -create -output cloud-backup ./cloud-backup-arm64 ./cloud-backup-amd64

compile-arm64:
	GOOS=darwin GOARCH=arm64 go build -o ./cloud-backup-arm64 ./cmd/

compile-amd64:
	GOOS=darwin GOARCH=amd64 go build -o ./cloud-backup-amd64 ./cmd/

compileLocal:
	go build -tags localDesination -o ./cloud-backup-local ./cmd/
