LDFLAGS=1.0.0
VERSION=v1.0
LIST=linux darwin
build:
	@for os in $(LIST); do \
		@CGO_ENABLED=0 GOOS=$$os GOARCH=amd64 go build -ldflags="-s -X main.version=${LDFLAGS}" . ; \
		@mkdir -p ~/Documents/git-push-bin/${VERSION}/$$os ; \
		@mv git-push  ~/Documents/git-push-bin/${VERSION}/$$os ; done \

	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -X main.version=${LDFLAGS}" .
	@mkdir -p ~/Documents/git-push-bin/${VERSION}/windows ;
	@mv git-push.exe  ~/Documents/git-push-bin/${VERSION}/windows

update:
	@for os in $(LIST); do \
		@CGO_ENABLED=0 GOOS=$$os GOARCH=amd64 go build -ldflags="-s -X main.version=${LDFLAGS}" . ; \
		@rm ~/Documents/git-push-bin/${VERSION}/$$os/git-push ; \
		@mkdir -p ~/Documents/git-push-bin/${VERSION}/$$os ; \
		@mv git-push  ~/Documents/git-push-bin/${VERSION}/$$os ; done \

	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -X main.version=${LDFLAGS}" .
	@rm ~/Documents/git-push-bin/${VERSION}/windows/git-push.exe
	@mkdir -p ~/Documents/git-push-bin/${VERSION}/windows ;
	@mv git-push.exe  ~/Documents/git-push-bin/${VERSION}/windows

upload:
	cd ~/go/src/github.com/urvil38/git-push-upload && ./git-push-upload -v ${VERSION}
