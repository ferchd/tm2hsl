package normalizer

import (
	"encoding/json"
	"testing"

	"github.com/ferchd/tm2hsl/internal/parser"
)

func TestNewNormalizer(t *testing.T) {
	n := NewNormalizer()
	if n == nil {
		t.Error("NewNormalizer() returned nil")
	}
	if n.supportedFeatures == nil {
		t.Error("supportedFeatures not initialized")
	}
	if !n.strictMode {
		t.Error("strictMode should be true by default")
	}
}

func TestNormalizer_Normalize(t *testing.T) {
	tests := []struct {
		name    string
		ast     *parser.TextMateAST
		wantErr bool
	}{
		{
			name: "valid simple grammar",
			ast: &parser.TextMateAST{
				ScopeName: "source.test",
				Patterns: []parser.GrammarRule{
					{
						Name:  "keyword",
						Match: "\\btest\\b",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid AST - missing scopeName",
			ast: &parser.TextMateAST{
				Name: "Test",
			},
			wantErr: true,
		},
		{
			name: "grammar with includes",
			ast: &parser.TextMateAST{
				ScopeName: "source.test",
				Patterns: []parser.GrammarRule{
					{
						Include: "#keywords",
					},
				},
				Repository: map[string]parser.GrammarRule{
					"keywords": {
						Patterns: []parser.GrammarRule{
							{
								Name:  "keyword",
								Match: "\\bif\\b",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNormalizer()
			_, err := n.Normalize(tt.ast)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalizer.Normalize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalizer_validateAST(t *testing.T) {
	tests := []struct {
		name    string
		ast     *parser.TextMateAST
		wantErr bool
	}{
		{
			name: "valid AST",
			ast: &parser.TextMateAST{
				ScopeName: "source.test",
			},
			wantErr: false,
		},
		{
			name:    "nil AST",
			ast:     nil,
			wantErr: true,
		},
		{
			name: "missing scopeName",
			ast: &parser.TextMateAST{
				Name: "Test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNormalizer()
			err := n.validateAST(tt.ast)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalizer.validateAST() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func FuzzNormalizer_Normalize(f *testing.F) {
	// Add seed ASTs
	f.Add(`{"scopeName": "source.test", "patterns": [{"match": "test"}]}`)

	f.Fuzz(func(t *testing.T, jsonData string) {
		// Try to parse as AST
		ast := &parser.TextMateAST{}
		if err := json.Unmarshal([]byte(jsonData), ast); err != nil {
			return // Skip invalid JSON
		}

		n := NewNormalizer()
		// Just ensure it doesn't panic
		n.Normalize(ast)
	})
}
