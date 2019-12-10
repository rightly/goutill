package config

import (
	"errors"
	"github.com/rightly/goutill"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/configor"
)

type config struct {
	parser     *configor.Configor
	path       string `required:"true"`
	prefix     string
	configName string `required:"true" `
	suffix     string
	indent     string
}

type Config interface {
	SetPrefix(prefix string) Config
	SetSuffix(suffix string) Config
	SetPath(path string) Config
	SetConfigName(configName string) Config
	SetIndent(path string) Config
	SetPrefixToENV() Config
	SetSuffixToENV() Config

	AutoReload(
		flag bool,
		interval time.Duration,
		callback func(interface{}),
	) Config
	Silent(flag bool) Config
	Verbose(flag bool) Config

	Load(v interface{}) error
}

func New() Config {
	c := new(config)
	c.parser = configor.New(
		&configor.Config{
			Debug:                false,
			Verbose:              false,
			Silent:               true,
			AutoReload:           false,
			ErrorOnUnmatchedKeys: true,
		})
	c.indent = "-"

	return c
}

func (c *config) SetPrefix(prefix string) Config {
	c.prefix = prefix
	return c
}

func (c *config) SetSuffix(suffix string) Config {
	c.suffix = suffix
	return c
}

func (c *config) SetPath(path string) Config {
	c.path = path
	return c
}

func (c *config) SetConfigName(configName string) Config {
	c.configName = configName
	return c
}

func (c *config) SetIndent(path string) Config {
	c.path = path
	return c
}

func (c *config) SetPrefixToENV() Config {
	c.prefix = goutill.OS.ENV()
	return c
}

func (c *config) SetSuffixToENV() Config {
	c.suffix = goutill.OS.ENV()
	return c
}

func (c *config) Debug(flag bool) Config {
	c.parser.Debug = flag
	return c
}

func (c *config) AutoReload(
	flag bool,
	interval time.Duration,
	callback func(interface{}),
) Config {

	c.parser.AutoReload = flag
	return c
}

func (c *config) Silent(flag bool) Config {
	c.parser.Silent = flag
	return c
}

func (c *config) Verbose(flag bool) Config {
	c.parser.Verbose = flag
	return c
}

func (c *config) Prefix() string {
	return c.prefix
}
func (c *config) Suffix() string {
	return c.suffix
}
func (c *config) ConfigName() string {
	return c.configName
}

func (c *config) Indent() string {
	return c.indent
}

func (c *config) Load(v interface{}) error {
	if err := validator(c); err != nil {
		return err
	}

	config := c.filename()
	ok, err := exist(config)
	if ok {
		c.load(v, config)
	}

	return err
}

func (c *config) load(v interface{}, filename string) {
	err := c.parser.Load(v, filename)
	if err != nil {
		v = nil
		panic(err)
	}
}

func validator(v interface{}) error {
	rv := reflect.ValueOf(v).Elem()

	for i := 0; i < rv.NumField(); i++ {
		v := rv.Field(i)
		t := rv.Type().Field(i)
		switch t.Type.String() {
		case "string":
			if v.String() == "" && t.Tag.Get("required") == "true" {
				return errors.New(t.Name + " needs to be set")
			}
		}
	}

	return nil
}

func extractExtension(file string) (filename string, ext string) {
	chunks := strings.Split(file, ".")
	for _, v := range chunks[:len(chunks)-1] {
		filename += v
	}
	ext = "." + chunks[len(chunks)-1]

	return
}

func exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, err
	}

	return false, err
}

func (c *config) filename() string {
	f := c.path
	name, ext := extractExtension(c.configName)

	if string(c.path[len(c.path)-1]) != "/" {
		f += "/"
	}
	if c.prefix != "" {
		f += c.prefix + c.indent
	}
	f += name
	if c.suffix != "" {
		f += c.indent + c.suffix
	}
	f += ext
	return f
}
