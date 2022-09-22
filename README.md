# ADCTL
## Overview

`auto deploy ctl` similar as helm, Adctl is a tool for managing Charts,Charts are package of pre-configured docker-compose yaml resource

To Use Adctl:

- auto deploy docker applications via docker-compose
- share your own applications as Adctl Charts
- Manage releases of Adctl packages

## Install 
`adctl` is a command line program to deploy  applications via docker-compose.

It can be installed by running:

```
go install github.com/shikingram/adctl@latest
```

## Usage

### Precondition
docker and docker-compose are required for this tool,The minimum required version is
```
$docker --version 
Docker version 20.10.11, build dea9396

$docker-compose --version  
docker-compose version 1.29.2, build 5becea4c
```
### adctl nstall
This command can render all yaml files of an application, and then deploy it to docker

eg:
```
adctl install name chart
```
The chart package should have the following structure:
```
example-chart
├── Chart.yaml
├── templates
│   ├── 01-app-mysql.yaml.gtpl
│   ├── 02-job-init-databases.yaml.gtpl
│   ├── 03-app-register-center.yaml.gtpl
│   ├── NOTES.txt
│   └── config
│       ├── init-databases
│       │   └── config.env.gtpl
│       └── mysql
│           └── config.gtpl
└── values.yaml
```
After executing `adctl install`, we will generate an instance in the current directory and start it in docker vi docker-compose

### adctl upgrade
This command upgrades a release to a new version of a chart

When we modify some configurations, we can execute this command for `up-to-date`

### adctl uninstall
This command takes a release name and uninstalls the release.

However, in order to protect data from being lost, the locally generated directory `instance`will not be deleted. If you really need to delete it, you can use flag `--clean-instance`

## More Informations

for more informations, try to use `adctl --help` to view details 

