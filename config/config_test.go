package config

import (
	"github.com/rightly/goutill/config"
	"testing"
)

func TestLoader_validator(t *testing.T) {
	c := config.New()
	c.SetPath("/test")

	if err := validator(c); err != nil {
		t.Error(err)
	}
}

type test struct {
	Name string `yaml:"name"`
}

func TestLoader_Load(t *testing.T) {
	c := config.New()
	c.SetSuffix("s")
	s := &test{}
	err := c.Load(s)
	if err != nil {
		t.Error(err)
	}
}
