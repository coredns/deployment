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

The following commands are supported:

- `default`: Default returns true if the Corefile is the default for a that version of Kubernetes. 
If the Kubernetes version is omitted, returns true if the Corefile is the default for any version.

    The following flags are accepted by the `default` command:
     
     - `k8sversion`: The Kubernetes version for which you are checking the default. 
     If the Kubernetes version is omitted, returns true if the Corefile is the default for any version.
     - `corefile`  : The path where your Corefile is located. This flag is mandatory. 

- `deprecated`    : Deprecated returns a list of deprecated plugins or directives present in the Corefile.
    
    The following flags are accepted and mandatory for the `deprecated` command:
    
    - `from`: The CoreDNS version you are migrating from.
    - `to`  : The CoreDNS version you are migrating to.
    - `corefile`  : The path where your Corefile is located.
    
- `migrate`       : Migrate your CoreDNS corefile.

    The following flags are accepted and mandatory for the `migrate` command:
        
    - `from`        : The CoreDNS version you are migrating from.
    - `to`          : The CoreDNS version you are migrating to. This flag is mandatory.
    - `corefile`    : The path where your Corefile is located. This flag is mandatory. 
    - `deprecations`: Specify whether you want to migrate all the deprecations that are present in the current Corefile.
    Specifying `false` will result in the `migrate` command not migrating the deprecated plugins present in the Corefile.
        
- `released`      : Released determines whether your Docker Image ID of a CoreDNS release is valid or not.

    The following flags are accepted and mandatory for the `released` command:
    
    - `dockerImageID` : The docker image ID you want to check.
    
- `removed`       : Removed returns a list of removed plugins or directives present in the Corefile.

    The following flags are accepted and mandatory for the `removed` command:
        
    - `from`        : The CoreDNS version you are migrating from.
    - `to`          : The CoreDNS version you are migrating to. This flag is mandatory.
    - `corefile`    : The path where your Corefile is located. This flag is mandatory. 
    
- `unsupported`   : Unsupported returns a list of plugins that are not recognized/supported by the migration tool (but may still be valid in CoreDNS).

    The following flags are accepted and mandatory for the `unsupported` command:
        
    - `from`        : The CoreDNS version you are migrating from.
    - `to`          : The CoreDNS version you are migrating to. This flag is mandatory.
    - `corefile`    : The path where your Corefile is located. This flag is mandatory. 
    
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
# Migrate CoreDNS from v1.4.0 to v1.5.0 and also migrate all the deprecations 
# that are present in the current Corefile. 
corefile-tool migrate --from 1.4.0 --to 1.5.0 --corefile /path/to/Corefile  --deprecations true

# Migrate CoreDNS from v1.2.2 to v1.3.1 and do not also migrate all the deprecations 
# that are present in the current Corefile.
corefile-tool migrate --from 1.2.2 --to 1.3.1 --corefile /path/to/Corefile  --deprecations false
```

