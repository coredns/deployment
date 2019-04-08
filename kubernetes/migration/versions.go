package migration

import (
	"github.com/coredns/deployment/kubernetes/migration/corefile"
)

type plugin struct {
	status     string
	replacedBy string
	additional string
	action     pluginActionFn
	options    map[string]option
}

type option struct {
	status     string
	replacedBy string
	additional string
	action     optionActionFn
}

type release struct {
	k8sRelease    string
	nextVersion   string
	dockerImageID string
	plugins       map[string]plugin

	// defaultConf hold the default Corefile template packaged with the corresponding k8sRelease.
	// Wildcards are used for fuzzy matching:
	//   "*"   matches exactly one token
	//   "***" matches 0 all remaining tokens on the line
	// Order of server blocks, plugins, and options does not matter.
	// Order of arguments does matter.
	defaultConf string
}

type pluginActionFn func(*corefile.Plugin) (*corefile.Plugin, error)
type optionActionFn func(*corefile.Option) (*corefile.Option, error)

func removePlugin(*corefile.Plugin) (*corefile.Plugin, error) { return nil, nil }
func removeOption(*corefile.Option) (*corefile.Option, error) { return nil, nil }

func renamePlugin(p *corefile.Plugin, to string) (*corefile.Plugin, error) {
	p.Name = to
	return p, nil
}

var Versions = map[string]release{
	"1.5.0": {
		dockerImageID: "7987f0908caf",
		plugins: map[string]plugin{
			"errors": {
				options: map[string]option{
					"consolidate": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod": {
						status: deprecated,
						action: removeOption,
					},
					"endpoint":           {},
					"tls":                {},
					"kubeconfig":         {},
					"namespaces":         {},
					"labels":             {},
					"pods":               {},
					"endpoint_pod_names": {},
					"ttl":                {},
					"noendpoints":        {},
					"transfer":           {},
					"fallthrough":        {},
					"ignore":             {},
				},
			},
			"k8s_external": {
				options: map[string]option{
					"apex": {},
					"ttl":  {},
				},
			},
			"prometheus": {},
			"proxy": {
				status:     removed,
				replacedBy: "forward",
				action:     proxyToForwardPluginAction,
				options:    proxyToForwardOptionsMigrations,
			},
			"forward": {},
			"cache": {
				options: map[string]option{
					"success":  {},
					"denial":   {},
					"prefetch": {},
				},
			},
			"loop":        {},
			"reload":      {},
			"loadbalance": {},
		},
	},
	"1.4.0": {
		nextVersion:   "1.5.0",
		dockerImageID: "a9e015907f63",
		plugins: map[string]plugin{
			"errors": {
				options: map[string]option{
					"consolidate": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod":       {},
					"endpoint":           {}, // TODO: Multiple endpoint deprecation
					"tls":                {},
					"kubeconfig":         {},
					"namespaces":         {},
					"labels":             {},
					"pods":               {},
					"endpoint_pod_names": {},
					"upstream": {
						status: deprecated,
						action: removeOption,
					},
					"ttl":         {},
					"noendpoints": {},
					"transfer":    {},
					"fallthrough": {},
					"ignore":      {},
				},
			},
			"k8s_external": {
				options: map[string]option{
					"apex": {},
					"ttl":  {},
				},
			},
			"prometheus": {},
			"proxy": {
				status:     deprecated,
				replacedBy: "forward",
				action:     proxyToForwardPluginAction,
				options:    proxyToForwardOptionsMigrations,
			},
			"forward": {},
			"cache": {
				options: map[string]option{
					"success":  {},
					"denial":   {},
					"prefetch": {},
				},
			},
			"loop":        {},
			"reload":      {},
			"loadbalance": {},
		},
	},
	"1.3.1": {
		nextVersion:   "1.4.0",
		k8sRelease:    "1.14",
		dockerImageID: "eb516548c180",
		defaultConf: `.:53 {
    errors
    health
    kubernetes * *** {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . *
    cache 30
    loop
    reload
    loadbalance
}`,
		plugins: map[string]plugin{
			"errors": {
				options: map[string]option{
					"consolidate": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod":       {},
					"endpoint":           {}, // TODO: Multiple endpoint deprecation
					"tls":                {},
					"kubeconfig":         {},
					"namespaces":         {},
					"labels":             {},
					"pods":               {},
					"endpoint_pod_names": {},
					"upstream":           {},
					"ttl":                {},
					"noendpoints":        {},
					"transfer":           {},
					"fallthrough":        {},
					"ignore":             {},
				},
			},
			"k8s_external": {
				options: map[string]option{
					"apex": {},
					"ttl":  {},
				},
			},
			"prometheus": {},
			"proxy":      {},
			"forward":    {},
			"cache": {
				options: map[string]option{
					"success":  {},
					"denial":   {},
					"prefetch": {},
				},
			},
			"loop":        {},
			"reload":      {},
			"loadbalance": {},
		},
	},
}

var proxyToForwardOptionsMigrations = map[string]option{
	"policy": {
		action: func(o *corefile.Option) (*corefile.Option, error) {
			if len(o.Args) == 2 && o.Args[1] == "least_conn" {
				o.Name = "force_tcp"
				o.Args = nil
			}
			return o, nil
		},
	},
	"except":       {},
	"fail_timeout": {action: removeOption},
	"max_fails":    {action: removeOption},
	"health_check": {action: removeOption},
	"spray":        {action: removeOption},
	"protocol": {
		action: func(o *corefile.Option) (*corefile.Option, error) {
			if len(o.Args) >= 2 && o.Args[1] == "force_tcp" {
				o.Name = "force_tcp"
				o.Args = nil
				return o, nil
			}
			return nil, nil
		},
	},
}

var proxyToForwardPluginAction = func(p *corefile.Plugin) (*corefile.Plugin, error) { return renamePlugin(p, "forward") }
