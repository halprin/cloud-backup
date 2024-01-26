compile:
	go build -o ./cloud-backup ./cmd/

compileLocal:
	go build -tags localDesination -o ./cloud-backup-local ./cmd/
