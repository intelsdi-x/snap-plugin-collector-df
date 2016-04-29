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
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
)

const (
	// PluginName df collector plugin name
	PluginName = "df"
	// Version of plugin
	Version = 2
	// Type of plugin
	Type = plugin.CollectorPluginType
)

var (
	optionsKB       = []string{"--no-sync", "-P", "-T"}
	optionsINode    = []string{"--no-sync", "-P", "-T", "-i"}
	namespacePrefix = []string{"intel", "procfs", "filesystem"}
	metricsKind     = []string{
		"space_free",
		"space_reserved",
		"space_used",
		"space_percent_free",
		"space_percent_reserved",
		"space_percent_used",
		"inodes_free",
		"inodes_reserved",
		"inodes_used",
		"inodes_percent_free",
		"inodes_percent_reserved",
		"inodes_percent_used",
		"device_name",
		"device_type",
	}
)

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (p *dfCollector) GetMetricTypes(_ plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	dfms, err := p.stats.collect()

	if err != nil {
		return mts, fmt.Errorf(fmt.Sprintf("Unable to get available metrics from df: %s", err))
	}

	for _, dfm := range dfms {
		for _, kind := range metricsKind {
			mt := plugin.MetricType{Namespace_: core.NewNamespace(makeNamespace(dfm, kind)...)}
			mts = append(mts, mt)
		}
	}

	return mts, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (p *dfCollector) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	metrics := []plugin.MetricType{}
	dfms, err := p.stats.collect()
	if err != nil {
		return metrics, fmt.Errorf(fmt.Sprintf("Unable to collect metrics from df: %s", err))
	}

	hostname, _ := os.Hostname()
	for _, mt := range mts {
		tags := mt.Tags()
		if tags == nil {
			tags = map[string]string{}
		}
		tags["hostname"] = hostname

		namespace := mt.Namespace().Strings()
		if len(namespace) < 5 {
			return nil, fmt.Errorf("Wrong namespace length %d", len(namespace))
		}

		for _, dfm := range dfms {

			if validateMetric(namespace[3:], dfm) {

				kind := namespace[4]
				metric := plugin.MetricType{
					Timestamp_: time.Now(),
					Tags_:      tags,
					Namespace_: mt.Namespace(),
				}
				switch kind {
				case "space_free":
					metric.Data_ = dfm.Available
				case "space_reserved":
					metric.Data_ = dfm.Blocks - (dfm.Used + dfm.Available)
				case "space_used":
					metric.Data_ = dfm.Used
				case "space_percent_free":
					metric.Data_ = 100 * float64(dfm.Available) / float64(dfm.Blocks)
				case "space_percent_reserved":
					metric.Data_ = 100 * float64(dfm.Blocks-(dfm.Used+dfm.Available)) / float64(dfm.Blocks)
				case "space_percent_used":
					metric.Data_ = 100 * float64(dfm.Used) / float64(dfm.Blocks)
				case "device_name":
					metric.Data_ = dfm.Filesystem
				case "device_type":
					metric.Data_ = dfm.FsType
				case "inodes_free":
					metric.Data_ = dfm.IFree
				case "inodes_reserved":
					metric.Data_ = dfm.Inodes - (dfm.IUsed + dfm.IFree)
				case "inodes_used":
					metric.Data_ = dfm.IUsed
				case "inodes_percent_free":
					metric.Data_ = 100 * float64(dfm.IFree) / float64(dfm.Inodes)
				case "inodes_percent_reserved":
					metric.Data_ = 100 * float64(dfm.Inodes-(dfm.IUsed+dfm.IFree)) / float64(dfm.Inodes)
				case "inodes_percent_used":
					metric.Data_ = 100 * float64(dfm.IUsed) / float64(dfm.Inodes)
				}
				metrics = append(metrics, metric)
			}
		}
	}
	return metrics, nil
}

// GetConfigPolicy returns config policy
// It returns error in case retrieval was not successful
func (p *dfCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

// NewDfCollector creates new instance of plugin and returns pointer to initialized object.
func NewDfCollector() *dfCollector {
	return &dfCollector{stats: &dfStats{}}
}

// Meta returns plugin's metadata
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		PluginName,
		Version,
		Type,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
	)
}

type dfCollector struct {
	stats collector
}

type dfMetric struct {
	Filesystem              string
	Used, Available, Blocks uint64
	Capacity                float64
	FsType                  string
	MountPoint              string
	Inodes, IUsed, IFree    uint64
	IUse                    float64
}

type collector interface {
	collect() ([]dfMetric, error)
}

type dfStats struct{}

func (dfs *dfStats) collect() ([]dfMetric, error) {
	dfms := []dfMetric{}
	kBOutputB, err := exec.Command("df", optionsKB...).Output()
	if err != nil {
		return dfms, err
	}
	InodeOutput, err := exec.Command("df", optionsINode...).Output()
	if err != nil {
		return dfms, err
	}
	stringkB := string(kBOutputB)
	stringInode := string(InodeOutput)
	lineskB := strings.Split(stringkB, "\n")
	linesInode := strings.Split(stringInode, "\n")
	if len(linesInode) != len(lineskB) {
		return nil, fmt.Errorf("Inodes stats not comparable to space stats!")
	}

	for i := 0; i < len(linesInode); i++ {
		inodeEntry := strings.Fields(linesInode[i])
		spaceEntry := strings.Fields(lineskB[i])
		if len(spaceEntry) < 8 && len(spaceEntry) > 0 { //check if not line with columns description or not empty
			var dfm dfMetric
			dfm.Filesystem = spaceEntry[0]
			dfm.FsType = spaceEntry[1]
			dfm.Blocks, _ = strconv.ParseUint(spaceEntry[2], 10, 64)
			dfm.Used, _ = strconv.ParseUint(spaceEntry[3], 10, 64)
			dfm.Available, _ = strconv.ParseUint(spaceEntry[4], 10, 64)
			dfm.Capacity, _ = parsePerc(spaceEntry[5])
			dfm.Inodes, _ = strconv.ParseUint(inodeEntry[2], 10, 64)
			dfm.IUsed, _ = strconv.ParseUint(inodeEntry[3], 10, 64)
			dfm.IFree, _ = strconv.ParseUint(inodeEntry[4], 10, 64)
			dfm.IUse, _ = parsePerc(inodeEntry[5])
			if spaceEntry[6] == "/" {
				dfm.MountPoint = "rootfs"
			} else {
				dfm.MountPoint = strings.Replace(spaceEntry[6][1:], "/", "_", -1)
			}
			dfms = append(dfms, dfm)
		}

	}
	return dfms, nil
}

func makeNamespace(dfm dfMetric, kind string) []string {
	ns := []string{}

	ns = append(ns, namespacePrefix...)
	ns = append(ns, dfm.MountPoint, kind)

	return ns
}

// validate if metric should be exposed
func validateMetric(namespace []string, dfm dfMetric) bool {
	mountPoint := namespace[0]
	if mountPoint == dfm.MountPoint {
		return true
	}

	return false
}

func parsePerc(s string) (float64, error) {
	length := len(s)
	if string(s[length-1]) != "%" {
		return 0, fmt.Errorf("Wrong format")
	}
	ret, err := strconv.ParseFloat(s[:length-1], 10)
	if err != nil {
		return 0, err
	}
	return ret / 100, nil
}
