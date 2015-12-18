# snap-plugin-collector-df

snap plugin for collecting free space metrics from df linux tool

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

 Plugin collects specified metrics in-band on OS level

### System Requirements

 - Linux system with df command

### Installation
#### Download df plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [Github Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-df
Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-df
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `/build/rootfs`

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description (optional)
----------|-----------|-----------------------
/intel/disk/{file_system}/space/available_space_inodes | uint64 | Available inodes
/intel/disk/{file_system}/space/available_space_kB | uint64 | Available space in kB
/intel/disk/{file_system}/space/percentage_space_inodes | float64 | Available percentage of inodes
/intel/disk/{file_system}/space/percentage_space_kB | float64 | Available percentage of space in kB
/intel/disk/{file_system}/space/used_space_inodes | uint64 | Used inodes
/intel/disk/{file_system}/space/used_space_kB | uint64 | Used space in kB
/intel/disk/{file_system}/fs_type | string | File system type
/intel/disk/{file_system}/mount_point | string | Mount point of device

### Examples
Example task manifest to use df plugin:
```
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "5s"
    },
    "workflow": {
        "collect": {
            "metrics": {
		        "/intel/disk/dev/sda1/space/available_space_inodes": {},
                "/intel/disk/dev/sda1/space/available_space_kB": {},
                "/intel/disk/dev/sda1/space/percentage_space_inodes": {},
                "/intel/disk/dev/sda1/space/percentage_space_kB": {},
                "/intel/disk/dev/sda1/space/used_space_inodes": {},
                "/intel/disk/dev/sda1/space/used_space_kB": {},
                "/intel/disk/dev/sda2/space/available_space_inodes": {},
                "/intel/disk/dev/sda2/space/available_space_kB": {},
                "/intel/disk/dev/sda2/space/percentage_space_inodes": {},
                "/intel/disk/dev/sda2/space/percentage_space_kB": {},
                "/intel/disk/dev/sda2/space/used_space_inodes": {},
                "/intel/disk/dev/sda2/space/used_space_kB": {},
                "/intel/disk/tmpfs/space/available_space_inodes": {},
                "/intel/disk/tmpfs/space/available_space_kB": {},
                "/intel/disk/tmpfs/space/percentage_space_inodes": {},
                "/intel/disk/tmpfs/space/percentage_space_kB": {},
                "/intel/disk/tmpfs/space/used_space_inodes": {},
                "/intel/disk/tmpfs/space/used_space_kB": {}
           },
            "config": {
            },
            "process": null,
            "publish": null
        }
    }
}

```


### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Patryk Matyjasek](https://github.com/PatrykMatyjasek)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.