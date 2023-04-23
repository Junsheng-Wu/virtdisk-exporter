package collector

import (
	"encoding/json"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"time"
)

func init() {
	registerCollector("filesystem", defaultEnabled, NewVirtFilesystemCollector)
}

const (
	virtfilesystem = "fs"
)

var (
	virtFilesystemTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, virtfilesystem, "total"),
		"virtual machine filesystem total",
		[]string{"uuid", "disk", "mountpoint"}, nil,
	)

	virtFilesystemUsed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, virtfilesystem, "used"),
		"virtual machine filesystem used",
		[]string{"uuid", "disk", "mountpoint"}, nil,
	)
)

type Response struct {
	Return []DiskFsInfo `json:"return,omitempty"`
}

type DiskFsInfo struct {
	Name       string  `json:"name,omitempty"`
	TotalBytes float64 `json:"total-bytes,omitempty"`
	MountPoint string  `json:"mountPoint,omitempty"`
	UsedBytes  float64 `json:"used-bytes,omitempty"`
	Type       string  `json:"type,omitempty"`
}

type VirtFilesystemCollector struct {
	logger              log.Logger
	VirtFilesystemTotal *prometheus.Desc
	VirtFilesystemUsed  *prometheus.Desc
}

func NewVirtFilesystemCollector(log log.Logger) (Collector, error) {
	return &VirtFilesystemCollector{
		logger:              log,
		VirtFilesystemTotal: virtFilesystemTotal,
		VirtFilesystemUsed:  virtFilesystemUsed,
	}, nil
}

func (v *VirtFilesystemCollector) Update(ch chan<- prometheus.Metric) error {
	return v.UpdateInfo(ch)
}

func (v *VirtFilesystemCollector) UpdateInfo(ch chan<- prometheus.Metric) error {
	begin := time.Now()
	var (
		uuids []string
	)
	execCmd := []string{"list", "--uuid"}
	output, stdErr, err := ExecCommand("virsh", execCmd)
	duration := time.Since(begin)
	if err != nil {
		level.Error(v.logger).Log("msg", "get uuid"+stdErr, "name", execCmd, "duration_seconds", duration.Seconds(), "err", err)
		return err
	} else {
		uuids = strings.Split(output, "\n")
	}

	for _, uuid := range uuids {
		if uuid == "" {
			continue
		}
		begin := time.Now()
		disksFsInfo, stdErr, err := getVirtFsInfo(uuid)
		duration := time.Since(begin)
		if err != nil {
			level.Error(v.logger).Log("msg", uuid+": getVirtFsInfo "+stdErr, "name", execCmd, "duration_seconds", duration.Seconds(), "err", err)
		}
		for _, fs := range disksFsInfo {
			ch <- prometheus.MustNewConstMetric(v.VirtFilesystemTotal, prometheus.GaugeValue, fs.TotalBytes, uuid, fs.Name, fs.MountPoint)
			ch <- prometheus.MustNewConstMetric(v.VirtFilesystemUsed, prometheus.GaugeValue, fs.UsedBytes, uuid, fs.Name, fs.MountPoint)
		}
	}

	return nil
}

func getVirtFsInfo(uuid string) ([]DiskFsInfo, string, error) {
	execCmd := []string{"qemu-agent-command", uuid, "--pretty", `{"execute": "guest-get-fsinfo"}`}
	output, stdErr, err := ExecCommand("virsh", execCmd)
	if err != nil {
		return nil, stdErr, err
	}
	res := Response{}

	json.Unmarshal([]byte(output), &res)

	return res.Return, stdErr, nil
}
