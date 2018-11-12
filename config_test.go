package envconfig

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type HTTPConfig struct {
	LivenessCheckPath  string
	ReadinessCheckPath string
	HTTPPort           int
	EnablePProf        bool
}

type PubSubConfig struct {
	Topic        string        `envconfig:"topic"`
	MaxExtension time.Duration `default:"15m" envconfig:"max_extension"`
}

type PubSubConfigWithFunnyTags struct {
	Topic        string        `envconfig:"toto"`
	MaxExtension time.Duration `default:"15m" envconfig:"patati_patata"`
}

type Config struct {
	HTTPConfig
	PubSubA PubSubConfig
	PubSubB PubSubConfig
}

type Config2 struct {
	Config Config
	Old    Config
	Name   string
}

type Config3 struct {
	Config
	Old  Config
	Name string
}

type Env map[string]string

func (e Env) setEnv() {
	for k, v := range e {
		os.Setenv(k, v)
	}
}

func (e Env) clearEnv() {
	/*for k := range e {
		os.Setenv(k, "")
	}*/
	os.Clearenv()
}

var expected = &Config{
	HTTPConfig: HTTPConfig{
		HTTPPort:           8088,
		LivenessCheckPath:  "/test/live",
		ReadinessCheckPath: "/test/ready",
		EnablePProf:        true,
	},
	PubSubA: PubSubConfig{
		Topic:        "topicA",
		MaxExtension: time.Duration(11 * time.Hour),
	},
	PubSubB: PubSubConfig{
		Topic:        "topicB",
		MaxExtension: time.Duration(22 * time.Minute),
	},
}

func TestConfigs(t *testing.T) {
	tests := []struct {
		prefix   string
		env      Env
		cfg      interface{}
		expected interface{}
	}{
		{
			"test",
			Env{
				"HTTPPORT":           "8088",
				"LIVENESSCHECKPATH":  "/test/live",
				"READINESSCHECKPATH": "/test/ready",
				"ENABLEPPROF":        "true",
			},
			&HTTPConfig{},
			&expected.HTTPConfig,
		},
		{
			"pubsubb",
			Env{
				"PUBSUBB_TOPIC": "topicB",
				"MAX_EXTENSION": "22m",
			},
			&PubSubConfig{},
			&expected.PubSubB,
		},
		{
			"pubsubb",
			Env{
				"PUBSUBB_TOTO":  "topicB",
				"PATATI_PATATA": "22m",
			},
			&PubSubConfigWithFunnyTags{},
			&PubSubConfigWithFunnyTags{
				Topic:        "topicB",
				MaxExtension: time.Duration(22 * time.Minute),
			},
		},
		{
			"test",
			Env{
				"HTTPPORT":              "8088",
				"LIVENESSCHECKPATH":     "/test/live",
				"READINESSCHECKPATH":    "/test/ready",
				"ENABLEPPROF":           "true",
				"TOPIC":                 "topicA",
				"PUBSUBA_MAX_EXTENSION": "11h",
				"PUBSUBB_TOPIC":         "topicB",
				"MAX_EXTENSION":         "22m",
			},
			&Config{},
			expected,
		},
		{
			"test",
			Env{
				"TEST_HTTPPORT":           "8088",
				"LIVENESSCHECKPATH":       "/test/live",
				"TEST_READINESSCHECKPATH": "/test/ready",
				"TEST_ENABLEPPROF":        "true",
				"PUBSUBA_TOPIC":           "topicA",
				"MAX_EXTENSION":           "11h",
				"TOPIC":                   "topicB",
				"TEST_PUBSUBB_MAX_EXTENSION": "22m",
			},
			&Config{},
			expected,
		},
		{
			"test",
			Env{
				"TEST_CONFIG_HTTPPORT":           "8080",
				"CONFIG_LIVENESSCHECKPATH":       "/cfg/live",
				"TEST_CONFIG_READINESSCHECKPATH": "/cfg/ready",
				"TEST_CONFIG_ENABLEPPROF":        "true",
				"CONFIG_PUBSUBA_TOPIC":           "topicA-prod",
				"MAX_EXTENSION":                  "11h",
				"TOPIC":                          "topicB",
				"TEST_CONFIG_PUBSUBB_MAX_EXTENSION": "33m",
				"NAME":                           "test-multi",
				"TEST_OLD_HTTPPORT":              "8088",
				"OLD_LIVENESSCHECKPATH":          "/test/live",
				"TEST_OLD_READINESSCHECKPATH":    "/test/ready",
				"TEST_OLD_ENABLEPPROF":           "true",
				"OLD_PUBSUBA_TOPIC":              "topicA",
				"TEST_OLD_PUBSUBB_MAX_EXTENSION": "22m",
			},
			&Config2{},
			&Config2{
				Config: Config{
					HTTPConfig: HTTPConfig{
						HTTPPort:           8080,
						LivenessCheckPath:  "/cfg/live",
						ReadinessCheckPath: "/cfg/ready",
						EnablePProf:        true,
					},
					PubSubA: PubSubConfig{
						Topic:        "topicA-prod",
						MaxExtension: time.Duration(11 * time.Hour),
					},
					PubSubB: PubSubConfig{
						Topic:        "topicB",
						MaxExtension: time.Duration(33 * time.Minute),
					},
				},
				Old:  *expected,
				Name: "test-multi",
			},
		},
		{
			"test",
			Env{
				"TEST_HTTPPORT":           "8080",
				"LIVENESSCHECKPATH":       "/cfg/live",
				"TEST_READINESSCHECKPATH": "/cfg/ready",
				"TEST_ENABLEPPROF":        "true",
				"PUBSUBA_TOPIC":           "topicA-prod",
				"MAX_EXTENSION":           "11h",
				"TOPIC":                   "topicB",
				"TEST_PUBSUBB_MAX_EXTENSION": "33m",
				"NAME":                           "test-multi",
				"TEST_OLD_HTTPPORT":              "8088",
				"OLD_LIVENESSCHECKPATH":          "/test/live",
				"TEST_OLD_READINESSCHECKPATH":    "/test/ready",
				"TEST_OLD_ENABLEPPROF":           "true",
				"OLD_PUBSUBA_TOPIC":              "topicA",
				"TEST_OLD_PUBSUBB_MAX_EXTENSION": "22m",
			},
			&Config3{},
			&Config3{
				Config: Config{
					HTTPConfig: HTTPConfig{
						HTTPPort:           8080,
						LivenessCheckPath:  "/cfg/live",
						ReadinessCheckPath: "/cfg/ready",
						EnablePProf:        true,
					},
					PubSubA: PubSubConfig{
						Topic:        "topicA-prod",
						MaxExtension: time.Duration(11 * time.Hour),
					},
					PubSubB: PubSubConfig{
						Topic:        "topicB",
						MaxExtension: time.Duration(33 * time.Minute),
					},
				},
				Old:  *expected,
				Name: "test-multi",
			},
		},
	}
	for i, test := range tests {
		t.Log("test envconfig with ", i, test)
		test.env.setEnv()
		if err := Process(test.prefix, test.cfg); err != nil {
			t.Error(err)
		}
		assert.Equal(t, test.expected, test.cfg)
		test.env.clearEnv()
	}
}
