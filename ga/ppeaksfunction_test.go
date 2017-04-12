package ga

import (
	"testing"
)

// P-Peaks function utility.
func benchmarkPPeaksFitnessFunction(peaksNumber int, chromosomeSize int, b *testing.B) {
	peaks := make([]ByteVectorChromosome, peaksNumber)
	for i := 0; i < peaksNumber; i++ {
		peaks[i] = ByteVectorChromosomeInitialization(chromosomeSize, 0, 1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		chromosome := ByteVectorChromosomeInitialization(chromosomeSize, 0, 1)
		b.StartTimer()
		PPeaksFitnessFunction(chromosome, peaks)
	}
}

// P-Peaks P64.
func BenchmarkPPeaksFunctionFitnessEvaluation_P64_C64(b *testing.B) {
	benchmarkPPeaksFitnessFunction(64, 64, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P64_C128(b *testing.B) {
	benchmarkPPeaksFitnessFunction(64, 128, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P64_C256(b *testing.B) {
	benchmarkPPeaksFitnessFunction(64, 256, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P64_C512(b *testing.B) {
	benchmarkPPeaksFitnessFunction(64, 512, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P64_C1024(b *testing.B) {
	benchmarkPPeaksFitnessFunction(64, 1024, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P64_C2048(b *testing.B) {
	benchmarkPPeaksFitnessFunction(64, 2048, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P64_C4096(b *testing.B) {
	benchmarkPPeaksFitnessFunction(64, 4096, b)
}

// P-Peaks P128.
func BenchmarkPPeaksFunctionFitnessEvaluation_P128_C64(b *testing.B) {
	benchmarkPPeaksFitnessFunction(128, 64, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P128_C128(b *testing.B) {
	benchmarkPPeaksFitnessFunction(128, 128, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P128_C256(b *testing.B) {
	benchmarkPPeaksFitnessFunction(128, 256, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P128_C512(b *testing.B) {
	benchmarkPPeaksFitnessFunction(128, 512, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P128_C1024(b *testing.B) {
	benchmarkPPeaksFitnessFunction(128, 1024, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P128_C2048(b *testing.B) {
	benchmarkPPeaksFitnessFunction(128, 2048, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P128_C4096(b *testing.B) {
	benchmarkPPeaksFitnessFunction(128, 4096, b)
}

// P-Peaks P256.
func BenchmarkPPeaksFunctionFitnessEvaluation_P256_C64(b *testing.B) {
	benchmarkPPeaksFitnessFunction(256, 64, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P256_C128(b *testing.B) {
	benchmarkPPeaksFitnessFunction(256, 128, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P256_C256(b *testing.B) {
	benchmarkPPeaksFitnessFunction(256, 256, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P256_C512(b *testing.B) {
	benchmarkPPeaksFitnessFunction(256, 512, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P256_C1024(b *testing.B) {
	benchmarkPPeaksFitnessFunction(256, 1024, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P256_C2048(b *testing.B) {
	benchmarkPPeaksFitnessFunction(256, 2048, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P256_C4096(b *testing.B) {
	benchmarkPPeaksFitnessFunction(256, 4096, b)
}

// P-Peaks P512.
func BenchmarkPPeaksFunctionFitnessEvaluation_P512_C64(b *testing.B) {
	benchmarkPPeaksFitnessFunction(512, 64, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P512_C128(b *testing.B) {
	benchmarkPPeaksFitnessFunction(512, 128, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P512_C256(b *testing.B) {
	benchmarkPPeaksFitnessFunction(512, 256, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P512_C512(b *testing.B) {
	benchmarkPPeaksFitnessFunction(512, 512, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P512_C1024(b *testing.B) {
	benchmarkPPeaksFitnessFunction(512, 1024, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P512_C2048(b *testing.B) {
	benchmarkPPeaksFitnessFunction(512, 2048, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P512_C4096(b *testing.B) {
	benchmarkPPeaksFitnessFunction(512, 4096, b)
}

// P-Peaks P1024.
func BenchmarkPPeaksFunctionFitnessEvaluation_P1024_C64(b *testing.B) {
	benchmarkPPeaksFitnessFunction(1024, 64, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P1024_C128(b *testing.B) {
	benchmarkPPeaksFitnessFunction(1024, 128, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P1024_C256(b *testing.B) {
	benchmarkPPeaksFitnessFunction(1024, 256, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P1024_C512(b *testing.B) {
	benchmarkPPeaksFitnessFunction(1024, 512, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P1024_C1024(b *testing.B) {
	benchmarkPPeaksFitnessFunction(1024, 1024, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P1024_C2048(b *testing.B) {
	benchmarkPPeaksFitnessFunction(1024, 2048, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P1024_C4096(b *testing.B) {
	benchmarkPPeaksFitnessFunction(1024, 4096, b)
}

// P-Peaks P2048.
func BenchmarkPPeaksFunctionFitnessEvaluation_P2048_C64(b *testing.B) {
	benchmarkPPeaksFitnessFunction(2048, 64, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P2048_C128(b *testing.B) {
	benchmarkPPeaksFitnessFunction(2048, 128, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P2048_C256(b *testing.B) {
	benchmarkPPeaksFitnessFunction(2048, 256, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P2048_C512(b *testing.B) {
	benchmarkPPeaksFitnessFunction(2048, 512, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P2048_C1024(b *testing.B) {
	benchmarkPPeaksFitnessFunction(2048, 1024, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P2048_C2048(b *testing.B) {
	benchmarkPPeaksFitnessFunction(2048, 2048, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P2048_C4096(b *testing.B) {
	benchmarkPPeaksFitnessFunction(2048, 4096, b)
}

// P-Peaks P4096.
func BenchmarkPPeaksFunctionFitnessEvaluation_P4096_C64(b *testing.B) {
	benchmarkPPeaksFitnessFunction(4096, 64, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P4096_C128(b *testing.B) {
	benchmarkPPeaksFitnessFunction(4096, 128, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P4096_C256(b *testing.B) {
	benchmarkPPeaksFitnessFunction(4096, 256, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P4096_C512(b *testing.B) {
	benchmarkPPeaksFitnessFunction(4096, 512, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P4096_C1024(b *testing.B) {
	benchmarkPPeaksFitnessFunction(4096, 1024, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P4096_C2048(b *testing.B) {
	benchmarkPPeaksFitnessFunction(4096, 2048, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P4096_C4096(b *testing.B) {
	benchmarkPPeaksFitnessFunction(4096, 4096, b)
}

// P-Peaks P8192.
func BenchmarkPPeaksFunctionFitnessEvaluation_P8192_C64(b *testing.B) {
	benchmarkPPeaksFitnessFunction(8192, 64, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P8192_C128(b *testing.B) {
	benchmarkPPeaksFitnessFunction(8192, 128, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P8192_C256(b *testing.B) {
	benchmarkPPeaksFitnessFunction(8192, 256, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P8192_C512(b *testing.B) {
	benchmarkPPeaksFitnessFunction(8192, 512, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P8192_C1024(b *testing.B) {
	benchmarkPPeaksFitnessFunction(8192, 1024, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P8192_C2048(b *testing.B) {
	benchmarkPPeaksFitnessFunction(8192, 2048, b)
}
func BenchmarkPPeaksFunctionFitnessEvaluation_P8192_C4096(b *testing.B) {
	benchmarkPPeaksFitnessFunction(8192, 4096, b)
}
