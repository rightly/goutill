package goutill

import "testing"

func TestLoader_validator(t *testing.T) {
	c := New()
	SetPath("/test")

	if err := validator(c); err != nil {
		t.Error(err)
	}
}

type test struct {
	Name string `yaml:"name"`
}

func TestLoader_Load(t *testing.T) {
	c := New()
	SetSuffix("s")
	s := &test{}
	err := Load(s)
	if err != nil {
		t.Error(err)
	}
}