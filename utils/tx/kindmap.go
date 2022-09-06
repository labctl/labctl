package tx

import log "github.com/sirupsen/logrus"

type KindDef struct {
	Platform string   `yaml:"platform"` // the scrapligo SSH platform definition
	Comment  string   `yaml:"comment"`  // character to indicate line is a comment
	Commit   []string `yaml:"commit"`   // commands to commit and compare sessions
	Compare  []string `yaml:"compare"`  // commands  compare sessions
}

var KindMap = map[string]KindDef{
	"vr-sros": {
		Platform: "nokia_sros",
		Comment:  "#",
		Commit:   []string{"validate", "commit", "discard /"},
		Compare:  []string{"validate", "compare /", "discard /"},
	},
	"srl": {
		Platform: "nokia_srl",
		Comment:  "#",
		Commit:   []string{"commit validate", "commit now"},
		Compare:  []string{"commit validate", "diff", "discard now"},
	},
	"ceos": {
		Platform: "arista_eos",
	},
	"crpd": {
		Platform: "juniper_junos",
	},
	"vr-csr": {
		Platform: "cisco_iosxe",
	},
	"vr-n9kv": {
		Platform: "cisco_nxos",
	},
	"vr-nxos": {
		Platform: "cisco_nxos",
	},
	"vr-pan": {
		Platform: "paloalto_panos",
	},
	"vr-veos": {
		Platform: "arista_eos",
	},
	"vr-vmx": {
		Platform: "juniper_junos",
	},
	"vr-vqfx": {
		Platform: "juniper_junos",
	},
	"vr-xrv": {
		Platform: "cisco_iosxr",
	},
	"vr-xrv9k": {
		Platform: "cisco_iosxr",
	},
}

// Override KindMap with values from settings
func KindMapOverride(kindmap map[string]KindDef) {
	for k := range kindmap {
		if DebugCount > 2 {
			log.Debugf("Override KindMap[%s]. Please consider raising a PR to include this in labctl", k)
		}
		KindMap[k] = kindmap[k]
	}
}
