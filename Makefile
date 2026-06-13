SINGBOX_VERSION = 1.13.12
SINGBOX_TAGS = with_v2ray_api,with_gvisor,with_quic,with_dhcp,with_wireguard,with_utls,with_clash_api
SINGBOX_LDFLAGS = -X 'github.com/sagernet/sing-box/constant.Version=$(SINGBOX_VERSION)'

VERSION := $(shell git describe --tags 2>/dev/null || echo "dev")

DOCKER = docker run --rm \
	-v $(shell pwd):/app \
	-v $(HOME)/go/pkg/mod:/root/go/pkg/mod \
	-v $(HOME)/.cache/go-build:/root/.cache/go-build \
	-v $(HOME)/go/src/sing-box:/sing-box \
	slinx-builder sh -c

build: slinx core

slinx: slinx-amd64 slinx-arm64

core: core-amd64 core-arm64

slinx-amd64:
	@$(DOCKER) "mkdir -p /app/dist && \
		echo \"\$$(date '+%H:%M:%S') >>> 编译 slinx amd64...\" && \
		cd /app && GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-gnu-gcc go build -ldflags='-s -w -X main.Version=$(VERSION)' -o dist/slinx_linux_amd64 . && \
		echo \"\$$(date '+%H:%M:%S') >>> 压缩 slinx amd64...\" && \
		upx dist/slinx_linux_amd64 && \
		echo \"\$$(date '+%H:%M:%S') >>> 完成 slinx_linux_amd64\""

slinx-arm64:
	@$(DOCKER) "mkdir -p /app/dist && \
		echo \"\$$(date '+%H:%M:%S') >>> 编译 slinx arm64...\" && \
		cd /app && GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build -ldflags='-s -w -X main.Version=$(VERSION)' -o dist/slinx_linux_arm64 . && \
		echo \"\$$(date '+%H:%M:%S') >>> 压缩 slinx arm64...\" && \
		upx dist/slinx_linux_arm64 && \
		echo \"\$$(date '+%H:%M:%S') >>> 完成: slinx_linux_arm64\""

core-amd64:
	@$(DOCKER) "mkdir -p /app/dist && \
		echo \"\$$(date '+%H:%M:%S') >>> 编译 sing-box amd64...\" && \
		cd /sing-box && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags '$(SINGBOX_TAGS)' -ldflags='$(SINGBOX_LDFLAGS) -s -w' -o /app/dist/sing-box_linux_amd64 ./cmd/sing-box && \
		echo \"\$$(date '+%H:%M:%S') >>> 压缩 sing-box amd64...\" && \
		gzip -f -k /app/dist/sing-box_linux_amd64 && rm /app/dist/sing-box_linux_amd64 && \
		echo \"\$$(date '+%H:%M:%S') >>> 完成: sing-box_linux_amd64.gz\""

core-arm64:
	@$(DOCKER) "mkdir -p /app/dist && \
		echo \"\$$(date '+%H:%M:%S') >>> 编译 sing-box arm64...\" && \
		cd /sing-box && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags '$(SINGBOX_TAGS)' -ldflags='$(SINGBOX_LDFLAGS) -s -w' -o /app/dist/sing-box_linux_arm64 ./cmd/sing-box && \
		echo \"\$$(date '+%H:%M:%S') >>> 压缩 sing-box arm64...\" && \
		gzip -f -k /app/dist/sing-box_linux_arm64 && rm /app/dist/sing-box_linux_arm64 && \
		echo \"\$$(date '+%H:%M:%S') >>> 完成: sing-box_linux_arm64.gz\""