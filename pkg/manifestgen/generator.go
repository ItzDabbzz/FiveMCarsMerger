package manifestgen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/dft"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/flags"
	"github.com/charmbracelet/log"
)

type Manifest struct {
	HasCarcols          bool
	HasCarvariations    bool
	HasContentUnlocks   bool
	HasHandling         bool
	HasVehicleLayouts   bool
	HasVehicleModelsets bool
	HasVehicles         bool
	HasWeaponsFile      bool
	HasAudio            bool
	AudioConfigs        []string // Store unique audio config names
	AudioWavePacks      []string // Store unique wavepack folder names
}

type Generator interface {
	Generate() error
}

type generator struct {
	Flags flags.Flags
}

func New(_flags flags.Flags) Generator {
	return &generator{Flags: _flags}
}

func (g *generator) Generate() error {
	log.Debug("Starting manifest generation")

	tmpl, err := template.New("manifestTemplate").Parse(manifestTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	manifest := Manifest{
		AudioConfigs:   make([]string, 0),
		AudioWavePacks: make([]string, 0),
	}

	// Check if data directory exists
	dataPath := filepath.Join(g.Flags.OutputPath, "data")
	folders, err := ioutil.ReadDir(dataPath)
	if err != nil {
		return fmt.Errorf("failed to read data directory: %w", err)
	}

	log.Debug("Processing data folders", "count", len(folders))

	// Process data folders
	for _, folder := range folders {
		folderName := strings.ToLower(folder.Name())
		switch folderName {
		case strings.ToLower(dft.CARCOLS.String()):
			manifest.HasCarcols = true
			log.Debug("Found carcols folder", "folder", folderName)
		case strings.ToLower(dft.CARVARIATIONS.String()):
			manifest.HasCarvariations = true
			log.Debug("Found carvariations folder", "folder", folderName)
		case strings.ToLower(dft.CONTENTUNLOCKS.String()):
			manifest.HasContentUnlocks = true
			log.Debug("Found contentunlocks folder", "folder", folderName)
		case strings.ToLower(dft.HANDLING.String()):
			manifest.HasHandling = true
			log.Debug("Found handling folder", "folder", folderName)
		case strings.ToLower(dft.VEHICLELAYOUTS.String()):
			manifest.HasVehicleLayouts = true
			log.Debug("Found vehiclelayouts folder", "folder", folderName)
		case strings.ToLower(dft.VEHICLEMODELSETS.String()):
			manifest.HasVehicleModelsets = true
			log.Debug("Found vehiclemodelsets folder", "folder", folderName)
		case strings.ToLower(dft.VEHICLES.String()):
			manifest.HasVehicles = true
			log.Debug("Found vehicles folder", "folder", folderName)
		case strings.ToLower(dft.WEAPONSFILE.String()):
			manifest.HasWeaponsFile = true
			log.Debug("Found weaponsfile folder", "folder", folderName)
		default:
			log.Debug("Skipping unknown folder", "folder", folderName)
		}
	}

	// Process audio files if they exist
	// Process audio config files
	audioConfigPath := filepath.Join(g.Flags.OutputPath, "audioconfig")
	if _, err := os.Stat(audioConfigPath); err == nil {
		manifest.HasAudio = true
		audioFiles, _ := ioutil.ReadDir(audioConfigPath)

		uniqueConfigs := make(map[string]struct{})

		for _, file := range audioFiles {
			fileName := file.Name()
			// Example: hondaf20c_game.dat151.rel -> hondaf20c
			if strings.Contains(fileName, "_game") || strings.Contains(fileName, "_sounds") {
				parts := strings.Split(fileName, "_")
				if len(parts) > 0 {
					engineName := parts[0]
					uniqueConfigs[engineName] = struct{}{}
					log.Debug("Found audio config", "engine", engineName, "file", fileName)
				}
			}
		}

		// Convert to sorted slice for consistent output
		for config := range uniqueConfigs {
			manifest.AudioConfigs = append(manifest.AudioConfigs, config)
		}
		sort.Strings(manifest.AudioConfigs)
	}

	// Create manifest file
	manifestPath := filepath.Join(g.Flags.OutputPath, "fxmanifest.lua")
	log.Debug("Creating manifest file", "path", manifestPath)

	fxManifest, err := os.Create(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to create manifest file: %w", err)
	}
	defer fxManifest.Close()

	// Execute template
	if err := tmpl.Execute(fxManifest, manifest); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	log.Info("Successfully generated fxmanifest.lua")
	return nil

}
