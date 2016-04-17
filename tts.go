package termux

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

var (
	// ErrNoTTSEngineProvided is returned if there's no engine specified
	ErrNoTTSEngineProvided = errors.New("No engine specified for TTS")

	// ErrNoTTSContentProvided is returned if no content is given to speak
	ErrNoTTSContentProvided = errors.New("No content given to TTS engine")
)

// TTSEngine is the data returned by TTSEngines. TODO: Need example output.
type TTSEngine struct {
}

// TTSEngines returns the list of available TTS engines
func TTSEngines() ([]TTSEngine, error) {
	var resp []TTSEngine
	respBytes, err := toolExec(nil, "TextToSpeech", "--es", "engine", "LIST_AVAILABLE")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBytes, &resp)
	return resp, err
}

// TTSSpeak speaks the provided content.
func TTSSpeak(content, engine, language string, pitch, rate float64) error {
	var args []string
	if engine == "" {
		return ErrNoTTSEngineProvided
	}
	args = append(args, []string{"--es", "engine", engine}...)
	if content == "" {
		return ErrNoTTSContentProvided
	}
	if language != "" {
		args = append(args, []string{"--es", "language", language}...)
	}
	if pitch > 0.001 {
		args = append(args, []string{"--ef", "pitch", strconv.FormatFloat(pitch, 'f', -1, 32)}...)
	}
	if rate > 0.001 {
		args = append(args, []string{"--ef", "rate", strconv.FormatFloat(rate, 'f', -1, 32)}...)
	}
	_, err := toolExec(bytes.NewBuffer([]byte(content)), "TextToSpeech", args...)
	return err
}
