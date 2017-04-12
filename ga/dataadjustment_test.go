package ga

import (
	"sort"
	"testing"

	"github.com/pasqualesalza/amqpga/util"
)

// Utility function.
func benchmarkDataAdjustment(chromosomeSize int, individualsNumber int, b *testing.B) {
	individuals := make([]*Individual, individualsNumber)
	for i := 0; i < individualsNumber; i++ {
		individuals[i] = new(Individual)
		individuals[i].Chromosome = ByteVectorChromosomeInitialization(chromosomeSize, 0, 1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		for i := 0; i < individualsNumber; i++ {
			individuals[i].Id = util.RandomInt64InRange(0, 1)
		}

		b.StartTimer()

		individualsCopy := make(SortByIdIndividuals, len(individuals))
		copy(individualsCopy, individuals)
		sort.Sort(individualsCopy)
		individuals = individualsCopy
	}
}

// Chromosome size 64.
func BenchmarkDataAdjustment_C64_I1000(b *testing.B) {
	benchmarkDataAdjustment(64, 1000, b)
}
func BenchmarkDataAdjustment_C128_I1000(b *testing.B) {
	benchmarkDataAdjustment(128, 1000, b)
}
func BenchmarkDataAdjustment_C256_I1000(b *testing.B) {
	benchmarkDataAdjustment(256, 1000, b)
}
func BenchmarkDataAdjustment_C512_I1000(b *testing.B) {
	benchmarkDataAdjustment(512, 1000, b)
}
func BenchmarkDataAdjustment_C1024_I1000(b *testing.B) {
	benchmarkDataAdjustment(1024, 1000, b)
}
func BenchmarkDataAdjustment_C2048_I1000(b *testing.B) {
	benchmarkDataAdjustment(2048, 1000, b)
}
func BenchmarkDataAdjustment_C4096_I1000(b *testing.B) {
	benchmarkDataAdjustment(4096, 1000, b)
}
