package migration

import "fmt"

// Notice is a migration warning
type Notice struct {
	Plugin     string
	Option     string
	Severity   string // 'deprecated', 'removed', or 'unsupported'
	ReplacedBy string
	Additional string
	Version    string
}

func (n *Notice) ToString() string {
	s := ""
	if n.Option == "" {
		s += fmt.Sprintf(`Plugin "%v" `, n.Plugin)
	} else {
		s += fmt.Sprintf(`Option "%v" in plugin "%v" `, n.Option, n.Plugin)
	}
	if n.Severity == unsupported {
		s += "is unsupported by this migration tool in " + n.Version + "."
	} else if n.Severity == newdefault {
		s += "is added as a default in " + n.Version + "."
	} else {
		s += "is " + n.Severity + " in " + n.Version + "."
	}
	if n.ReplacedBy != "" {
		s += fmt.Sprintf(` It is replaced by "%v".`, n.ReplacedBy)
	}
	if n.Additional != "" {
		s += " " + n.Additional
	}
	return s
}

const (
	deprecated  = "deprecated"  // plugin/option is deprecated in CoreDNS
	ignored     = "ignored"     // plugin/option is ignored by CoreDNS
	removed     = "removed"     // plugin/option has been removed from CoreDNS
	unsupported = "unsupported" // plugin/option is not supported by the migration tool
	newdefault  = "newdefault"  // plugin/option was added to the default corefile
	all			= "all"         // all plugin/option that are deprecated, ignored and removed.
)
