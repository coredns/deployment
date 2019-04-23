package migration

// This package provides a set of functions to help handle migrations of CoreDNS Corefiles to be compatible with new
// versions of CoreDNS. The task of upgrading CoreDNS is the responsibility of a variety of Kubernetes management tools
// (e.g. kubeadm and others), and the precise behavior may be different for each one. This library abstracts some basic
// helper functions that make this easier to implement.

import (
	"fmt"
	"github.com/coredns/deployment/kubernetes/migration/corefile"
)

// Deprecated returns a list of deprecated plugins or directives present in the Corefile. Each Notice returned is a
// warning, e.g. "plugin 'foo' is deprecated." An empty list returned means there are no deprecated plugins/options
// present in the Corefile.
func Deprecated(fromCoreDNSVersion, toCoreDNSVersion, corefileStr string) ([]Notice, error) {
	return getStatus(fromCoreDNSVersion, toCoreDNSVersion, corefileStr, deprecated)
}

// Removed returns a list of removed plugins or directives present in the Corefile. Each Notice returned is a warning,
// e.g. "plugin 'foo' is no longer supported." An empty list returned means there are no removed plugins/options
// present in the Corefile.
func Removed(fromCoreDNSVersion, toCoreDNSVersion, corefileStr string) ([]Notice, error) {
	return getStatus(fromCoreDNSVersion, toCoreDNSVersion, corefileStr, removed)
}

// Unsupported returns a list of plugins that are not recognized/supported by the migration tool (but may still be valid in CoreDNS).
func Unsupported(fromCoreDNSVersion, toCoreDNSVersion, corefileStr string) ([]Notice, error) {
	return getStatus(fromCoreDNSVersion, toCoreDNSVersion, corefileStr, unsupported)
}

func getStatus(fromCoreDNSVersion, toCoreDNSVersion, corefileStr, status string) ([]Notice, error) {
	notices := []Notice{}
	err := validateVersions(fromCoreDNSVersion, toCoreDNSVersion)
	if err != nil {
		return notices, err
	}
	cf, err := corefile.New(corefileStr)
	if err != nil {
		return notices, err
	}
	v := fromCoreDNSVersion
	for {
		v = Versions[v].nextVersion
		for _, s := range cf.Servers {
			for _, p := range s.Plugins {
				vp, present := Versions[v].plugins[p.Name]
				if status == unsupported {
					if present {
						continue
					}
					notices = append(notices, Notice{Plugin: p.Name, Severity: status, Version: v})
					continue
				}
				if !present {
					continue
				}
				if vp.status == status {
					notices = append(notices, Notice{
						Plugin:     p.Name,
						Severity:   status,
						Version:    v,
						ReplacedBy: vp.replacedBy,
						Additional: vp.additional,
					})
					continue
				}
				for _, o := range p.Options {
					vo, present := Versions[v].plugins[p.Name].options[o.Name]
					if status == unsupported {
						if present {
							continue
						}
						notices = append(notices, Notice{
							Plugin:     p.Name,
							Option:     o.Name,
							Severity:   status,
							Version:    v,
							ReplacedBy: vo.replacedBy,
							Additional: vo.additional,
						})
						continue
					}
					if !present {
						continue
					}
					if vo.status == status {
						notices = append(notices, Notice{Plugin: p.Name, Option: o.Name, Severity: status, Version: v})
						continue
					}
				}
			}
		}
		if v == toCoreDNSVersion {
			break
		}
	}
	return notices, nil
}

// Migrate returns version of the Corefile migrated to toCoreDNSVersion, or an error if it cannot.
// TODO: add newdefault bool parameter?
func Migrate(fromCoreDNSVersion, toCoreDNSVersion, corefileStr string, deprecations bool) (string, error) {
	err := validateVersions(fromCoreDNSVersion, toCoreDNSVersion)
	if err != nil {
		return "", err
	}
	cf, err := corefile.New(corefileStr)
	if err != nil {
		return "", err
	}
	v := fromCoreDNSVersion
	for {
		v = Versions[v].nextVersion
		newSrvs := []*corefile.Server{}
		for i, s := range cf.Servers {
			newPlugs := []*corefile.Plugin{}
			for _, p := range s.Plugins {
				vp, present := Versions[v].plugins[p.Name]
				if !present {
					newPlugs = append(newPlugs, p)
					continue
				}
				if !deprecations && vp.status == deprecated {
					newPlugs = append(newPlugs, p)
					continue
				}
				if vp.action != nil {
					p, err := vp.action(p)
					if err != nil {
						return "", err
					}
					if p == nil {
						// remove plugin, skip options processing
						continue
					}
				}
				newOpts := []*corefile.Option{}
				for _, o := range p.Options {
					vo, present := Versions[v].plugins[p.Name].options[o.Name]
					if !present {
						newOpts = append(newOpts, o)
						continue
					}
					if !deprecations && vo.status == deprecated {
						newOpts = append(newOpts, o)
						continue
					}
					if vo.action == nil {
						newOpts = append(newOpts, o)
						continue
					}
					o, err := vo.action(o)
					if err != nil {
						return "", err
					}
					if o == nil {
						// remove option
						continue
					}
					newOpts = append(newOpts, o)
				}
			CheckForNewOptions:
				for name, vo := range Versions[v].plugins[p.Name].options {
					if vo.status != newdefault {
						continue
					}
					for _, o := range p.Options {
						if name == o.Name {
							continue CheckForNewOptions
						}
					}
					newOpt, err := vo.action(nil)
					if err != nil {
						return "", err
					}
					newOpts = append(newOpts, newOpt)
				}
				newPlugs = append(newPlugs,
					&corefile.Plugin{
						Name:    p.Name,
						Args:    p.Args,
						Options: newOpts,
					})
			}
			if i == 0 {
			CheckForNewPlugins:
				for name, vp := range Versions[v].plugins {
					if vp.status != newdefault {
						continue
					}
					for _, p := range s.Plugins {
						if name == p.Name {
							continue CheckForNewPlugins
						}
					}
					newPlug, err := vp.action(nil)
					if err != nil {
						return "", err
					}
					newPlugs = append(newPlugs, newPlug)
				}
				newSrvs = append(newSrvs,
					&corefile.Server{
						DomPorts: s.DomPorts,
						Plugins:  newPlugs,
					},
				)
			}
		}
		cf = corefile.Corefile{Servers: newSrvs}
		if v == toCoreDNSVersion {
			break
		}
	}
	return cf.ToString(), nil
}

// Default returns true if the Corefile is the default for a given version of Kubernetes.
// Or, if k8sVersion is empty, Default returns true if the Corefile is the default for any version of Kubernetes.
func Default(k8sVersion, corefileStr string) bool {
	cf, err := corefile.New(corefileStr)
	if err != nil {
		return false
	}
NextVersion:
	for _, v := range Versions {
		if k8sVersion != "" && k8sVersion != v.k8sRelease {
			continue
		}
		defCf, err := corefile.New(v.defaultConf)
		if err != nil {
			continue
		}
		// check corefile against k8s release default
		if len(cf.Servers) != len(defCf.Servers) {
			continue NextVersion
		}
		for _, s := range cf.Servers {
			defS, found := s.FindMatch(defCf.Servers)
			if !found {
				continue NextVersion
			}
			if len(s.Plugins) != len(defS.Plugins) {
				continue NextVersion
			}
			for _, p := range s.Plugins {
				defP, found := p.FindMatch(defS.Plugins)
				if !found {
					continue NextVersion
				}
				if len(p.Options) != len(defP.Options) {
					continue NextVersion
				}
				for _, o := range p.Options {
					_, found := o.FindMatch(defP.Options)
					if !found {
						continue NextVersion
					}
				}
			}
		}
		return true
	}
	return false
}

// Released returns true if dockerImageID matches any released image of CoreDNS.
func Released(dockerImageID string) bool {
	for _, v := range Versions {
		if v.dockerImageID == dockerImageID {
			return true
		}
	}
	return false
}

// ValidVersions returns a list of all versions defined
func ValidVersions() []string {
	var vStrs []string
	for vStr, _ := range Versions {
		vStrs = append(vStrs, vStr)
	}
	return vStrs
}

func validateVersion(fromCoreDNSVersion string) error {
	if _, ok := Versions[fromCoreDNSVersion]; !ok {
		return fmt.Errorf("start version '%v' not supported", fromCoreDNSVersion)
	}
	return nil
}

func validateVersions(fromCoreDNSVersion, toCoreDNSVersion string) error {
	err := validateVersion(fromCoreDNSVersion)
	if err != nil {
		return err
	}
	for next := Versions[fromCoreDNSVersion].nextVersion; next != ""; next = Versions[next].nextVersion {
		if next != toCoreDNSVersion {
			continue
		}
		return nil
	}
	return fmt.Errorf("cannot migrate to '%v' from '%v'", toCoreDNSVersion, fromCoreDNSVersion)
}
