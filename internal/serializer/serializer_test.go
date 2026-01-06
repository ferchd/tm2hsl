package serializer

import (
	"bytes"
	"testing"

	"github.com/ferchd/tm2hsl/pkg/hsl"
)

func TestNewSerializer(t *testing.T) {
	s := NewSerializer()
	if s == nil {
		t.Error("NewSerializer() returned nil")
	}
	if s.byteOrder == nil {
		t.Error("byteOrder not initialized")
	}
}

func TestSerializer_Serialize(t *testing.T) {
	s := NewSerializer()

	tests := []struct {
		name     string
		bytecode *hsl.Bytecode
		wantErr  bool
	}{
		{
			name: "nil bytecode",
			bytecode: &hsl.Bytecode{
				Header: hsl.Header{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := s.Serialize(tt.bytecode, &buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Serializer.Serialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
