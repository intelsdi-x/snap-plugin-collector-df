# snap plugin collector - df

## Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description
----------|-----------|-----------------------
/intel/procfs/filesystem/\<mount_point\>/inodes_free | uint64 | the number of free inodes on the file system
/intel/procfs/filesystem/\<mount_point\>/inodes_reserved | uint64 | the number of reserved inodes
/intel/procfs/filesystem/\<mount_point\>/inodes_used | uint64 | the number of used inodes
/intel/procfs/filesystem/\<mount_point\>/space_free | uint64 | the number of free bytes
/intel/procfs/filesystem/\<mount_point\>/space_reserved | uint64 | the number of reserved bytes
/intel/procfs/filesystem/\<mount_point\>/space_used | uint64 | the number of used bytes
/intel/procfs/filesystem/\<mount_point\>/inodes_percent_free | float64 | the percentage of free inodes on the file system
/intel/procfs/filesystem/\<mount_point\>/inodes_percent_reserved | float64 | the percentage of reserved inodes
/intel/procfs/filesystem/\<mount_point\>/inodes_percent_used | float64 | the percentage of used inodes
/intel/procfs/filesystem/\<mount_point\>/space_percent_free | float64 | the percentage of free bytes
/intel/procfs/filesystem/\<mount_point\>/space_percent_reserved | float64 | the percentage of reserved bytes
/intel/procfs/filesystem/\<mount_point\>/space_percent_used | float64 | the percentage of used bytes
/intel/procfs/filesystem/\<mount_point\>/device_name | string | device name as presented in filesystem (eg. /dev/sda1)
/intel/procfs/filesystem/\<mount_point\>/device_type | string | device type as presented in filesystem (eg. ext4)
