.PHONY: install-dependencies-local install-dependencies-ci codegen format lint

BUF_VERSION = v1.3.0
BUF_FILENAME = buf-$(shell uname -s)-$(shell uname -m)
BIN_DIR = /usr/local/bin

install-dependencies-local:
	brew install bufbuild/buf/buf
	brew install clang-format

install-dependencies-ci:
	curl -OL https://github.com/bufbuild/buf/releases/download/${BUF_VERSION}/${BUF_FILENAME}
	sudo mv ${BUF_FILENAME} ${BIN_DIR}/buf
	sudo chmod +x ${BIN_DIR}/buf
	apt-get --no-install-recommends install -y clang-format

codegen:
	buf generate

format:
	find ./proto -name "*.proto" | xargs clang-format -i

lint:
	buf lint
