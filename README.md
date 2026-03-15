# Odoo Dev Tool (odv)

All in one tool for Odoo development.

## Usage

### Installation

To install, run: `go install`. Make sure to have go installed locally and have the proper path variables set.

### Configuration

The configuration is determined by the `~/.odvrc` file. The default configuration is as follows:

```toml
odoo_home = "$ODOO_HOME"
db_prefix = "rd-"
odoo_port = 8069

[repositories]
.workspace = ".workspace"
community  = "community"
enterprise = "enterprise"
upgrade    = "upgrade"
```

The `odoo_home` variable is the path to your Odoo installation. The `repositories` section defines the repositories that will be used for development. The key is the name of the repository and the value is the path to the repository relative to the Odoo home directory.

## Features

### Git

The git features assume that your Community, Enterprise, Upgrade and Workspace repositories are all in `odoo_home`. Pull only accepts version branches, while rebase doesn't accept them. Version branches are the base versions for development, such as master, saas-19.2, 18.0...

### Database

The database module of odv provides a list, duplicate and drop commands for databases. List and drop --all work with a prefix system, where only databases with the specified prefix are list/dropped. The default prefix is `rd-`. This is a trick to avoid operating on the system Postgres databases. You can change the prefix in the configuration or by writing it in the command.

### Utils

The kill-odoo command will find the process running on the Odoo default port (8069) and kill it (port can be changed in the configuration).
