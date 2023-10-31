package hashring

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/rand"
	"sort"

	"slashslinging/hasher/models/cluster"
)

type HashRingStrategy struct {
	mappings []Mapping
}

type Mapping struct {
	id  int
	pos int
}

func (hrs *HashRingStrategy) Len() int {
	return len(hrs.mappings)
}

func (hrs *HashRingStrategy) Less(i, j int) bool {
	return hrs.mappings[i].pos < hrs.mappings[j].pos
}

func (hrs *HashRingStrategy) Swap(i, j int) {
	temp := hrs.mappings[i]
	hrs.mappings[i] = hrs.mappings[j]
	hrs.mappings[j] = temp
}

func New(clus cluster.Cluster) *HashRingStrategy {
	strategy := HashRingStrategy{make([]Mapping, 0)}
	generatedMappings := make([]Mapping, 0)

	serverSegmentLength := math.MaxInt32 / len(clus.Servers())
	serverPos := serverSegmentLength

	for _, server := range clus.Servers() {
		generatedMappings = append(generatedMappings, Mapping{server.Id(), serverPos})
		serverPos += serverSegmentLength
	}
	strategy.mappings = generatedMappings
	return &strategy
}

func (hrs *HashRingStrategy) Del(id int) {
	newMappings := make([]Mapping, 0)
	for _, mapping := range hrs.mappings {
		if mapping.id != id {
			newMappings = append(newMappings, mapping)
		}
	}
	hrs.mappings = newMappings
}

func (hrs *HashRingStrategy) Add(id int) {
	pick := rand.Int31n(int32(len(hrs.mappings)))

	newMappings := make([]Mapping, 0)
	for idx, mapping := range hrs.mappings {
		if pick == int32(idx) {
			nextMapping := hrs.findNextServer(mapping.pos)
			serverPos := int64(float64(mapping.pos)+math.Abs(float64(nextMapping.pos-mapping.pos))/2) % int64(math.MaxInt32)
			newMappings = append(newMappings, mapping, Mapping{id, int(serverPos)})
		} else {
			newMappings = append(newMappings, mapping)
		}
	}
	hrs.mappings = newMappings
	sort.Sort(hrs)
}

func (hrs HashRingStrategy) findNextServer(pos int) Mapping {
	for _, mapping := range hrs.mappings {
		if mapping.pos > pos {
			return mapping
		}
	}
	return hrs.mappings[0]
}

func (hrs HashRingStrategy) GetPartitionIndex(clus cluster.Cluster, key string) int32 {
	keyBytes := make([]byte, 0)
	for _, char := range key {
		keyBytes = append(keyBytes, byte(char))
	}
	sum := sha256.Sum256(keyBytes)
	val := int32(0)
	for i := 0; i < 4; i++ {
		val += int32(sum[i]) << i * 8
	}

	rng := rand.New(rand.NewSource(int64(val)))
	pos := rng.Int31()
	mapping := hrs.findNextServer(int(pos))
	for idx, server := range clus.Servers() {
		if server.Id() == mapping.id {
			return int32(idx)
		}
	}
	return -1
}

func (hms HashRingStrategy) Info() {
	println("[HashRing Strategy Info]")
	for _, mapping := range hms.mappings {
		fmt.Printf("(id, pos): (%d, %d)\n", mapping.id, mapping.pos)
	}
}
