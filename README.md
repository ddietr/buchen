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

```shell
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
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
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

```shell
---
date: 29 Jul
time: 0,12
startedAt: 2021-07-29T07:41:19+02:00
project: ...
description: |
  - TICKET-146
```

## Edit

Edit date entries manually in your $EDITOR

```shell
buchen e
# or
buchen edit
```
