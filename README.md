DISCONTINUATION OF PROJECT. 

This project will no longer be maintained by Intel.

This project has been identified as having known security escapes.

Intel has ceased development and contributions including, but not limited to, maintenance, bug fixes, new releases, or updates, to this project.  

Intel no longer accepts patches to this project.
# DISCONTINUATION OF PROJECT 

**This project will no longer be maintained by Intel.  Intel will not provide or guarantee development of or support for this project, including but not limited to, maintenance, bug fixes, new releases or updates.  Patches to this project are no longer accepted by Intel. If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the community, please create your own fork of the project.**

# Snap plugin collector - df

Snap plugin for collecting free space metrics from df linux tool

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
All OSs currently supported by Snap:
* Linux/amd64

### Installation
#### Download the plugin binary:

You can get the pre-built binaries for your OS and architecture from the plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-collector-df/releases) page. Download the plugin from the latest release and load it into `snapteld` (`/opt/snap/plugins` is the default location for snap packages).

#### To build the plugin binary:

Fork https://github.com/intelsdi-x/snap-plugin-collector-df
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-df.git
```

Build the snap df plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build/`

### Configuration and Usage

* Set up the [Snap framework](https://github.com/intelsdi-x/snap#getting-started).
* Load the plugin and create a task, see example in [Examples](#examples).
* Available configuration:

| Namespace                    | Data Type | Default Value | Description |
|-----------------------------|----------|-------------------------|------|
| **proc_path**                | string    | `/proc` | Path to `/proc` filesystem |
| **excluded_fs_names**        | []string  | <ul><li>`/proc/sys/fs/binfmt_misc`</li><li>`/var/lib/docker/aufs`</li></ul> | List of excluded mount points |
| **excluded_fs_types**        | []string  | <ul><li>`proc`</li><li>`binfmt_misc`</li><li>`fuse.gvfsd-fuse`</li><li>`sysfs`</li><li>`cgroup`</li><li>`fusectl`</li><li>`pstore`</li><li>`debugfs`</li><li>`securityfs`</li><li>`devpts`</li><li>`mqueue`</li><li>`hugetlbfs`</li><li>`nsfs`</li><li>`rpc_pipefs`</li><li>`devtmpfs`</li><li>`none`</li><li>`tmpfs`</li><li>`aufs`</li></ul> | List of excluded filesystem types |
| **keep_original_mountpoint** | bool      | `true` | Whether original mount point names should be retained |

## Documentation

### Collected Metrics

List of collected metrics is described in [METRICS.md](METRICS.md).


### Examples

Example running snap-plugin-collector-df plugin and writing data to a file.

Ensure [snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `service snap-telemetry start`
* systemd: `sysctl start snap-telemetry`
* command line: `snapteld -l 1 -t 0 &`

Download and load snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-df/latest/linux/x86_64/snap-plugin-collector-df
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ chmod 755 snap-plugin-*
$ snaptel plugin load snap-plugin-collector-df
$ snaptel plugin load snap-plugin-publisher-file
```

See all available metrics:
```
$ snaptel metric list
```

Download an [example task file](examples/tasks/df-file.json) and load it:
```
$ curl -sfLO https://raw.githubusercontent.com/intelsdi-x/snap-plugin-collector-df/master/examples/tasks/df-file.json
$ snaptel task create -t df-file.json
Using task manifest to create task
Task created
ID: 480323af-15b0-4af8-a526-eb2ca6d8ae67
Name: Task-480323af-15b0-4af8-a526-eb2ca6d8ae67
State: Running
```

See realtime output from `snaptel task watch <task_id>` (CTRL+C to exit)
```
$ snaptel task watch 480323af-15b0-4af8-a526-eb2ca6d8ae67
```

This data is published to a file `/tmp/published_df.log` per task specification

Stop task:
```
$ snaptel task stop 480323af-15b0-4af8-a526-eb2ca6d8ae67
Task stopped:
ID: 480323af-15b0-4af8-a526-eb2ca6d8ae67
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-df/issues) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-df/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [Slack](http://slack.snap-telemetry.io).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[Snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Patryk Matyjasek](https://github.com/PatrykMatyjasek)
* Author: [Marcin Krolik](https://github.com/marcin-krolik)
