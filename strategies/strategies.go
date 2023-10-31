package strategies

import (
	"fmt"

	"slashslinging/hasher/models/cluster"
)

type DistributionStrategy interface {
	GetPartitionIndex(clus cluster.Cluster, key string) int32
	Info()
}

type RedistributionData struct {
	key      string
	value    string
	oldIndex int
	newIndex int
}

func Redistribute(ds DistributionStrategy, clus *cluster.Cluster, deleted map[string]string) {
	redistributed := 0
	moved := 0

	redistributionList := make([]RedistributionData, 0)
	for oldIndex, server := range clus.Servers() {
		for _, key := range server.Keys() {
			newIndex := ds.GetPartitionIndex(*clus, key)
			if oldIndex != int(newIndex) {
				redistributionList = append(redistributionList,
					RedistributionData{key, server.Get(key), oldIndex, int(newIndex)})
				redistributed += 1
			}
		}
	}

	for k, v := range deleted {
		idx := ds.GetPartitionIndex(*clus, k)
		clus.Servers()[idx].Put(k, v)
		moved += 1
	}

	for _, data := range redistributionList {
		clus.Servers()[data.oldIndex].Del(data.key)
		clus.Servers()[data.newIndex].Put(data.key, data.value)
	}

	fmt.Printf("Moved %d keys\n", moved)
	fmt.Printf("Redistributed %d keys\n", redistributed)
}
