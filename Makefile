all: clido

WARNINGS = -Wall
DEBUG = -ggdb -fno-omit-frame-pointer
OPTIMIZE = -O2

passgen: Makefile clido.cpp
	$(CPP) -o $@ $(WARNINGS) $(DEBUG) $(OPTIMIZE) clido.cpp

clean:
	rm -f clido

install:
	echo "Installing is not supported"

run:
	./clido

