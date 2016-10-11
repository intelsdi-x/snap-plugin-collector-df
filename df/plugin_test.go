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
	"errors"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
)

type DfPluginSuite struct {
	suite.Suite
	cfg           plugin.ConfigType
	mockCollector *MockCollector
}

func (dfp *DfPluginSuite) SetupSuite() {
	dfms := []dfMetric{
		dfMetric{
			Blocks:     100,
			Used:       50,
			Available:  40,
			Capacity:   0.5,
			FsType:     "ext4",
			Filesystem: "/dev/sda1",
			MountPoint: "rootfs",
			Inodes:     1000,
			IUsed:      500,
			IFree:      400,
			IUse:       0.5,
		},
		dfMetric{
			Blocks:     200,
			Used:       110,
			Available:  80,
			Capacity:   0.3,
			FsType:     "ext4",
			Filesystem: "/dev/sda2",
			MountPoint: "big",
			Inodes:     2000,
			IUsed:      1000,
			IFree:      800,
			IUse:       0.5,
		},
	}
	mc := &MockCollector{}
	mc.On("collect", "/proc", dfltExcludedFSNames, dfltExcludedFSTypes).Return(dfms, nil)
	mc.On("collect", "/dummy", dfltExcludedFSNames, dfltExcludedFSTypes).Return(dfms, errors.New("Fake error"))
	dfp.mockCollector = mc
	dfp.cfg = plugin.ConfigType{}
}

func (dfp *DfPluginSuite) TestGetMetricTypes() {
	Convey("Given df plugin is badly initialized", dfp.T(), func() {
		dfPlg := NewDfCollector()
		dfPlg.stats = dfp.mockCollector

		node := cdata.NewNode()
		node.AddItem(ProcPath, ctypes.ConfigValueStr{Value: "/dummy"})

		Convey("When list of available metrics is requested", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace("intel", "procfs", "filesystem", "rootfs", "space_free"),
					Config_:    node,
				},
			}
			metrics, err := dfPlg.CollectMetrics(mts)
			Convey("Then error should be reported", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "no such file or directory")
				So(metrics, ShouldBeNil)
			})
		})
	})
	Convey("Given df plugin is initialized", dfp.T(), func() {
		dfPlg := NewDfCollector()
		dfPlg.stats = dfp.mockCollector

		Convey("When values for given metrics are requested", func() {

			mts, err := dfPlg.GetMetricTypes(dfp.cfg)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then proper metrics are returned", func() {
				ns := []string{}
				for _, m := range mts {
					ns = append(ns, m.Namespace().String())
				}
				So(len(mts), ShouldEqual, 14)
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/space_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/space_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/space_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/space_percent_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/space_percent_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/space_percent_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/inodes_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/inodes_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/inodes_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/inodes_percent_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/inodes_percent_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/inodes_percent_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/device_name")
				So(ns, ShouldContain, "/intel/procfs/filesystem/*/device_type")
			})
		})
	})
}

func (dfp *DfPluginSuite) TestCollectMetrics() {
	Convey("Given df plugin is initialized", dfp.T(), func() {
		dfPlg := NewDfCollector()
		dfPlg.stats = dfp.mockCollector

		Convey("When list of metrics is requested with too short namespace", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace("filesystem", "rootfs", "space_free"),
				},
			}
			metrics, err := dfPlg.CollectMetrics(mts)

			Convey("Then error should be reported", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "Wrong namespace length")
				So(metrics, ShouldBeNil)
			})
		})

		Convey("When list of metrics is requested with bad namespace", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace("intel", "procfs", "filesystem", "rootfs"),
				},
			}
			metrics, err := dfPlg.CollectMetrics(mts)

			Convey("Then error should be reported", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "Namespace should contain wildcard")
				So(metrics, ShouldBeNil)
			})
		})

		Convey("When list of specific metrics is requested", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace("intel", "procfs", "filesystem", "rootfs", "space_free"),
				},
				plugin.MetricType{
					Namespace_: core.NewNamespace("intel", "procfs", "filesystem", "big", "inodes_percent_free"),
				},
			}
			metrics, err := dfPlg.CollectMetrics(mts)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeNil)
			})

			Convey("Then proper metrics are returned", func() {
				metvals := map[string]interface{}{}
				for _, m := range metrics {
					stat := strings.Join(m.Namespace().Strings()[3:], "/")
					metvals[stat] = m.Data()
				}
				So(len(metrics), ShouldEqual, 2)

				val, ok := metvals["rootfs/space_free"]
				So(ok, ShouldBeTrue)
				So(val, ShouldNotBeNil)

				val, ok = metvals["big/inodes_percent_free"]
				So(ok, ShouldBeTrue)
				So(val, ShouldNotBeNil)
			})
		})

		Convey("When one single available dynamic metrics is requested", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace("intel", "procfs", "filesystem", "*", "space_free"),
				},
			}
			metrics, err := dfPlg.CollectMetrics(mts)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeNil)
			})

			Convey("Then proper metrics are returned", func() {
				metvals := map[string]interface{}{}
				for _, m := range metrics {
					stat := strings.Join(m.Namespace().Strings()[3:], "/")
					metvals[stat] = m.Data()
				}

				So(len(metrics), ShouldEqual, 2)

				val, ok := metvals["rootfs/space_free"]
				So(ok, ShouldBeTrue)
				So(val, ShouldNotBeNil)

				val, ok = metvals["big/space_free"]
				So(ok, ShouldBeTrue)
				So(val, ShouldNotBeNil)
			})
		})

		Convey("When all available dynamic metrics are requested for given mountpoint", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace("intel", "procfs", "filesystem", "rootfs", "*"),
				},
			}
			metrics, err := dfPlg.CollectMetrics(mts)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeNil)
			})

			Convey("Then proper metrics are returned", func() {
				metvals := map[string]interface{}{}
				for _, m := range metrics {
					stat := strings.Join(m.Namespace().Strings()[3:], "/")
					So(stat, ShouldStartWith, "rootfs")
					metvals[stat] = m.Data()
				}
				So(len(metrics), ShouldEqual, 14)

				val, ok := metvals["rootfs/space_free"]
				So(ok, ShouldBeTrue)
				So(val, ShouldNotBeNil)
			})
		})

		Convey("When all available dynamic metrics are requested for all mountpoints", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace("intel", "procfs", "filesystem", "*", "*"),
				},
			}
			metrics, err := dfPlg.CollectMetrics(mts)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeNil)
			})

			Convey("Then proper metrics are returned", func() {
				metvals := map[string]interface{}{}
				for _, m := range metrics {
					stat := strings.Join(m.Namespace().Strings()[3:], "/")
					metvals[stat] = m.Data()
				}

				So(len(metrics), ShouldEqual, 28)

				val, ok := metvals["rootfs/space_free"]
				So(ok, ShouldBeTrue)
				So(val, ShouldNotBeNil)

				val, ok = metvals["big/space_free"]
				So(ok, ShouldBeTrue)
				So(val, ShouldNotBeNil)
			})
		})

		Convey("When list of all available dynamic metrics is requested", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace("intel", "procfs", "filesystem", "*"),
				},
			}
			metrics, err := dfPlg.CollectMetrics(mts)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeNil)
			})

			Convey("Then proper metrics are returned", func() {
				metvals := map[string]interface{}{}
				for _, m := range metrics {
					stat := strings.Join(m.Namespace().Strings()[3:], "/")
					metvals[stat] = m.Data()
				}

				So(len(metrics), ShouldEqual, 28)

				val, ok := metvals["rootfs/space_free"]
				So(ok, ShouldBeTrue)
				So(val, ShouldNotBeNil)

				val, ok = metvals["big/space_free"]
				So(ok, ShouldBeTrue)
				So(val, ShouldNotBeNil)
			})
		})

		Convey("When calling twice", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace("intel", "procfs", "filesystem", "rootfs", "space_free"),
				},
			}
			metrics, err := dfPlg.CollectMetrics(mts)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeNil)
			})

			dfPlg.proc_path = "/dummy"
			_, err = dfPlg.CollectMetrics(mts)

			Convey("Then error should be reported", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "Fake error")
			})
		})
	})
}

func (dfp *DfPluginSuite) TestCollect() {
	Convey("Given df plugin is initialized", dfp.T(), func() {
		dfPlg := NewDfCollector()

		Convey("When called with non existing path", func() {
			metrics, err := dfPlg.stats.collect("/dummy", []string{}, []string{})
			Convey("Then error should be reported", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "no such file or directory")
				So(metrics, ShouldBeNil)
			})
		})

		Convey("When called with existing path and different exclusion lists", func() {
			metrics, err := dfPlg.stats.collect("/proc", []string{"dummy"}, []string{"dummy"})
			Convey("Then no error should be reported with dummy exclusion lists", func() {
				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeNil)
			})
			nbMetrics := len(metrics)
			exclusions := false
			for _, m := range metrics {
				if excludedFSFromList(m.UnchangedMountPoint, dfltExcludedFSNames) ||
					excludedFSFromList(m.FsType, dfltExcludedFSTypes) {
					exclusions = true
				}
			}
			Convey("Then no exclusion occured", func() {
				So(exclusions, ShouldEqual, true)
			})

			metrics, err = dfPlg.stats.collect("/proc", dfltExcludedFSNames, dfltExcludedFSTypes)
			Convey("Then error should be reported", func() {
				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeNil)
				So(len(metrics), ShouldBeLessThan, nbMetrics)
			})
		})
	})
}

func (dfp *DfPluginSuite) TestHelperRoutines() {
	Convey("Basically initialized", dfp.T(), func() {

		Convey("Namespace", func() {

			ns := createNamespace("element", "name")

			Convey("Then no error should be reported", func() {
				So(ns, ShouldNotBeNil)
				So(strings.Join(ns, ","), ShouldStartWith,
					strings.Join(namespacePrefix, ","))
				So(strings.Join(ns, ","), ShouldEndWith,
					"element,name")
			})

			dfm := dfMetric{
				MountPoint: "/mount",
			}
			ns = makeNamespace(dfm, "mykind")

			Convey("Then namespace shoud be reported", func() {
				So(ns, ShouldNotBeNil)
				So(len(ns), ShouldNotEqual, 0)
				So(ns[len(ns)-1], ShouldEqual, "mykind")
				So(ns[len(ns)-2], ShouldEqual, "/mount")
			})
		})

		Convey("Set config variables", func() {

			node := cdata.NewNode()
			node.AddItem(ProcPath, ctypes.ConfigValueStr{Value: "/etc/hosts"})
			cfg := plugin.ConfigType{ConfigDataNode: node}
			dfPlg := NewDfCollector()
			dfPlg.stats = dfp.mockCollector
			err := dfPlg.setProcPath(cfg)

			Convey("Then error should be reported (not a directory)", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "is not a directory")
			})

			node = cdata.NewNode()
			node.AddItem(ProcPath, ctypes.ConfigValueStr{Value: "/proc"})
			cfg = plugin.ConfigType{ConfigDataNode: node}
			dfPlg = NewDfCollector()
			dfPlg.stats = dfp.mockCollector
			err = dfPlg.setProcPath(cfg)

			Convey("Then no error should be reported (proc_path)", func() {
				So(err, ShouldBeNil)
			})

			err = dfPlg.setProcPath(cfg)

			Convey("Then no error should be reported (already called)", func() {
				So(err, ShouldBeNil)
				So(len(dfPlg.excluded_fs_types), ShouldEqual, 18)
				So(len(dfPlg.excluded_fs_names), ShouldEqual, 2)
			})

			node = cdata.NewNode()
			node.AddItem(ExcludedFSNames, ctypes.ConfigValueStr{Value: ""})
			dfPlg = NewDfCollector()
			dfPlg.stats = dfp.mockCollector
			cfg = plugin.ConfigType{ConfigDataNode: node}
			err = dfPlg.setProcPath(cfg)

			Convey("Then no error should be reported (excluded_fs_names with proper value for empty string)", func() {
				So(err, ShouldBeNil)
				So(len(dfPlg.excluded_fs_names), ShouldEqual, 0)
			})

			node = cdata.NewNode()
			node.AddItem(ExcludedFSNames, ctypes.ConfigValueStr{Value: "n1,n2,n3"})
			dfPlg = NewDfCollector()
			dfPlg.stats = dfp.mockCollector
			cfg = plugin.ConfigType{ConfigDataNode: node}
			err = dfPlg.setProcPath(cfg)

			Convey("Then no error should be reported (excluded_fs_names with proper value)", func() {
				So(err, ShouldBeNil)
				So(len(dfPlg.excluded_fs_names), ShouldEqual, 3)
			})

			node = cdata.NewNode()
			node.AddItem(ExcludedFSTypes, ctypes.ConfigValueStr{Value: ""})
			dfPlg = NewDfCollector()
			dfPlg.stats = dfp.mockCollector
			cfg = plugin.ConfigType{ConfigDataNode: node}
			err = dfPlg.setProcPath(cfg)

			Convey("Then no error should be reported (excluded_fs_types with proper value for empty string)", func() {
				So(err, ShouldBeNil)
				So(len(dfPlg.excluded_fs_types), ShouldEqual, 0)
			})

			node = cdata.NewNode()
			node.AddItem(ExcludedFSTypes, ctypes.ConfigValueStr{Value: "fs1,fs2"})
			dfPlg = NewDfCollector()
			dfPlg.stats = dfp.mockCollector
			cfg = plugin.ConfigType{ConfigDataNode: node}
			err = dfPlg.setProcPath(cfg)

			Convey("Then no error should be reported (excluded_fs_types with proper value)", func() {
				So(err, ShouldBeNil)
				So(len(dfPlg.excluded_fs_types), ShouldEqual, 2)
			})
		})

		Convey("Set get config policy", func() {

			dfPlg := NewDfCollector()
			dfPlg.stats = dfp.mockCollector
			cp, err := dfPlg.GetConfigPolicy()

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
				So(cp, ShouldNotBeNil)
			})
		})

		Convey("Miscellaneous", func() {

			dfm := dfMetric{
				MountPoint: "/mount",
			}
			v := validateMetric([]string{"/mount"}, dfm)

			Convey("Then true value should be reported", func() {
				So(v, ShouldEqual, true)
			})

			v = validateMetric([]string{"/test"}, dfm)

			Convey("Then false value should be reported", func() {
				So(v, ShouldEqual, false)
			})

			l := []string{"1", "2", "3", "4"}
			v = excludedFSFromList("3", l)

			Convey("Then value should be reported as found", func() {
				So(v, ShouldEqual, true)
			})

			v = excludedFSFromList("error", l)

			Convey("Then value should be reported as not found", func() {
				So(v, ShouldEqual, false)
			})

			m := Meta()

			Convey("Then meta value should be reported as not nil", func() {
				So(m, ShouldNotBeNil)
			})
		})

		Convey("ceilPercent", func() {

			v := ceilPercent(1, 0)

			Convey("Then 0.0 value should be reported", func() {
				So(v, ShouldEqual, 0.0)
			})

			v = ceilPercent(80, 111)

			Convey("Then x value should be reported", func() {
				So(v, ShouldEqual, 73.0)
			})
		})
	})
}

func TestDfPluginSuite(t *testing.T) {
	suite.Run(t, &DfPluginSuite{})
}

type MockCollector struct {
	mock.Mock
}

func (mc *MockCollector) collect(p string, n []string, f []string) ([]dfMetric, error) {
	ret := mc.Mock.Called(p, n, f)
	return ret.Get(0).([]dfMetric), ret.Error(1)
}
