#!/bin/sh

bname="translate-mcp"
bdir="./cmd/${bname}"
oname="${bdir}/${bname}"

go \
	build \
	-v \
	./...

go \
	build \
	-v \
	-o "${oname}" \
	"${bdir}"
