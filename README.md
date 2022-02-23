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
   csv      Print entries as csv
   add      Add description to current entry
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
# start new timer
buchen start
🏃 Start timer for a new day.

# edit entry description to TICKET-1
buchen e
# or use add command
buchen add TICKET-1

buchen stop
💤 Stopped "TICKET-1" at 0,12

buchen start
🏃 Restart "TICKET-1" at 0,12

buchen new "TICKET-2"
💤 Stopped "TICKET-1" at 0,23
🏃 Start "TICKET-2"

# switch to TICKET-1
buchen start
? Choose task:  [Use arrows to move, type to filter]
> TICKET-1
  TICKET-2 🏃
# press enter
💤 Stopped "TICKET-2" at 0,01
🏃 Start "TICKET-1" at 0,02
```

## View

```shell
buchen
# or
buchen view
```

```
+------------+------+---------+--------------------+
|    DATE    | TIME | PROJECT |    DESCRIPTION     |
+------------+------+---------+--------------------+
| 06.02.2022 | 0,35 | ...     | TICKET-1, TICKET-2 |
+------------+------+---------+--------------------+
```

## Export CSV

Print date entries in CSV format

```shell
# All entries
buchen csv > all.csv
# Filtered by project
buchen csv -p sso > sso.csv

# Output
Datum;Beschreibung;Aufwand
18.10.2021;TICKET-358 TICKET-387;8,10
```
