# Buchen Tool

## Install

```shell
go build
sudo cp buchen /usr/local/bin/buchen
```

## Usage

```shell
buchen help
```

```
NAME:
   buchen - time tracking from cli

USAGE:
   buchen [global options] command [command options] [arguments...]

COMMANDS:
   view     Print entries
   edit, e  Open entries file in EDITOR
   start    Start the timer
   stop     Stop the timer
   new      Create new date entry
   init     Init
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

## Init

```shell
# creates a ~/buchen.yaml file
buchen init
```

## Start/Stop Timer

```shell
buchen start
buchen stop
Time: 0,12
```

## View

```shell
buchen
# or
buchen view
```

```
+--------+------+---------+----------------------------+
|  DATE  | TIME | PROJECT |        DESCRIPTION         |
+--------+------+---------+----------------------------+
| 13 Oct | 8,77 | ...     | - TICKET-387               |
|        |      |         |                            |
| 14 Oct | 1,06 | ...     | - ...                      |
|        |      |         |                            |
+--------+------+---------+----------------------------+
```

## Edit

Edit date entries manually in your $EDITOR

```shell
buchen e
# or
buchen edit
```

## Export CSV

Print date entries in CSV format

```shell
# All entries
buchen csv > all.csv
# Filtered by project
buchen csv -p sso > sso.csv
```
