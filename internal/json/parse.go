package json

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	basictypes "github.com/BIQ-Cat/easyserver/config/types"
)

func ParseFiles() error {
	for key, cfg := range Configurations {
		if !(*cfg).HasExternalFile() {
			continue
		}

		raw, err := os.ReadFile((*cfg).GetJSONPath())
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return fmt.Errorf("error while reading file: %w", err)
		}

		res := reflect.ValueOf(*cfg).Interface()

		err = json.Unmarshal(raw, &res)
		*Configurations[key] = res.(basictypes.JSONConfig)
		if err != nil {
			return err
		}
	}
	return nil
}
