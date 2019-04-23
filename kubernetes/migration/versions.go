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
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health": {},
			"ready": {
				status: newdefault,
				action: func(*corefile.Plugin) (*corefile.Plugin, error) { return &corefile.Plugin{Name: "ready"}, nil },
			},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod": {
						status: deprecated,
						action: removeOption,
					},
					"endpoint": {
						status: ignored,
						action: useFirstArgumentOnly,
					},
					"tls":                {},
					"kubeconfig":         {},
					"namespaces":         {},
					"labels":             {},
					"pods":               {},
					"endpoint_pod_names": {},
					"upstream": {
						status: ignored,
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
				status:     removed,
				replacedBy: "forward",
				action:     proxyToForwardPluginAction,
				options:    proxyToForwardOptionsMigrations,
			},
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"prefer_udp":     {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
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
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod": {},
					"endpoint": {
						status: ignored,
						action: useFirstArgumentOnly,
					},
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
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"prefer_udp":     {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
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
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod": {},
					"endpoint": {
						status: deprecated,
						action: useFirstArgumentOnly,
					},
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
			"proxy": {
				options: map[string]option{
					"policy":       {},
					"fail_timeout": {},
					"max_fails":    {},
					"health_check": {},
					"except":       {},
					"spray":        {},
					"protocol":     {},
				},
			},
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"prefer_udp":     {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
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
	"1.3.0": {
		nextVersion:   "1.3.1",
		dockerImageID: "2ee68ed074c6",
		plugins: map[string]plugin{
			"errors": {
				options: map[string]option{
					"consolidate": {},
				},
			},
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod":       {},
					"endpoint":           {},
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
			"proxy": {
				options: map[string]option{
					"policy":       {},
					"fail_timeout": {},
					"max_fails":    {},
					"health_check": {},
					"except":       {},
					"spray":        {},
					"protocol":     {},
				},
			},
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"prefer_udp":     {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
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
	"1.2.6": {
		nextVersion:   "1.3.0",
		k8sRelease:    "1.13",
		dockerImageID: "f59dcacceff4",
		defaultConf: `.:53 {
    errors
    health
    kubernetes * *** {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . *
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
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod":       {},
					"endpoint":           {},
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
			"prometheus": {},
			"proxy": {
				options: map[string]option{
					"policy":       {},
					"fail_timeout": {},
					"max_fails":    {},
					"health_check": {},
					"except":       {},
					"spray":        {},
					"protocol":     {},
				},
			},
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"prefer_udp":     {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
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
	"1.2.5": {
		nextVersion:   "1.2.6",
		dockerImageID: "bd254cf72111",
		plugins: map[string]plugin{
			"errors": {},
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod":       {},
					"endpoint":           {},
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
			"prometheus": {},
			"proxy": {
				options: map[string]option{
					"policy":       {},
					"fail_timeout": {},
					"max_fails":    {},
					"health_check": {},
					"except":       {},
					"spray":        {},
					"protocol":     {},
				},
			},
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"prefer_udp":     {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
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
	"1.2.4": {
		nextVersion:   "1.2.5",
		dockerImageID: "d35fe8670379",
		plugins: map[string]plugin{
			"errors": {},
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod":       {},
					"endpoint":           {},
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
			"prometheus": {},
			"proxy": {
				options: map[string]option{
					"policy":       {},
					"fail_timeout": {},
					"max_fails":    {},
					"health_check": {},
					"except":       {},
					"spray":        {},
					"protocol":     {},
				},
			},
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"prefer_udp":     {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
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
	"1.2.3": {
		nextVersion:   "1.2.4",
		dockerImageID: "d46263e07d7a",
		plugins: map[string]plugin{
			"errors": {},
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod":       {},
					"endpoint":           {},
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
			"prometheus": {},
			"proxy": {
				options: map[string]option{
					"policy":       {},
					"fail_timeout": {},
					"max_fails":    {},
					"health_check": {},
					"except":       {},
					"spray":        {},
					"protocol":     {},
				},
			},
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"prefer_udp":     {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
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
	"1.2.2": {
		nextVersion:   "1.2.3",
		k8sRelease:    "1.12",
		dockerImageID: "367cdc8433a4",
		defaultConf: `.:53 {
    errors
    health
    kubernetes * *** {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . *
    cache 30
    loop
    reload
    loadbalance
}`,
		plugins: map[string]plugin{
			"errors": {},
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod":       {},
					"endpoint":           {},
					"tls":                {},
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
			"prometheus": {},
			"proxy": {
				options: map[string]option{
					"policy":       {},
					"fail_timeout": {},
					"max_fails":    {},
					"health_check": {},
					"except":       {},
					"spray":        {},
					"protocol":     {},
				},
			},
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"prefer_udp":     {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
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
	"1.2.1": {
		nextVersion:   "1.2.2",
		dockerImageID: "a575d86d4058",
		plugins: map[string]plugin{
			"errors": {},
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod":       {},
					"endpoint":           {},
					"tls":                {},
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
			"prometheus": {},
			"proxy": {
				options: map[string]option{
					"policy":       {},
					"fail_timeout": {},
					"max_fails":    {},
					"health_check": {},
					"except":       {},
					"spray":        {},
					"protocol":     {},
				},
			},
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"prefer_udp":     {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
			"cache": {
				options: map[string]option{
					"success":  {},
					"denial":   {},
					"prefetch": {},
				},
			},
			"loop": {
				status: newdefault,
				action: func(*corefile.Plugin) (*corefile.Plugin, error) { return &corefile.Plugin{Name: "loop"}, nil },
			},
			"reload":      {},
			"loadbalance": {},
		},
	},
	"1.2.0": {
		nextVersion:   "1.2.1",
		dockerImageID: "da1adafc0e78",
		plugins: map[string]plugin{
			"errors": {},
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod":       {},
					"endpoint":           {},
					"tls":                {},
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
			"prometheus": {},
			"proxy": {
				options: map[string]option{
					"policy":       {},
					"fail_timeout": {},
					"max_fails":    {},
					"health_check": {},
					"except":       {},
					"spray":        {},
					"protocol": {
						status: removed,
						action: proxyRemoveHttpsGoogleProtocol,
					},
				},
			},
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"prefer_udp":     {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
			"cache": {
				options: map[string]option{
					"success":  {},
					"denial":   {},
					"prefetch": {},
				},
			},
			"reload":      {},
			"loadbalance": {},
		},
	},
	"1.1.4": {
		nextVersion:   "1.2.0",
		dockerImageID: "9919f8566026",
		plugins: map[string]plugin{
			"errors": {},
			"log": {
				options: map[string]option{
					"class": {},
				},
			},
			"health":   {},
			"autopath": {},
			"kubernetes": {
				options: map[string]option{
					"resyncperiod":       {},
					"endpoint":           {},
					"tls":                {},
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
			"prometheus": {},
			"proxy": {
				options: map[string]option{
					"policy":       {},
					"fail_timeout": {},
					"max_fails":    {},
					"health_check": {},
					"except":       {},
					"spray":        {},
					"protocol": {
						status: ignored,
						action: proxyRemoveHttpsGoogleProtocol,
					},
				},
			},
			"forward": {
				options: map[string]option{
					"except":         {},
					"force_tcp":      {},
					"expire":         {},
					"max_fails":      {},
					"tls":            {},
					"tls_servername": {},
					"policy":         {},
					"health_check":   {},
				},
			},
			"cache": {
				options: map[string]option{
					"success":  {},
					"denial":   {},
					"prefetch": {},
				},
			},
			"reload":      {},
			"loadbalance": {},
		},
	},
	"1.1.3": {
		nextVersion:   "1.1.4",
		k8sRelease:    "1.11",
		dockerImageID: "b3b94275d97c",
		defaultConf: `.:53 {
    errors
    health
    kubernetes * *** {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . *
    cache 30
    reload
}`},
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

var proxyToForwardPluginAction = func(p *corefile.Plugin) (*corefile.Plugin, error) {
	return renamePlugin(p, "forward")
}

var useFirstArgumentOnly = func(o *corefile.Option) (*corefile.Option, error) {
	if len(o.Args) < 1 {
		return o, nil
	}
	o.Args = o.Args[:1]
	return o, nil
}

var proxyRemoveHttpsGoogleProtocol = func(o *corefile.Option) (*corefile.Option, error) {
	if len(o.Args) > 0 && o.Args[0] == "https_google" {
		return nil, nil
	}
	return o, nil
}
