package application

import (
	"context"
	"fmt"
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
	cpuInstanceData := make(map[string]map[int64]float64)
	memInstanceData := make(map[string]map[int64]float64)
	allTimestamps := make(map[int64]bool)
	instanceNameMap := make(map[string]string)
	var bytesToMB float64 = 1024 * 1024

	for _, db := range metrics {
		if _, ok := instanceNameMap[db.InstanceName]; !ok {
			instanceNameMap[db.InstanceName] = fmt.Sprintf("%s-%d", application.Name, len(instanceNameMap))
		}
		instanceName, _ := instanceNameMap[db.InstanceName]
		if _, ok := cpuInstanceData[instanceName]; !ok {
			cpuInstanceData[instanceName] = make(map[int64]float64)
			memInstanceData[instanceName] = make(map[int64]float64)
		}
		cpuInstanceData[instanceName][db.BucketTimestamp] = db.CPU
		memInstanceData[instanceName][db.BucketTimestamp] = db.Memory
		allTimestamps[db.BucketTimestamp] = true
	}

	// Create sorted timestamp array
	timestamps := make([]int64, 0, len(allTimestamps))
	for t := range allTimestamps {
		timestamps = append(timestamps, t)
	}

	// Sort timestamps
	for i := 0; i < len(timestamps)-1; i++ {
		for j := i + 1; j < len(timestamps); j++ {
			if timestamps[i] > timestamps[j] {
				timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
			}
		}
	}

	// Build CPU series
	cpuSeries := make([]*types.MetricsSeriesViewModel, 0, len(cpuInstanceData))
	for instanceName, data := range cpuInstanceData {
		seriesValue := make([]float64, len(timestamps))
		for i, ts := range timestamps {
			if val, ok := data[ts]; ok {
				seriesValue[i] = val
			} else {
				seriesValue[i] = 0
			}
		}

		cpuSeries = append(cpuSeries, &types.MetricsSeriesViewModel{
			Label:      instanceName,
			Timestamps: timestamps,
			Values:     seriesValue,
		})
	}

	// Memory series
	memSeries := make([]*types.MetricsSeriesViewModel, 0, len(memInstanceData))
	for instanceName, data := range memInstanceData {
		seriesValue := make([]float64, len(timestamps))
		for i, ts := range timestamps {
			if val, ok := data[ts]; ok {
				seriesValue[i] = val / float64(bytesToMB)
			} else {
				seriesValue[i] = 0
			}
		}

		memSeries = append(memSeries, &types.MetricsSeriesViewModel{
			Label:      instanceName,
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
