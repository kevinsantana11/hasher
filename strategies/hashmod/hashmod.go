package hashmod

import (
	"crypto/sha256"

	"slashslinging/hasher/models/cluster"
)

type HashModStrategy struct{}

func New() *HashModStrategy {
	return &HashModStrategy{}
}

func (hms HashModStrategy) GetPartitionIndex(clus cluster.Cluster, key string) int32 {
	keyBytes := make([]byte, 0)
	for _, char := range key {
		keyBytes = append(keyBytes, byte(char))
	}
	sum := sha256.Sum256(keyBytes)
	val := int32(0)
	for i := 0; i < 4; i++ {
		val += int32(sum[i]) << i * 8
	}
	return int32(val % int32(len(clus.Servers())))
}

func (hms HashModStrategy) Info() {}
