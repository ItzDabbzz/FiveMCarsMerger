package validator

import "strings"

type Validator interface {
	IsValidDataFile(file string) bool
	IsValidStreamFile(file string) bool
	IsValidAudioFile(file string) bool
	IsValidAudioDataFile(file string) bool
}

type validator struct {
}

func New() Validator {
	return &validator{}
}

func (v *validator) IsValidDataFile(file string) bool {
	validDataExtensions := []string{".meta", ".dat"}
	return v.HasAnyFileExtension(file, validDataExtensions)
}

func (v *validator) IsValidStreamFile(file string) bool {
	validStreamExtensions := []string{".yft", ".ytd"}
	return v.HasAnyFileExtension(file, validStreamExtensions)
}

func (v *validator) IsValidAudioFile(file string) bool {
	validAudioExtensions := []string{".awc", ".awc2"}
	return v.HasAnyFileExtension(file, validAudioExtensions)
}

func (v *validator) IsValidAudioDataFile(file string) bool {
	validAudioExtensions := []string{".dat151", ".dat151.nametable", ".dat151.rel", ".dat10", ".dat10.nametable", ".dat10.rel", ".dat54", ".dat54.nametable", ".dat54.rel", ".dat"}
	return v.HasAnyFileExtension(file, validAudioExtensions)
}

func (v *validator) HasAnyFileExtension(file string, extensions []string) bool {
	for _, extension := range extensions {
		if strings.HasSuffix(file, extension) {
			return true
		}
	}
	return false
}
