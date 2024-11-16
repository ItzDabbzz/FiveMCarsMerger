package rpf

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

type Extractor struct {
	CachePath string
	Header    *RPFHeader
	Entries   []interface{} // Can be DirectoryEntry or FileEntry
	Names     map[int32]string
}

var resourceTypeExtensions = map[byte]string{
	0x01: ".ytd",    // Texture Dictionary
	0x02: ".yft",    // Fragment
	0x0B: ".ydr",    // Drawable
	0x0C: ".ydd",    // Drawable Dictionary
	0x0D: ".ybn",    // Bounds
	0x0F: ".ycd",    // Clips Dictionary
	0x12: ".ymf",    // Map Types
	0x16: ".ytyp",   // Types
	0x17: ".ynd",    // Nav Mesh
	0x1B: ".ydd",    // Expression Dictionary
	0x1C: ".yet",    // Expression
	0x20: ".meta",   // Meta files
	0x22: ".gxt2",   // GXT2 Text files
	0x23: ".dat151", // Audio data
	0x24: ".dat10",  // Audio data
	0x25: ".dat54",  // Audio data
	0x26: ".awc",    // Audio wave container
	0x27: ".awc2",   // Audio wave container v2
}

// Add additional handling for related files
func getRelatedExtensions(baseExt string) []string {
	switch baseExt {
	case ".dat151", ".dat10", ".dat54":
		return []string{
			baseExt + ".nametable",
			baseExt + ".rel",
		}
	default:
		return nil
	}
}

func (e *Extractor) Extract(rpfPath string) error {
	log.Debug("Starting RPF extraction", "file", rpfPath)

	file, err := os.Open(rpfPath)
	if err != nil {
		log.Error("Failed to open RPF file", "error", err)
		return err
	}
	defer file.Close()

	if err := os.MkdirAll(e.CachePath, 0755); err != nil {
		log.Error("Failed to create output directory", "path", e.CachePath)
		return err
	}

	// Parse header first
	e.Header, err = ParseRPF(rpfPath)
	if err != nil {
		return err
	}

	// Read TOC entries
	tocSize := e.Header.EntryCount * 16
	if tocSize <= 0 || tocSize > 1<<30 { // Max 1GB for safety
		log.Error("Invalid TOC size calculated", "size", tocSize)
		return fmt.Errorf("invalid TOC size calculated: %d", tocSize)
	}

	toc := make([]byte, tocSize)
	if _, err := file.Seek(TOCOffset, 0); err != nil {
		return err
	}
	if _, err := io.ReadFull(file, toc); err != nil {
		return err
	}

	// Read string table
	stringTableSize := e.Header.TOCSize - tocSize
	stringTable := make([]byte, stringTableSize)
	if _, err := io.ReadFull(file, stringTable); err != nil {
		return err
	}

	// Parse string table
	e.Names = make(map[int32]string)
	currentOffset := int32(0)
	currentString := ""
	for i := 0; i < len(stringTable); i++ {
		if stringTable[i] == 0 {
			if currentString != "" {
				e.Names[currentOffset] = currentString
				log.Debug("Found filename", "offset", currentOffset, "name", currentString)
			}
			currentOffset = int32(i + 1)
			currentString = ""
		} else {
			currentString += string(stringTable[i])
		}
	}

	// Parse entries
	for i := 0; i < int(e.Header.EntryCount); i++ {
		offset := i * 16
		nameOffset := binary.LittleEndian.Uint32(toc[offset:])

		if (nameOffset & 0x80000000) != 0 {
			dir := &DirectoryEntry{}
			binary.Read(bytes.NewReader(toc[offset:offset+16]), binary.LittleEndian, dir)
			e.Entries = append(e.Entries, dir)
		} else {
			file := &FileEntry{}
			binary.Read(bytes.NewReader(toc[offset:offset+16]), binary.LittleEndian, file)
			e.Entries = append(e.Entries, file)
		}
	}

	// Extract files
	for _, entry := range e.Entries {
		if fileEntry, ok := entry.(*FileEntry); ok {
			if err := e.extractFile(file, fileEntry); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Extractor) extractFile(rpf *os.File, entry *FileEntry) error {
	// Get appropriate extension for the file
	extension := resourceTypeExtensions[entry.ResourceType]
	if extension == "" {
		// Skip unknown file types
		return nil
	}

	// Calculate absolute file offset
	fileInfo, err := rpf.Stat()
	if err != nil {
		log.Error("Failed to get RPF file info", "error", err)
		return err
	}

	fileOffset := int64(entry.Offset[0])<<16 | int64(entry.Offset[1])<<8 | int64(entry.Offset[2])
	expectedSize := int64(entry.Size)

	// Validate that claimed size doesn't exceed file bounds
	if fileOffset+expectedSize > fileInfo.Size() {
		log.Debug("Adjusting file size to match actual bounds",
			"originalSize", entry.Size,
			"adjustedSize", fileInfo.Size()-fileOffset)
		entry.Size = uint32(fileInfo.Size() - fileOffset)
	}

	data := make([]byte, entry.Size)
	if _, err := rpf.Seek(fileOffset, 0); err != nil {
		log.Error("Failed to seek to file offset", "error", err)
		return err
	}

	if _, err := io.ReadFull(rpf, data); err != nil {
		log.Error("Failed to read file data", "error", err)
		return err
	}

	if e.IsCompressed(entry.Flags) {
		log.Debug("Decompressing file data")
		// Check if data actually has zlib header before attempting decompression
		if len(data) > 2 && (data[0] == 0x78 || data[0] == 0x1f) {
			decompressedData, err := e.decompressData(data)
			if err != nil {
				log.Error("Failed to decompress data", "error", err)
				return err
			}
			data = decompressedData
		} else {
			log.Debug("File marked as compressed but no valid compression header found, using raw data")
		}
	}

	// Preserve directory structure
	outPath := filepath.Join(e.CachePath, filepath.Dir(e.Names[entry.NameOffset]), filepath.Base(e.Names[entry.NameOffset])+extension)
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}

	// Handle related files (like .nametable and .rel)
	if relatedExts := getRelatedExtensions(extension); len(relatedExts) > 0 {
		for _, relExt := range relatedExts {
			relPath := filepath.Join(filepath.Dir(outPath), filepath.Base(outPath)+relExt)
			if err := os.WriteFile(relPath, data, 0644); err != nil {
				return err
			}
		}
	}

	// Extract nested RPFs
	if extension == ".rpf" {
		nestedExtractor := &Extractor{
			CachePath: filepath.Join(filepath.Dir(outPath), "extracted_"+filepath.Base(outPath)),
		}
		if err := os.WriteFile(outPath, data, 0644); err != nil {
			return err
		}
		return nestedExtractor.Extract(outPath)
	}

	return os.WriteFile(outPath, data, 0644)
}

func (e *Extractor) decompressData(compressedData []byte) ([]byte, error) {
	reader := bytes.NewReader(compressedData)
	zlibReader, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer zlibReader.Close()

	return io.ReadAll(zlibReader)
}

func (e *Extractor) IsCompressed(flags uint32) bool {
	return (flags & 1) != 0
}
