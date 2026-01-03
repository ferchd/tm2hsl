package normalizer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ferchd/tm2hsl/pkg/textmate"
)

type Normalizer struct {
	grammar      *textmate.Grammar
	stateMachine *StateMachine
	stateCounter int
	ruleCounter  int

	// Mapeos para seguimiento
	ruleToState  map[string]int  // include name -> state ID
	stateStack   []int           // Para detección de ciclos
	visitedRules map[string]bool // Para evitar recursión infinita
}

type StateMachine struct {
	InitialState int
	States       map[int]*State
	RegexCache   map[string]*regexp.Regexp
}

type State struct {
	ID       int
	Rules    []StateRule
	IsFinal  bool
	Metadata map[string]string
}

type StateRule struct {
	ID        int
	Pattern   string // Regex compilada (referencia)
	Action    RuleAction
	NextState int // -1 para pop, -2 para stay
	Priority  int
	Captures  map[int]string
}

type RuleAction uint8

const (
	ActionPush RuleAction = iota
	ActionPop
	ActionStay
	ActionScope
)

func NewNormalizer(grammar *textmate.Grammar) *Normalizer {
	return &Normalizer{
		grammar: grammar,
		stateMachine: &StateMachine{
			States:     make(map[int]*State),
			RegexCache: make(map[string]*regexp.Regexp),
		},
		ruleToState:  make(map[string]int),
		visitedRules: make(map[string]bool),
	}
}

func (n *Normalizer) Normalize() (*StateMachine, error) {
	// Crear estado inicial
	initialState := n.createState()
	n.stateMachine.InitialState = initialState.ID

	// Procesar patrones principales
	for _, rule := range n.grammar.Patterns {
		stateRule, err := n.processRule(rule, initialState.ID)
		if err != nil {
			return nil, fmt.Errorf("procesando patrón: %w", err)
		}
		initialState.Rules = append(initialState.Rules, stateRule)
	}

	// Procesar repositorio
	for name, rule := range n.grammar.Repository {
		n.stateStack = []int{}
		if err := n.processRepositoryRule(name, rule); err != nil {
			return nil, fmt.Errorf("repositorio '%s': %w", name, err)
		}
	}

	// Validar máquina de estados
	if err := n.validateStateMachine(); err != nil {
		return nil, fmt.Errorf("máquina de estados inválida: %w", err)
	}

	return n.stateMachine, nil
}

func (n *Normalizer) processRule(rule textmate.Rule, currentState int) (StateRule, error) {
	stateRule := StateRule{
		ID:        n.ruleCounter,
		Captures:  make(map[int]string),
		NextState: -2, // Por defecto: stay
	}
	n.ruleCounter++

	// Procesar match simple
	if rule.Match != "" {
		if err := n.compileRegex(rule.Match); err != nil {
			return stateRule, err
		}
		stateRule.Pattern = rule.Match

		if rule.Name != "" {
			stateRule.Action = ActionScope
			// Convertir nombre de scope a ID
		}
	}

	// Procesar begin/end (estados push/pop)
	if rule.Begin != "" && rule.End != "" {
		if err := n.compileRegex(rule.Begin); err != nil {
			return stateRule, err
		}
		if err := n.compileRegex(rule.End); err != nil {
			return stateRule, err
		}

		// Crear nuevo estado para contenido entre begin/end
		contentState := n.createState()
		stateRule.Pattern = rule.Begin
		stateRule.Action = ActionPush
		stateRule.NextState = contentState.ID

		// Procesar patrones internos
		for _, subRule := range rule.Patterns {
			contentRule, err := n.processRule(subRule, contentState.ID)
			if err != nil {
				return stateRule, err
			}
			contentState.Rules = append(contentState.Rules, contentRule)
		}

		// Crear regla de end
		endRule := StateRule{
			ID:        n.ruleCounter,
			Pattern:   rule.End,
			Action:    ActionPop,
			NextState: -1,
		}
		n.ruleCounter++

		if rule.Name != "" {
			// Asignar scope al end
		}

		contentState.Rules = append(contentState.Rules, endRule)
	}

	// Procesar includes
	if rule.Include != "" {
		return n.processInclude(rule.Include, currentState)
	}

	// Procesar capturas
	for idx, capture := range rule.Captures {
		if capture.Name != "" {
			stateRule.Captures[idx] = capture.Name
		}
	}

	return stateRule, nil
}

func (n *Normalizer) processRepositoryRule(name string, rule textmate.Rule) error {
	// Detectar ciclos en la pila de estados
	if n.visitedRules[name] {
		return fmt.Errorf("ciclo detectado en regla de repositorio: %s", name)
	}

	n.visitedRules[name] = true
	defer delete(n.visitedRules, name)

	// Crear estado para esta regla de repositorio
	ruleState := n.createState()
	n.ruleToState[name] = ruleState.ID

	// Procesar la regla
	stateRule, err := n.processRule(rule, ruleState.ID)
	if err != nil {
		return fmt.Errorf("procesando regla '%s': %w", name, err)
	}

	ruleState.Rules = append(ruleState.Rules, stateRule)

	return nil
}

func (n *Normalizer) processInclude(include string, currentState int) (StateRule, error) {
	// Detectar ciclos
	if n.visitedRules[include] {
		return StateRule{}, fmt.Errorf("ciclo detectado en include: %s", include)
	}

	n.visitedRules[include] = true
	defer delete(n.visitedRules, include)

	// Include a repositorio ($self, $base, o nombre)
	if strings.HasPrefix(include, "#") {
		ruleName := strings.TrimPrefix(include, "#")
		if rule, exists := n.grammar.Repository[ruleName]; exists {
			return n.processRule(rule, currentState)
		}
		return StateRule{}, fmt.Errorf("regla de repositorio no encontrada: %s", ruleName)
	}

	// Include a gramática externa
	return StateRule{}, fmt.Errorf("includes externos no soportados aún")
}

func (n *Normalizer) compileRegex(pattern string) error {
	if _, exists := n.stateMachine.RegexCache[pattern]; !exists {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("regex inválida '%s': %w", pattern, err)
		}
		n.stateMachine.RegexCache[pattern] = re
	}
	return nil
}

func (n *Normalizer) createState() *State {
	state := &State{
		ID:       n.stateCounter,
		Rules:    []StateRule{},
		Metadata: make(map[string]string),
	}
	n.stateMachine.States[n.stateCounter] = state
	n.stateCounter++
	return state
}

func (n *Normalizer) validateStateMachine() error {
	// Verificar estados inalcanzables
	reachable := make(map[int]bool)
	n.markReachable(n.stateMachine.InitialState, reachable)

	for id, state := range n.stateMachine.States {
		if !reachable[id] && len(state.Rules) > 0 {
			return fmt.Errorf("estado inalcanzable: %d", id)
		}
	}

	// Verificar ciclos infinitos
	for id := range n.stateMachine.States {
		if n.hasInfiniteLoop(id, make(map[int]bool)) {
			return fmt.Errorf("ciclo infinito detectado desde estado %d", id)
		}
	}

	return nil
}

func (n *Normalizer) markReachable(stateID int, visited map[int]bool) {
	if visited[stateID] {
		return
	}
	visited[stateID] = true

	state, exists := n.stateMachine.States[stateID]
	if !exists {
		return
	}

	for _, rule := range state.Rules {
		if rule.NextState >= 0 {
			n.markReachable(rule.NextState, visited)
		}
	}
}

func (n *Normalizer) hasInfiniteLoop(stateID int, path map[int]bool) bool {
	if path[stateID] {
		return true
	}

	path[stateID] = true
	defer delete(path, stateID)

	state := n.stateMachine.States[stateID]
	for _, rule := range state.Rules {
		if rule.NextState >= 0 {
			if n.hasInfiniteLoop(rule.NextState, path) {
				return true
			}
		}
	}

	return false
}
