## Corefile-tool

Corefile-tool is a command line tool which helps you to evaluate and migrate your CoreDNS Corefile Configuration.
It uses the [CoreDNS migration tool library](https://github.com/coredns/deployment/tree/master/kubernetes/migration).

## Usage

Use the following syntax to run the `corefile-tool` command:

```
Usage:
    corefile-tool default --corefile <path> [--k8sversion <k8s-ver>]
    corefile-tool deprecated --from <coredns-ver> --to <coredns-ver> --corefile <path>
    corefile-tool migrate --from <coredns-ver> --to <coredns-ver> --corefile <path> [--deprecations <true|false>]
    corefile-tool released --dockerImageId <id>
    corefile-tool removed --from <coredns-ver> --to <coredns-ver> --corefile <path>
    corefile-tool unsupported --from <coredns-ver> --to <coredns-ver> --corefile <path>
    corefile-tool validversions
```


### Operations

The following operations are supported:

- `default`: returns true if the Corefile is the default for the given version of Kubernetes. If `--k8sversion` is not specified, then this will return true if the Corefile is the default for any version of Kubernetes supported by the tool.

- `deprecated`: returns a list of plugins/options in the Corefile that have been deprecated.

- `migrate`: updates your CoreDNS corefile to be compatible with the `-to` version. Setting the `--deprecations` flag to `true` will migrate plugins/options as soon as they are announced as deprecated.  Setting the `--deprecations` flag to `false` will migrate plugins/options only once they are removed (or made a no-op).  The default is `false`. 

- `released`: determines if the `--dockerImageID` was an official CoreDNS release or not.  Only official releases of CoreDNS are supported by the tool.

- `removed`: returns a list plugins/options in the Corefile that have been removed from CoreDNS.

- `unsupported`: returns a list of plugins/options in the Corefile that are not supported by the migration tool (but may still be valid in CoreDNS).

- `validversions`: Shows the list of CoreDNS versions supported by the this tool.


### Examples

The following examples will help you understand the basic usage of the migration tool.

```bash
# See if the Corefile is the default in CoreDNS v1.4.0. 
corefile-tool default --k8sversion 1.4.0 --corefile /path/to/Corefile
```

```bash
# See deprecated plugins CoreDNS from v1.4.0 to v1.5.0. 
corefile-tool deprecated --from 1.4.0 --to 1.5.0 --corefile /path/to/Corefile
```

```bash
# See unsupported plugins CoreDNS from v1.4.0 to v1.5.0. 
corefile-tool unsupported --from 1.4.0 --to 1.5.0 --corefile /path/to/Corefile
```

```bash
# See removed plugins CoreDNS from v1.4.0 to v1.5.0. 
corefile-tool removed --from 1.4.0 --to 1.5.0 --corefile /path/to/Corefile
```

```bash
# Migrate CoreDNS from v1.4.0 to v1.5.0 and also migrate all the deprecations 
# that are present in the current Corefile. 
corefile-tool migrate --from 1.4.0 --to 1.5.0 --corefile /path/to/Corefile  --deprecations true

# Migrate CoreDNS from v1.2.2 to v1.3.1 and do not also migrate all the deprecations 
# that are present in the current Corefile.
corefile-tool migrate --from 1.2.2 --to 1.3.1 --corefile /path/to/Corefile  --deprecations false
```

