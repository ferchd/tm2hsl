package parser

import (
	"strings"
	"testing"
)

func TestLoadGrammar_JSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected *TextMateAST
	}{
		{
			name: "valid JSON grammar",
			input: `{
				"scopeName": "source.test",
				"name": "Test Grammar",
				"patterns": [
					{
						"name": "keyword.control",
						"match": "\\b(if|else)\\b"
					}
				]
			}`,
			wantErr: false,
			expected: &TextMateAST{
				ScopeName: "source.test",
				Name:      "Test Grammar",
				Patterns: []GrammarRule{
					{
						Name:  "keyword.control",
						Match: "\\b(if|else)\\b",
					},
				},
			},
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := LoadGrammar(reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGrammar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if got.ScopeName != tt.expected.ScopeName {
					t.Errorf("LoadGrammar() ScopeName = %v, want %v", got.ScopeName, tt.expected.ScopeName)
				}
				if got.Name != tt.expected.Name {
					t.Errorf("LoadGrammar() Name = %v, want %v", got.Name, tt.expected.Name)
				}
			}
		})
	}
}

func TestLoadGrammar_XML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "valid XML plist",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>scopeName</key>
	<string>source.test</string>
	<key>name</key>
	<string>Test Grammar</string>
</dict>
</plist>`,
			wantErr: false,
		},
		{
			name:    "invalid XML",
			input:   `<invalid>`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			_, err := LoadGrammar(reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGrammar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTextMateAST_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ast     TextMateAST
		wantErr bool
	}{
		{
			name: "valid AST",
			ast: TextMateAST{
				ScopeName: "source.test",
			},
			wantErr: false,
		},
		{
			name: "missing scopeName",
			ast: TextMateAST{
				Name: "Test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ast.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("TextMateAST.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func FuzzLoadGrammar(f *testing.F) {
	// Add seed inputs
	f.Add(`{"scopeName": "source.test", "name": "Test"}`)
	f.Add(`<?xml version="1.0"?><plist><dict><key>scopeName</key><string>test</string></dict></plist>`)

	f.Fuzz(func(t *testing.T, input string) {
		reader := strings.NewReader(input)
		// Just ensure it doesn't panic
		LoadGrammar(reader)
	})
}
