package serializer

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/ferchd/tm2hsl/pkg/hsl"
)

type Serializer struct {
	byteOrder binary.ByteOrder
}

func NewSerializer() *Serializer {
	return &Serializer{
		byteOrder: binary.LittleEndian,
	}
}

func (s *Serializer) Serialize(bytecode *hsl.Bytecode, w io.Writer) error {
	// 1. Cabecera
	if err := binary.Write(w, s.byteOrder, &bytecode.Header); err != nil {
		return err
	}

	// 2. Escribir offsets dinámicos (se actualizarán después)
	offsets := struct {
		StringTable uint32
		RegexTable  uint32
		ScopeTable  uint32
		StateTable  uint32
		RuleTable   uint32
	}{}

	offsetPos, _ := w.(io.Seeker).Seek(0, io.SeekCurrent)
	if err := binary.Write(w, s.byteOrder, &offsets); err != nil {
		return err
	}

	// 3. Tabla de strings
	offsets.StringTable = uint32(offsetPos)
	if err := s.writeStringTable(w, &bytecode.StringTable); err != nil {
		return err
	}

	// 4. Tabla de regex
	offsets.RegexTable = uint32(offsetPos)
	if err := s.writeRegexTable(w, &bytecode.RegexTable); err != nil {
		return err
	}

	// 5. Tabla de scopes
	offsets.ScopeTable = uint32(offsetPos)
	if err := s.writeScopeTable(w, &bytecode.ScopeTable); err != nil {
		return err
	}

	// 6. Tabla de estados
	offsets.StateTable = uint32(offsetPos)
	if err := s.writeStateTable(w, &bytecode.StateTable); err != nil {
		return err
	}

	// 7. Tabla de reglas
	offsets.RuleTable = uint32(offsetPos)
	if err := s.writeRuleTable(w, &bytecode.RuleTable); err != nil {
		return err
	}

	// 8. Volver atrás y escribir offsets reales
	if seeker, ok := w.(io.Seeker); ok {
		seeker.Seek(offsetPos, io.SeekStart)
		binary.Write(w, s.byteOrder, &offsets)
	}

	// 9. Escribir tamaño total
	if seeker, ok := w.(io.Seeker); ok {
		endPos, _ := seeker.Seek(0, io.SeekCurrent)
		seeker.Seek(int64(binary.Size(bytecode.Header)-8), io.SeekStart)
		binary.Write(w, s.byteOrder, uint32(endPos))
	}

	return nil
}

func (s *Serializer) WriteToFile(bytecode *hsl.Bytecode, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.Serialize(bytecode, file)
}

func (s *Serializer) writeStringTable(w io.Writer, table *hsl.StringTable) error {
	if err := binary.Write(w, s.byteOrder, table.Count); err != nil {
		return err
	}

	// Escribir offsets
	for _, offset := range table.Offsets {
		if err := binary.Write(w, s.byteOrder, offset); err != nil {
			return err
		}
	}

	// Escribir datos de strings
	_, err := w.Write(table.Data)
	return err
}

func (s *Serializer) writeRegexTable(w io.Writer, table *hsl.RegexTable) error {
	if err := binary.Write(w, s.byteOrder, table.Count); err != nil {
		return err
	}

	for _, entry := range table.Entries {
		if err := binary.Write(w, s.byteOrder, entry.ID); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.PatternHash); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, uint32(len(entry.Bytecode))); err != nil {
			return err
		}
		if _, err := w.Write(entry.Bytecode); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.Flags); err != nil {
			return err
		}
	}

	return nil
}

func (s *Serializer) writeScopeTable(w io.Writer, table *hsl.ScopeTable) error {
	if err := binary.Write(w, s.byteOrder, table.Count); err != nil {
		return err
	}

	for _, entry := range table.Entries {
		if err := binary.Write(w, s.byteOrder, entry.ID); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.NameID); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.ParentID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Serializer) writeStateTable(w io.Writer, table *hsl.StateTable) error {
	if err := binary.Write(w, s.byteOrder, table.Count); err != nil {
		return err
	}

	for _, entry := range table.Entries {
		if err := binary.Write(w, s.byteOrder, entry.ID); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.RuleOffset); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.RuleCount); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.Flags); err != nil {
			return err
		}
	}

	return nil
}

func (s *Serializer) writeRuleTable(w io.Writer, table *hsl.RuleTable) error {
	if err := binary.Write(w, s.byteOrder, table.Count); err != nil {
		return err
	}

	for _, entry := range table.Entries {
		if err := binary.Write(w, s.byteOrder, entry.RegexID); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.Action); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.NextState); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.ScopeID); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.Priority); err != nil {
			return err
		}
		if err := binary.Write(w, s.byteOrder, entry.CaptureCount); err != nil {
			return err
		}

		// Escribir capturas
		for i := uint8(0); i < entry.CaptureCount; i++ {
			cap := entry.Captures[i]
			if err := binary.Write(w, s.byteOrder, cap.Group); err != nil {
				return err
			}
			if err := binary.Write(w, s.byteOrder, cap.ScopeID); err != nil {
				return err
			}
		}
	}

	return nil
}
