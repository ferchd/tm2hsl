package parser

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
    "regexp"
    
    "github.com/ferchd/tm2hsl/pkg/textmate"
)

type Parser struct {
    grammar *textmate.Grammar
    errors  []error
}

func NewParser() *Parser {
    return &Parser{
        grammar: &textmate.Grammar{
            Repository: make(map[string]textmate.Rule),
        },
    }
}

func (p *Parser) ParseFile(path string) (*textmate.Grammar, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("no se pudo abrir gramática: %w", err)
    }
    defer file.Close()
    
    return p.Parse(file)
}

func (p *Parser) Parse(r io.Reader) (*textmate.Grammar, error) {
    var rawGrammar map[string]interface{}
    
    decoder := json.NewDecoder(r)
    if err := decoder.Decode(&rawGrammar); err != nil {
        return nil, fmt.Errorf("JSON inválido: %w", err)
    }
    
    // Extraer metadatos básicos
    if scopeName, ok := rawGrammar["scopeName"].(string); ok {
        p.grammar.ScopeName = scopeName
    }
    
    // Parsear repositorio
    if repo, ok := rawGrammar["repository"].(map[string]interface{}); ok {
        for name, ruleData := range repo {
            rule, err := p.parseRule(ruleData)
            if err != nil {
                p.errors = append(p.errors, 
                    fmt.Errorf("repositorio '%s': %w", name, err))
                continue
            }
            p.grammar.Repository[name] = rule
        }
    }
    
    // Parsear patrones principales
    if patterns, ok := rawGrammar["patterns"].([]interface{}); ok {
        for i, patternData := range patterns {
            rule, err := p.parseRule(patternData)
            if err != nil {
                p.errors = append(p.errors, 
                    fmt.Errorf("patrón %d: %w", i, err))
                continue
            }
            p.grammar.Patterns = append(p.grammar.Patterns, rule)
        }
    }
    
    if len(p.errors) > 0 {
        return p.grammar, fmt.Errorf("errores de parseo: %v", p.errors)
    }
    
    return p.grammar, nil
}

func (p *Parser) parseRule(data interface{}) (textmate.Rule, error) {
    rule := textmate.Rule{}
    
    m, ok := data.(map[string]interface{})
    if !ok {
        return rule, fmt.Errorf("regla no es objeto")
    }
    
    // Extraer campos básicos
    if name, ok := m["name"].(string); ok {
        rule.Name = name
    }
    if match, ok := m["match"].(string); ok {
        rule.Match = match
        if err := p.validateRegex(match); err != nil {
            return rule, fmt.Errorf("regex inválida: %w", err)
        }
    }
    if begin, ok := m["begin"].(string); ok {
        rule.Begin = begin
        if err := p.validateRegex(begin); err != nil {
            return rule, fmt.Errorf("regex begin inválida: %w", err)
        }
    }
    if end, ok := m["end"].(string); ok {
        rule.End = end
        if err := p.validateRegex(end); err != nil {
            return rule, fmt.Errorf("regex end inválida: %w", err)
        }
    }
    
    // Parsear sub-patrones
    if patterns, ok := m["patterns"].([]interface{}); ok {
        for _, patternData := range patterns {
            subRule, err := p.parseRule(patternData)
            if err != nil {
                return rule, err
            }
            rule.Patterns = append(rule.Patterns, subRule)
        }
    }
    
    // Parsear capturas
    if captures, ok := m["captures"].(map[string]interface{}); ok {
        rule.Captures = make(map[int]textmate.Capture)
        for key, capData := range captures {
            var index int
            if _, err := fmt.Sscanf(key, "%d", &index); err != nil {
                continue
            }
            
            capMap, ok := capData.(map[string]interface{})
            if !ok {
                continue
            }
            
            capture := textmate.Capture{}
            if name, ok := capMap["name"].(string); ok {
                capture.Name = name
            }
            rule.Captures[index] = capture
        }
    }
    
    return rule, nil
}

func (p *Parser) validateRegex(pattern string) error {
    if pattern == "" {
        return nil
    }
    
    // Validación básica de regex
    _, err := regexp.Compile(pattern)
    return err
}