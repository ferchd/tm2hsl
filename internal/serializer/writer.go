package serializer

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"

	"github.com/ferchd/tm2hsl/internal/ir"
)

// BytecodeWriter - Writes bytecode with alignment
type BytecodeWriter struct {
	w      io.Writer
	offset uint32
	buffer *bytes.Buffer
	layout *BytecodeLayout
}

func NewBytecodeWriter(w io.Writer) *BytecodeWriter {
	return &BytecodeWriter{
		w:      w,
		offset: 0,
		buffer: &bytes.Buffer{},
		layout: &BytecodeLayout{},
	}
}

// WriteBytecode - Writes entire bytecode
func (w *BytecodeWriter) WriteBytecode(machine *ir.StateMachine) error {
	// 1. Preparar layout
	if err := w.prepareLayout(machine); err != nil {
		return err
	}

	// 2. Escribir header (offset 0)
	w.writeHeader()

	// 3. Escribir secciones alineadas
	sections := []struct {
		offsetPtr *uint32
		sizePtr   *uint32
		write     func() error
	}{
		{offsetPtr: &w.layout.Header.StateTable, sizePtr: &w.layout.Header.StateTableSize, write: w.writeStateTable},
		{offsetPtr: &w.layout.Header.TokenTable, sizePtr: &w.layout.Header.TokenTableSize, write: w.writeTokenTable},
		{offsetPtr: &w.layout.Header.ActionTable, sizePtr: &w.layout.Header.ActionTableSize, write: w.writeActionTable},
		{offsetPtr: &w.layout.Header.StringTable, sizePtr: &w.layout.Header.StringTableSize, write: w.writeStringTable},
		{offsetPtr: &w.layout.Header.RegexTable, sizePtr: &w.layout.Header.RegexTableSize, write: w.writeRegexTable},
	}

	for _, section := range sections {
		// Alinear a 8 bytes
		w.alignTo(8)

		*section.offsetPtr = w.offset
		if err := section.write(); err != nil {
			return err
		}
		*section.sizePtr = w.offset - *section.offsetPtr
	}

	// 4. Calcular y escribir checksum
	w.writeChecksum()

	// 5. Volcar buffer al writer real
	_, err := w.w.Write(w.buffer.Bytes())
	return err
}

// alignTo - Añade padding para alineación
func (w *BytecodeWriter) alignTo(alignment uint32) {
	padding := (alignment - (w.offset % alignment)) % alignment
	if padding > 0 {
		w.buffer.Write(make([]byte, padding))
		w.offset += padding
	}
}

// writeHeader - Escribe header con campos calculados
func (w *BytecodeWriter) writeHeader() {
	// Magic number
	copy(w.layout.Header.Magic[:], "HSLv1\x00\x00\x00")

	// Versión
	w.layout.Header.Version = 0x0001

	// Tamaño fijo del header
	w.layout.Header.HeaderSize = 64

	// Write header
	binary.Write(w.buffer, binary.LittleEndian, w.layout.Header)
	w.offset += 64
}

// writeStateTable - Escribe tabla de estados
func (w *BytecodeWriter) writeStateTable() error {
	// TODO: Implement state table writing
	return nil
}

// writeStringTable - Escribe pool de strings
func (w *BytecodeWriter) writeStringTable() error {
	// TODO: Implement string table writing
	return nil
}

// writeChecksum - Calcula y escribe checksum
func (w *BytecodeWriter) writeChecksum() {
	data := w.buffer.Bytes()

	// Calcular CRC32 desde después del header
	start := 64 // Saltar header completo
	checksum := crc32.ChecksumIEEE(data[start:])

	// Actualizar en el buffer
	binary.LittleEndian.PutUint32(data[56:60], checksum)
}

// Helper para preparar layout
func (w *BytecodeWriter) prepareLayout(machine *ir.StateMachine) error {
	// Convertir máquina a layout estructurado
	layout := &BytecodeLayout{}
	// TODO: Implement conversion logic from machine to layout
	w.layout = layout
	return nil
}

// Stub methods - TODO: Implement properly
func (w *BytecodeWriter) writeTokenTable() error {
	// TODO: Implement token table writing
	return nil
}

func (w *BytecodeWriter) writeActionTable() error {
	// TODO: Implement action table writing
	return nil
}

func (w *BytecodeWriter) writeRegexTable() error {
	// TODO: Implement regex table writing
	return nil
}
