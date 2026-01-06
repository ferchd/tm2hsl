package serializer

type BytecodeFileLayout struct {
	Header      BytecodeHeader
	StateTable  StateTableSection
	TokenTable  TokenTableSection
	ActionTable ActionTableSection
	StringTable StringTableSection
	RegexTable  RegexTableSection
}

// StateTableSection - Tabla de estados
type StateTableSection struct {
	Entries []StateEntry
}

type StateEntry struct {
	Flags           uint16
	NumTransitions  uint16
	FirstTransition uint32 // Offset en la tabla de transiciones
	OnEntryAction   uint32 // Índice en ActionTable (-1 si no hay)
	OnExitAction    uint32 // Índice en ActionTable (-1 si no hay)
	Reserved        [12]byte
}

// TransitionEntry - Entrada de transición (24 bytes)
type TransitionEntry struct {
	PredicateType uint8
	PredicateData uint32 // Índice/offset según tipo
	TargetState   uint32
	Priority      uint8
	ConsumeChars  uint8
	NumActions    uint16
	FirstAction   uint32 // Índice en ActionTable
	Reserved      [8]byte
}

// TokenTableSection - Tabla de tokens
type TokenTableSection struct {
	Entries []TokenEntry
}

type TokenEntry struct {
	NameOffset  uint32 // Offset en StringTable
	ScopeDepth  uint16
	Flags       uint16
	ParentToken uint32 // Índice del token padre (-1 si raíz)
	FirstChild  uint32 // Primer hijo en jerarquía
	NextSibling uint32 // Siguiente hermano
}

// ActionTableSection - Tabla de acciones
type ActionTableSection struct {
	Entries []ActionEntry
}

type ActionEntry struct {
	Type       uint8
	DataSize   uint16
	DataOffset uint32 // Offset en datos de acción
	NextAction uint32 // Índice de siguiente acción en secuencia
}

// ActionData - Datos específicos de acción
type PushScopeData struct {
	TokenID      uint32
	CaptureIndex int32 // -1 si no es captura
}

type PopScopeData struct {
	Count uint16
}

// StringTableSection - Pool de strings
type StringTableSection struct {
	Strings []StringEntry
}

type StringEntry struct {
	Length uint32
	Offset uint32 // Offset desde inicio de sección
	Hash   uint32 // Hash para búsqueda rápida
}

// RegexTableSection - Pool de regex compiladas
type RegexTableSection struct {
	Entries []RegexEntry
}

type RegexEntry struct {
	PatternOffset  uint32 // Offset en StringTable
	Flags          uint16
	NumGroups      uint16
	BytecodeSize   uint32
	BytecodeOffset uint32 // Offset al bytecode nativo
}
