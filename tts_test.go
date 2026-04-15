package tts

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	JSONFilePath = "response.json"
	ModelPath    = "ggml-small.bin"
)

type JSONWithAudioB64 struct {
	Audio string `json:"audio"`
}

func TestTranscribeFromBase64(t *testing.T) {
	data, err := os.ReadFile(JSONFilePath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", JSONFilePath, err)
	}

	var resp JSONWithAudioB64
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if resp.Audio == "" {
		t.Fatal("audio string is empty")
	}

	if _, err := os.Stat(ModelPath); os.IsNotExist(err) {
		t.Skipf("bypass test: model file %s not found", ModelPath)
	}

	transcriber, err := NewTranscriber(ModelPath, "ru", false)
	if err != nil {
		t.Fatalf("failed to initialize transcriber: %v", err)
	}
	defer transcriber.Close()

	result, err := transcriber.TranscribeFromBase64(resp.Audio)
	if err != nil {
		t.Errorf("TranscribeFromBase64 error: %v", err)
	}

	if result == "" {
		t.Error("result is empty; expected result")
	}

	t.Logf("transcribed text: %s", result)
}
