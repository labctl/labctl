package config

import (
	"net/netip"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	clabcore "github.com/srl-labs/containerlab/core"
)

func TestFarEndIP(t *testing.T) {
	lst := map[string]string{
		"10.0.0.1/32": "",

		"10.0.0.0/31": "10.0.0.1/31",
		"10.0.0.1/31": "10.0.0.0/31",
		"10.0.0.2/31": "10.0.0.3/31",
		"10.0.0.3/31": "10.0.0.2/31",

		"10.0.0.1/30": "10.0.0.2/30",
		"10.0.0.2/30": "10.0.0.1/30",
		"10.0.0.0/30": "",
		"10.0.0.3/30": "",
		"10.0.0.4/30": "",
		"10.0.0.5/30": "10.0.0.6/30",
		"10.0.0.6/30": "10.0.0.5/30",
	}

	for k, v := range lst {
		p, err := netip.ParsePrefix(k)
		if err != nil {
			t.Errorf("not a valid IP prefix %s", k)
		}

		n := ipFarEnd(p)
		n2, _ := ipFarEndS(k)

		if !n.IsValid() && v == "" && n2 == "" {
			continue
		}

		if n2 != n.String() {
			t.Errorf("far end %s != %s, expected %s", n, n2, v)
		}

		if v != n.String() {
			t.Errorf("far end of %s, got %s, expected %s", k, n, v)
		}
	}
}

func TestIPLastOctect(t *testing.T) {
	lst := map[string]int{
		"10.0.0.1/32": 1,
		"::1/32":      1,
	}
	for k, v := range lst {
		n, err := netip.ParsePrefix(k)
		if err != nil {
			t.Error(err)
		}
		lo := ipLastOctet(n.Addr())
		if v != lo {
			t.Errorf("far end of %s, got %d, expected %d", k, lo, v)
		}
	}
}

func gettestLink() (*ConfigLink, Dict, Dict) {
	return &ConfigLink{
			Link: clabcore.Link{
				Source: "a",
				Target: "b",
			},
			Vars: map[string]any{},
		},
		map[string]any{
			vkSystemIP: "10.0.0.1/32",
		},
		map[string]any{
			vkSystemIP: "10.0.0.2/32",
		}
}

// func assert(t *testing.T, val, exp any) {
// 	if !cmp.Equal(val, exp) {
// 		_, fn, line, _ := runtime.Caller(1)
// 		t.Errorf("assert failed on line %v in %s\n%s", line, fn, cmp.Diff(val, exp))
// 	}/
// }

func TestLinkName(t *testing.T) {
	l, sv, tv := gettestLink()
	n1, n2, _ := linkName(l, sv, tv)
	assert.Equal(t, n1, "to_b")
	assert.Equal(t, n2, "to_a")

	l.Vars[vkLinkNum] = 1
	n1, n2, _ = linkName(l, sv, tv)
	assert.Equal(t, n1, "to_b_1")
	assert.Equal(t, n2, "to_a_1")
}

func TestLinkIP(t *testing.T) {
	l, sv, tv := gettestLink()
	n1, n2, _ := linkIP(l, sv, tv)
	assert.Equal(t, n1, "1.1.2.0/31")
	assert.Equal(t, n2, "1.1.2.1/31")

	l.Vars[vkLinkNum] = 1
	n1, n2, _ = linkIP(l, sv, tv)
	assert.Equal(t, n1, "1.1.2.2/31")
	assert.Equal(t, n2, "1.1.2.3/31")
}

func TestPrepareLinkVars(t *testing.T) {
	a := make(Dict)
	b := make(Dict)
	l, sv, tv := gettestLink()
	nc := map[string]*NodeConfig{
		l.Source: {Vars: sv},
		l.Target: {Vars: tv},
	}
	_ = prepareLinkVars(l, a, b, nc)
	assert.Equal(t, a, Dict{
		vkFarEnd:   Dict{vkLinkIP: "1.1.2.1/31", vkLinkName: "to_a", vkNodeName: "b"},
		vkLinkIP:   "1.1.2.0/31",
		vkLinkName: "to_b",
	})
	assert.Equal(t, b, Dict{
		vkFarEnd:   Dict{vkLinkIP: "1.1.2.0/31", vkLinkName: "to_b", vkNodeName: "a"},
		vkLinkIP:   "1.1.2.1/31",
		vkLinkName: "to_a",
	})

	l.Vars[vkLinkIP] = []string{"1.1.2.0/16", "1.1.2.1/16"}
	l.Vars[vkLinkName] = "the_same"

	_ = prepareLinkVars(l, a, b, nc)
	assert.Equal(t, a, Dict{
		vkFarEnd:   Dict{vkLinkIP: "1.1.2.1/16", vkLinkName: "the_same", vkNodeName: "b"},
		vkLinkIP:   "1.1.2.0/16",
		vkLinkName: "the_same",
	})
	assert.Equal(t, b, Dict{
		vkFarEnd:   Dict{vkLinkIP: "1.1.2.0/16", vkLinkName: "the_same", vkNodeName: "a"},
		vkLinkIP:   "1.1.2.1/16",
		vkLinkName: "the_same",
	})
}

func TestIPfarEndS(t *testing.T) {
	ipA := "10.0.3.0/31"
	feA, err := ipFarEndS(ipA)
	assert.Equal(t, err, nil)
	assert.Equal(t, feA, "10.0.3.1/31")

	ipA = "10.0.3.0/30"
	feA, err = ipFarEndS(ipA)
	assert.Equal(t, err.Error(), "invalid ip 10.0.3.0/30 - invalid Prefix")
	assert.Equal(t, feA, "")
}

func TestPrepareVarsFromConfig(t *testing.T) {
	// load a topo file
	tests := map[string]struct {
		got          string
		node         string
		envvar       map[string]string
		wantNodeVars map[string]any
	}{
		"env_defined_at_node_level": {
			got:  "test_data/topoa.clab.yml",
			node: "PE1",
			wantNodeVars: map[string]any{
				"clab_kind":            "nokia_srsim",
				"clab_management_ipv4": "",
				"clab_management_ipv6": "",
				"clab_node":            "PE1",
				"clab_role":            "nokia_srsim",
				"clab_system_ip":       "192.0.2.1/32",
				"clab_type":            "SR-1",
				"clab_links": []any{
					Dict{
						"cflow_parameters": string("sampling unicast type interface"),
						"clab_far": Dict{
							"cflow_parameters": "sampling unicast type interface",
							"clab_link_ip":     "192.168.0.2/30",
							"clab_link_name":   "to_PE1",
							"clab_node":        "PE2",
							"isis_iid":         []any{int(0)},
							"port":             "1/1/c1/1",
							"vlan":             100,
						},
						"clab_link_ip":   "192.168.0.1/30",
						"clab_link_name": "to_PE2",
						"isis_iid":       []any{int(0)},
						"port":           "1/1/c1/1",
						"vlan":           100,
					},
				},
				"clab_nodes": Dict{
					"PE1": 123,
					"PE2": 123,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := []clabcore.ClabOption{
				clabcore.WithTopoPath(tc.got, ""),
			}
			c, err := clabcore.NewContainerLab(opts...)
			if err != nil {
				t.Fatal(err)
			}
			nc := PrepareVars(c)
			assert.Equal(t, len(nc), 2)
			comparewith := make(map[string]any)
			theVars := nc[tc.node].Vars
			for k, v := range theVars {
				if k == "clab_nodes" {
					nodes := make(Dict)
					comparewith[k] = nodes
					for k1 := range v.(Dict) {
						nodes[k1] = 123
					}
					continue
				}
				if strings.HasPrefix(k, "clab_") {
					comparewith[k] = v
				}
			}
			assert.Equal(t, comparewith, tc.wantNodeVars)

		})
	}

}
