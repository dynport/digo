default: test install

test:
	DIGITAL_OCEAN_CLIENT_ID=client_id DIGITAL_OCEAN_API_KEY=api_key go test -v

ctags:
	gotags *.go > tags
	
rebuild: clean
	go build -a -o ./godo

clean:
	rm -f ./godo

install:
	go install github.com/dynport/digo 

all: clean
	go build -o ./bin/digo
