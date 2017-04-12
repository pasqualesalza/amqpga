package ga

import (
	"math"
	"time"

	"github.com/pasqualesalza/amqpga/util"
)

const (
	SphereFunctionMinBound = -5.12
	SphereFunctionMaxBound = 5.12

	RastriginFunctionMinBound = -5.12
	RastriginFunctionMaxBound = 5.12

	AckleyFunctionMinBound = -32.768
	AckleyFunctionMaxBound = 32.768

	SchwefelFunctionMinBound = -500.0
	SchwefelFunctionMaxBound = 500.0

	RosenbrockFunctionMinBound = -2.048
	RosenbrockFunctionMaxBound = 2.048

	PPeaksFunctionMinBound = 0
	PPeaksFunctionMaxBound = 1
)

func SphereFunctionFitnessEvaluation(vector Float64VectorChromosome) Float64FitnessValue {
	result := 0.0

	sum := 0.0
	for _, x := range vector {
		sum += math.Pow(x, 2.0)
	}

	result += sum

	return Float64FitnessValue(result)
}

func RastriginFunctionFitnessEvaluation(vector Float64VectorChromosome) Float64FitnessValue {
	result := 10.0 * float64(len(vector))

	sum := 0.0
	for _, x := range vector {
		sum += math.Pow(x, 2.0) - 10.0*math.Cos(2.0*math.Pi*x)
	}

	result += sum

	return Float64FitnessValue(result)
}

func AckleyFunctionFitnessEvaluation(vector Float64VectorChromosome) Float64FitnessValue {
	result := 0.0

	sum1 := 0.0
	for _, x := range vector {
		sum1 += math.Pow(x, 2.0)
	}

	sum2 := 0.0
	for _, x := range vector {
		sum2 += math.Cos(2.0 * math.Pi * x)
	}

	result += -20.0*math.Exp(-0.2*math.Sqrt(1.0/float64(len(vector))*sum1)) - math.Exp(1.0/float64(len(vector))*sum2) + 20.0 + math.E

	return Float64FitnessValue(result)
}

func SchwefelFunctionFitnessEvaluation(vector Float64VectorChromosome) Float64FitnessValue {
	result := 0.0

	sum := 0.0
	for _, x := range vector {
		sum += -x * math.Sin(math.Sqrt(math.Abs(x)))
	}

	result += sum

	return Float64FitnessValue(result)
}

func RosenbrockFunctionFitnessEvaluation(vector Float64VectorChromosome) Float64FitnessValue {
	result := 0.0

	sum := 0.0
	for i := 0; i < len(vector)-1; i++ {
		sum += 100.0*math.Pow(vector[i+1]-math.Pow(vector[i], 2.0), 2.0) + math.Pow(1-vector[i], 2.0)
	}

	result += sum

	return Float64FitnessValue(result)
}

func ShiftFloat64Vector(vector []float64, optimumDecisionVector []float64) []float64 {
	shiftedVector := make([]float64, len(vector))
	for i := 0; i < len(vector); i++ {
		shiftedVector[i] = vector[i] - optimumDecisionVector[i]
	}
	return shiftedVector
}

func RotateFloat64Vector(vector []float64, rotationMatrix [][]float64) []float64 {
	rotatedVector := make([]float64, len(vector))
	for i := 0; i < len(vector); i++ {
		rotatedVector[i] = 0.0
		for j := 0; j < len(vector); j++ {
			rotatedVector[i] += rotationMatrix[i][j] * vector[j]
		}
	}
	return rotatedVector
}

func PPeaksFitnessFunction(vector ByteVectorChromosome, peaks []ByteVectorChromosome) Float64FitnessValue {
	n := len(vector)
	p := len(peaks)

	max := n - util.Hamming(vector, peaks[0])

	for i := 1; i < p; i++ {
		value := n - util.Hamming(vector, peaks[i])
		if value > max {
			max = value
		}
	}

	result := float64(max) / float64(n)

	return Float64FitnessValue(result)
}

func EuclideanDistanceFitnessFunction(tour IntVectorChromosome, nodes [][2]int) IntFitnessValue {
	distance := 0

	for i := 1; i < len(nodes)-1; i++ {
		node1 := nodes[tour[i]]
		node2 := nodes[tour[i+1]]

		distance += util.EuclideanDistance(node1, node2)
	}

	return IntFitnessValue(distance)
}

func MaxMakespanFitnessFunction(schedule IntVectorChromosome, instance [][][2]int) IntFitnessValue {
	numberOfJobs := len(instance)
	numberOfMachines := len(instance[0])

	currentOperations := make([]int, numberOfJobs)

	for i := 0; i < numberOfJobs; i++ {
		currentOperations[i] = 0
	}

	currentMachine := 0
	currentMakespan := 0
	maxMakespan := 0

	for i := 0; i < len(schedule); i++ {
		job := schedule[i]

		currentOperation := currentOperations[job]

		currentMakespan += instance[job][currentOperation][1]
		currentOperations[job]++

		if (i+1)%numberOfMachines == 0 {
			if currentMakespan > maxMakespan {
				maxMakespan = currentMakespan
			}

			currentMachine++
			currentMakespan = 0
		}

		i++
	}

	return IntFitnessValue(maxMakespan)
}

func SleepFitnessFunction(duration time.Duration) {
	time.Sleep(duration)
}
