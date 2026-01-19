#!/bin/sh

man echo |
	./cmd/translate-cli/translate-cli \
	-model translategemma:4b-it-q4_K_M \
	-source en \
	-target ja
