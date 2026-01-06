// bytecode.go - ABI estable, no solo dumping de structs
package serializer

import (
	"bytes"
	"encoding/binary"

	"github.com/ferchd/tm2hsl/internal/ir"
)

// Constants
const stateEntrySize = 32 // Size of StateEntry in bytes

// Stub functions - TODO: Implement properly
func serializeStates(states map[ir.StateID]*ir.State) []StateEntry {
	// TODO: Implement state serialization
	return []StateEntry{}
}

func serializeTokens(tokens map[ir.TokenID]ir.TokenDef) []TokenEntry {
	// TODO: Implement token serialization
	return []TokenEntry{}
}

func serializeActions(actions map[ir.ActionID]ir.Action) []ActionEntry {
	// TODO: Implement action serialization
	return []ActionEntry{}
}

func computeChecksum(data []byte) uint32 {
	// Simple checksum - TODO: Use proper CRC32
	var sum uint32
	for _, b := range data {
		sum += uint32(b)
	}
	return sum
}

// BytecodeHeader - Fixed header
type BytecodeHeader struct {
	Magic           [4]byte // 'HSL1'
	Version         uint16  // 0x0001
	Flags           uint16
	NumStates       uint32
	NumTokens       uint32
	NumActions      uint32
	EntryPoint      uint32 // Offset del estado inicial
	StateTable      uint32 // Offset to state table
	StateTableSize  uint32 // Size of state table
	TokenTable      uint32 // Offset to token table
	TokenTableSize  uint32 // Size of token table
	ActionTable     uint32 // Offset to action table
	ActionTableSize uint32 // Size of action table
	StringTable     uint32 // Offset to string table
	StringTableSize uint32 // Size of string table
	RegexTable      uint32 // Offset to regex table
	RegexTableSize  uint32 // Size of regex table
	HeaderSize      uint32 // Total header size
	Checksum        uint32
}

// BytecodeLayout - Predefined layout
type BytecodeLayout struct {
	Header      BytecodeHeader
	StateTable  []StateEntry  // Tabla ordenada
	TokenTable  []TokenEntry  // Tabla ordenada por ID
	ActionTable []ActionEntry // Tabla ordenada
	StringPool  []StringEntry // Pool de strings
	RegexPool   []RegexEntry  // Pool de regex compiladas
}

// Serialize - Stable format writer
func Serialize(machine *ir.StateMachine) ([]byte, error) {
	layout := BytecodeLayout{}

	// 1. Escribir header con offsets calculados
	header := BytecodeHeader{
		Magic:     [4]byte{'H', 'S', 'L', '1'},
		Version:   0x0001,
		NumStates: uint32(len(machine.States)),
		NumTokens: uint32(len(machine.Tokens)),
	}

	// 2. Ordenar tablas por ID (orden determinista)
	layout.StateTable = serializeStates(machine.States)
	layout.TokenTable = serializeTokens(machine.Tokens)
	layout.ActionTable = serializeActions(machine.Actions)

	// 3. Calcular offsets y checksum
	buf := &bytes.Buffer{}

	// Header at fixed offset 0
	binary.Write(buf, binary.LittleEndian, header)

	// Tables at predictable offsets
	stateOffset := binary.Size(header)
	_ = stateOffset // TODO: Use stateOffset
	// etc.

	// 4. Escribir con padding para alineación
	writeAligned(buf, layout.StateTable, 8)
	writeAligned(buf, layout.TokenTable, 8)
	writeAligned(buf, layout.ActionTable, 8)

	// 5. Calcular checksum final
	data := buf.Bytes()
	checksum := computeChecksum(data)
	binary.LittleEndian.PutUint32(data[20:], checksum)

	return data, nil
}

// writeAligned - Garantiza alineación para acceso rápido
func writeAligned(buf *bytes.Buffer, data interface{}, align int) {
	current := buf.Len()
	padding := (align - (current % align)) % align
	buf.Write(make([]byte, padding))
	binary.Write(buf, binary.LittleEndian, data)
}
