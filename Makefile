GIT_COMMIT = $(shell git rev-parse --short HEAD)
GIT_STATUS = $(shell test -n "`git status --porcelain`" && echo "+CHANGES")
GIT_REV = $(GIT_COMMIT)$(GIT_STATUS)

default: test install

test:
	DIGITAL_OCEAN_CLIENT_ID=client_id DIGITAL_OCEAN_API_KEY=api_key go test -v

release:
	GOOS=linux  GOARCH=amd64 bash ./scripts/release.sh
	GOOS=darwin GOARCH=amd64 bash ./scripts/release.sh

ctags:
	gotags *.go > tags
	
rebuild: clean
	go build -a -o ./godo

clean:
	rm -f ./godo

install:
	go install -ldflags "-X main.GITCOMMIT $(GIT_REV)" github.com/dynport/digo 
	go install github.com/dynport/digo/digo
