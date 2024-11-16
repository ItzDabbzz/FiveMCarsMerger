package main

import (
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/flags"
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/merger"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/lipgloss"
	flag "github.com/spf13/pflag"
)

var Flags flags.Flags

func init() {
	ParseFlags()
	ConfigureLogger()
}

func ParseFlags() {
	flag.BoolVar(&Flags.Verbose, "verbose", false, "Enable verbose logging")
	flag.StringVar(&Flags.InputPath, "input-path", ".", "Path to all cars")
	flag.StringVar(&Flags.OutputPath, "output-path", "out", "Output path")
	flag.BoolVar(&Flags.Clean, "clean", false, "Clear output directory")
	flag.Parse()
}

func ConfigureLogger() {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7EBC89"))

	log.SetFormatter(log.TextFormatter)
	log.SetReportCaller(true)
	log.SetPrefix(style.Render("FiveMCarsMerger"))

	if Flags.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// log.SetFormatter(&log.TextFormatter{
	// 	ForceColors: true,
	// })

	// if Flags.Verbose {
	// 	log.SetLevel(log.DebugLevel)
	// } else {
	// 	log.SetLevel(log.InfoLevel)
	// }
}

func main() {
	carsMerger := merger.New(Flags)
	err := carsMerger.Merge()
	if err != nil {
		log.Fatal(err)
	}
}
