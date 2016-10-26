# snap-plugin-collector-df

snap plugin for collecting free space metrics from df linux tool

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

 Plugin collects specified metrics in-band on OS level.

### System Requirements

* Linux system with df command

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

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

### Configuration and Usage

* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-df/blob/master/README.md#examples).

## Documentation

### Collected Metrics

List of collected metrics is described in [METRICS.md](https://github.com/intelsdi-x/snap-plugin-collector-df/blob/master/METRICS.md).


### Examples

Example running snap-plugin-collector-df plugin and writing data to a file.

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=<snapDirectoryPath>/build
```
Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

In one terminal window, open the snap daemon (in this case with logging set to 1,  trust disabled):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0
```
In another terminal window:

Load snap-plugin-collector-df plugin:
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-df
```
Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-publisher-file
```
See available metrics:

```
$ $SNAP_PATH/bin/snapctl metric list
```

Create a task manifest file to use snap-plugin-collector-df plugin (exemplary file in [examples/task/] (https://github.com/intelsdi-x/snap-plugin-collector-df/blob/master/examples/task/)):
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
                "/intel/procfs/filesystem/*/space_free": {},
                "/intel/procfs/filesystem/*/space_reserved": {},
                "/intel/procfs/filesystem/*/inodes_percent_free": {},
                "/intel/procfs/filesystem/*/inodes_percent_used": {},
                "/intel/procfs/filesystem/*/device_name": {},
                "/intel/procfs/filesystem/*/inodes_used": {}
            },
            "config": {
                "/intel/procfs/filesystem": {
                    "proc_path": "/proc",
                    "keep_original_mountpoint": true,
                    "excluded_fs_names": "/proc/sys/fs/binfmt_misc,/var/lib/docker/aufs",
                    "excluded_fs_types": "proc,binfmt_misc,fuse.gvfsd-fuse,sysfs,cgroup,fusectl,pstore,debugfs,securityfs,devpts,mqueue,hugetlbfs,nsfs,rpc_pipefs,devtmpfs,none,tmpfs,aufs"
                }
            },
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_df.log"
                    }
                }
            ]
        }
    }
}

```
Create a task:
```
$ $SNAP_PATH/bin/snapctl task create -t df-file.json
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-users/issues) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-users/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [snap Gitter channel](https://gitter.im/intelsdi-x/snap).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Patryk Matyjasek](https://github.com/PatrykMatyjasek)
* Author: [Marcin Krolik](https://github.com/marcin-krolik)
