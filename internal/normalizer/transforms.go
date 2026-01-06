package normalizer

// SemanticTransformer - Aplica transformaciones semánticas a la IR
type SemanticTransformer struct {
    machine *ir.StateMachine
}

func NewSemanticTransformer(machine *ir.StateMachine) *SemanticTransformer {
    return &SemanticTransformer{machine: machine}
}

// ApplyAllTransforms - Aplica todas las transformaciones
func (t *SemanticTransformer) ApplyAllTransforms() error {
    transforms := []func() error{
        t.normalizeCaptureNames,
        t.mergeEquivalentScopes,
        t.removeEmptyStates,
        t.optimizeTransitionOrder,
        t.inlineSimpleStates,
    }
    
    for _, transform := range transforms {
        if err := transform(); err != nil {
            return err
        }
    }
    
    return nil
}

// normalizeCaptureNames - Normaliza nombres de capturas
func (t *SemanticTransformer) normalizeCaptureNames() error {
    // Asegurar que los nombres de scope siguen convenciones HSL
    for _, state := range t.machine.States {
        for _, trans := range state.Transitions {
            for i, action := range trans.Actions {
                if push, ok := action.(*ir.PushScopeAction); ok {
                    push.Scope = t.normalizeScopeName(push.Scope)
                }
            }
        }
    }
    return nil
}

func (t *SemanticTransformer) normalizeScopeName(name string) string {
    // Convertir nombres TextMate a formato HSL
    name = strings.TrimSpace(name)
    
    // Remover prefijos comunes
    name = strings.TrimPrefix(name, "source.")
    name = strings.TrimPrefix(name, "text.")
    
    // Reemplazar caracteres no válidos
    name = strings.ReplaceAll(name, ".", "_")
    name = strings.ReplaceAll(name, "-", "_")
    
    return name
}

// mergeEquivalentScopes - Fusiona scopes equivalentes
func (t *SemanticTransformer) mergeEquivalentScopes() error {
    // Identificar scopes que son equivalentes semánticamente
    scopeMap := make(map[string]string)
    
    for _, state := range t.machine.States {
        for _, trans := range state.Transitions {
            for _, action := range trans.Actions {
                if push, ok := action.(*ir.PushScopeAction); ok {
                    // Buscar scope equivalente (mismo nombre normalizado)
                    normalized := t.normalizeScopeName(push.Scope)
                    if existing, found := scopeMap[push.Scope]; found {
                        push.Scope = existing
                    } else {
                        scopeMap[push.Scope] = normalized
                        push.Scope = normalized
                    }
                }
            }
        }
    }
    
    return nil
}

// removeEmptyStates - Elimina estados sin transiciones útiles
func (t *SemanticTransformer) removeEmptyStates() error {
    toRemove := []ir.StateID{}
    
    for id, state := range t.machine.States {
        // Estado vacío (sin transiciones y no final)
        if len(state.Transitions) == 0 && !state.IsFinal {
            toRemove = append(toRemove, id)
        }
    }
    
    for _, id := range toRemove {
        delete(t.machine.States, id)
    }
    
    // Reenviar transiciones a estados eliminados
    t.redirectTransitions(toRemove)
    
    return nil
}

// optimizeTransitionOrder - Ordena transiciones por prioridad
func (t *SemanticTransformer) optimizeTransitionOrder() error {
    for _, state := range t.machine.States {
        sort.Slice(state.Transitions, func(i, j int) bool {
            return state.Transitions[i].Priority > state.Transitions[j].Priority
        })
    }
    return nil
}

// inlineSimpleStates - Inlinea estados con solo una transición
func (t *SemanticTransformer) inlineSimpleStates() error {
    for id, state := range t.machine.States {
        if len(state.Transitions) == 1 && len(state.OnEntry) == 0 && len(state.OnExit) == 0 {
            t.inlineState(id, state)
        }
    }
    return nil
}

func (t *SemanticTransformer) inlineState(id ir.StateID, state *ir.State) {
    // Conectar transiciones entrantes directamente a la salida
    for _, s := range t.machine.States {
        for i, trans := range s.Transitions {
            if trans.Target == id {
                // Reemplazar destino y combinar acciones
                trans.Target = state.Transitions[0].Target
                trans.Actions = append(trans.Actions, state.Transitions[0].Actions...)
                s.Transitions[i] = trans
            }
        }
    }
    
    delete(t.machine.States, id)
}