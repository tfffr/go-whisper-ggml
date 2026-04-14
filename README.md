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

## 3. Edit Makefile
Replace `WHISPER_DIR` in Makefile with your actual path to `whisper.cpp`

## 4. Run the Go project through Makefile

```bash
make run
```
or to build:
```bash
make build
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
