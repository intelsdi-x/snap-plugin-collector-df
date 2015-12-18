/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2015 Intel Corporation
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package df

import (
	"fmt"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	// Name of plugin
	Name = "df"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

type DfCollector struct {
}

type Metric struct {
	Filesystem      string
	Used, Available uint64
	Percentage      float64
	FsType          string
	MountPoint      string
	Inode           bool
}

var optionsKB = []string{"--no-sync", "-P", "-T"}
var optionsINode = []string{"--no-sync", "-P", "-T", "-i"}
var namespacePrefix = []string{"intel", "disk"}
var metricsKind = []string{"fs_type", "mount_point", "available_space", "used_space", "percentage_space"}

func makeNamespace(metrics Metric, kind string) []string {
	ns := []string{}
	ns = append(ns, namespacePrefix...)
	if strings.Contains(metrics.Filesystem, "/") {
		ns = append(ns, strings.Split(metrics.Filesystem, "/")[1:]...) // drop first element from array
	} else {
		ns = append(ns, metrics.Filesystem)
	}
	if kind == "mount_point" || kind == "fs_type" {
		ns = append(ns, kind)
	} else {
		ns = append(ns, "space")
		metric := ""
		metric += kind
		if metrics.Inode {
			metric += "_inodes"
		} else {
			metric += "_kB"
		}
		ns = append(ns, metric)
	}
	return ns
}

func collect() ([]Metric, error) {
	metrics := make([]Metric, 0)
	kBOutputB, err := exec.Command("df", optionsKB...).Output()
	if err != nil {
		return metrics, err
	}
	InodeOutput, err := exec.Command("df", optionsINode...).Output()
	if err != nil {
		return metrics, err
	}
	stringkB := string(kBOutputB)
	stringInode := string(InodeOutput)
	lineskB := strings.Split(stringkB, "\n")
	linesInode := strings.Split(stringInode, "\n")
	data := [][]string{lineskB, linesInode}
	for _, dat := range data {
		for _, line := range dat {
			columns := strings.Fields(line)
			if len(columns) < 8 && len(columns) > 0 { //check if not line with columns description or not empty
				if columns[0] != "none" { //check if no none device in df
					var metric Metric
					if strings.Contains(dat[0], "IUsed") { //check header of command output
						metric.Inode = true
					} else {
						metric.Inode = false
					}
					// fill struct fields
					metric.Filesystem = columns[0]
					metric.FsType = columns[1]
					metric.Used, _ = strconv.ParseUint(columns[3], 10, 64)
					metric.Available, _ = strconv.ParseUint(columns[4], 10, 64)
					metric.Percentage, _ = strconv.ParseFloat(columns[5], 10)
					metric.MountPoint = columns[6]
					metrics = append(metrics, metric)
				}
			}
		}
	}
	return metrics, nil
}

// validate if metric should be exposed
func validateMetric(namespace []string, metrics []plugin.PluginMetricType) bool {
	for _, metric := range metrics {
		if strings.Join(namespace, "/") == strings.Join(metric.Namespace_, "/") { //check if namespace is in mts
			return true
		}
	}
	return false
}

func (p *DfCollector) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	metrics := []plugin.PluginMetricType{}
	data, err := collect()
	if err != nil {
		return metrics, fmt.Errorf(fmt.Sprintf("Unable to collect metrics from df: %s", err))
	}
	timestamp := time.Now()
	hostname, _ := os.Hostname()
	metric := plugin.PluginMetricType{}
	for _, record := range data { //data is array of structs which contains all metrics per line
		metric = plugin.PluginMetricType{
			Namespace_: makeNamespace(record, "available_space"),
			Data_:      record.Available,
			Timestamp_: timestamp,
			Source_:    hostname,
		}
		if validateMetric(metric.Namespace_, mts) {
			metrics = append(metrics, metric)
		}
		metric = plugin.PluginMetricType{
			Namespace_: makeNamespace(record, "used_space"),
			Data_:      record.Used,
			Timestamp_: timestamp,
			Source_:    hostname,
		}
		if validateMetric(metric.Namespace_, mts) {
			metrics = append(metrics, metric)
		}
		metric = plugin.PluginMetricType{
			Namespace_: makeNamespace(record, "percentage_space"),
			Data_:      record.Percentage,
			Timestamp_: timestamp,
			Source_:    hostname,
		}
		if validateMetric(metric.Namespace_, mts) {
			metrics = append(metrics, metric)
		}
		metric = plugin.PluginMetricType{
			Namespace_: makeNamespace(record, "mount_point"),
			Data_:      record.MountPoint,
			Timestamp_: timestamp,
			Source_:    hostname,
		}
		if validateMetric(metric.Namespace_, mts) && !validateMetric(metric.Namespace_, metrics) { //validate if metric is not doubled
			metrics = append(metrics, metric)
		}
		metric = plugin.PluginMetricType{
			Namespace_: makeNamespace(record, "fs_type"),
			Data_:      record.FsType,
			Timestamp_: timestamp,
			Source_:    hostname,
		}
		if validateMetric(metric.Namespace_, mts) && !validateMetric(metric.Namespace_, metrics) { //validate if metric is not doubled
			metrics = append(metrics, metric)
		}
	}
	return metrics, nil
}

func (p *DfCollector) GetMetricTypes(_ plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	mts := []plugin.PluginMetricType{}
	data, err := collect()
	if err != nil {
		return mts, fmt.Errorf(fmt.Sprintf("Unable to get available metrics from df: %s", err))
	}
	for _, c := range data {
		mt := plugin.PluginMetricType{Namespace_: makeNamespace(c, "available_space")}
		mts = append(mts, mt)
		mt = plugin.PluginMetricType{Namespace_: makeNamespace(c, "used_space")}
		mts = append(mts, mt)
		mt = plugin.PluginMetricType{Namespace_: makeNamespace(c, "percentage_space")}
		mts = append(mts, mt)
		mt = plugin.PluginMetricType{Namespace_: makeNamespace(c, "mount_point")}
		mts = append(mts, mt)
		mt = plugin.PluginMetricType{Namespace_: makeNamespace(c, "fs_type")}
		mts = append(mts, mt)
	}
	return mts, nil
}

// GetConfigPolicy
func (p *DfCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

// Creates new instance of plugin and returns pointer to initialized object.
func NewDfCollector() *DfCollector {
	return &DfCollector{}
}

// Returns plugin's metadata
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(Name, Version, Type, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}
