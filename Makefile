BIN=fitness-backend
PORT=3000
PATH=./

all:
	CGO_ENABLED=0 /usr/local/go/bin/go build -o ./${BIN}

run:
	@./${BIN} --port=${PORT} --path=${PATH}

clean:
	rm -f ./${BIN}
