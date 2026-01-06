package parser

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

// TextMateAST - Exact representation of the TextMate file
type TextMateAST struct {
	ScopeName          string                 `json:"scopeName" xml:"scopeName"`
	FileTypes          []string               `json:"fileTypes" xml:"fileTypes>string"`
	UUID               string                 `json:"uuid" xml:"uuid"`
	Name               string                 `json:"name" xml:"name"`
	Patterns           []GrammarRule          `json:"patterns" xml:"array>dict"`
	Repository         map[string]GrammarRule `json:"repository" xml:"dict>key"`
	FirstLineMatch     string                 `json:"firstLineMatch" xml:"firstLineMatch"`
	FoldingStartMarker string                 `json:"foldingStartMarker" xml:"foldingStartMarker"`
	FoldingStopMarker  string                 `json:"foldingStopMarker" xml:"foldingStopMarker"`

	// Fields ignored in normalization but preserved
	HiddenFields map[string]interface{} `json:"-" xml:"-"`
}

// GrammarRule - Exact grammatical rule
type GrammarRule struct {
	Name           string          `json:"name,omitempty" xml:"key>name,omitempty"`
	Match          string          `json:"match,omitempty" xml:"string>match,omitempty"`
	Begin          string          `json:"begin,omitempty" xml:"string>begin,omitempty"`
	End            string          `json:"end,omitempty" xml:"string>end,omitempty"`
	ContentName    string          `json:"contentName,omitempty" xml:"string>contentName,omitempty"`
	Captures       map[int]Capture `json:"captures,omitempty" xml:"dict>captures,omitempty"`
	BeginCaptures  map[int]Capture `json:"beginCaptures,omitempty" xml:"dict>beginCaptures,omitempty"`
	EndCaptures    map[int]Capture `json:"endCaptures,omitempty" xml:"dict>endCaptures,omitempty"`
	Include        string          `json:"include,omitempty" xml:"string>include,omitempty"`
	Patterns       []GrammarRule   `json:"patterns,omitempty" xml:"array>patterns,omitempty"`
	RepositoryName string          `json:"-" xml:"-"` // Para tracking interno
}

// Capture - Captura exacta
type Capture struct {
	Name string `json:"name" xml:"string"`
}

// RawPattern - Para preservar patrones no soportados
type RawPattern struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

func LoadGrammar(r io.Reader) (*TextMateAST, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read grammar: %w", err)
	}

	var ast TextMateAST

	if err := json.Unmarshal(data, &ast); err == nil {
		return &ast, nil
	}

	if err := parsePlist(data, &ast); err == nil {
		return &ast, nil
	}

	return nil, fmt.Errorf("unsupported grammar format")
}

func parsePlist(data []byte, ast *TextMateAST) error {
	return xml.Unmarshal(data, ast)
}

func (ast *TextMateAST) Validate() error {
	if ast.ScopeName == "" {
		return fmt.Errorf("scopeName is required")
	}
	return nil
}
