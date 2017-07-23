.PHONY: clean build run

ARTIFACT := owlclicker

clean:
	-rm -f ${ARTIFACT}
	-mkdir -p functionmetrics

build:
	make clean
	go build -o ${ARTIFACT} .

run:
	make build
	./${ARTIFACT}
