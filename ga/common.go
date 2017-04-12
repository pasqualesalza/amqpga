package ga

import (
	"encoding/gob"

	"github.com/pasqualesalza/amqpga/util"
)

// byte

type ByteFitnessValue int

func (fitnessValue ByteFitnessValue) Less(other FitnessValue) bool {
	return fitnessValue < other.(ByteFitnessValue)
}

type ByteVectorChromosome []byte

func ByteVectorChromosomeInitialization(size int, min, max byte) ByteVectorChromosome {
	chromosome := make(ByteVectorChromosome, size)
	for i := 0; i < size; i++ {
		chromosome[i] = util.RandomByteInRange(min, max)
	}
	return chromosome
}

// int

type IntFitnessValue int

func (fitnessValue IntFitnessValue) Less(other FitnessValue) bool {
	return fitnessValue < other.(IntFitnessValue)
}

type IntVectorChromosome []int

func IntVectorChromosomeInitialization(size int, min, max int) IntVectorChromosome {
	chromosome := make(IntVectorChromosome, size)
	for i := 0; i < size; i++ {
		chromosome[i] = util.RandomIntInRange(min, max)
	}
	return chromosome
}

// int64

type Int64FitnessValue int64

func (fitnessValue Int64FitnessValue) Less(other FitnessValue) bool {
	return fitnessValue < other.(Int64FitnessValue)
}

type Int64VectorChromosome []int64

func Int64VectorChromosomeInitialization(size int, min, max int64) Int64VectorChromosome {
	chromosome := make(Int64VectorChromosome, size)
	for i := 0; i < size; i++ {
		chromosome[i] = util.RandomInt64InRange(min, max)
	}
	return chromosome
}

// float32

type Float32FitnessValue float32

func (fitnessValue Float32FitnessValue) Less(other FitnessValue) bool {
	return fitnessValue < other.(Float32FitnessValue)
}

type Float32VectorChromosome []float32

func Float32VectorChromosomeInitialization(size int, min, max float32) Float32VectorChromosome {
	chromosome := make(Float32VectorChromosome, size)
	for i := 0; i < size; i++ {
		chromosome[i] = util.RandomFloat32InRange(min, max)
	}
	return chromosome
}

// float64

type Float64FitnessValue float64

func (fitnessValue Float64FitnessValue) Less(other FitnessValue) bool {
	return fitnessValue < other.(Float64FitnessValue)
}

type Float64VectorChromosome []float64

func Float64VectorChromosomeInitialization(size int, min float64, max float64) Float64VectorChromosome {
	chromosome := make(Float64VectorChromosome, size)
	for i := 0; i < size; i++ {
		chromosome[i] = util.RandomFloat64InRange(min, max)
	}
	return chromosome
}

func init() {
	gob.Register(ByteFitnessValue(0))
	gob.Register(ByteVectorChromosome{})

	gob.Register(IntFitnessValue(0))
	gob.Register(IntVectorChromosome{})

	gob.Register(Int64FitnessValue(0))
	gob.Register(Int64VectorChromosome{})

	gob.Register(Float32FitnessValue(0.0))
	gob.Register(Float32VectorChromosome{})

	gob.Register(Float64FitnessValue(0.0))
	gob.Register(Float64VectorChromosome{})
}
