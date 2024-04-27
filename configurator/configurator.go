package configurator

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"text/template"

	"github.com/jinzhu/copier"
	"github.com/kelseyhightower/envconfig"
)

const (
	//nolint: lll
	defaultTableFormatSplit = `{{range .}}{{usage_key .}},{{usage_type .}},{{usage_default .}},{{usage_required .}},{{usage_description .}}
{{end}}`

	mdFormat = `| KEY	| TYPE	| DEFAULT	| REQUIRED	| DESCRIPTION	|
| 	| 	| 	| 	| 	|
{{range .}}| {{ usage_key . }}	| {{usage_type .}}	| {{usage_default .}}	| {{usage_required .}}	| {{usage_description .}}	|
{{end}}`
)

var copierOption = copier.Option{IgnoreEmpty: true, DeepCopy: true}

type config struct {
	prefix string
	conf   any
}

// Configurator is a struct for getting/setting a configuration
type Configurator struct {
	mu        sync.Mutex
	configMap map[string]config
}

// NewConfigurator returns a new instance of Configurator
func NewConfigurator(configs ...Configuration) *Configurator {
	c := &Configurator{
		configMap: map[string]config{},
	}

	for _, cfg := range configs {
		c.Set(cfg.Name, cfg.Config)
	}

	return c
}

// Set sets a new configuration to a config map
func (c *Configurator) Set(key string, conf any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.configMap[key] = config{conf: conf}
}

// Get returns or creates a Configuration
func (c *Configurator) Get(key string, conf interface{}) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	configuration, exists := c.configMap[key]
	if exists {
		if err := copier.CopyWithOption(conf, configuration.conf, copierOption); err != nil {
			panic(err)
		}

		return conf
	}

	if err := NewConfiguration(conf); err != nil {
		panic(err)
	}

	c.configMap[key] = config{conf: conf}

	return conf
}

// New sets a new configuration
func (c *Configurator) New(key string, conf interface{}, prefixes ...string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := NewConfiguration(conf, prefixes...); err != nil {
		panic(err)
	}

	prefix := ""

	if len(prefixes) > 0 {
		prefix = prefixes[0]
	}

	c.configMap[key] = config{conf: conf, prefix: prefix}

	return conf
}

func (c *Configurator) Print() {
	c.mu.Lock()
	defer c.mu.Unlock()

	var b bytes.Buffer
	tabs := tabwriter.NewWriter(&b, 1, 0, 4, ' ', 0)

	s := ""
	buf := bytes.NewBufferString(s)

	for _, v := range c.configMap {
		envconfig.Usagef(v.prefix, v.conf, buf, defaultTableFormatSplit)
	}

	var newSlice [][]string

	for _, line := range strings.Split(buf.String(), "\n") {
		arrays := strings.Split(line, ",")
		if len(arrays) < 1 {
			continue
		}
		newSlice = append(newSlice, arrays)
	}

	tmpl := prepareTemplate(mdFormat)

	if err := tmpl.Execute(tabs, newSlice); err != nil {
		panic(err)
	}
	if err := tabs.Flush(); err != nil {
		panic(err)
	}

	bytes := b.Bytes()
	lines := 0
	for i, b := range bytes {
		if b == '\n' {
			lines++
		}
		if lines == 1 && b == ' ' {
			bytes[i] = '-'
		}
	}

	io.Copy(os.Stdout, &b)
}

func findElementByIndex(slice []string, index int) string {
	for i, v := range slice {
		if i == index {
			return v
		}
	}

	return ""
}

func prepareTemplate(format string) *template.Template {
	// Specify the default usage template functions
	functions := template.FuncMap{
		"usage_key":         func(v []string) string { return findElementByIndex(v, 0) },
		"usage_type":        func(v []string) string { return findElementByIndex(v, 1) },
		"usage_default":     func(v []string) string { return findElementByIndex(v, 2) },
		"usage_required":    func(v []string) string { return findElementByIndex(v, 3) },
		"usage_description": func(v []string) string { return findElementByIndex(v, 4) },
	}

	tmpl, err := template.New("envconfig").Funcs(functions).Parse(format)
	if err != nil {
		panic(err)
	}

	return tmpl
}
