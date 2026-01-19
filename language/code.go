package language

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

//go:embed iso.list.jsonl
var languageList string

type langEntry struct {
	ID      string `json:"id"`
	IsISO   string `json:"isIso"`
	Display string `json:"display"`
}

var displayNameToTag = make(map[string]language.Tag)

var ErrUnknownLanguage = errors.New("unknown language")

func init() {
	decoder := json.NewDecoder(strings.NewReader(languageList))
	for decoder.More() {
		var entry langEntry
		err := decoder.Decode(&entry)
		if err != nil {
			fmt.Printf("Error decoding language entry: %v\n", err)
			continue
		}

		tag, err := language.Parse(entry.ID)
		if err != nil {
			fmt.Printf("Error parsing language ID '%s': %v\n", entry.ID, err)
			continue
		}

		displayNameToTag[strings.ToLower(entry.Display)] = tag
		displayNameToTag[strings.ToLower(entry.ID)] = tag // Also allow parsing by ID (e.g., "en")
	}
}

// Parse attempts to parse a language string.
// It first tries to parse it as a language tag (e.g., "en-US").
// If that fails, it tries to look it up as an English display name (e.g., "English").
func Parse(lang string) (language.Tag, error) {
	tag, err := language.Parse(lang)
	if err == nil {
		return tag, nil
	}

	lowerLang := strings.ToLower(lang)
	if tag, ok := displayNameToTag[lowerLang]; ok {
		return tag, nil
	}

	return language.Und, fmt.Errorf("%w: %s", ErrUnknownLanguage, lang)
}
