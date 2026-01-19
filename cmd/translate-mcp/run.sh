#!/bin/sh

export OLLAMA_HOST=127.0.0.1:11435

model=translategemma:4b-it-q4_K_M
port=11981

./translate-mcp \
	-model "${model}" \
	-port ${port}
