# go-whisper-ggml
Run:
```bash
git clone https://github.com/ggerganov/whisper.cpp.git
cd whisper.cpp
git clone https://github.com/openai/whisper.git ./openai-whisper
uv add huggingface_hub transformers torch
uv run python -c 'from huggingface_hub import snapshot_download; snapshot_download(repo_id="path/to/model/from/hugging/face", local_dir="./my_model")'
uv run models/convert-h5-to-ggml.py ./my_model ./openai-whisper ./my_model
```
Your model is now located at ./my_model/ggml-model.bin.
