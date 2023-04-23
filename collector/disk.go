package collector

import (
	"github.com/go-kit/log/level"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("disk", defaultEnabled, NewVirtDiskCollector)
}

const (
	virtdisk = "disk"
)

var (
	virtDiskCapacity = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, virtdisk, "capacity"),
		"virtual disk capacity.",
		[]string{"uuid", "diskname"}, nil,
	)

	virtDiskAllocation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, virtdisk, "allocation"),
		"virtual disk allocation.",
		[]string{"uuid", "diskname"}, nil,
	)

	virtDiskPhysical = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, virtdisk, "physical"),
		"virtual disk physical.",
		[]string{"uuid", "diskname"}, nil,
	)
)

type VirtDiskCollector struct {
	logger             log.Logger
	VirtDiskCapacity   *prometheus.Desc
	VirtDiskAllocation *prometheus.Desc
	VirtDiskPhysical   *prometheus.Desc
}

type VirtDiskInfo struct {
	Name       string
	Capacity   float64
	Allocation float64
	Physical   float64
}

func NewVirtDiskCollector(log log.Logger) (Collector, error) {
	return &VirtDiskCollector{
		logger:             log,
		VirtDiskCapacity:   virtDiskCapacity,
		VirtDiskAllocation: virtDiskAllocation,
		VirtDiskPhysical:   virtDiskPhysical,
	}, nil
}

func (v *VirtDiskCollector) Update(ch chan<- prometheus.Metric) error {
	return v.UpdateInfo(ch)
}

func (v *VirtDiskCollector) UpdateInfo(ch chan<- prometheus.Metric) error {
	begin := time.Now()
	var (
		uuids []string
	)
	execCmd := []string{"list", "--uuid"}
	output, stdErr, err := ExecCommand("virsh", execCmd)
	duration := time.Since(begin)
	if err != nil {
		level.Error(v.logger).Log("msg", stdErr, "name", execCmd, "duration_seconds", duration.Seconds(), "err", err)
		return err
	} else {
		uuids = strings.Split(output, "\n")
	}

	for _, uuid := range uuids {
		if uuid == "" {
			continue
		}
		begin = time.Now()
		disksInfo, stdErr, err := getVirtDisksInfo(uuid)
		duration = time.Since(begin)
		if err != nil {
			level.Error(v.logger).Log("msg", uuid+":"+stdErr, "name", execCmd, "duration_seconds", duration.Seconds(), "err", err)
		}
		for _, disk := range disksInfo {
			ch <- prometheus.MustNewConstMetric(v.VirtDiskCapacity, prometheus.GaugeValue, disk.Capacity, uuid, disk.Name)
			ch <- prometheus.MustNewConstMetric(v.VirtDiskAllocation, prometheus.GaugeValue, disk.Allocation, uuid, disk.Name)
			ch <- prometheus.MustNewConstMetric(v.VirtDiskPhysical, prometheus.GaugeValue, disk.Physical, uuid, disk.Name)
		}
	}
	return nil
}

func getVirtDisksInfo(uuid string) ([]VirtDiskInfo, string, error) {
	execCmd := []string{"domblkinfo", uuid, "--all"}
	output, stdErr, err := ExecCommand("virsh", execCmd)
	if err != nil {
		return nil, stdErr, err
	}
	lines := strings.Split(output, "\n")
	lines = lines[2:]
	var disksInfo []VirtDiskInfo
	for _, line := range lines {
		infoLine := strings.Fields(line)
		if len(infoLine) != 4 {
			continue
		}

		capa, err := strconv.ParseFloat(strings.TrimSpace(infoLine[1]), 64)
		if err != nil {
			continue
		}

		alloc, err := strconv.ParseFloat(strings.TrimSpace(infoLine[2]), 64)
		if err != nil {
			continue
		}

		phy, err := strconv.ParseFloat(strings.TrimSpace(infoLine[3]), 64)
		if err != nil {
			continue
		}

		disksInfo = append(disksInfo, VirtDiskInfo{
			Name:       infoLine[0],
			Capacity:   capa,
			Allocation: alloc,
			Physical:   phy,
		})
	}

	return disksInfo, stdErr, nil
}
