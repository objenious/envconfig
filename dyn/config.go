package config

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/reagere/envconfig"
)

// MapFactories maps the factories per type of config
type MapFactories map[reflect.Type]func(func(c interface{}) error) (interface{}, error)

// Process processes the environment variables with a dynamic layer of envconfig processing
func Process(prefix string, cfg interface{}) error {
	value := reflect.ValueOf(cfg)
	elem := value.Elem()
	for i := 0; i < elem.NumField(); i++ {
		fieldType := elem.Type().FieldByIndex([]int{i})
		name := prefix
		if !fieldType.Anonymous {
			if name != "" {
				name += "_"
			}
			name += fieldType.Name
		}
		cv := reflect.New(fieldType.Type)
		name = strings.ToLower(name)
		interf := cv.Interface()
		if err := envconfig.Process(name, interf); err != nil {
			return errors.Wrap(err, "envconfig process")
		}
		elem.Field(i).Set(cv.Elem())
	}
	return nil
}
