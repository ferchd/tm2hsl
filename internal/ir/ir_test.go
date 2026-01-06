package ir

import (
	"testing"

	"github.com/ferchd/tm2hsl/internal/parser"
)

func TestBuildFromAST(t *testing.T) {
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
			wantErr: false, // Currently doesn't error on nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm, err := BuildFromAST(tt.ast)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildFromAST() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && sm != nil {
				if tt.ast != nil && sm.Name != tt.ast.ScopeName {
					t.Errorf("BuildFromAST() Name = %v, want %v", sm.Name, tt.ast.ScopeName)
				}
			}
		})
	}
}

func TestPredicateBuilder_FromTextMateMatch(t *testing.T) {
	pb := NewPredicateBuilder(false)

	tests := []struct {
		name     string
		pattern  string
		wantErr  bool
		expected PredicateType
	}{
		{
			name:     "any character",
			pattern:  ".",
			wantErr:  false,
			expected: PredicateAny,
		},
		{
			name:     "regex pattern",
			pattern:  "\\w+",
			wantErr:  false,
			expected: PredicateRegex,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pred, err := pb.FromTextMateMatch(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("PredicateBuilder.FromTextMateMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && pred.Type() != tt.expected {
				t.Errorf("PredicateBuilder.FromTextMateMatch() type = %v, want %v", pred.Type(), tt.expected)
			}
		})
	}
}

func TestActionBuilder_FromTextMateCaptures(t *testing.T) {
	ab := NewActionBuilder()

	tests := []struct {
		name     string
		captures map[int]parser.Capture
		wantLen  int
	}{
		{
			name: "simple captures",
			captures: map[int]parser.Capture{
				0: {Name: "match"},
				1: {Name: "group1"},
			},
			wantLen: 2,
		},
		{
			name:     "empty captures",
			captures: map[int]parser.Capture{},
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actions, err := ab.FromTextMateCaptures(tt.captures)
			if err != nil {
				t.Errorf("ActionBuilder.FromTextMateCaptures() error = %v", err)
				return
			}
			if len(actions) != tt.wantLen {
				t.Errorf("ActionBuilder.FromTextMateCaptures() len = %v, want %v", len(actions), tt.wantLen)
			}
		})
	}
}

func TestPredicates(t *testing.T) {
	// Test AnyPredicate
	ap := &AnyPredicate{}
	if ap.Type() != PredicateAny {
		t.Errorf("AnyPredicate.Type() = %v, want %v", ap.Type(), PredicateAny)
	}
	if ap.String() != "any" {
		t.Errorf("AnyPredicate.String() = %v, want %v", ap.String(), "any")
	}

	// Test RegexPredicate
	rp := &RegexPredicate{Pattern: "\\w+"}
	if rp.Type() != PredicateRegex {
		t.Errorf("RegexPredicate.Type() = %v, want %v", rp.Type(), PredicateRegex)
	}
	if rp.Equal(&RegexPredicate{Pattern: "\\w+"}) != true {
		t.Error("RegexPredicate.Equal() should return true for same pattern")
	}
}

func TestActions(t *testing.T) {
	// Test PushScopeAction
	psa := &PushScopeAction{Scope: "test"}
	if psa.Type() != ActionPushScope {
		t.Errorf("PushScopeAction.Type() = %v, want %v", psa.Type(), ActionPushScope)
	}

	// Test EmitTokenAction
	eta := &EmitTokenAction{Text: "test"}
	if eta.Type() != ActionEmitToken {
		t.Errorf("EmitTokenAction.Type() = %v, want %v", eta.Type(), ActionEmitToken)
	}
}
