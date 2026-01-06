package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type LanguageConfig struct {
	Name       string            `toml:"name"`
	Version    string            `toml:"version"`
	Scope      string            `toml:"scope"`
	Grammar    string            `toml:"grammar"`
	Snippets   []string          `toml:"snippets,omitempty"`
	Repository map[string]string `toml:"repository,omitempty"`
	Metadata   map[string]string `toml:"metadata,omitempty"`

	// Campos calculados
	baseDir     string
	grammarPath string
}

func LoadConfig(configPath string) (*LanguageConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer config: %w", err)
	}

	var config LanguageConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("config inválida: %w", err)
	}

	config.baseDir = filepath.Dir(configPath)

	// Resolver rutas relativas
	if config.Grammar != "" {
		if !filepath.IsAbs(config.Grammar) {
			config.grammarPath = filepath.Join(config.baseDir, config.Grammar)
		} else {
			config.grammarPath = config.Grammar
		}

		if _, err := os.Stat(config.grammarPath); err != nil {
			return nil, fmt.Errorf("gramática no encontrada: %s", config.grammarPath)
		}
	}

	// Validaciones básicas
	if config.Name == "" {
		return nil, fmt.Errorf("nombre del lenguaje requerido")
	}
	if config.Scope == "" {
		return nil, fmt.Errorf("scope requerido")
	}

	return &config, nil
}

func (c *LanguageConfig) GrammarPath() string {
	return c.grammarPath
}

func (c *LanguageConfig) BaseDir() string {
	return c.baseDir
}
