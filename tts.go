package tts

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/go-audio/wav"
)

const (
	TargetSampleRate = 16000
	TargetChannels   = 1
)

type Transcriber struct {
	model     whisper.Model
	Language  string
	Translate bool
}

func NewTranscriber(modelPath string, lang string, translate bool) (*Transcriber, error) {
	model, err := whisper.New(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	return &Transcriber{
		model:     model,
		Language:  lang,
		Translate: translate,
	}, nil
}

func (t *Transcriber) Close() error {
	return t.model.Close()
}

func (t *Transcriber) TranscribeFromBase64(b64String string) (string, error) {
	prefix := "data:audio/wav;base64,"
	if strings.HasPrefix(b64String, prefix) {
		b64String = strings.TrimPrefix(b64String, prefix)
	}

	audioBytes, err := base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	return t.TranscribeFromBytes(audioBytes)
}

func (t *Transcriber) TranscribeFromFile(filePath string) (string, error) {
	audioBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return t.TranscribeFromBytes(audioBytes)
}

func (t *Transcriber) TranscribeFromBytes(audioBytes []byte) (string, error) {
	cleanAudioBytes, err := denoiseAndSpedUp(audioBytes)
	if err != nil {
		fmt.Printf("WARNING: denoise and spedup failed, using original audio: %v\n", err)
		cleanAudioBytes = audioBytes
	}

	audioData, err := parseWavData(cleanAudioBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse wav: %w", err)
	}

	if len(audioData) > TargetSampleRate*2 {
		audioData = audioData[TargetSampleRate : len(audioData)-TargetSampleRate]
	}

	transcription, err := t.transcribe(audioData)
	if err != nil {
		return "", fmt.Errorf("failed to transcribe: %w", err)
	}

	return transcription, nil
}

func (t *Transcriber) transcribe(audioData []float32) (string, error) {
	context, err := t.model.NewContext()
	if err != nil {
		return "", err
	}
	context.SetLanguage(t.Language)
	context.SetTranslate(t.Translate)

	if err := context.Process(audioData, nil, nil, nil); err != nil {
		return "", err
	}

	var b strings.Builder

	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		b.WriteString(segment.Text)
	}

	return b.String(), nil
}

// validate wav and convert to pcm float32
func parseWavData(data []byte) ([]float32, error) {
	reader := bytes.NewReader(data)
	decoder := wav.NewDecoder(reader)

	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid wav")
	}

	format := decoder.Format()

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	pcmData := buf.Data
	pcmData = convertWavFormat(pcmData, format.SampleRate, format.NumChannels)

	return float32Buffer(pcmData, buf.SourceBitDepth), nil
}

// to 16khz / 1ch
func convertWavFormat(pcmData []int, sampleRate, numChannels int) []int {
	var processedData []int

	// to 1ch
	if numChannels > TargetChannels {
		pcmDataLen := len(pcmData)
		processedData = make([]int, 0, pcmDataLen/numChannels)
		for i := 0; i < pcmDataLen; i += numChannels {
			sum := 0
			// normalizing
			for c := 0; c < numChannels; c++ {
				sum += pcmData[i+c]
			}
			processedData = append(processedData, sum/numChannels)
		}
	} else {
		processedData = pcmData
	}

	// to 16khz
	if sampleRate != TargetSampleRate {
		return resampleLinear(processedData, sampleRate)
	}

	return processedData
}

// to 16khz
func resampleLinear(data []int, srcRate int) []int {
	if srcRate == TargetSampleRate || len(data) == 0 {
		return data
	}

	ratio := float64(srcRate) / float64(TargetSampleRate)
	newLen := int(float64(len(data)) / ratio)
	out := make([]int, newLen)

	for i := 0; i < newLen; i++ {
		srcPos := float64(i) * ratio
		j := int(srcPos)
		if j >= len(data)-1 {
			out[i] = data[len(data)-1]
			continue
		}

		frac := srcPos - float64(j)
		v := float64(data[j])*(1-frac) + float64(data[j+1])*frac
		out[i] = int(v)
	}

	return out
}

// eg -32768..32767 -> -1..1
func float32Buffer(data []int, bitDepth int) []float32 {
	res := make([]float32, len(data))

	var scale float32
	switch bitDepth {
	case 8:
		scale = 128.0
	case 16:
		scale = 32768.0
	case 24:
		scale = 8388608.0
	case 32:
		scale = 2147483648.0
	default:
		scale = 32768.0
	}

	for i, v := range data {
		res[i] = float32(v) / scale
	}

	return res
}

// removes noise and speds up
func denoiseAndSpedUp(audioBytes []byte) ([]byte, error) {
	cmd := exec.Command("ffmpeg",
		"-i", "pipe:0",
		"-af", "highpass=f=200,lowpass=f=3000,afftdn=nf=-25,atempo=1.25",
		"-f", "wav",
		"pipe:1",
	)
	cmd.Stdin = bytes.NewReader(audioBytes)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg processing failed: %v; stderr: %s", err, stderr.String())
	}

	return out.Bytes(), nil
}
