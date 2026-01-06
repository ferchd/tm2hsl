package normalizer

// Expander - Expands references and includes in grammars
type Expander struct {
	// Expander configuration
}

// NewExpander - Creates a new expander
func NewExpander() *Expander {
	return &Expander{}
}

// Expand - Expands all references in the grammar
func (e *Expander) Expand(grammar interface{}) (interface{}, error) {
	// TODO: Implement reference expansion
	// For now, return unchanged
	return grammar, nil
}
