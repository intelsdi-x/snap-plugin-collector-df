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
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
)

type DfPluginSuite struct {
	suite.Suite
	cfg           plugin.ConfigType
	mockCollector *MockCollector
}

func (dfp *DfPluginSuite) SetupSuite() {
	mc := &MockCollector{}

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
	mc.On("collect", "/proc").Return(dfms, nil)
	dfp.mockCollector = mc
	dfp.cfg = plugin.ConfigType{}
}

func (dfp *DfPluginSuite) TestGetMetricTypes() {
	Convey("Given df plugin is initialized", dfp.T(), func() {
		//dfPlg := NewDfCollector()
		dfPlg := dfCollector{
			stats: dfp.mockCollector,
		}

		Convey("When list of available metrics is requested", func() {
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
	})
}

func (dfp *DfPluginSuite) TestCollectMetrics() {
	Convey("Given df plugin is initialized", dfp.T(), func() {
		//dfPlg := NewDfCollector()
		dfPlg := dfCollector{
			stats: dfp.mockCollector,
		}

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
				So(len(mts), ShouldEqual, 28)
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/space_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/space_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/space_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/space_percent_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/space_percent_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/space_percent_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/inodes_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/inodes_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/inodes_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/inodes_percent_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/inodes_percent_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/inodes_percent_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/device_name")
				So(ns, ShouldContain, "/intel/procfs/filesystem/rootfs/device_type")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/space_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/space_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/space_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/space_percent_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/space_percent_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/space_percent_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/inodes_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/inodes_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/inodes_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/inodes_percent_free")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/inodes_percent_reserved")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/inodes_percent_used")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/device_name")
				So(ns, ShouldContain, "/intel/procfs/filesystem/big/device_type")
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

func (mc *MockCollector) collect(p string) ([]dfMetric, error) {
	ret := mc.Mock.Called("/proc")
	return ret.Get(0).([]dfMetric), ret.Error(1)
}
