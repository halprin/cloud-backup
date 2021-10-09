compile: compileProtobuf
	go build -o ./cloud-backup ./cmd/

compileLocal: compileProtobuf
	go build -tags localDesination -o ./cloud-backup-local ./cmd/

compileProtobuf:
	protoc --go_out=. --go_opt=paths=source_relative \
		./external/pb/envelope_encryption_preamble.proto \
		./external/pb/envelope_encryption_v100.proto

installProtobufDependenciesForLinux:
	apt update
	apt -y install unzip
	curl -L -o /tmp/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip
	unzip /tmp/protoc.zip -d /usr/local/
	go version
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	protoc-gen-go --version
