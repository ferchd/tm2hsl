package optimizer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ferchd/tm2hsl/internal/ir"
)

// GraphAnalyzer - Analyzes the graph structure
type GraphAnalyzer struct {
	machine *ir.StateMachine
}

func NewGraphAnalyzer(machine *ir.StateMachine) *GraphAnalyzer {
	return &GraphAnalyzer{machine: machine}
}

// findReachableStates - Encuentra estados alcanzables desde el inicial
func findReachableStates(machine *ir.StateMachine) map[ir.StateID]bool {
	visited := make(map[ir.StateID]bool)
	queue := []ir.StateID{machine.Initial}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		state := machine.States[current]
		for _, trans := range state.Transitions {
			if !visited[trans.Target] {
				queue = append(queue, trans.Target)
			}
		}
	}

	return visited
}

// computeStateEquivalence - Calcula equivalencia de estados
func computeStateEquivalence(machine *ir.StateMachine) map[ir.StateID]int {
	// Inicializar con todos los estados en la misma clase
	class := make(map[ir.StateID]int)
	for id := range machine.States {
		class[id] = 0
	}

	// Algoritmo de minimización de DFA
	changed := true
	iteration := 0

	for changed {
		changed = false
		iteration++

		// Calcular nueva partición
		newClass := make(map[ir.StateID]int)
		signatures := make(map[string]ir.StateID)

		for id, state := range machine.States {
			// Crear firma del estado
			sig := stateSignature(state, class)

			if existingID, exists := signatures[sig]; exists {
				newClass[id] = newClass[existingID]
			} else {
				newClass[id] = len(signatures)
				signatures[sig] = id
			}
		}

		// Verificar si cambió
		if !mapsEqual(class, newClass) {
			class = newClass
			changed = true
		}
	}

	return class
}

func stateSignature(state *ir.State, class map[ir.StateID]int) string {
	var parts []string

	// Transiciones por clase de destino
	transMap := make(map[int][]string)
	for _, trans := range state.Transitions {
		destClass := class[trans.Target]
		predStr := trans.Predicate.String()
		transMap[destClass] = append(transMap[destClass], predStr)
	}

	// Ordenar para consistencia
	classes := make([]int, 0, len(transMap))
	for c := range transMap {
		classes = append(classes, c)
	}
	sort.Ints(classes)

	for _, c := range classes {
		preds := transMap[c]
		sort.Strings(preds)
		parts = append(parts, fmt.Sprintf("%d:[%s]", c, strings.Join(preds, ",")))
	}

	// Acciones de entrada/salida
	if len(state.OnEntry) > 0 {
		parts = append(parts, fmt.Sprintf("entry:%d", len(state.OnEntry)))
	}
	if len(state.OnExit) > 0 {
		parts = append(parts, fmt.Sprintf("exit:%d", len(state.OnExit)))
	}

	return strings.Join(parts, "|")
}

// mapsEqual compares two maps for equality
func mapsEqual(a, b map[ir.StateID]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}

// calculateCyclomaticComplexity - Complejidad ciclomática del autómata
func (a *GraphAnalyzer) calculateCyclomaticComplexity() int {
	// V(G) = E - N + 2P
	// donde E = aristas, N = nodos, P = componentes conexas
	edges := 0
	for _, state := range a.machine.States {
		edges += len(state.Transitions)
	}

	nodes := len(a.machine.States)
	components := a.countConnectedComponents()

	return edges - nodes + 2*components
}

func (a *GraphAnalyzer) countConnectedComponents() int {
	visited := make(map[ir.StateID]bool)
	components := 0

	for id := range a.machine.States {
		if !visited[id] {
			components++
			a.dfs(id, visited)
		}
	}

	return components
}

func (a *GraphAnalyzer) dfs(start ir.StateID, visited map[ir.StateID]bool) {
	stack := []ir.StateID{start}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[current] {
			continue
		}
		visited[current] = true

		for _, trans := range a.machine.States[current].Transitions {
			if !visited[trans.Target] {
				stack = append(stack, trans.Target)
			}
		}
	}
}
