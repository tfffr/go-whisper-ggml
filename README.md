# go-whisper-ggml

## What is needed on a new machine

* Go
* Git
* CMake
* C/C++ compiler toolchain for cgo (`gcc`, `g++`, `make` or `ninja`)
* `whisper.cpp` source tree
* ggml model file (`ggml-small.bin`, etc.)

## 1. Download and build `whisper.cpp`

```bash
git clone https://github.com/ggml-org/whisper.cpp.git
cd whisper.cpp

cmake -B build
cmake --build build -j --config Release
```

## 2. Download a ggml model

```bash
sh ./models/download-ggml-model.sh small
```

The model will be placed here:

```bash
./models/ggml-small.bin
```

## 3. Point Go/cgo to the built library

```bash
export WHISPER_DIR="/absolute/path/to/whisper.cpp"
export CGO_CFLAGS="-I${WHISPER_DIR}/include -I${WHISPER_DIR}/ggml/include"
export CGO_LDFLAGS="-L${WHISPER_DIR}/build/src -L${WHISPER_DIR}/build/ggml/src -lwhisper -lggml -lstdc++ -lm"
```

## 4. Run the Go project

```bash
cd /your/go/project
go run .
```

## Notes

### If you use the default model path in code

If your code contains:

```go
const ModelPath = "ggml-small.bin"
```

then either:

* copy `./models/ggml-small.bin` into the project directory, or
* change the code to the full/relative actual path, for example:

```go
const ModelPath = "/absolute/path/to/whisper.cpp/models/ggml-small.bin"
```
