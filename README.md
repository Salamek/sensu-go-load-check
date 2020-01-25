[![Bonsai Asset Badge](https://img.shields.io/badge/Sensu%20Go%20Load%20Check-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/Salamek/sensu-go-load-check) [![TravisCI Build Status](https://travis-ci.org/Salamek/sensu-go-load-check.svg?branch=master)](https://travis-ci.org/Salamek/sensu-go-load-check)

# Sensu Go Load Check
- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Asset definition](#asset-definition)
  - [Check definition](#resource-definition)
- [Installation from source and contributing](#installation-from-source-and-contributing)
- [Additional notes](#additional-notes)

## Overview

This plugin provides native load instrumentation including load health and per-core metrics. The `sensu-go-load-check` check takes the flags `-w` (warning) and `-c` (critical) and a desired L1/L5/L15. By default, these are a warning value of 2.75, 2.5, 2.0 and a critical value of 3.5, 3.25, 3.0. This check also outputs data as `nagios_perfdata`(for more information, see [this Nagios documentation article](https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/perfdata.html). This allows for the check to be used as both a status check and a metric check. You can see an example of this in the [example check definition](#check-definition) below.

## Usage Examples

### Command line help

```
The Sensu Go check for system Load usage

Usage:
  sensu-go-load-check [flags]

Flags:
  -c, --critical string   Critical value for system load (default "3.5, 3.25, 3.0")
  -h, --help              help for sensu-go-load-check
  -w, --warning string    Warning value for system load. (default "2.75, 2.5, 2.0")
```

### Example Output

```bash
./sensu-go-load-check
CheckLoad OK - value = 0.14, 0.10, 0.07 | core_load_1=0.14, core_load_5=0.10, core_load_15=0.07
```

## Configuration

### Asset registration

Assets are the best way to make use of this check. If you're not using this plugin as an asset, please consider doing so! If you're using Sensu 5.13 or later, you can install this plugin as an asset by running:

`sensuctl asset add Salamek/sensu-go-load-check`

Else, you can find this asset on the [Bonsai Asset Index](https://bonsai.sensu.io/assets/Salamek/sensu-go-load-check).

### Asset definition

You can download the asset definition there, or you can do a little bit of copy/pasta and use the one below:

```yml
---
type: Asset
api_version: core/v2
metadata:
  name: sensu-go-load-check
  namespace: CHANGEME
  labels: {}
  annotations: {}
spec:
  url: https://github.com/asachs01/sensu-go-load-check/releases/download/0.0.1/sensu-go-load-check_0.0.1_linux_amd64.tar.gz
  sha512:
  filters:
  - entity.system.os == 'linux'
  - entity.system.arch == 'amd64'
```

**NOTE**: ***PLEASE ENSURE YOU UPDATE YOUR URL AND SHA512 BEFORE USING THE ASSET***. If you don't, you might just be stuck on a super old version. Don't say I didn't warn you ¯\\_(ツ)_/¯

### Check definition

Example Sensu Go definition:

**sensu-go-load-check**
```yml
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-go-load-check
  namespace: CHANGEME
spec:
  command: sensu-go-load-check
  runtime_assets:
  - sensu-go-load-check
  interval: 60
  publish: true
  output_metric_format: nagios_perfdata
  output_metric_handlers:
  - influxdb
  handlers:
  - slack
  subscriptions:
  - system
```

## Installation from source and contributing

While it's generally recommended to use an asset, you can download a copy of the handler plugin from [releases][1],
or create an executable script from this source.

From the local path of the sensu-go-load-check repository:

**sensu-go-load-check**
```
go build -o /usr/local/bin/sensu-go-load-check main.go
```
To contribute to this repo, see https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md.

## Additional notes

### Supported operating systems

This project uses `gopsutil`, and is thus largely dependent on the systems that it supports. For this plugin, the following operating systems are supported:

* Linux
* FreeBSD
* OpenBSD
* Mac OS X
* Windows
* Solaris

[1]: https://github.com/Salamek/sensu-go-load-check/releases
