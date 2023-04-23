## virtdisk-exporter
此服务用于非侵入式监控openstack平台虚拟机磁盘容量和使用量，以及磁盘所挂载的文件系统的容量及使用量，
虚拟机内部需要提前安装好qemu-guest-agent客户端。

## 制作容器镜像
```make docker-build```

## 制作二进制文件
```make build```

## metris
虚拟机磁盘可用量： virtual_machine_disk_allocation{diskname="vda",uuid="10109125-65d2-4057-a70d-e90268ff3354"} 5.2896108544e+10
虚拟机磁盘容量： virtual_machine_disk_capacity{diskname="vda",uuid="10109125-65d2-4057-a70d-e90268ff3354"} 5.36870912e+10
虚拟机磁盘物理容量： virtual_machine_disk_physical{diskname="vda",uuid="10109125-65d2-4057-a70d-e90268ff3354"} 5.36870912e+10
虚拟机磁盘文件容量：virtual_machine_fs_total{disk="vda2",mountpoint="/boot",uuid="52edeb68-48b9-4d19-a9bd-a88fbaf9cf3c"} 5.18684672e+08
虚拟机磁盘文件已使用量：virtual_machine_fs_used{disk="vda2",mountpoint="/boot",uuid="52edeb68-48b9-4d19-a9bd-a88fbaf9cf3c"} 2.04091392e+08
