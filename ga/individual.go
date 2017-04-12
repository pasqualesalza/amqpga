package ga

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/golang/snappy"

	"github.com/pasqualesalza/amqpga/util"
)

type Individual struct {
	Id           int64
	Generation   int64
	Chromosome   Chromosome
	FitnessValue FitnessValue
}

type Chromosome interface{}

type FitnessValue interface {
	Less(other FitnessValue) bool
}

func (individual *Individual) String() string {
	return fmt.Sprintf("{Generation: %v, Id: %v, Chromosome: %v, FitnessValue: %v}", individual.Generation, individual.Id, individual.Chromosome, individual.FitnessValue)
}

func (individual *Individual) Encode() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(*individual)
	util.FailOnError(err, "Failed to encode individual")
	bytes := buffer.Bytes()
	bytes = snappy.Encode(nil, bytes)
	return bytes
}

func (individual *Individual) Decode(data []byte) {
	data, _ = snappy.Decode(nil, data)
	buffer := *bytes.NewBuffer(data)
	decoder := gob.NewDecoder(&buffer)
	err := decoder.Decode(individual)
	util.FailOnError(err, "Failed to decode individual")
}

type SortByIdIndividuals []*Individual

func (individuals SortByIdIndividuals) Len() int {
	return len(individuals)
}
func (individuals SortByIdIndividuals) Less(i, j int) bool {
	return individuals[i].Id < individuals[j].Id
}
func (individuals SortByIdIndividuals) Swap(i, j int) {
	individuals[i], individuals[j] = individuals[j], individuals[i]
}

type SortByMinFitnessValueIndividuals []*Individual

func (individuals SortByMinFitnessValueIndividuals) Len() int {
	return len(individuals)
}
func (individuals SortByMinFitnessValueIndividuals) Less(i, j int) bool {
	return individuals[i].FitnessValue.Less(individuals[j].FitnessValue)
}
func (individuals SortByMinFitnessValueIndividuals) Swap(i, j int) {
	individuals[i], individuals[j] = individuals[j], individuals[i]
}

type SortByMaxFitnessValueIndividuals []*Individual

func (individuals SortByMaxFitnessValueIndividuals) Len() int {
	return len(individuals)
}
func (individuals SortByMaxFitnessValueIndividuals) Less(i, j int) bool {
	return !individuals[i].FitnessValue.Less(individuals[j].FitnessValue)
}
func (individuals SortByMaxFitnessValueIndividuals) Swap(i, j int) {
	individuals[i], individuals[j] = individuals[j], individuals[i]
}
