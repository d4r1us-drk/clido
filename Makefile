all: tuido

WARNINGS = -Wall
DEBUG = -ggdb -fno-omit-frame-pointer
OPTIMIZE = -O2

passgen: Makefile tuido.c
	$(CC) -o $@ $(WARNINGS) $(DEBUG) $(OPTIMIZE) tuido.c

clean:
	rm -f tuido

install:
	echo "Installing is not supported"

run:
	./tuido

