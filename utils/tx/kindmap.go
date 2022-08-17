package tx

type KindDef struct {
	platform        string   // the scrapligo SSH platform definition
	comment         string   // character to indicate line is a comment
	commit, compare []string // commands to commit and compare sessions
}

var KindMap = map[string]KindDef{
	"vr-sros": {
		platform: "nokia_sros",
		comment:  "#",
		commit:   []string{"validate", "commit", "discard /"},
		compare:  []string{"validate", "compare /", "discard /"},
	},
	"srl": {
		platform: "nokia_srl",
		comment:  "#",
		commit:   []string{"commit now", "discard stay"},
		compare:  []string{"compare", "discard"},
	},
	"ceos": {
		platform: "arista_eos",
	},
	"crpd": {
		platform: "juniper_junos",
	},
	"vr-csr": {
		platform: "cisco_iosxe",
	},
	"vr-n9kv": {
		platform: "cisco_nxos",
	},
	"vr-nxos": {
		platform: "cisco_nxos",
	},
	"vr-pan": {
		platform: "paloalto_panos",
	},
	"vr-veos": {
		platform: "arista_eos",
	},
	"vr-vmx": {
		platform: "juniper_junos",
	},
	"vr-vqfx": {
		platform: "juniper_junos",
	},
	"vr-xrv": {
		platform: "cisco_iosxr",
	},
	"vr-xrv9k": {
		platform: "cisco_iosxr",
	},
}
