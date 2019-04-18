## Corefile-tool

Corefile-tool is a simple command line tool which helps you to evaluate and migrate your CoreDNS Corefile Configuration.
It is based on the [CoreDNS migration tool library](https://github.com/coredns/deployment/tree/master/kubernetes/migration).

## Usage

Use the following syntax to run the `corefile-tool` command:

`corefile-tool [command] [flags]`

where `command`, `flags` are:

- `command`: The operation you want to perform. 
- `flags`  : Specifies flags required to carry out the operations.


### Operations

The following operations are supported:
- `default`: Default returns true if the Corefile is the default for a that version of Kubernetes. 
If the Kubernetes version is omitted, returns true if the Corefile is the default for any version.

- `deprecated`    : Deprecated returns a list of deprecated plugins or directives present in the Corefile.
- `migrate`       : Migrate your CoreDNS corefile.
- `removed`       : Removed returns a list of removed plugins or directives present in the Corefile.
- `unsupported`   : Unsupported returns a list of plugins that are not recognized/supported by the migration tool (but may still be valid in CoreDNS).
- `validversions` : Shows valid versions of CoreDNS.


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
# Migrate CoreDNS from v1.4.0 to v1.5.0 and handle deprecations . 
corefile-tool migrate --from 1.4.0 --to 1.5.0 --corefile /path/to/Corefile  --deprecations true

# Migrate CoreDNS from v1.2.2 to v1.3.1 and do not handle deprecations .
corefile-tool migrate --from 1.2.2 --to 1.3.1 --corefile /path/to/Corefile  --deprecations false
```

