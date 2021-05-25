CFLAGS=-std=c11 -g -static

9cc: cmd/compiler/main.go
	go build -o 9cc cmd/compiler/main.go

test: 9cc
	./test.sh

clean:
	rm -f 9cc *.o *~ tmp*

.PHONY: test clean
