# K8s/CoreDNS Corefile Migration Tools

This Go library provides a set of functions to help handle migrations of CoreDNS Corefiles to be compatible with new versions of CoreDNS.  The task of upgrading CoreDNS is the responsibility of a variety of Kubernetes management tools (e.g. kubeadm and others), and the precise behavior may be different for each one.  This library abstracts some basic helper functions that make this easier to implement.

## Proposed functions:

`Deprecated(fromVersion, toVersion, corefile string) []string`: returns a list of deprecated plugins or directives present in the Corefile.

`Removed(fromVersion, toVersion, corefile string) []string`: returns a list of removed plugins or directives present in the Corefile.

`Migrate(fromVersion, toVersion, corefile string, deprecations boolean) (string, error)`: returns an automatically migrated version of the Corefile, or an error if it cannot. It should:
  * replace/convert any plugins/directives that have replacements (e.g. _proxy_ -> _forward_)
  * return an error if replacable plugins/directives cannot be converted (e.g. proxy _options_ not available in _forward_)
  * remove plugins/directives that do not have replacements (e.g. kubernetes `upstream`)
  * if _deprecations_ is set to true, also migrate deprecated plugins/directives.

`Unsupported(fromVersion, toVersion, corefile string) []string`: returns a list of plugins that are not recognized/supported by the migration tool.  We must handle all the default plugins, and we should make an effort to handle the most common non-default plugins. 

Although it would be useful, we cannot for example provide a function `Default(version, corefile string) boolean` that returns  `true` if the corefile is the default for a that version, because each management tool may deploy their own defaults.  So detection of default Corefiles must be implemented by each management tool that requires it.

## Command Line Converter

We should also write a simple command line tool that allows someone to use these functions via the command line.

E.g.

```
Usage:
  corefile-tool deprecated --from <version> --to <version> --corefile <path>
  corefile-tool removed --from <version> --to <version> --corefile <path> [--deprecations]
  corefile-tool migrate --from <version> --to <version> --corefile <path>
  etc ...
```
