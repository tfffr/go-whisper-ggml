# go-whisper-ggml
1. Build ggml binary:
```bash
git clone https://github.com/ggerganov/whisper.cpp.git
cd whisper.cpp
git clone https://github.com/openai/whisper.git ./openai-whisper
uv add transformers torch
bash ./models/download-ggml-model.sh small
```
Your model is now located at ./my_model/ggml-model.bin.

2. Compile whisper.cpp statically:
```bash
cmake -B build -DBUILD_SHARED_LIBS=OFF
cmake --build build --config Release
export WHISPER_DIR="/absolute/path/to/whisper.cpp"
export CGO_CFLAGS="-I${WHISPER_DIR}/include -I${WHISPER_DIR}/ggml/include"
export CGO_LDFLAGS="-L${WHISPER_DIR}/build/src -L${WHISPER_DIR}/build/ggml/src -lwhisper -lggml -lstdc++ -lm"
```

3. Run Go code with ggml binary:
```bash
cd /your/go/project/
go run .
```
