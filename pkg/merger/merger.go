package merger

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/carfinder"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/copier"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/dft"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/flags"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/manifestgen"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/typeidentifier"
	sliceutils "github.com/ItzDabbzz/FiveMCarsMerger/pkg/utils/slice"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/validator"
	"github.com/charmbracelet/log"
)

type Merger interface {
	Merge() error
}

type merger struct {
	Flags          flags.Flags
	Generator      manifestgen.Generator
	Validator      validator.Validator
	TypeIdentifier typeidentifier.TypeIdentifier
	CarFinder      carfinder.CarFinder
	Copier         copier.Copier
}

func New(_flags flags.Flags) Merger {
	return &merger{
		Flags:          _flags,
		Generator:      manifestgen.New(_flags),
		Validator:      validator.New(),
		TypeIdentifier: typeidentifier.New(),
		CarFinder:      carfinder.New(_flags),
		Copier:         copier.New(_flags),
	}
}

func (m *merger) Merge() error {
	var streamFiles []dft.StreamFile
	var dataFiles []dft.DataFile
	var audioFiles []dft.AudioFile // New slice for audio files

	log.Info("Creating Output Directory...")
	outputPath, err := filepath.Abs(m.Flags.OutputPath)
	if err != nil {
		return err
	}
	m.Flags.OutputPath = outputPath
	if err := m.CreateOutputDirectory(); err != nil {
		return err
	}

	log.Info("Identifying cars", "path", m.Flags.InputPath)

	err = filepath.Walk(m.Flags.InputPath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return err
		}
		if m.Validator.IsValidAudioFile(f.Name()) || m.Validator.IsValidAudioDataFile(f.Name()) {
			isConfig := m.Validator.IsValidAudioDataFile(f.Name())
			dlcFolder := filepath.Base(filepath.Dir(path))
			if strings.HasPrefix(dlcFolder, "dlc_") {
				dlcFolder = strings.TrimPrefix(dlcFolder, "dlc_")
			}

			audioFiles = append(audioFiles, dft.AudioFile{
				Path:      path,
				Name:      f.Name(),
				IsConfig:  isConfig,
				DLCFolder: dlcFolder,
			})
			return err
		}
		if m.Validator.IsValidStreamFile(f.Name()) {
			streamFiles = append(streamFiles, dft.StreamFile{
				Path: path,
				Name: f.Name(),
			})
			return err
		}
		if m.Validator.IsValidDataFile(f.Name()) {
			_type, err := m.TypeIdentifier.IdentifyDataFileType(path)
			if err != nil {
				return err
			}
			dataFile := dft.DataFile{
				Path: path,
				Name: f.Name(),
				Type: _type,
			}

			if dataFile.Type != dft.INVALID {
				dataFiles = append(dataFiles, dataFile)
			}
		}
		return err
	})
	if err != nil {
		return err
	}

	if len(dataFiles) == 0 || len(streamFiles) == 0 {
		log.Error("Cannot find any cars in the specified folder")

		if err := m.Cleanup(); err != nil {
			return err
		}
		return nil
	}

	// Add audio file copying after stream/data files
	if len(audioFiles) > 0 {
		log.Info("Copying Audio files...")
		err = m.Copier.CopyAudioFilesToOutputDirectory(audioFiles)
		if err != nil {
			return err
		}
	}

	log.Info("Copying Stream files...")
	err = m.Copier.CopyStreamFilesToOutputDirectory(streamFiles)
	if err != nil {
		return err
	}

	log.Info("Copying Data files...")
	err = m.Copier.CopyDataFilesToOutputDirectory(dataFiles)
	if err != nil {
		return err
	}

	log.Info("Generating fxmanifest.lua")
	err = m.Generator.Generate()
	if err != nil {
		return err
	}

	log.Info("Parsing Data Files For Cars...")
	dataFileCars, err := m.CarFinder.FindDataFileCars()
	if err != nil {
		return err
	}
	log.Info("Parsing Stream Files For Cars...")
	streamFileCars, err := m.CarFinder.FindStreamFileCars()
	if err != nil {
		return err
	}

	validCars := sliceutils.RemoveDuplicates(m.CarFinder.FindValidCars(dataFileCars, streamFileCars))

	log.Info("Valid cars in the car pack", "cars", validCars)
	log.Info("Success! Resource ready", "output_folder", m.Flags.OutputPath)

	return nil
}

func (m *merger) CreateOutputDirectory() error {
	if m.Flags.Clean {
		if err := m.Cleanup(); err != nil {
			return err
		}
	}
	err := os.Mkdir(m.Flags.OutputPath, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (m *merger) Cleanup() error {
	if _, err := os.Stat(m.Flags.OutputPath); !os.IsNotExist(err) {
		err = os.RemoveAll(m.Flags.OutputPath)
		if err != nil {
			return err
		}
	}
	return nil
}
