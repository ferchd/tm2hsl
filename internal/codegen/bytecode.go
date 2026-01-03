package codegen

import (
	"encoding/binary"
	"hash/crc32"

	"github.com/ferchd/tm2hsl/internal/ir"
	"github.com/ferchd/tm2hsl/pkg/hsl"
)

type BytecodeGenerator struct {
	program  *ir.Program
	bytecode *hsl.Bytecode
}

func NewGenerator(program *ir.Program) *BytecodeGenerator {
	return &BytecodeGenerator{
		program:  program,
		bytecode: &hsl.Bytecode{},
	}
}

func (g *BytecodeGenerator) Generate() (*hsl.Bytecode, error) {
	// 1. Cabecera
	g.generateHeader()

	// 2. Tabla de strings
	g.generateStringTable()

	// 3. Tabla de regex
	g.generateRegexTable()

	// 4. Tabla de scopes
	g.generateScopeTable()

	// 5. Tabla de estados
	g.generateStateTable()

	// 6. Tabla de reglas
	g.generateRuleTable()

	// 7. Calcular checksum
	g.calculateChecksum()

	return g.bytecode, nil
}

func (g *BytecodeGenerator) generateHeader() {
	g.bytecode.Header = hsl.Header{
		Magic:             [4]byte{'H', 'S', 'L', '1'},
		Version:           g.program.Version,
		HeaderSize:        uint16(binary.Size(hsl.Header{})),
		StringTableOffset: 0,
		RegexTableOffset:  0,
		ScopeTableOffset:  0,
		StateTableOffset:  0,
		RuleTableOffset:   0,
		TotalSize:         0,
		Checksum:          0,
		Flags:             hsl.FlagValidated | hsl.FlagOptimized,
	}

	// Metadata
	g.bytecode.Name = g.program.Name
	g.bytecode.Scope = g.program.Scope
}

func (g *BytecodeGenerator) generateStringTable() {
	// Serializar strings con offsets
	var data []byte
	offsets := make([]uint32, len(g.program.StringTable))

	for i, str := range g.program.StringTable {
		offsets[i] = uint32(len(data))
		data = append(data, []byte(str)...)
		data = append(data, 0) // Null terminator
	}

	g.bytecode.StringTable = hsl.StringTable{
		Count:   uint32(len(g.program.StringTable)),
		Offsets: offsets,
		Data:    data,
	}
}

func (g *BytecodeGenerator) generateRegexTable() {
	regexes := make([]hsl.RegexEntry, len(g.program.RegexTable))

	for i, re := range g.program.RegexTable {
		// Compilar regex a bytecode interno si es posible
		bytecode := g.compileRegexToBytecode(re.Pattern)

		regexes[i] = hsl.RegexEntry{
			ID:          re.ID,
			PatternHash: g.hashString(re.Pattern),
			Bytecode:    bytecode,
			Flags:       0,
		}
	}

	g.bytecode.RegexTable = hsl.RegexTable{
		Count:   uint32(len(regexes)),
		Entries: regexes,
	}
}

func (g *BytecodeGenerator) compileRegexToBytecode(pattern string) []byte {
	// Implementación simplificada
	// En producción, usar un compilador de regex a DFA
	return []byte(pattern) // Placeholder
}

func (g *BytecodeGenerator) generateScopeTable() {
	scopes := make([]hsl.ScopeEntry, len(g.program.ScopeTable))

	for i, scope := range g.program.ScopeTable {
		// Buscar string ID en tabla
		strID := g.findStringID(scope.Name)

		scopes[i] = hsl.ScopeEntry{
			ID:       scope.ID,
			NameID:   strID,
			ParentID: 0xFFFF, // Ninguno por defecto
		}
	}

	g.bytecode.ScopeTable = hsl.ScopeTable{
		Count:   uint32(len(scopes)),
		Entries: scopes,
	}
}

func (g *BytecodeGenerator) generateStateTable() {
	states := make([]hsl.StateEntry, len(g.program.StateTable))

	for i, state := range g.program.StateTable {
		states[i] = hsl.StateEntry{
			ID:         state.ID,
			RuleOffset: state.RuleOffset,
			RuleCount:  state.RuleCount,
			Flags:      uint8(state.Flags),
		}
	}

	g.bytecode.StateTable = hsl.StateTable{
		Count:   uint32(len(states)),
		Entries: states,
	}
}

func (g *BytecodeGenerator) generateRuleTable() {
	rules := make([]hsl.RuleEntry, len(g.program.RuleTable))

	for i, rule := range g.program.RuleTable {
		rules[i] = hsl.RuleEntry{
			RegexID:      rule.RegexID,
			Action:       uint8(rule.Action),
			NextState:    int16(rule.NextState),
			ScopeID:      rule.ScopeID,
			Priority:     rule.Priority,
			CaptureCount: uint8(len(rule.CaptureMap)),
		}

		// Añadir mapeos de captura si existen
		if len(rule.CaptureMap) > 0 {
			for _, cap := range rule.CaptureMap {
				rules[i].Captures = append(rules[i].Captures, hsl.CaptureMapping{
					Group:   cap.Group,
					ScopeID: cap.ScopeID,
				})
			}
		}
	}

	g.bytecode.RuleTable = hsl.RuleTable{
		Count:   uint32(len(rules)),
		Entries: rules,
	}
}

func (g *BytecodeGenerator) calculateChecksum() {
	// Calcular CRC32 de todo el contenido excepto el checksum mismo
	// Implementación simplificada
	data := g.serializeForChecksum()
	g.bytecode.Header.Checksum = crc32.ChecksumIEEE(data)
}

func (g *BytecodeGenerator) serializeForChecksum() []byte {
	// Serializar para cálculo de checksum
	// Implementación simplificada
	return []byte{}
}

func (g *BytecodeGenerator) findStringID(str string) uint32 {
	for i, s := range g.program.StringTable {
		if s == str {
			return uint32(i)
		}
	}
	return 0xFFFFFFFF
}

func (g *BytecodeGenerator) hashString(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}
