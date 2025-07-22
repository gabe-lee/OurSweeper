package env_loader

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/joho/godotenv"
)

type (
	Writer = io.Writer
)

const (
	ENV_NAME_TAG  = "env"
	FILE_NAME_TAG = "file"
	DEFAULT_TAG   = "default"
)

func LoadAndFill(env any, errWriter Writer, files ...string) bool {
	success := true
	envMap, err := godotenv.Read(files...)
	if err != nil {
		fmt.Fprintf(errWriter, "%s", err)
		return false
	}
	T := reflect.TypeOf(env)
	if T.Kind() != reflect.Pointer || T.Elem().Kind() != reflect.Struct {
		fmt.Fprintf(errWriter, "parameter `env` must be a pointer to a struct type")
		return false
	}
	TT := T.Elem()
	V := reflect.ValueOf(env)
	for i := range TT.NumField() {
		field := TT.Field(i)
		envName, isEnvField := field.Tag.Lookup(ENV_NAME_TAG)
		if isEnvField {
			envVal, exists := envMap[envName]
			if !exists {
				envDefault, hasDefault := field.Tag.Lookup(DEFAULT_TAG)
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
						fmt.Fprintf(errWriter, "cannot parse env string value to field `%s` type `%s`", field.Name, field.Type.Name())
						success = false
						continue
					}
					*(*[]byte)(fieldPtrRaw) = []byte(envVal)
				case reflect.Bool:
					val, err := strconv.ParseBool(envVal)
					if err != nil {
						fmt.Fprintf(errWriter, "cannot parse env string value to field `%s`: %s", field.Name, err)
						success = false
						continue
					}
					*(*bool)(fieldPtrRaw) = val
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					val, err := strconv.ParseInt(envVal, 0, 64)
					if err != nil {
						fmt.Fprintf(errWriter, "cannot parse env string value to field `%s`: %s", field.Name, err)
						success = false
						continue
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
						fmt.Fprintf(errWriter, "cannot parse env string value to field `%s`: %s", field.Name, err)
						success = false
						continue
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
					fmt.Fprintf(errWriter, "cannot parse env string value to field `%s` type `%s`", field.Name, field.Type.Name())
					success = false
					continue
				}
			}
		} else {
			fileName, isFileField := field.Tag.Lookup(FILE_NAME_TAG)
			if isFileField {
				_, err := os.Stat(fileName)
				fileExists := err == nil
				var data []byte
				exists := false
				if fileExists {
					file, err := os.Open(fileName)
					if err != nil {
						fmt.Fprintf(errWriter, "cannot load file contents `%s` to env field `%s`: %s", fileName, field.Name, err)
						success = false
					} else {
						data, err = io.ReadAll(file)
						if err != nil {
							fmt.Fprintf(errWriter, "cannot load file contents `%s` to env field `%s`: %s", fileName, field.Name, err)
							success = false
						} else {
							exists = true
						}
						file.Close()
					}
				}
				if !exists {
					fileDefault, hasDefault := field.Tag.Lookup(DEFAULT_TAG)
					if hasDefault {
						_, err := os.Stat(fileDefault)
						defaultFileExists := err == nil
						if defaultFileExists {
							file, err := os.Open(fileDefault)
							if err != nil {
								fmt.Fprintf(errWriter, "cannot load file contents (default) `%s` to env field `%s`: %s", fileDefault, field.Name, err)
								success = false
							} else {
								data, err = io.ReadAll(file)
								if err != nil {
									fmt.Fprintf(errWriter, "cannot load file contents (default) `%s` to env field `%s`: %s", fileDefault, field.Name, err)
									success = false
								} else {
									exists = true
								}
								file.Close()
							}
						}
					}
				}
				if exists {
					fieldPtrRaw := unsafe.Pointer(uintptr(V.UnsafePointer()) + field.Offset)
					kind := field.Type.Kind()
					switch kind {
					case reflect.String:
						*(*string)(fieldPtrRaw) = string(data)
					case reflect.Slice:
						E := field.Type.Elem()
						if E.Kind() != reflect.Uint8 {
							fmt.Fprintf(errWriter, "cannot parse file byte contents to field `%s` type `%s`", field.Name, field.Type.Name())
							success = false
							continue
						}
						*(*[]byte)(fieldPtrRaw) = data
					default:
						fmt.Fprintf(errWriter, "cannot parse file byte contents to field `%s` type `%s`", field.Name, field.Type.Name())
						success = false
						continue
					}
				} else {
					success = false
				}
			}
		}
	}
	return success
}
