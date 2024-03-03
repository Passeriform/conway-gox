package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/passeriform/conway-gox/internal/cell_map"
	"github.com/passeriform/conway-gox/internal/utility"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type PrimitiveInfo struct {
	Type     string
	Name     string
	SeedName string
	SeedPath string
}

func beautifySeedName(input string) string {
	caseChunks := make([]string, 0)
	strChunks := strings.FieldsFunc(input, func(r rune) bool {
		return r == ':' ||
			r == '-' ||
			r == '_' ||
			r == ' ' ||
			r == ';' ||
			r == '|' ||
			r == ','
	})

	for _, chunk := range strChunks {
		caseChunks = append(caseChunks, cases.Title(language.English).String(chunk))
	}

	return strings.Join(caseChunks, " ")
}

func LoadFromFile(fp string, padding int) (cell_map.Map, error) {
	dat, err := os.ReadFile(fp)
	if err != nil {
		return cell_map.Map{}, fmt.Errorf("unable to load map from file %v: %v", fp, err)
	}
	var messageObject [][2]int
	json.Unmarshal(dat, &messageObject)

	return cell_map.DecodeJson(messageObject, padding), nil
}

func SaveToFile(cm cell_map.Map, fp string, padding int) error {
	wb, err := json.Marshal(cm.EncodeJson(0))

	if err != nil {
		return fmt.Errorf("unable to marshal map: %v", err)
	}

	if _, err := os.Stat(fp); err == nil {
		fmt.Printf("Existing save file found at %v. Will be overwritten.", fp)
	}

	os.WriteFile(fp, wb, 0644)

	return nil
}

func LoadFromPrimitive(primitive string, padding int) (cell_map.Map, error) {
	pGlob := fmt.Sprintf("%v/*/%v.json", filepath.Join(os.Getenv("GOPATH"), "assets", "primitives"), primitive)
	matches, err := filepath.Glob(pGlob)
	if err != nil {
		return cell_map.Map{}, fmt.Errorf("could not prepare glob to search for primitives: %v", err)
	}

	if len(matches) == 0 {
		return cell_map.Map{}, fmt.Errorf("could not find the requested primitive: %v", primitive)
	}

	return LoadFromFile(matches[0], padding)
}

func ScanPrimitives() ([]PrimitiveInfo, error) {
	pGlob := fmt.Sprintf("%v/*/*.json", filepath.Join(os.Getenv("GOPATH"), "assets", "primitives"))
	matches, err := filepath.Glob(pGlob)
	if err != nil {
		return nil, fmt.Errorf("could not prepare glob to search for primitives: %v", err)
	}

	primitives := make([]PrimitiveInfo, len(matches))

	for idx, match := range matches {
		fileName := filepath.Base(match)
		seedName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		primitiveType := filepath.Base(filepath.Dir(match))
		primitives[idx] = PrimitiveInfo{primitiveType, beautifySeedName(seedName), seedName, match}
	}

	return primitives, nil
}

func ScanPrimitivesByType() (map[string][]PrimitiveInfo, error) {
	primitives, err := ScanPrimitives()
	if err != nil {
		return nil, fmt.Errorf("unable to scan primitives: %v", err)
	}
	pMap := utility.PartitionMany(primitives, func(el PrimitiveInfo) string { return el.Type })
	return pMap, nil
}
