package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/labctl/labctl/utils"
	log "github.com/sirupsen/logrus"
	"gotest.tools/v3/assert"
)

func AssertHasPrefix(t *testing.T, str, substr string) {
	if len(str) < len(substr) {
		assert.Equal(t, str, substr)
	}
	assert.DeepEqual(t, str[:len(substr)], substr)
}

func TestLoad(t *testing.T) {
	c := DiscoverTemplate{}

	p := utils.KindPath{
		Kind: "vr-sros",
		Name: "interfaces",
		Ext:  "yml",
	}
	ps, _ := p.Resolve("../../templates")
	c.Load(ps)

	assert.DeepEqual(t, c.Command, "/show router interface")
	AssertHasPrefix(t, c.ParseTemplate, "Value Required,")
	AssertHasPrefix(t, c.OutputTemplate, "{{")
}

func TestAll(t *testing.T) {
	fp, _ := filepath.Abs("../../templates/*__*.yml")
	files, err := filepath.Glob(fp)
	wd, _ := os.Getwd()
	log.Errorf("%s %s", wd, fp)
	assert.Assert(t, err)
	assert.Assert(t, len(files) > 1)
	for _, file := range files {
		log.Errorf("%s ", file)
		c := DiscoverTemplate{}
		err := c.Load(file)
		assert.Assert(t, err)
	}
}

func TestParse(t *testing.T) {
	show := `
    | to_SR-42-6                       Up        Up/Up       Network 1/1/c8/1:6
    |    39.42.6.39/24                                               n/a
    |    2010::272a:627/120                                          PREFERRED`
	c, err := NewDiscoveryConfig("interfaces", "vr-sros")
	assert.Assert(t, err)

	r, err := c.ParseShow(show)
	assert.Assert(t, err)

	expect := []map[string]interface{}{{"address": []string{"39.42.6.39/24", "2010::272a:627/120"},
		"admin":     string("Up"),
		"interface": string("to_SR-42-6"),
		"mode":      string("Network"),
		"oper_v4":   string("Up"),
		"oper_v6":   string("Up"),
		"pfx_state": []string{"n/a", "PREFERRED"},
		"port_sap":  string("1/1/c8/1:6")}}

	assert.DeepEqual(t, r, expect)

	rp, err := c.ProcessShow(show, "")
	assert.Assert(t, err)
	exp2 := []map[string]interface{}{{
		"i":   0,
		"if":  "to_SR-42-6",
		"ip":  "39.42.6.39/24",
		"sap": "1/1/c8/1:6"},
	}
	assert.DeepEqual(t, rp, exp2)

	// and a custom template
	cust := "{{ range $i, $n := . }}- if: {{ $n.interface }}{{end}}"

	rp, err = c.ProcessShow(show, cust)
	assert.Assert(t, err)
	exp2 = []map[string]interface{}{{
		"if": "to_SR-42-6",
	}}
	assert.DeepEqual(t, rp, exp2)

}
