package discovery

import (
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/charmbracelet/log"
	"github.com/chigopher/pathlib"
	"github.com/labctl/labctl/utils"
)

func TestLoad(t *testing.T) {
	c := DiscoverTemplate{}
	p := utils.KindPath{
		Kind: "nokia_srsim",
		Name: "interfaces",
		Ext:  "yml",
	}
	ps, err := p.Resolve("../../templates")
	assert.NoError(t, err)
	err = c.Load(ps)
	assert.NoError(t, err)

	assert.Equal(t, c.Command, "/show router interface")
	assert.HasPrefix(t, c.ParseTemplate, "Value Required,")
	assert.HasPrefix(t, c.OutputTemplate, "{{")
}

func TestAll(t *testing.T) {
	fp := pathlib.NewPath("../../templates")
	files, err := fp.Glob("*__*.y*ml")
	wd, _ := os.Getwd()
	log.Errorf("%s %s", wd, fp)
	assert.NoError(t, err)
	assert.True(t, len(files) > 1)
	for _, file := range files {
		log.Errorf("%s ", file)
		c := DiscoverTemplate{}
		err := c.Load(file.String())
		assert.NoError(t, err)
	}
}

func TestParse(t *testing.T) {
	show := `
    | to_SR-42-6                       Up        Up/Up       Network 1/1/c8/1:6
    |    39.42.6.39/24                                               n/a
    |    2010::272a:627/120                                          PREFERRED`

	c, err := NewDiscoveryConfig("interfaces", "nokia_srsim", "../../templates")
	assert.NoError(t, err)

	r, err := c.ParseShow(show)
	assert.NoError(t, err)

	expect := []map[string]any{{
		"address":   []string{"39.42.6.39/24", "2010::272a:627/120"},
		"admin":     string("Up"),
		"interface": string("to_SR-42-6"),
		"mode":      string("Network"),
		"oper_v4":   string("Up"),
		"oper_v6":   string("Up"),
		"pfx_state": []string{"n/a", "PREFERRED"},
		"port_sap":  string("1/1/c8/1:6"),
	}}

	assert.Equal(t, r, expect)

	rp, err := c.ProcessShow(show, "")
	assert.NoError(t, err)
	exp2 := []map[string]any{
		{
			"i":   0,
			"if":  "to_SR-42-6",
			"ip":  "39.42.6.39/24",
			"sap": "1/1/c8/1:6",
		},
	}
	assert.Equal(t, rp, exp2)

	// and a custom template
	cust := "{{ range $i, $n := . }}- if: {{ $n.interface }}{{end}}"

	rp, err = c.ProcessShow(show, cust)
	assert.NoError(t, err)
	exp2 = []map[string]any{{
		"if": "to_SR-42-6",
	}}
	assert.Equal(t, rp, exp2)
}
