package copier

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/dft"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/flags"
	fileutils "github.com/ItzDabbzz/FiveMCarsMerger/pkg/utils/file"
	"github.com/charmbracelet/log"
)

type Copier interface {
	CopyStreamFilesToOutputDirectory(streamFiles []dft.StreamFile) error
	CopyDataFilesToOutputDirectory(dataFiles []dft.DataFile) error
	CopyAudioFilesToOutputDirectory(audioFiles []dft.AudioFile) error
}

type copier struct {
	Flags flags.Flags
}

func New(_flags flags.Flags) Copier {
	return &copier{Flags: _flags}
}

func (c *copier) CopyDataFilesToOutputDirectory(dataFiles []dft.DataFile) error {
	// First ensure the base output directory exists
	log.Info("Creating base output directory", "path", c.Flags.OutputPath)
	if err := os.MkdirAll(c.Flags.OutputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create data directory with explicit path
	baseDataPath := filepath.Join(c.Flags.OutputPath, "data")
	log.Info("Creating data directory", "path", baseDataPath)
	if err := os.MkdirAll(baseDataPath, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create all subdirectories with explicit paths and verification
	dataDirs := []string{"vehicles", "carcols", "carvariations", "handling", "vehiclelayouts", "contentunlocks"}
	for _, dir := range dataDirs {
		fullPath := filepath.Join(baseDataPath, dir)
		log.Debug("Creating subdirectory", "path", fullPath)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", dir, err)
		}
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("directory creation failed for: %s", fullPath)
		}
	}

	// Map to store vehicle names
	vehicleNames := make(map[string]string)

	// First pass - get vehicle names from vehicles.meta
	for _, dataFile := range dataFiles {
		if dataFile.Type == dft.VEHICLES {
			dirPath := filepath.Dir(dataFile.Path)
			content, err := os.ReadFile(dataFile.Path)
			if err != nil {
				log.Debug("Failed to read vehicles.meta", "path", dataFile.Path)
				continue
			}

			re := regexp.MustCompile(`<modelName.*?>(.*)</modelName>`)
			matches := re.FindStringSubmatch(string(content))
			if len(matches) > 1 {
				vehicleNames[dirPath] = strings.TrimSpace(matches[1])
				log.Debug("Found vehicle name", "name", matches[1], "dir", dirPath)
			}
		}
	}

	// Second pass - copy files with correct names
	for _, dataFile := range dataFiles {
		dirPath := filepath.Dir(dataFile.Path)
		vehicleName := vehicleNames[dirPath]
		if vehicleName == "" {
			log.Debug("No vehicle name found for directory", "dir", dirPath)
			continue
		}

		typeDir := strings.ToLower(dataFile.Type.String())
		destPath := filepath.Join(baseDataPath, typeDir, fmt.Sprintf("%s_%s.meta", typeDir, vehicleName))

		log.Debug("Copying file", "from", dataFile.Path, "to", destPath)
		if _, err := fileutils.CopyFile(dataFile.Path, destPath); err != nil {
			return fmt.Errorf("failed to copy file %s: %w", dataFile.Path, err)
		}
	}

	return nil
}

func (c *copier) CopyStreamFilesToOutputDirectory(streamFiles []dft.StreamFile) error {
	err := c.CreateDirectoryInOutput("stream")
	if err != nil {
		return err
	}

	streamPath := c.Flags.OutputPath + "/stream/"

	for _, streamFile := range streamFiles {
		log.Debug("Copying file", "name", streamFile.Name)
		_, err := fileutils.CopyFile(streamFile.Path, streamPath+streamFile.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *copier) CopyAudioFilesToOutputDirectory(audioFiles []dft.AudioFile) error {
	// Create audio directories
	if err := c.CreateDirectoryInOutput("audioconfig"); err != nil {
		return err
	}
	if err := c.CreateDirectoryInOutput("sfx"); err != nil {
		return err
	}

	for _, audio := range audioFiles {
		if audio.IsConfig {
			destPath := filepath.Join(c.Flags.OutputPath, "audioconfig", audio.Name)
			if _, err := fileutils.CopyFile(audio.Path, destPath); err != nil {
				return err
			}
		} else {
			// When creating the destination path:
			dlcPath := filepath.Join(c.Flags.OutputPath, "sfx", "dlc_"+audio.DLCFolder)

			if err := os.MkdirAll(dlcPath, 0755); err != nil {
				return err
			}
			destPath := filepath.Join(dlcPath, audio.Name)
			if _, err := fileutils.CopyFile(audio.Path, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *copier) CreateDirectoryInOutput(name string) error {
	return os.MkdirAll(c.Flags.OutputPath+"/"+name, 0755)
}
