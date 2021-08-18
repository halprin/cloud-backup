compile: compileProtobuf
	go build -o ./cloud-backup ./cmd/

compileLocal: compileProtobuf
	go build -tags localDesination -o ./cloud-backup-local ./cmd/

compileProtobuf:
	protoc --go_out=. --go_opt=paths=source_relative \
		./external/pb/envelope_encryption_preamble.proto \
		./external/pb/envelope_encryption_v100.proto
