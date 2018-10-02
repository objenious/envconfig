package test

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/objenious/kitty"
	"github.com/pkg/errors"
	"github.com/reagere/envconfig"
	"github.com/stretchr/testify/assert"
)

type PubSubConfig struct {
	Topic        string
	MaxExtension time.Duration `default:"15m"`
}

type Config struct {
	kitty.Config
	PubSubA PubSubConfig
	PubSubB PubSubConfig
}

type MapFactories map[reflect.Type]func(func(c interface{}) error) (interface{}, error)

func Process(prefix string, cfg interface{}, factories MapFactories) error {
	v := reflect.ValueOf(cfg).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().FieldByIndex([]int{i})
		name := prefix
		if !field.Anonymous {
			if name != "" {
				name += "_"
			}
			name += field.Name
		}
		name = strings.ToLower(name)
		factory := factories[v.Field(i).Type()]
		if factory != nil {
			c, err := factory(func(c interface{}) error {
				if err := envconfig.Process(name, c); err != nil {
					return errors.Wrap(err, "envconfig process")
				}
				return nil
			})
			if err == nil {
				v.Field(i).Set(reflect.ValueOf(c))
			} else {
				return err
			}
		}
		/*
			el := v.Field(i).Interface()
			switch el.(type) {
			case PubSubConfig:
				c := el.(PubSubConfig)
				if err := envconfig.Process(name, &c); err != nil {
					return err
				}
				v.Field(i).Set(reflect.ValueOf(c))
			case kitty.Config:
				c := el.(kitty.Config)
				if err := envconfig.Process(name, &c); err != nil {
					return err
				}
				v.Field(i).Set(reflect.ValueOf(c))
			}
		*/
	}
	return nil
}

func TestEmbeddedConfigs(t *testing.T) {
	os.Setenv("HTTPPORT", "8088")
	os.Setenv("LIVENESSCHECKPATH", "/test/live")
	os.Setenv("READINESSCHECKPATH", "/test/ready")
	os.Setenv("ENABLEPPROF", "true")
	os.Setenv("PUBSUBA_TOPIC", "topicA")
	os.Setenv("PUBSUBA_MAXEXTENSION", "11h")
	os.Setenv("PUBSUBB_TOPIC", "topicB")
	os.Setenv("PUBSUBB_MAXEXTENSION", "22m")
	cfg := &Config{}
	if err := Process("", cfg, MapFactories{
		reflect.TypeOf(kitty.Config{}): func(ec func(c interface{}) error) (interface{}, error) {
			c := kitty.Config{}
			err := ec(&c)
			return c, err
		},
		reflect.TypeOf(PubSubConfig{}): func(ec func(c interface{}) error) (interface{}, error) {
			c := PubSubConfig{}
			err := ec(&c)
			return c, err
		},
	}); err != nil {
		panic(err)
	}
	assert.Equal(t, 8088, cfg.HTTPPort)
	assert.Equal(t, "/test/live", cfg.LivenessCheckPath)
	assert.Equal(t, "/test/ready", cfg.ReadinessCheckPath)
	assert.Equal(t, true, cfg.EnablePProf)
	assert.Equal(t, 11*time.Hour, cfg.PubSubA.MaxExtension)
	assert.Equal(t, "topicA", cfg.PubSubA.Topic)
	assert.Equal(t, 22*time.Minute, cfg.PubSubB.MaxExtension)
	assert.Equal(t, "topicB", cfg.PubSubB.Topic)
}
