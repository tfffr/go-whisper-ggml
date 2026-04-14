# go-whisper-ggml

## What is needed on a new machine

* Go
* Git
* CMake
* C/C++ compiler toolchain for cgo (`gcc`, `g++`, `make` or `ninja`)
* (optional) `ffmpeg` — if you plan to work with formats other than WAV

> This project uses CGO via `whisper.cpp`, so a working C/C++ toolchain is required.

---

## 1. Download and build `whisper.cpp`

```bash
git clone https://github.com/ggml-org/whisper.cpp.git
cd whisper.cpp

cmake -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build -j
```

---

## 2. Download a ggml model

`small` model is used as an example. You can choose any available model.

List of models:
https://github.com/ggml-org/whisper.cpp/tree/master/models

```bash
sh ./models/download-ggml-model.sh small
```

The model will be placed at:

```bash
./models/ggml-small.bin
```

---

## 3. Configure Makefile

This project relies on CGO, so you must provide paths to `whisper.cpp`.

Edit `WHISPER_DIR` in Makefile:

```make
WHISPER_DIR ?= /absolute/path/to/whisper.cpp
```

Or override it at runtime:

```bash
make run WHISPER_DIR=/absolute/path/to/whisper.cpp
```

---

## 4. Copy Makefile to your main project

Any project that imports this package must use the same CGO configuration.

---

## 5. Run the project

Make sure you provide a valid path to the `.bin` model.

```bash
make run WHISPER_DIR=/absolute/path/to/whisper.cpp MODEL_PATH=/absolute/path/to/ggml-small.bin
```

Or if your code reads model path from a constant or config, update it accordingly.

Build binary:

```bash
make build
```

Example usage:

```go
transcriber, _ := tts.NewTranscriber("ggml-small.bin")
defer transcriber.Close()

text, _ := transcriber.TranscribeFromFile("example.wav")
fmt.Println(text)
```

---

## Notes

* CGO flags (`CGO_CFLAGS`, `CGO_LDFLAGS`) are required for linking with `whisper.cpp`
* These flags are defined inside the Makefile for convenience
* Running `go run` without them will fail
