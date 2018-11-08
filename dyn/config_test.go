package config

import (
	"os"
	"testing"
	"time"

	"github.com/objenious/kitty"
	"github.com/stretchr/testify/assert"
)

type PubSubConfig struct {
	Topic        string        `envconfig:"topic"`
	MaxExtension time.Duration `default:"15m" envconfig:"max_extension"`
}

type Config struct {
	kitty.Config
	PubSubA PubSubConfig
	PubSubB PubSubConfig
}

func TestEmbeddedConfigsNoPrefix(t *testing.T) {
	os.Setenv("HTTPPORT", "8088")
	os.Setenv("LIVENESSCHECKPATH", "/test/live")
	os.Setenv("READINESSCHECKPATH", "/test/ready")
	os.Setenv("ENABLEPPROF", "true")
	os.Setenv("PUBSUBA_TOPIC", "topicA")
	os.Setenv("MAX_EXTENSION", "11h")
	os.Setenv("TOPIC", "topicB")
	os.Setenv("PUBSUBB_MAX_EXTENSION", "22m")
	cfg := &Config{}
	if err := Process("", cfg); err != nil {
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

func TestEmbeddedConfigs(t *testing.T) {
	os.Setenv("TEST_HTTPPORT", "8088")
	os.Setenv("LIVENESSCHECKPATH", "/test/live")
	os.Setenv("TEST_READINESSCHECKPATH", "/test/ready")
	os.Setenv("TEST_ENABLEPPROF", "true")
	os.Setenv("PUBSUBA_TOPIC", "topicA")
	os.Setenv("MAX_EXTENSION", "11h")
	os.Setenv("TOPIC", "topicB")
	os.Setenv("TEST_PUBSUBB_MAX_EXTENSION", "22m")
	cfg := &Config{}
	if err := Process("test", cfg); err != nil {
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
