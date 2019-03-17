BINARY = ssh-hammer
HOST = localhost
PORT = 2022

all: $(BINARY)

$(BINARY): *.go
	go build .

deps:
	go get .

build: $(BINARY)

clean:
	rm $(BINARY)

run: $(BINARY)
	./$(BINARY) -vv $(HOST):$(PORT) --num 2

test:
	go test .
	golint
