package termux

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// TTSEngine is the data returned by TTSEngines. TODO: Need example output.
type TTSEngine map[string]interface{}

// TTSEngines returns the list of available TTS engines
func TTSEngines() ([]TTSEngine, error) {
	return ttsEngines(toolExec)
}

// TTSEngines returns the list of available TTS engines
func ttsEngines(execF toolExecFunc) ([]TTSEngine, error) {
	var resp []TTSEngine
	respBytes, err := execF(nil, "TextToSpeech", "--es", "engine", "LIST_AVAILABLE")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBytes, &resp)
	return resp, err
}

// TTSSpeak speaks the provided content.
func TTSSpeak(content, engine, language string, pitch, rate float64) error {
	return ttsSpeak(toolExec, content, engine, language, pitch, rate)
}

func ttsSpeak(execF toolExecFunc, content, engine, language string, pitch, rate float64) error {
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
	_, err := execF(bytes.NewBuffer([]byte(content)), "TextToSpeech", args...)
	return err
}
