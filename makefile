PROTO_PATH := $(CURDIR)/protos/SSO
PKG_PROTO_PATH := $(CURDIR)/sso/pkg
PROTOC := protoc

# Файл сгенерированного кода
PROTO_FILES := $(wildcard $(PROTO_PATH)/*.proto)

install_protoc:
	@echo "🔧 Устанавливаем protoc ..."
	apt install -y protobuf-compiler
	protoc --version

install_go_plugins:
	@echo "🔧 Устанавливаем плагины для Go..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Плагины установлены. Добавьте GOPATH/bin в PATH:"
	@echo 'export PATH=$$PATH:$(go env GOPATH)/bin'

generate_proto:
	@echo "Генерируем proto ..."
	$(PROTOC) -I=$(PROTO_PATH) --go_out=$(PKG_PROTO_PATH) --go-grpc_out=$(PKG_PROTO_PATH) $(PROTO_FILES)