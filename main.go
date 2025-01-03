package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/config"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/flags"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/merger"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var Flags flags.Flags

func main() {
	appFlags, err := config.LoadConfig()
	logFile := "merger.log"
	// Clear existing log file by recreating it
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fileWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(fileWriter)

	if appFlags == nil {
		// First time setup
		appFlags = &flags.Flags{}
		if err := initialSetup(appFlags); err != nil {
			log.Fatal(err)
		}
		if err := config.SaveConfig(appFlags); err != nil {
			log.Fatal(err)
		}
	}

	if absPath, err := filepath.Abs(appFlags.OutputPath); err == nil {
		appFlags.OutputPath = absPath
	}
	if absPath, err := filepath.Abs(appFlags.InputPath); err == nil {
		appFlags.InputPath = absPath
	}

	for {
		mainMenu := []string{"Start Merge Process", "Edit Settings", "Exit"}
		var selected string

		form := huh.NewSelect[string]().
			Title("FiveM Cars Merger").
			Options(huh.NewOptions(mainMenu...)...).
			Value(&selected)

		if err := form.Run(); err != nil {
			log.Error("Menu error:", err)
			continue
		}

		switch selected {
		case "Start Merge Process":
			ConfigureLogger()
			carsMerger := merger.New(*appFlags)
			if err := carsMerger.Merge(); err != nil {
				log.Error("Merge failed:", err)
				continue
			}
		case "Edit Settings":
			if err := editSettings(appFlags); err != nil {
				log.Fatal(err)
			}
			if err := config.SaveConfig(appFlags); err != nil {
				log.Fatal(err)
			}
		// case "Extract RPF File":
		// 	extractor := rpf.Extractor{
		// 		CachePath: filepath.Join("assets", "out"), // Add 'extracted' subfolder
		// 	}

		// 	if err := extractor.Extract(filepath.Join("assets", "dlc.rpf")); err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	log.Info("RPF extraction completed successfully")
		case "Exit":
			return
		}
	}
}

func initialSetup(flags *flags.Flags) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input Path").
				Description("Path to all cars").
				Placeholder("cars").
				Value(&flags.InputPath),
			huh.NewInput().
				Title("Output Path").
				Description("Output path for merged cars").
				Placeholder("merged-cars").
				Value(&flags.OutputPath),
			huh.NewConfirm().
				Title("Enable Verbose Logging").
				Value(&flags.Verbose),
			huh.NewConfirm().
				Title("Clean Output Directory").
				Value(&flags.Clean),
		).Title("Initial Configuration"),
	).Run()
}

func editSettings(flags *flags.Flags) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input Path").
				Description("Path to all cars").
				Value(&flags.InputPath),
			huh.NewInput().
				Title("Output Path").
				Description("Output path for merged cars").
				Value(&flags.OutputPath),
			huh.NewConfirm().
				Title("Enable Verbose Logging").
				Value(&flags.Verbose),
			huh.NewConfirm().
				Title("Clean Output Directory").
				Value(&flags.Clean),
		).Title("Edit Settings"),
	).Run()
}

func ConfigureLogger() {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ff7df9"))

	log.SetFormatter(log.TextFormatter)
	log.SetReportCaller(true)
	log.SetPrefix(style.Render("FiveMCarsMerger"))

	if Flags.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
