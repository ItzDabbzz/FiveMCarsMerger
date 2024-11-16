package typeidentifier

import (
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/dft"
	xmlutils "github.com/iLLeniumStudios/FiveMCarsMerger/pkg/utils/xml"
	"github.com/charmbracelet/log"
	"io/ioutil"
	"os"
)

type TypeIdentifier interface {
	IdentifyDataFileType(path string) (dft.DataFileType, error)
}

type typeIdentifier struct {
}

func New() TypeIdentifier {
	return &typeIdentifier{}
}

func (ti *typeIdentifier) IdentifyDataFileType(path string) (dft.DataFileType, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return dft.INVALID, err
	}
	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return dft.INVALID, err
	}

	defer xmlFile.Close()

	startTag, err := xmlutils.GetStartTag(byteValue)
	if err != nil {
		return dft.INVALID, err
	}

	dataFileType := dft.INVALID

	switch startTag {
	case "CVehicleModelInfo__InitDataList":
		dataFileType = dft.VEHICLES
		break
	case "CAmbientModelSets":
		dataFileType = dft.VEHICLEMODELSETS
		break
	case "CVehicleMetadataMgr":
		dataFileType = dft.VEHICLELAYOUTS
		break
	case "SContentUnlocks":
		dataFileType = dft.CONTENTUNLOCKS
		break
	case "CHandlingDataMgr":
		dataFileType = dft.HANDLING
		break
	case "CVehicleModelInfoVariation":
		dataFileType = dft.CARVARIATIONS
		break
	case "CVehicleModelInfoVarGlobal":
		dataFileType = dft.CARCOLS
		break
	case "CWeaponInfoBlob":
		dataFileType = dft.WEAPONSFILE
	case "":
		log.Debug("Invalid XML file", "file", path)
		break
	default:
		log.Debug("Unknown tag detected", "file", path, "tag", startTag)
		break
	}

	return dataFileType, nil
}
