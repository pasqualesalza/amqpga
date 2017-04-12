package ga

import (
	"testing"
	"time"

	"github.com/pasqualesalza/amqpga/ga/data/jss"
	"github.com/pasqualesalza/amqpga/ga/data/tsp"
)

// Utility function.
func benchmarkFloat64FitnessEvaluation(fitnessFunction func(vector Float64VectorChromosome) Float64FitnessValue, chromosomeSize int, minBound float64, maxBound float64, b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		chromosome := Float64VectorChromosomeInitialization(chromosomeSize, minBound, maxBound)
		b.StartTimer()
		fitnessFunction(chromosome)
	}
}

// Sphere function.
func BenchmarkSphereFunctionFitnessEvaluation_C1000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(SphereFunctionFitnessEvaluation, 1000, SphereFunctionMinBound, SphereFunctionMaxBound, b)
}
func BenchmarkSphereFunctionFitnessEvaluation_C10000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(SphereFunctionFitnessEvaluation, 10000, SphereFunctionMinBound, SphereFunctionMaxBound, b)
}
func BenchmarkSphereFunctionFitnessEvaluation_C100000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(SphereFunctionFitnessEvaluation, 100000, SphereFunctionMinBound, SphereFunctionMaxBound, b)
}

// Rastrigin function.
func BenchmarkRastriginFunctionFitnessEvaluation_C1000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(RastriginFunctionFitnessEvaluation, 1000, RastriginFunctionMinBound, RastriginFunctionMaxBound, b)
}
func BenchmarkRastriginFunctionFitnessEvaluation_C10000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(RastriginFunctionFitnessEvaluation, 10000, RastriginFunctionMinBound, RastriginFunctionMaxBound, b)
}
func BenchmarkRastriginFunctionFitnessEvaluation_C100000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(RastriginFunctionFitnessEvaluation, 100000, RastriginFunctionMinBound, RastriginFunctionMaxBound, b)
}

// Ackley function.
func BenchmarkAckleyFunctionFitnessEvaluation_C1000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(AckleyFunctionFitnessEvaluation, 1000, AckleyFunctionMinBound, AckleyFunctionMaxBound, b)
}
func BenchmarkAckleyFunctionFitnessEvaluation_C10000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(AckleyFunctionFitnessEvaluation, 10000, AckleyFunctionMinBound, AckleyFunctionMaxBound, b)
}
func BenchmarkAckleyFunctionFitnessEvaluation_C100000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(AckleyFunctionFitnessEvaluation, 100000, AckleyFunctionMinBound, AckleyFunctionMaxBound, b)
}

// Schwefel function.
func BenchmarkSchwefelFunctionFitnessEvaluation_C1000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(SchwefelFunctionFitnessEvaluation, 1000, SchwefelFunctionMinBound, SchwefelFunctionMaxBound, b)
}
func BenchmarkSchwefelFunctionFitnessEvaluation_C10000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(SchwefelFunctionFitnessEvaluation, 10000, SchwefelFunctionMinBound, SchwefelFunctionMaxBound, b)
}
func BenchmarkSchwefelFunctionFitnessEvaluation_C100000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(SchwefelFunctionFitnessEvaluation, 100000, SchwefelFunctionMinBound, SchwefelFunctionMaxBound, b)
}

// Rosenbrock function.
func BenchmarkRosenbrockFunctionFitnessEvaluation_C1000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(RosenbrockFunctionFitnessEvaluation, 1000, RosenbrockFunctionMinBound, RosenbrockFunctionMaxBound, b)
}
func BenchmarkRosenbrockFunctionFitnessEvaluation_C10000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(RosenbrockFunctionFitnessEvaluation, 10000, RosenbrockFunctionMinBound, RosenbrockFunctionMaxBound, b)
}
func BenchmarkRosenbrockFunctionFitnessEvaluation_C100000(b *testing.B) {
	benchmarkFloat64FitnessEvaluation(RosenbrockFunctionFitnessEvaluation, 100000, RosenbrockFunctionMinBound, RosenbrockFunctionMaxBound, b)
}

// TSP utility function.
func benchmarkTSPFitnessEvaluation(data [][2]int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		chromosome := make(IntVectorChromosome, len(data))
		for j := 0; j < len(data); j++ {
			chromosome[j] = j
		}
		b.StartTimer()

		EuclideanDistanceFitnessFunction(chromosome, data)
	}
}

// TSP functions.
func BenchmarkA280TSPFitnessEvaluation(b *testing.B) {
	benchmarkTSPFitnessEvaluation(tsp.A280TSP, b)
}
func BenchmarkD15112TSPFitnessEvaluation(b *testing.B) {
	benchmarkTSPFitnessEvaluation(tsp.D15112TSP, b)
}

// JSS utility function.
func benchmarkJSSFitnessEvaluation(data [][][2]int, b *testing.B) {
	numberOfJobs := len(data)
	numberOfMachines := len(data[0])

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		chromosome := make(IntVectorChromosome, numberOfJobs*numberOfMachines)
		currentJob := 0
		for j := 0; j < len(chromosome); j++ {
			chromosome[j] = currentJob
			if (j+1)%numberOfMachines == 0 {
				currentJob = 0
			} else {
				currentJob++
			}
		}
		b.StartTimer()

		MaxMakespanFitnessFunction(chromosome, data)
	}
}

// JSS functions.
func BenchmarkABZ5JSSFitnessEvaluation(b *testing.B) {
	benchmarkJSSFitnessEvaluation(jss.ABZ5JSS, b)
}
func BenchmarkYN4JSSFitnessEvaluation(b *testing.B) {
	benchmarkJSSFitnessEvaluation(jss.YN4JSS, b)
}

// Sleep utility function.
func benchmarkSleepFitnessEvaluation(duration time.Duration, b *testing.B) {
	for i := 0; i < b.N; i++ {
		SleepFitnessFunction(duration)
	}
}

// Sleep functions.
func BenchmarkSleepFitnessEvaluation_10000ns(b *testing.B) {
	benchmarkSleepFitnessEvaluation(10000*time.Nanosecond, b)
}

func BenchmarkSleepFitnessEvaluation_100000ns(b *testing.B) {
	benchmarkSleepFitnessEvaluation(100000*time.Nanosecond, b)
}

func BenchmarkSleepFitnessEvaluation_1ms(b *testing.B) {
	benchmarkSleepFitnessEvaluation(1*time.Millisecond, b)
}

func BenchmarkSleepFitnessEvaluation_10ms(b *testing.B) {
	benchmarkSleepFitnessEvaluation(10*time.Millisecond, b)
}

func BenchmarkSleepFitnessEvaluation_100ms(b *testing.B) {
	benchmarkSleepFitnessEvaluation(100*time.Millisecond, b)
}

func BenchmarkSleepFitnessEvaluation_1000ms(b *testing.B) {
	benchmarkSleepFitnessEvaluation(1000*time.Millisecond, b)
}
