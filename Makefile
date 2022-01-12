ABIDIR=./abi
GENDIR=./ERC721

all: $(ABIDIR)/*
	for file in $^; do \
		f=$${file##*/}; \
		filename=$${f%.*}; \
		abigen  --abi=$^ --pkg=$${filename} --out=${GENDIR}/$${filename}.go; \
	done \

