package rpf

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

const (
	HeaderSize        = 20
	TOCOffset         = 2048
	StringTableOffset = TOCOffset + 16 // String table follows TOC entries
)

func ParseRPF(filepath string) (*RPFHeader, error) {
	log.Debug("Starting RPF parsing", "file", filepath)

	file, err := os.Open(filepath)
	if err != nil {
		log.Error("Failed to open RPF file", "error", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	header := &RPFHeader{}
	if err := binary.Read(file, binary.LittleEndian, header); err != nil {
		log.Error("Failed to read RPF header", "error", err)
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Validate header values
	if header.TOCSize <= 0 || header.TOCSize > 1<<30 {
		log.Error("Invalid TOC size in header", "size", header.TOCSize)
		return nil, fmt.Errorf("invalid TOC size: %d", header.TOCSize)
	}

	if header.EntryCount <= 0 || header.EntryCount > 1<<20 {
		log.Error("Invalid entry count in header", "count", header.EntryCount)
		return nil, fmt.Errorf("invalid entry count: %d", header.EntryCount)
	}

	log.Debug("RPF header parsed successfully",
		"version", header.Version,
		"tocSize", header.TOCSize,
		"entryCount", header.EntryCount)

	return header, nil
}
