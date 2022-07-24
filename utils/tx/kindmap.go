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
		commit:   []string{"validate", "commit", "discard"},
		compare:  []string{"validate", "compare /", "discard"},
	},
	"srl": {
		platform: "nokia_sros",
		comment:  "#",
		commit:   []string{"commit now", "discard stay"},
		compare:  []string{"compare", "discard"},
	},
}
