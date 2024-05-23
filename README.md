# Koksmat

The name `koksmat` is a Danish word that translates into the assistant to a chef on a ship. The Koksmat is responsible for preparing the food and keeping the kitchen clean.

So see `koksmat` as you assistanct in your kitchen.

`koksmat` like to have a clear set of rules on how the kitchen is organized, where the secret ingredients are stored, and how to prepare the food.

As a chef, you can tell `koksmat` what you like to have done, and `koksmat` will do it for you.

The digital manifestation of `koksmat` is a Command Line Interface (CLI) written in Go (golang). You can install that on any operating system that supports Go.

Currently, `koksmat` is only being developed for Linux, but it should be possible to run it on Windows and MacOS as well.

## Installation

To install `koksmat` you need to have Go installed on your computer. If you don't have Go installed, you can download it from [golang.org](https://golang.org/dl/).

Once you have Go installed, you can install `koksmat` by running the following command in your terminal:

```bash
go install github.com/koksmat-com/koksmat@latest
```

## Usage

Once installed, you can run `koksmat` by typing `koksmat` in your terminal.

```bash
koksmat
```

Running `koksmat` will show you a list of available commands.

```bash
koksmat


##  ###   ## ##   ##  ###   ## ##   ##   ##    ##     #### ##
##  ##   ##   ##  ##  ##   ##   ##   ## ##      ##    # ## ##
## ##    ##   ##  ## ##    ####     # ### #   ## ##     ##
## ##    ##   ##  ## ##     #####   ## # ##   ##  ##    ##
## ###   ##   ##  ## ###       ###  ##   ##   ## ###    ##
##  ##   ##   ##  ##  ##   ##   ##  ##   ##   ##  ##    ##
##  ###   ## ##   ##  ###   ## ##   ##   ##  ###  ##   ####


Usage:
  koksmat [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  connector   connector
  help        Help about any command
  kitchen     kitchen
  sail        Auto pilot mode

Flags:
  -h, --help     help for koksmat
  -t, --toggle   Help message for toggle

Use "koksmat [command] --help" for more information about a command.
```

## Command: `kitchen`

The `kitchen` command is the main command for `koksmat`. It is the command that you will use to interact with `koksmat`.

```bash
koksmat kitchen --help
```

Outputs a list of available subcommands for the `kitchen` command.

```bash
kitchen

Usage:
  koksmat kitchen [[service]] [flags]
  koksmat kitchen [command]

Available Commands:
  boot        Boot kitchens
  build       Build kitchen
  create      Create a new kitchen
  launch      Launch kitchen
  open        Open kitchen
  script      Working with scripts
  stations    List stations in kitchen
  status      Get status of kitchen

Flags:
  -h, --help   help for kitchen

Use "koksmat kitchen [command] --help" for more information about a command.
```

Install-Module -Name PnP.PowerShell -Force -AllowClobber
Install-Module -Name ExchangeOnlineManagement -force
koksmat ship get mate
cp .koksmat/kitchenroot/\* .. -r
cp .koksmat/kitchenroot/.koksmat .. -r
echo #x > ../.env
koksmat sail
