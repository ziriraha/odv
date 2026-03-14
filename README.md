# Odoo Dev Tool (odv)

All in one tool for Odoo development.

## Install Guide

To install, run: `go install`. Make sure to have go installed locally and have the proper path variables set.

## Configuration

The configuration is determined by the `~/.odvrc` file. The default configuration is as follows:

```toml
odoo_home = "$ODOO_HOME"

[repositories]
.workspace = ".workspace"
community  = "community"
enterprise = "enterprise"
upgrade    = "upgrade"
```

The `odoo_home` variable is the path to your Odoo installation. The `repositories` section defines the repositories that will be used for development. The key is the name of the repository and the value is the path to the repository relative to the Odoo home directory.
