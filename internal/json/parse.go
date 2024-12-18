package json

import (
	"encoding/json"
	"fmt"
	"os"
)

func ParseFiles() error {
	for _, cfg := range Configurations {
		if !(*cfg).HasExternalFile() {
			continue
		}

		raw, err := os.ReadFile((*cfg).GetJSONPath())
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return fmt.Errorf("error while reading file: %w", err)
		}

		err = json.Unmarshal(raw, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}
