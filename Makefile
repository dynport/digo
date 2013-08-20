default: test

test:
	DIGITAL_OCEAN_CLIENT_ID=client_id DIGITAL_OCEAN_API_KEY=api_key go test -v

ctags:
	gotags *.go > tags
	
rebuild: clean
	go build -a -o ./godo

clean:
	rm -f ./godo

install:
	cp ./bin/godo /usr/local/bin/

all: clean
	go build -o ./bin/digo
