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
/intel/linux/filesystem/\<mount_point\>/inodes_free | uint64 | the number of free inodes on the file system
/intel/linux/filesystem/\<mount_point\>/inodes_reserved | uint64 | the number of reserved inodes
/intel/linux/filesystem/\<mount_point\>/inodes_used | uint64 | the number of used inodes
/intel/linux/filesystem/\<mount_point\>/space_free | uint64 | the number of free bytes
/intel/linux/filesystem/\<mount_point\>/space_reserved | uint64 | the number of reserved bytes
/intel/linux/filesystem/\<mount_point\>/space_used | uint64 | the number of used bytes
/intel/linux/filesystem/\<mount_point\>/inodes_percent_free | float64 | the percentage of free inodes on the file system
/intel/linux/filesystem/\<mount_point\>/inodes_percent_reserved | float64 | the percentage of reserved inodes
/intel/linux/filesystem/\<mount_point\>/inodes_percent_used | float64 | the percentage of used inodes
/intel/linux/filesystem/\<mount_point\>/space_percent_free | float64 | the percentage of free bytes
/intel/linux/filesystem/\<mount_point\>/space_percent_reserved | float64 | the percentage of reserved bytes
/intel/linux/filesystem/\<mount_point\>/space_percent_used | float64 | the percentage of used bytes
/intel/linux/filesystem/\<mount_point\>/device_name | string | device name as presented in filesystem (eg. /dev/sda1)
/intel/linux/filesystem/\<mount_point\>/device_type | string | device type as presented in filesystem (eg. ext4)

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
		        "/intel/linux/filesystem/rootfs/space_free": {},
                "/intel/linux/filesystem/rootfs/space_reserved": {},
                "/intel/linux/filesystem/rootfs/inodes_percent_free": {},
                "/intel/linux/filesystem/rootfs/inodes_percent_used": {},
                "/intel/linux/filesystem/rootfs/device_name": {},
                "/intel/linux/filesystem/sys_fs_cgroup/inodes_used": {}
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
* Author: [Marcin Krolik](https://github.com/marcin-krolik)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.