package main

import (
	"hash/fnv"
	"math"
	"os"
)

const (
	BloomLength = 29
	HashFuncNum = 10
)

type BloomFilter struct {
	size      uint
	hashFuncs []func(data []byte) uint
	bitset    []bool
}

// NewBloomFilter creates a new Bloom filter with the given size and number of hash functions.
func NewBloomFilter(size uint, numHashFuncs uint) *BloomFilter {
	hashFuncs := make([]func(data []byte) uint, numHashFuncs)
	for i := uint(0); i < numHashFuncs; i++ {
		hashFuncs[i] = createHashFunc(i)
	}

	return &BloomFilter{
		size:      size,
		hashFuncs: hashFuncs,
		bitset:    make([]bool, size),
	}
}

func helperByteToBoolSlice(data []byte) []bool {
	boolSlice := make([]bool, len(data))

	for i, b := range data {
		if b == 0x01 {
			boolSlice[i] = true
		}
	}

	return boolSlice
}

func CreateBloomFilter(bitset []byte) *BloomFilter {
	bloom := NewBloomFilter(uint(len(bitset)), HashFuncNum)
	bloom.bitset = helperByteToBoolSlice(bitset)
	return bloom
}

// Add adds an element to the Bloom filter.
func (bf *BloomFilter) Add(data []byte) {
	for _, hashFunc := range bf.hashFuncs {
		index := hashFunc(data) % bf.size
		bf.bitset[index] = true
	}
}

// Test checks if an element is possibly in the Bloom filter.
func (bf *BloomFilter) Test(data []byte) bool {
	for _, hashFunc := range bf.hashFuncs {
		index := hashFunc(data) % bf.size
		if !bf.bitset[index] {
			return false
		}
	}

	return true
}

func createHashFunc(seed uint) func(data []byte) uint {
	return func(data []byte) uint {
		hasher := fnv.New32a()
		hasher.Write(data)
		hashValue := hasher.Sum32()
		return (uint(hashValue) + 19*seed) % math.MaxUint32
	}
}

func (bf *BloomFilter) WriteToFile(writer *os.File) error {
	var b byte
	for _, value := range bf.bitset {
		if value {
			b = 0x01
		} else {
			b = 0x00
		}
		if _, err := writer.Write([]byte{b}); err != nil {
			return nil
		}
	}
	return nil
}
