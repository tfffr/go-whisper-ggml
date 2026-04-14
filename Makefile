WHISPER_DIR ?= /home/jeff/dev/python/whisper.cpp

export CGO_CFLAGS := -I$(WHISPER_DIR)/include -I$(WHISPER_DIR)/ggml/include
export CGO_LDFLAGS := -L$(WHISPER_DIR)/build/src -L$(WHISPER_DIR)/build/ggml/src -lwhisper -lggml -lstdc++ -lm

run:
	go run .

build:
	go build -o app .
