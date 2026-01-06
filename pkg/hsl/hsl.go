package hsl

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Header del archivo HSL
type Header struct {
	Magic             [4]byte
	Version           uint16
	HeaderSize        uint16
	StringTableOffset uint32
	RegexTableOffset  uint32
	ScopeTableOffset  uint32
	StateTableOffset  uint32
	RuleTableOffset   uint32
	TotalSize         uint32
	Checksum          uint32
	Flags             uint32
}

// Flags de cabecera
const (
	FlagValidated = 1 << iota
	FlagOptimized
	FlagDeterministic
	FlagLinearTime
)

// Bytecode completo
type Bytecode struct {
	Header      Header
	Name        string
	Scope       string
	StringTable StringTable
	RegexTable  RegexTable
	ScopeTable  ScopeTable
	StateTable  StateTable
	RuleTable   RuleTable
}

// Tablas
type StringTable struct {
	Count   uint32
	Offsets []uint32
	Data    []byte
}

type RegexTable struct {
	Count   uint32
	Entries []RegexEntry
}

type ScopeTable struct {
	Count   uint32
	Entries []ScopeEntry
}

type StateTable struct {
	Count   uint32
	Entries []StateEntry
}

type RuleTable struct {
	Count   uint32
	Entries []RuleEntry
}

// Entradas
type RegexEntry struct {
	ID          uint32
	PatternHash uint32
	Bytecode    []byte
	Flags       uint8
}

type ScopeEntry struct {
	ID       uint16
	NameID   uint32
	ParentID uint16
}

type StateEntry struct {
	ID         uint32
	RuleOffset uint32
	RuleCount  uint16
	Flags      uint8
}

type RuleEntry struct {
	RegexID      uint32
	Action       uint8
	NextState    int16
	ScopeID      uint16
	Priority     uint8
	CaptureCount uint8
	Captures     []CaptureMapping
}

type CaptureMapping struct {
	Group   uint8
	ScopeID uint16
}

// Validación
func (h *Header) Validate() error {
	if string(h.Magic[:]) != "HSL1" {
		return fmt.Errorf("magic number inválido")
	}
	if h.Version != 1 {
		return fmt.Errorf("versión no soportada: %d", h.Version)
	}
	return nil
}

// Serialización/Deserialización helpers
func ReadHeader(data []byte) (*Header, error) {
	if len(data) < binary.Size(Header{}) {
		return nil, fmt.Errorf("datos insuficientes para cabecera")
	}

	var header Header
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	return &header, nil
}
