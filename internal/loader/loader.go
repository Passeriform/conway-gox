package loader

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/passeriform/conway-gox/internal/cell_map"
)

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
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return cell_map.Map{}, errors.New("could not fetch runtime caller to get primitives directory")
	}
	exPath := filepath.Dir(filename)

	primitiveFp := filepath.Join(exPath, "primitives", fmt.Sprintf("%v.json", primitive))

	return LoadFromFile(primitiveFp, padding)
}
