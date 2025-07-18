package env_loader

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/joho/godotenv"
)

const (
	ENV_NAME_TAG    = "env"
	ENV_DEFAULT_TAG = "default"
)

func LoadInto(env any, files ...string) error {
	envMap, err := godotenv.Read(files...)
	if err != nil {
		return err
	}
	T := reflect.TypeOf(env)
	if T.Kind() != reflect.Pointer || T.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("parameter `env` must be a pointer to a struct type")
	}
	TT := T.Elem()
	V := reflect.ValueOf(env)
	for i := range TT.NumField() {
		field := TT.Field(i)
		envName, isEnvField := field.Tag.Lookup(ENV_NAME_TAG)
		if isEnvField {
			envVal, exists := envMap[envName]
			if !exists {
				envDefault, hasDefault := field.Tag.Lookup(ENV_DEFAULT_TAG)
				if hasDefault {
					envVal = envDefault
					exists = true
				}
			}
			if exists {
				fieldPtrRaw := unsafe.Pointer(uintptr(V.UnsafePointer()) + field.Offset)
				kind := field.Type.Kind()
				switch kind {
				case reflect.String:
					*(*string)(fieldPtrRaw) = envVal
				case reflect.Slice:
					E := field.Type.Elem()
					if E.Kind() != reflect.Uint8 {
						return fmt.Errorf("cannot parse env string value to field `%s` type `%s`", field.Name, field.Type.Name())
					}
					*(*[]byte)(fieldPtrRaw) = []byte(envVal)
				case reflect.Bool:
					val, err := strconv.ParseBool(envVal)
					if err != nil {
						return err
					}
					*(*bool)(fieldPtrRaw) = val
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					val, err := strconv.ParseInt(envVal, 0, 64)
					if err != nil {
						return err
					}
					switch kind {
					case reflect.Int:
						*(*int)(fieldPtrRaw) = int(val)
					case reflect.Int8:
						*(*int8)(fieldPtrRaw) = int8(val)
					case reflect.Int16:
						*(*int8)(fieldPtrRaw) = int8(val)
					case reflect.Int32:
						*(*int8)(fieldPtrRaw) = int8(val)
					case reflect.Int64:
						*(*int8)(fieldPtrRaw) = int8(val)
					default:
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					val, err := strconv.ParseUint(envVal, 0, 64)
					if err != nil {
						return err
					}
					switch kind {
					case reflect.Uint:
						*(*int)(fieldPtrRaw) = int(val)
					case reflect.Uint8:
						*(*int8)(fieldPtrRaw) = int8(val)
					case reflect.Uint16:
						*(*int8)(fieldPtrRaw) = int8(val)
					case reflect.Uint32:
						*(*int8)(fieldPtrRaw) = int8(val)
					case reflect.Uint64:
						*(*int8)(fieldPtrRaw) = int8(val)
					default:
					}
				default:
					return fmt.Errorf("cannot parse env string value to field `%s` type `%s`", field.Name, field.Type.Name())
				}
			}
		}
	}
	return nil
}
