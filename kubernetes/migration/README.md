# K8s/CoreDNS Corefile Migration Tools

This Go library provides a set of functions to help handle migrations of CoreDNS Corefiles to be compatible
with new versions of CoreDNS.  The task of upgrading CoreDNS is the responsibility of a variety of Kubernetes
management tools (e.g. kubeadm and others), and the precise behavior may be different for each one.  This
library abstracts some basic helper functions that make this easier to implement.

## Notifications

Several functions in the library return a list of Notices.  Each Notice is a warning of a feature deprecation,
an unsupported plugin/option, or a new required plugin/option added to the Corefile.  A Notice has a `ToString()`
For display to an end user.  e.g.

```
Plugin "foo" is deprecated in <version>. It is replaced by "bar".
Plugin "bar" is removed in <version>. It is replaced by "qux".
Option "foo" in plugin "bar" is added as a default in <version>.
Plugin "baz" is unsupported by this migration tool in <version>.
```


## Functions

`Deprecated(fromCoreDNSVersion, toCoreDNSVersion, corefileStr string) ([]Notice, error)`

Deprecated returns a list of deprecation notices affecting the given Corefile.  Notices are returned for
any deprecated, removed, or ignored plugins/directives present in the Corefile.  Notices are also returned for
any new default plugins that would be added in a migration.  Notices


`Migrate(fromCoreDNSVersion, toCoreDNSVersion, corefileStr string, deprecations bool) (string, error)`

Migrate returns a migrated version of the Corefile, or an error if it cannot. It will:
  * replace/convert any plugins/directives that have replacements (e.g. _proxy_ -> _forward_)
  * return an error if replacable plugins/directives cannot be converted (e.g. proxy _options_ not available in _forward_)
  * remove plugins/directives that do not have replacements (e.g. kubernetes `upstream`)
  * add in any new default plugins where applicable if they are not already present (e.g. adding _loop_ plugin when it was added).
    This will have to be case by case, and could potentially get complicated.
  * If deprecations is true, deprecated plugins/options will be migrated as soon as they are deprecated.
  * If deprecations is false, deprecated plugins/options will be migrated only once they become removed or ignored.


`Unsupported(fromCoreDNSVersion, toCoreDNSVersion, corefileStr string) ([]Notice, error)`

Unsupported returns a list Notices for plugins/options that are unhandled by this migration tool,
but may still be valid in CoreDNS.  Currently, only a subset of plugins included by default in CoreDNS are supported
by this tool.


`Default(k8sVersion, corefileStr string) bool`

Default returns true if the Corefile is the default for a given version of Kubernetes.
Or, if k8sVersion is empty, Default returns true if the Corefile is the default for any version of Kubernetes.


`Released(dockerImageSHA string) bool`

Released returns true if dockerImageSHA matches any released image of CoreDNS.


`ValidVersions() []string`

ValidVersions returns a list of all versions supported by this tool.


## Command Line Converter Example

An example use of this library is provided [here](../corefile-tool/).


## Kubernetes Cluster Managemnt Tool Usage

This is an example flow of how this library could be used by a Kubernetes cluster management tool to perform a
Corefile migration during an upgrade...

0. Check `Released()` to verify that the installed version of CoreDNS is an official release.
1. Check `Default()`, if the Corefile is a default, simply re-deploy the new default over top the old one. No migration needed.
   If the Corefile is not a default, continue...
2. Check `Deprecated()`, if anything is deprecated, warn user, but continue install.
3. Check `Unsupported()`, if anything is unsupported, abort and warn user (allow user to override to pass this).
4. Call `Migrate()`, if there is an error, abort and warn user.
5. If there is no error, pause and ask user if they want to continue with the migration.  If the starting Corefile was at defaults,
   proceed use the migrated Corefile.



