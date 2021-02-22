package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"TestDemoApp/service/web/graph/generated"
	"TestDemoApp/service/web/graph/model"
	"context"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/process"
)

func (r *queryResolver) HostInfo(_ context.Context) (*model.HostInfo, error) {
	result := model.HostInfo{
		Uptime:               0,
		HostName:             "unknown",
		HostID:               "unknown",
		Os:                   "unknown",
		Platform:             "unknown",
		PlatformFamily:       "unknown",
		PlatformVersion:      "unknown",
		KernelVersion:        "unknown",
		KernelArch:           "unknown",
		VirtualizationSystem: "unknown",
		VirtualizationRole:   "unknown",
	}

	hostInfo, err := host.Info()

	if err != nil {
		r.log.Warnf("не удалось получить данных о хостах: %s", err)
		return &result, nil
	}

	result.HostName = hostInfo.Hostname
	result.Uptime = int(hostInfo.Uptime)
	result.HostName = hostInfo.Hostname
	result.HostID = hostInfo.HostID
	result.Os = hostInfo.OS
	result.Platform = hostInfo.Platform
	result.PlatformFamily = hostInfo.PlatformFamily
	result.PlatformVersion = hostInfo.PlatformVersion
	result.KernelVersion = hostInfo.KernelVersion
	result.KernelArch = hostInfo.KernelArch
	result.VirtualizationSystem = hostInfo.VirtualizationSystem
	result.VirtualizationRole = hostInfo.VirtualizationRole

	return &result, nil
}

func (r *queryResolver) ProcessInfo(_ context.Context) ([]*model.ProcessItem, error) {
	result := make([]*model.ProcessItem, 0)

	processes, err := process.Processes()
	if err != nil {
		r.log.Warnf("ошибка получения списка процессов: %s", err.Error())
		return result, nil
	}

	for _, processItem := range processes {
		// Получаем параметры процессов
		pName, err := processItem.Name()
		if err != nil {
			pName = "unknown"
		}
		pPID := int(processItem.Pid)

		pMemoryInfo, err := processItem.MemoryInfo()
		if err != nil {
			//r.log.Warnf("ошибка получения данных о памяти процесса %s: %s", pName, err)
			pMemoryInfo = &process.MemoryInfoStat{
				RSS:    0,
				VMS:    0,
				HWM:    0,
				Data:   0,
				Stack:  0,
				Locked: 0,
				Swap:   0,
			}
		}
		pMemHWM := int(pMemoryInfo.HWM)
		pMemVMS := int(pMemoryInfo.VMS)

		newItem := model.ProcessItem{
			Name: pName,
			Pid:  pPID,
			Vms:  pMemVMS,
			Hwm:  pMemHWM,
		}
		result = append(result, &newItem)
	}

	return result, nil
}

func (r *queryResolver) CPUInfo(_ context.Context) (*model.CPUInfo, error) {
	result := model.CPUInfo{
		CPU: 0,
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		r.log.Warnf("ошибка получения информации о CPU: %s", err.Error())
		return &result, nil
	}

	for _, cpuItem := range cpuInfo {
		result.CPU += int(cpuItem.CPU)
	}

	return &result, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
