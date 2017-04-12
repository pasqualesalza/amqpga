package ga

import (
	"math/rand"
	"reflect"
	"sort"

	"github.com/pasqualesalza/amqpga/util"
)

const (
	BLXAlpha = 0.5
)

func TournamentSelection(individuals []*Individual, size int, minimization bool) *Individual {
	// Select individuals for tournament.
	randomSelectionIndices := rand.Perm(len(individuals))
	randomSelection := make([]*Individual, size)
	for i := 0; i < size; i++ {
		randomSelection[i] = individuals[randomSelectionIndices[i]]
	}

	// Sort the individuals by fitness value.
	switch minimization {
	case true:
		sortedSelection := make(SortByMinFitnessValueIndividuals, len(randomSelection))
		copy(sortedSelection, randomSelection)
		sort.Sort(sortedSelection)

		// Returns the first individual.
		return sortedSelection[0]
	case false:
		sortedSelection := make(SortByMaxFitnessValueIndividuals, len(randomSelection))
		copy(sortedSelection, randomSelection)
		sort.Sort(sortedSelection)

		// Returns the first individual.
		return sortedSelection[0]
	}

	return nil
}

func RouletteWheelSelection(individuals []*Individual, minimization bool) *Individual {
	// Extracts the weights.
	weights := make([]float64, len(individuals))
	for i, individual := range individuals {
		weights[i] = float64(individual.FitnessValue.(Float64FitnessValue))
	}

	// Finds min and max weight.
	minWeight := weights[0]
	maxWeight := weights[0]
	for _, weight := range weights {
		if weight < minWeight {
			minWeight = weight
		}

		if weight > maxWeight {
			maxWeight = weight
		}
	}

	// Normalizes the weights.
	for i := 0; i < len(weights); i++ {
		difference := (maxWeight - minWeight)
		if difference == 0 {
			weights[i] = 0
		} else {
			weights[i] = (weights[i] - minWeight) / (maxWeight - minWeight)
		}
	}

	// Adjustes the weights for minimization problems.
	if minimization {
		for i := 0; i < len(weights); i++ {
			weights[i] = 1.0 - weights[i]
		}
	}

	// Calculates the total weight.
	weightSum := 0.0
	for _, weight := range weights {
		weightSum += weight
	}

	// Gets random value.
	value := rand.Float64() * weightSum

	// Locates the random value based on the weights
	for i := 0; i < len(weights); i++ {
		value -= weights[i]
		if value <= 0 {
			return individuals[i]
		}
	}

	return individuals[len(individuals)-1]
}

func BLXCrossover(parent1, parent2 *Individual, crossoverRate float64) (Individual, Individual) {
	if rand.Float64() <= crossoverRate {
		parent1Chromosome := parent1.Chromosome.(Float64VectorChromosome)
		parent2Chromosome := parent2.Chromosome.(Float64VectorChromosome)

		child1Chromosome := make(Float64VectorChromosome, len(parent1Chromosome))
		child2Chromosome := make(Float64VectorChromosome, len(parent1Chromosome))

		for i := 0; i < len(parent1Chromosome); i++ {
			x1 := parent1Chromosome[i]
			x2 := parent2Chromosome[i]
			if x2 < x1 {
				x1, x2 = x2, x1
			}
			partial := BLXAlpha * (x2 - x1)

			minBound := x1 - partial
			maxBound := x2 - partial

			child1Chromosome[i] = util.RandomFloat64InRange(minBound, maxBound)
			child2Chromosome[i] = util.RandomFloat64InRange(minBound, maxBound)
		}

		var child1 Individual
		child1.Generation = parent1.Generation
		child1.Chromosome = child1Chromosome

		var child2 Individual
		child2.Generation = parent2.Generation
		child2.Chromosome = child2Chromosome

		return child1, child2
	}

	return *parent1, *parent2
}

func SinglePointCrossover(parent1, parent2 *Individual, crossoverRate float64) (Individual, Individual) {
	if rand.Float64() <= crossoverRate {
		parent1Chromosome := reflect.ValueOf(parent1.Chromosome)
		parent2Chromosome := reflect.ValueOf(parent2.Chromosome)

		chromosomeType := reflect.TypeOf(parent1.Chromosome)

		child1Chromosome := reflect.MakeSlice(chromosomeType, parent1Chromosome.Len(), parent1Chromosome.Len())
		child2Chromosome := reflect.MakeSlice(chromosomeType, parent1Chromosome.Len(), parent1Chromosome.Len())

		point := rand.Intn(parent1Chromosome.Len()-1) + 1

		for i := 0; i < parent1Chromosome.Len(); i++ {
			if i < point {
				child1Chromosome.Index(i).Set(parent1Chromosome.Index(i))
				child2Chromosome.Index(i).Set(parent2Chromosome.Index(i))
			} else {
				child1Chromosome.Index(i).Set(parent2Chromosome.Index(i))
				child2Chromosome.Index(i).Set(parent1Chromosome.Index(i))
			}
		}

		var child1 Individual
		child1.Generation = parent1.Generation
		child1.Chromosome = child1Chromosome.Interface()

		var child2 Individual
		child2.Generation = parent2.Generation
		child2.Chromosome = child2Chromosome.Interface()

		return child1, child2
	}

	return *parent1, *parent2
}

func TwoPointsCrossover(parent1, parent2 *Individual, crossoverRate float64) (Individual, Individual) {
	if rand.Float64() <= crossoverRate {
		parent1Chromosome := reflect.ValueOf(parent1.Chromosome)
		parent2Chromosome := reflect.ValueOf(parent2.Chromosome)

		chromosomeType := reflect.TypeOf(parent1.Chromosome)

		child1Chromosome := reflect.MakeSlice(chromosomeType, parent1Chromosome.Len(), parent1Chromosome.Len())
		child2Chromosome := reflect.MakeSlice(chromosomeType, parent1Chromosome.Len(), parent1Chromosome.Len())

		point1 := util.RandomIntInRange(1, parent1Chromosome.Len()-2)
		point2 := util.RandomIntInRange(point1+1, parent1Chromosome.Len()-1)

		for i := 0; i < parent1Chromosome.Len(); i++ {
			switch {
			case i < point1:
				child1Chromosome.Index(i).Set(parent1Chromosome.Index(i))
				child2Chromosome.Index(i).Set(parent2Chromosome.Index(i))
			case i >= point1 && i < point2:
				child1Chromosome.Index(i).Set(parent2Chromosome.Index(i))
				child2Chromosome.Index(i).Set(parent1Chromosome.Index(i))
			case i >= point2:
				child1Chromosome.Index(i).Set(parent1Chromosome.Index(i))
				child2Chromosome.Index(i).Set(parent2Chromosome.Index(i))
			}
		}

		var child1 Individual
		child1.Generation = parent1.Generation
		child1.Chromosome = child1Chromosome.Interface()

		var child2 Individual
		child2.Generation = parent2.Generation
		child2.Chromosome = child2Chromosome.Interface()

		return child1, child2
	}

	return *parent1, *parent2
}

func Float64RandomMutation(individual *Individual, min float64, max float64, mutationRate float64) {
	chromosome := individual.Chromosome.(Float64VectorChromosome)
	for i := 0; i < len(chromosome); i++ {
		if rand.Float64() <= mutationRate {
			chromosome[i] = util.RandomFloat64InRange(min, max)
		}
	}
}

func ByteRandomMutation(individual *Individual, min, max byte, mutationRate float64) {
	chromosome := individual.Chromosome.(ByteVectorChromosome)
	for i := 0; i < len(chromosome); i++ {
		if rand.Float64() <= mutationRate {
			chromosome[i] = util.RandomByteInRange(min, max)
		}
	}
}

func IntRandomMutation(individual *Individual, min, max int, mutationRate float64) {
	chromosome := individual.Chromosome.(IntVectorChromosome)
	for i := 0; i < len(chromosome); i++ {
		if rand.Float64() <= mutationRate {
			chromosome[i] = util.RandomIntInRange(min, max)
		}
	}
}
