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
	k8sRelease     string
	nextVersion    string
	dockerImageSHA string
	plugins        map[string]plugin

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
		dockerImageSHA: "e83beb5e43f8513fa735e77ffc5859640baea30a882a11cc75c4c3244a737d3c",
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
		nextVersion:    "1.5.0",
		dockerImageSHA: "70a92e9f6fc604f9b629ca331b6135287244a86612f550941193ec7e12759417",
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
		nextVersion:    "1.4.0",
		k8sRelease:     "1.14",
		dockerImageSHA: "02382353821b12c21b062c59184e227e001079bb13ebd01f9d3270ba0fcbf1e4",
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
		nextVersion:    "1.3.1",
		dockerImageSHA: "e030773c7fee285435ed7fc7623532ee54c4c1c4911fb24d95cd0170a8a768bc",
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
		nextVersion:    "1.3.0",
		k8sRelease:     "1.13",
		dockerImageSHA: "81936728011c0df9404cb70b95c17bbc8af922ec9a70d0561a5d01fefa6ffa51",
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
		nextVersion:    "1.2.6",
		dockerImageSHA: "33c8da20b887ae12433ec5c40bfddefbbfa233d5ce11fb067122e68af30291d6",
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
		nextVersion:    "1.2.5",
		dockerImageSHA: "a0d40ad961a714c699ee7b61b77441d165f6252f9fb84ac625d04a8d8554c0ec",
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
		nextVersion:    "1.2.4",
		dockerImageSHA: "12f3cab301c826978fac736fd40aca21ac023102fd7f4aa6b4341ae9ba89e90e",
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
		nextVersion:    "1.2.3",
		k8sRelease:     "1.12",
		dockerImageSHA: "3e2be1cec87aca0b74b7668bbe8c02964a95a402e45ceb51b2252629d608d03a",
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
		nextVersion:    "1.2.2",
		dockerImageSHA: "fb129c6a7c8912bc6d9cc4505e1f9007c5565ceb1aa6369750e60cc79771a244",
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
	"1.2.0": {
		nextVersion:    "1.2.1",
		dockerImageSHA: "ae69a32f8cc29a3e2af9628b6473f24d3e977950a2cb62ce8911478a61215471",
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
		nextVersion:    "1.2.0",
		dockerImageSHA: "463c7021141dd3bfd4a75812f4b735ef6aadc0253a128f15ffe16422abe56e50",
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
		nextVersion:    "1.1.4",
		k8sRelease:     "1.11",
		dockerImageSHA: "a5dd18e048983c7401e15648b55c3ef950601a86dd22370ef5dfc3e72a108aaa",
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
