package optimizer

import (
	"sort"

	"github.com/ferchd/tm2hsl/internal/ir"
)

type Optimizer struct {
	program *ir.Program
}

func NewOptimizer(program *ir.Program) *Optimizer {
	return &Optimizer{program: program}
}

func (o *Optimizer) Optimize() *ir.Program {
	// Aplicar optimizaciones en orden
	o.deduplicateRegexes()
	o.removeDeadStates()
	o.mergeEquivalentStates()
	o.sortRulesByPriority()
	o.compressStateTable()

	return o.program
}

func (o *Optimizer) deduplicateRegexes() {
	// Ya manejado por IR, pero podemos optimizar más
	// identificando regex equivalentes pero con diferentes strings
}

func (o *Optimizer) removeDeadStates() {
	// Identificar estados inalcanzables
	reachable := make(map[uint32]bool)
	o.markReachable(0, reachable) // Estado inicial

	// Reconstruir tabla de estados sin los inalcanzables
	newStateTable := []ir.StateEntry{}
	oldToNew := make(map[uint32]uint32)

	for _, state := range o.program.StateTable {
		if reachable[state.ID] || state.ID == 0 {
			newID := uint32(len(newStateTable))
			oldToNew[state.ID] = newID
			state.ID = newID
			newStateTable = append(newStateTable, state)
		}
	}

	// Actualizar referencias en reglas
	for i := range o.program.RuleTable {
		rule := &o.program.RuleTable[i]
		if rule.NextState >= 0 {
			if newID, exists := oldToNew[uint32(rule.NextState)]; exists {
				rule.NextState = int32(newID)
			} else {
				rule.NextState = -2 // stay (estado eliminado)
			}
		}
	}

	o.program.StateTable = newStateTable
}

func (o *Optimizer) markReachable(stateID uint32, visited map[uint32]bool) {
	if visited[stateID] {
		return
	}
	visited[stateID] = true

	state := o.program.StateTable[stateID]
	for i := uint32(0); i < uint32(state.RuleCount); i++ {
		rule := o.program.RuleTable[state.RuleOffset+uint32(i)]
		if rule.NextState >= 0 {
			o.markReachable(uint32(rule.NextState), visited)
		}
	}
}

func (o *Optimizer) mergeEquivalentStates() {
	// Identificar estados con reglas equivalentes
	equivalence := make(map[string][]uint32)

	for _, state := range o.program.StateTable {
		// Crear firma del estado basada en sus reglas
		signature := o.stateSignature(state)
		equivalence[signature] = append(equivalence[signature], state.ID)
	}

	// Fusionar estados equivalentes
	for _, group := range equivalence {
		if len(group) > 1 {
			o.mergeStateGroup(group)
		}
	}
}

func (o *Optimizer) stateSignature(state ir.StateEntry) string {
	// Crear una firma única para el estado basada en sus reglas
	var rules []string
	for i := uint32(0); i < uint32(state.RuleCount); i++ {
		rule := o.program.RuleTable[state.RuleOffset+uint32(i)]
		rules = append(rules, o.ruleSignature(rule))
	}

	sort.Strings(rules) // Orden independiente
	signature := ""
	for _, r := range rules {
		signature += r + "|"
	}
	return signature
}

func (o *Optimizer) ruleSignature(rule ir.RuleEntry) string {
	return string(rule.Action) + ":" +
		string(rune(rule.NextState)) + ":" +
		string(rune(rule.ScopeID)) + ":" +
		string(rune(rule.Priority))
}

func (o *Optimizer) mergeStateGroup(states []uint32) {
	if len(states) < 2 {
		return
	}

	// Mantener el primer estado, redirigir referencias a los otros
	keeper := states[0]

	for i := 1; i < len(states); i++ {
		o.redirectStateReferences(states[i], keeper)
	}
}

func (o *Optimizer) redirectStateReferences(oldID, newID uint32) {
	// Actualizar todas las reglas que apuntan al estado viejo
	for i := range o.program.RuleTable {
		rule := &o.program.RuleTable[i]
		if rule.NextState == int32(oldID) {
			rule.NextState = int32(newID)
		}
	}
}

func (o *Optimizer) sortRulesByPriority() {
	for i := range o.program.StateTable {
		state := &o.program.StateTable[i]

		// Extraer reglas del estado
		rules := o.program.RuleTable[state.RuleOffset : state.RuleOffset+uint32(state.RuleCount)]

		// Ordenar por prioridad (descendente)
		sort.Slice(rules, func(i, j int) bool {
			return rules[i].Priority > rules[j].Priority
		})

		// Copiar de vuelta
		copy(o.program.RuleTable[state.RuleOffset:], rules)
	}
}

func (o *Optimizer) compressStateTable() {
	// Compactar offsets de reglas después de optimizaciones
	currentOffset := uint32(0)

	for i := range o.program.StateTable {
		state := &o.program.StateTable[i]
		state.RuleOffset = currentOffset
		currentOffset += uint32(state.RuleCount)
	}
}
