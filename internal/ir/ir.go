package ir

import (
	"regexp"
)

type Program struct {
	Version     uint16
	Name        string
	Scope       string
	RegexTable  []RegexEntry
	StateTable  []StateEntry
	RuleTable   []RuleEntry
	ScopeTable  []ScopeEntry
	StringTable []string
}

type RegexEntry struct {
	ID       uint32
	Pattern  string
	Compiled *regexp.Regexp
	Bytecode []byte // Para regex compiladas a bytecode
}

type StateEntry struct {
	ID         uint32
	RuleOffset uint32
	RuleCount  uint16
	Flags      StateFlags
}

type StateFlags uint8

const (
	StateFinal StateFlags = 1 << iota
	StatePush
	StatePop
)

type RuleEntry struct {
	RegexID    uint32
	Action     RuleAction
	NextState  int32 // -1: pop, -2: stay, >=0: state ID
	ScopeID    uint16
	Priority   uint8
	CaptureMap []CaptureMapping
}

type CaptureMapping struct {
	Group   uint8
	ScopeID uint16
}

type ScopeEntry struct {
	ID   uint16
	Name string
}

type RuleAction uint8

const (
	RuleActionMatch RuleAction = iota
	RuleActionPushScope
	RuleActionPopScope
	RuleActionTransition
)

func NewProgram(name, scope string) *Program {
	return &Program{
		Version:     1,
		Name:        name,
		Scope:       scope,
		RegexTable:  []RegexEntry{},
		StateTable:  []StateEntry{},
		RuleTable:   []RuleEntry{},
		ScopeTable:  []ScopeEntry{},
		StringTable: []string{},
	}
}

func (p *Program) AddRegex(pattern string) uint32 {
	// Buscar regex duplicada
	for i, entry := range p.RegexTable {
		if entry.Pattern == pattern {
			return uint32(i)
		}
	}

	id := uint32(len(p.RegexTable))
	compiled, _ := regexp.Compile(pattern)

	p.RegexTable = append(p.RegexTable, RegexEntry{
		ID:       id,
		Pattern:  pattern,
		Compiled: compiled,
	})

	return id
}

func (p *Program) AddState(rules []RuleEntry, flags StateFlags) uint32 {
	stateID := uint32(len(p.StateTable))

	ruleOffset := uint32(len(p.RuleTable))
	p.RuleTable = append(p.RuleTable, rules...)

	p.StateTable = append(p.StateTable, StateEntry{
		ID:         stateID,
		RuleOffset: ruleOffset,
		RuleCount:  uint16(len(rules)),
		Flags:      flags,
	})

	return stateID
}

func (p *Program) AddScope(name string) uint16 {
	// Buscar scope duplicado
	for i, scope := range p.ScopeTable {
		if scope.Name == name {
			return uint16(i)
		}
	}

	id := uint16(len(p.ScopeTable))
	p.ScopeTable = append(p.ScopeTable, ScopeEntry{
		ID:   id,
		Name: name,
	})

	// Añadir a tabla de strings si es nuevo
	p.AddString(name)

	return id
}

func (p *Program) AddString(str string) uint32 {
	// Buscar string duplicada
	for i, s := range p.StringTable {
		if s == str {
			return uint32(i)
		}
	}

	id := uint32(len(p.StringTable))
	p.StringTable = append(p.StringTable, str)
	return id
}

func (p *Program) Statistics() ProgramStats {
	return ProgramStats{
		RegexCount:  len(p.RegexTable),
		StateCount:  len(p.StateTable),
		RuleCount:   len(p.RuleTable),
		ScopeCount:  len(p.ScopeTable),
		StringCount: len(p.StringTable),
	}
}

type ProgramStats struct {
	RegexCount  int
	StateCount  int
	RuleCount   int
	ScopeCount  int
	StringCount int
}
