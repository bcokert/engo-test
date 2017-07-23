.PHONY: clean build run

ARTIFACT := owlclicker

clean:
	-rm -f ${ARTIFACT}
	-rm -rf functionmetrics
	-mkdir -p functionmetrics

build:
	make clean
	go build -o ${ARTIFACT} .

run:
	make build
	./${ARTIFACT}
