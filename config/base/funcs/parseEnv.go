package funcs

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/BIQ-Cat/easyserver/config/base/types"
	"github.com/joho/godotenv"
)

var ErrEnvNotSet = errors.New("not all environment variables are set")

func ParseEnv(debug bool, defaults *types.EnvConfig) (types.EnvConfig, error) {
	if err := godotenv.Load(); err != nil && debug {
		fmt.Println(fmt.Errorf("WARNING: %w", err))
	}

	err := UnmarshalEnv(defaults)
	return *defaults, err
}

func UnmarshalEnv(cfg *types.EnvConfig) error {
	fields := make(map[string]reflect.Value)
	defaults := make(map[string]reflect.Value)
	cfgValue := reflect.ValueOf(cfg).Elem()
	for i := 0; i < cfgValue.NumField(); i++ {
		fieldInfo := cfgValue.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("name")
		if name == "" {
			name = strings.ToUpper(fieldInfo.Name)
		}

		v := cfgValue.Field(i)
		fields[name] = v

		if reflect.Zero(fieldInfo.Type).Interface() == v.Interface() {
			zero, ok := tag.Lookup("default")
			if ok {
				zeroV, err := parseString(zero, fieldInfo.Type.Kind(), fieldInfo.Name, "default value of "+fieldInfo.Name)
				if err != nil {
					return err
				}

				defaults[name] = zeroV
			}
		} else {
			defaults[name] = v
		}

	}

	for name, v := range fields {
		value, ok := os.LookupEnv(name)
		if !ok {
			if zero, hasZero := defaults[name]; hasZero {
				v.Set(zero)
				continue
			}

			return ErrEnvNotSet
		}
		newV, err := parseString(value, v.Kind(), name, "environment variable "+name)
		if err != nil {
			return err
		}

		v.Set(newV)
	}

	return nil
}

func parseString(raw string, kind reflect.Kind, name string, errorStart string) (value reflect.Value, err error) {
	switch kind {
	case reflect.Int:
		value, err = parseInt(raw, 0, errorStart)
	case reflect.Int8:
		value, err = parseInt(raw, 8, errorStart)
	case reflect.Int16:
		value, err = parseInt(raw, 16, errorStart)
	case reflect.Int32:
		value, err = parseInt(raw, 32, errorStart)
	case reflect.Int64:
		value, err = parseInt(raw, 64, errorStart)
	case reflect.Uint:
		value, err = parseUint(raw, 0, errorStart)
	case reflect.Uint8:
		value, err = parseUint(raw, 8, errorStart)
	case reflect.Uint16:
		value, err = parseUint(raw, 16, errorStart)
	case reflect.Uint32:
		value, err = parseUint(raw, 32, errorStart)
	case reflect.Uint64:
		value, err = parseUint(raw, 64, errorStart)
	case reflect.Float32:
		value, err = parseFloat(raw, 32, errorStart)
	case reflect.Float64:
		value, err = parseFloat(raw, 64, errorStart)
	case reflect.String:
		value = reflect.ValueOf(raw)
		err = nil
	default:
		panic(fmt.Sprintf("unsupported type of %s", name))
	}
	return
}

func parseInt(raw string, bits int, errorStart string) (reflect.Value, error) {
	value, err := strconv.ParseInt(raw, 0, bits)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("%s cannot be converted to integer: %w", errorStart, err)
	}
	switch bits {
	case 0:
		return reflect.ValueOf(int(value)), nil // #nosec G115
	case 8:
		return reflect.ValueOf(int8(value)), nil // #nosec G115
	case 16:
		return reflect.ValueOf(int16(value)), nil // #nosec G115
	case 32:
		return reflect.ValueOf(int32(value)), nil // #nosec G115
	case 64:
		return reflect.ValueOf(int64(value)), nil // #nosec G115
	default:
		panic("bad bits")
	}

}

func parseUint(raw string, bits int, errorStart string) (reflect.Value, error) {
	value, err := strconv.ParseUint(raw, 0, bits)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("%s cannot be converted to unsigned integer: %w", errorStart, err)
	}
	switch bits {
	case 0:
		return reflect.ValueOf(uint(value)), nil // #nosec G115
	case 8:
		return reflect.ValueOf(uint8(value)), nil // #nosec G115
	case 16:
		return reflect.ValueOf(uint16(value)), nil // #nosec G115
	case 32:
		return reflect.ValueOf(uint32(value)), nil // #nosec G115
	case 64:
		return reflect.ValueOf(uint64(value)), nil // #nosec G115
	default:
		panic("bad bits")
	}
}

func parseFloat(raw string, bits int, errorStart string) (reflect.Value, error) {
	value, err := strconv.ParseFloat(raw, bits)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("%s cannot be converted to unsigned integer: %w", errorStart, err)
	}
	switch bits {
	case 32:
		return reflect.ValueOf(float32(value)), nil
	case 64:
		return reflect.ValueOf(float64(value)), nil
	default:
		panic("bad bits")
	}
}
