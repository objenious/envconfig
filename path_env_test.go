package envconfig

import "testing"

type cfPath struct {
	JSONPath string `envconfig:"json_path" default:"test"`
}

func TestPath(t *testing.T) {
	config := &cfPath{}
	MustProcess("", config)
	if config.JSONPath != "test" {
		t.Logf("should be default, not path : %+v", config)
		t.Fail()
	}
}
