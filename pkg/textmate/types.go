package textmate

// Grammar representa una gramática TextMate completa
type Grammar struct {
	ScopeName  string
	Patterns   []Rule
	Repository map[string]Rule
}

// Rule representa una regla en la gramática TextMate
type Rule struct {
	Name          string
	Match         string
	Begin         string
	End           string
	ContentName   string
	Patterns      []Rule
	Include       string
	Captures      map[int]Capture
	BeginCaptures map[int]Capture
	EndCaptures   map[int]Capture
}

// Capture representa una captura en una regla
type Capture struct {
	Name string
}

// RepositoryRule es un alias para Rule, usado en el repositorio
type RepositoryRule = Rule
