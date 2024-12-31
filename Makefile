.PHONY: bridgeaddr deploy

bridgeaddr: $(shell find . -name "*.go")
	CGO_ENABLED=1 CC=$$(which musl-gcc) go build -tags netgo,osusergo -ldflags='-s -w -linkmode external -extldflags "-static"' -o ./bridgeaddr

deploy: bridgeaddr
	systemctl stop bridgeaddr
	cp ./bridgeaddr /opt/bridgeaddr/bridgeaddr
	chown bridgeaddr:bridgeaddr /opt/bridgeaddr/bridgeaddr
	systemctl start bridgeaddr
