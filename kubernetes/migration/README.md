# K8s/CoreDNS Corefile Migration Tools

This Go library provides a set of functions to help handle migrations of CoreDNS Corefiles to be compatible with new versions of CoreDNS.  The task of upgrading CoreDNS is the responsibility of a variety of Kubernetes management tools (e.g. kubeadm and others), and the precise behavior may be different for each one.  This library abstracts some basic helper functions that make this easier to implement.

## Proposed functions:

`Deprecated(fromCoreDNSVersion, toCoreDNSVersion, corefile string) []string`: returns a list of deprecated plugins or directives present in the Corefile. Each string returned is a warning, e.g. "plugin 'foo' is deprecated." An empty list returned means there are no deprecated plugins/options present in the Corefile.

`Removed(fromCoreDNSVersion, toCoreDNSVersion, corefile string) []string`: returns a list of removed plugins or directives present in the Corefile. Each string returned is a warning, e.g. "plugin 'foo' is no longer supported." An empty list returned means there are no removed plugins/options present in the Corefile.

`Migrate(fromCoreDNSVersion, toCoreDNSVersion, corefile string, deprecations boolean) (string, error)`: returns an automatically migrated version of the Corefile, or an error if it cannot. It should:
  * replace/convert any plugins/directives that have replacements (e.g. _proxy_ -> _forward_)
  * return an error if replacable plugins/directives cannot be converted (e.g. proxy _options_ not available in _forward_)
  * remove plugins/directives that do not have replacements (e.g. kubernetes `upstream`)
  * add in any new default plugins where applicable if they are not already present (e.g. adding _loop_ plugin when it was added).  This will have to be case by case, and could potentially get complicated.
  * if _deprecations_ is set to true, also migrate deprecated plugins/directives.

`Unsupported(fromCoreDNSVersion, toCoreDNSVersion, corefile string) []string`: returns a list of plugins that are not recognized/supported by the migration tool (but may still be valid in CoreDNS).  We must handle all the default plugins, and we should make an effort to handle the most common non-default plugins. Each string returned is a warning, e.g. "plugin 'foo' is not supported by the migration tool." An empty list returned means there are no unsupported plugins/options present in the Corefile.

Although it would be useful, we cannot for example provide a function `Default(coreDNSVersion, corefile string) boolean` that returns  `true` if the corefile is the default for a that version, because each management tool may deploy their own defaults.  So detection of default Corefiles must be implemented by each management tool that requires it.

## Command Line Converter

We should also write a simple command line tool that allows someone to use these functions via the command line.

E.g.

```
Usage:
  corefile-tool deprecated --from <coredns-version> --to <coredns-version> --corefile <path>
  corefile-tool removed --from <coredns-version> --to <coredns-version> --corefile <path> [--deprecations]
  corefile-tool migrate --from <coredns-version> --to <coredns-version> --corefile <path>
  etc ...
```

## Example of Usage

This is an example of how these tools could be used by an installer/upgrader... 

1. check Deprecated(), if anything is deprecated, warn user, but continue install. 
2. check Unsupported(), if anything is unsupported, abort and warn user (allow user to override to pass this).
3. call Migrate(), if there is an error, abort and warn user.
4. If there is no error, and the starting Corefile was not a default, pause and ask user if they want to continue with the migration.  If the starting Corefile was at defaults, proceed use the migrated corefile.


