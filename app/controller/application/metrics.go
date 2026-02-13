package application

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

func (c *Controller) GetMetrics(ctx context.Context, application *types.Application, span enum.MetricsSpan) (*types.AppMetricsViewModel, error) {
	from, to, bucket := getTimeFromSpan(span)

	metrics, err := c.metricsStore.ListByApplicationUID(ctx, application.UID, from, to, bucket)
	if err != nil {
		return nil, err
	}

	return convertToMetricsView(application, metrics), nil
}

func convertToMetricsView(application *types.Application, metrics []*types.AppMetricsAggregate) *types.AppMetricsViewModel {
	var bytesToMB float64 = 1024 * 1024
	allTimestamps := make(map[int64]bool)

	// Collect per-instance data and time ranges.
	type instanceInfo struct {
		name    string
		minTime int64
		maxTime int64
		cpu     map[int64]float64
		mem     map[int64]float64
	}

	instances := make(map[string]*instanceInfo)

	for _, db := range metrics {
		inst, ok := instances[db.InstanceName]
		if !ok {
			inst = &instanceInfo{
				name:    db.InstanceName,
				minTime: db.BucketTimestamp,
				maxTime: db.BucketTimestamp,
				cpu:     make(map[int64]float64),
				mem:     make(map[int64]float64),
			}
			instances[db.InstanceName] = inst
		}
		inst.cpu[db.BucketTimestamp] = db.CPU
		inst.mem[db.BucketTimestamp] = db.Memory
		if db.BucketTimestamp < inst.minTime {
			inst.minTime = db.BucketTimestamp
		}
		if db.BucketTimestamp > inst.maxTime {
			inst.maxTime = db.BucketTimestamp
		}
		allTimestamps[db.BucketTimestamp] = true
	}

	// Sort instances by their first appearance so we can merge
	// non-overlapping pods (e.g. old pod → new pod after deployment)
	// into the same series slot.
	sortedInstances := make([]*instanceInfo, 0, len(instances))
	for _, inst := range instances {
		sortedInstances = append(sortedInstances, inst)
	}
	sort.Slice(sortedInstances, func(i, j int) bool {
		return sortedInstances[i].minTime < sortedInstances[j].minTime
	})

	type replicaSlot struct {
		label   string
		maxTime int64
		cpu     map[int64]float64
		mem     map[int64]float64
	}

	var slots []*replicaSlot

	for _, inst := range sortedInstances {
		merged := false
		for _, slot := range slots {
			if inst.minTime > slot.maxTime {
				// Non-overlapping — same replica, different pod name after deploy.
				slot.maxTime = inst.maxTime
				for ts, val := range inst.cpu {
					slot.cpu[ts] = val
				}
				for ts, val := range inst.mem {
					slot.mem[ts] = val
				}
				merged = true
				break
			}
		}
		if !merged {
			slot := &replicaSlot{
				label:   fmt.Sprintf("%s-%d", application.Name, len(slots)),
				maxTime: inst.maxTime,
				cpu:     make(map[int64]float64),
				mem:     make(map[int64]float64),
			}
			for ts, val := range inst.cpu {
				slot.cpu[ts] = val
			}
			for ts, val := range inst.mem {
				slot.mem[ts] = val
			}
			slots = append(slots, slot)
		}
	}

	// Create sorted timestamp array
	timestamps := make([]int64, 0, len(allTimestamps))
	for t := range allTimestamps {
		timestamps = append(timestamps, t)
	}
	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i] < timestamps[j]
	})

	// Build CPU series
	cpuSeries := make([]*types.MetricsSeriesViewModel, 0, len(slots))
	for _, slot := range slots {
		seriesValue := make([]float64, len(timestamps))
		for i, ts := range timestamps {
			if val, ok := slot.cpu[ts]; ok {
				seriesValue[i] = val
			}
		}
		cpuSeries = append(cpuSeries, &types.MetricsSeriesViewModel{
			Label:      slot.label,
			Timestamps: timestamps,
			Values:     seriesValue,
		})
	}

	// Memory series
	memSeries := make([]*types.MetricsSeriesViewModel, 0, len(slots))
	for _, slot := range slots {
		seriesValue := make([]float64, len(timestamps))
		for i, ts := range timestamps {
			if val, ok := slot.mem[ts]; ok {
				seriesValue[i] = val / float64(bytesToMB)
			}
		}
		memSeries = append(memSeries, &types.MetricsSeriesViewModel{
			Label:      slot.label,
			Timestamps: timestamps,
			Values:     seriesValue,
		})
	}

	return &types.AppMetricsViewModel{
		CPU:    cpuSeries,
		Memory: memSeries,
	}
}

func getTimeFromSpan(span enum.MetricsSpan) (from time.Time, to time.Time, bucket int64) {
	to = time.Now().UTC()
	switch span {
	case enum.MetricsSpan1h:
		from = to.Add(-1 * time.Hour)
		bucket = 60 // 1 min
	case enum.MetricsSpan6h:
		from = to.Add(-6 * time.Hour)
		bucket = 2 * 60 // 2 min
	case enum.MetricsSpan1d:
		from = to.Add(-24 * time.Hour)
		bucket = 5 * 60 // 5 min
	case enum.MetricsSpan7d:
		from = to.Add(-7 * 24 * time.Hour)
		bucket = 15 * 60 // 15 min
		// case enum.MetricsSpan30d:
		// 	from = to - 30*24*60*60
		// 	bucket = 30 * 24 * 300
	}
	return
}
