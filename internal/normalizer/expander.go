package normalizer

// Expander - Expande referencias y includes en gramáticas
type Expander struct {
	// Configuración del expansor
}

// NewExpander - Crea un nuevo expansor
func NewExpander() *Expander {
	return &Expander{}
}

// Expand - Expande todas las referencias en la gramática
func (e *Expander) Expand(grammar interface{}) (interface{}, error) {
	// TODO: Implementar expansión de referencias
	// Por ahora, devolver sin cambios
	return grammar, nil
}
