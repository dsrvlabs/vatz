# How to use CLI of VATZ

## Init and start VATZ

Visit [Installation](./installation.md).

## Plugin

**VATZ** binary also supports several plugin commands. In this document, usage of plugin command will be described.

Currently, there are four subcommands under the plugin as subcommand and the more subcommands will be added.
For more details, you can see the help text by adding `--help` flag.

```
~$ ./vatz plugin --help
Plugin commands

Usage:
   plugin [command]

Available Commands:
  install     Install new plugin
  list        List installed plugin
  start       Start installed plugin
  status      Get statuses of Plugin

Flags:
  -h, --help   help for plugin

Use " plugin [command] --help" for more information about a command.
```

First, you can install the plugin by CLI.
For installing the plugin, you have to know the repository URL where the plugin is implemented.

Then, you can try below to install the plugin.

```
~$ ./vatz plugin install github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_active_status cosmos-status
```

The last argument, `cosmos-status` is a simple name that is used for binary name on your machine.
So, you can choose any name as your convenience.

If you install some plugin, you can query installed plugins by list subcommand.

```
~$ ./vatz plugin list
2023-01-06T05:59:27Z INF List plugins module=plugin
2023-01-06T05:59:27Z INF List module=plugin
2023-01-06T05:59:27Z INF newReader /home/rootwarp/.vatz/vatz.db module=db
2023-01-06T05:59:27Z INF Create DB Instance module=db
2023-01-06T05:59:27Z INF List Plugin module=db
+---------------+---------------------+------------------------------------------------------------------------------+---------+
| NAME          | INSTALL DATA        | REPOSITORY                                                                   | VERSION |
+---------------+---------------------+------------------------------------------------------------------------------+---------+
| cosmos-status | 2022-11-29 08:12:28 | https://github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_active_status | latest  |
+---------------+---------------------+------------------------------------------------------------------------------+---------+
```

Then, installed plugin can be started like below.
Arguments that are required for plugin's binary can be provided by `--args` flag and it should be surrounded by quotes like as one string.

```
~$ ./vatz plugin start --plugin cosmos-status --args "--valoperAddr=<YOUR VALOPER ADDRESS>"
2022-11-21T02:57:14Z INF Start plugin cosmos-status%!(EXTRA string=--valoperAddr=<HIDDEN>) module=plugin
2022-11-21T02:57:14Z INF Start plugin cosmos-status module=plugin
2022-11-21T02:57:14Z INF newReader module=db
2022-11-21T02:57:14Z INF Create DB Instance module=db
```

This CLI feature is implementing actively but there are lots of more features should be fixed and implemented.
In near future, below subcommands will be added as soon as possible.
- Stop: Stop executing plugin.
- Status: Show running plugin's status. There is already subcommand but should be enhanced more.
- Remove: Remove installed plugin.
- Update: Update the pluging to new version.

## Help

For more details, you can query helps by adding `--help` flag.
Try like below if you need more detail instructions.

```
~$ ./vatz --help
Usage:
   [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Init
  plugin      Plugin commands
  start       start VATZ
  version     VATZ Version

Flags:
  -h, --help   help for this command

Use " [command] --help" for more information about a command.
```
