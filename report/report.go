package report

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/pasqualesalza/amqpga/util"
)

func Connect(host string) *mgo.Session {
	mongoSession, err := mgo.Dial(host)
	util.FailOnError(err, "Failed to connect to MongoDB.")
	return mongoSession
}

func MillisecondsSince(t time.Time) int64 {
	return time.Since(t).Nanoseconds() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

type Experiment struct {
	Id                      bson.ObjectId "_id,omitempty"
	RandomId                string        "randomId"
	Type                    string        "type"
	ClusterSize             int64         "clusterSize"
	RandomSeed              int64         "randomSeed"
	FitnessFunctionName     string        "fitnessFunctionName"
	PopulationSize          int           "populationSize"
	GenerationsNumber       int64         "generationsNumber"
	ChromosomeSize          int           "chromosomeSize"
	TournamentSelectionSize int           "tournamentSelectionSize"
	CrossoverRate           float64       "crossoverRate"
	MutationRate            float64       "mutationRate"
	PeaksNumber             int64         "peaksNumber"
	SleepTime               int64         "sleepTime"
}

type Time struct {
	Id         bson.ObjectId "_id,omitempty"
	Experiment mgo.DBRef     "experiment"
	Generation int64         "generation"
	Type       string        "type"
	Time       int64         "time"
}

type Individual struct {
	Id           bson.ObjectId "_id,omitempty"
	Experiment   mgo.DBRef     "experiment"
	Generation   int64         "generation"
	Type         string        "type"
	Chromosome   string        "chromosome"
	FitnessValue string        "fitnessValue"
}

type Latency struct {
	Id                 bson.ObjectId "_id,omitempty"
	NodeId             string        "nodeId"
	ExperimentRandomId string        "experimentRandomId"
	Generation         int64         "generation"
	IndividualId       int64         "individualId"
	Type               string        "type"
	Time               int64         "time"
}

func ReportExperiment(experiment *Experiment, collection *mgo.Collection) bson.ObjectId {
	experiment.Id = bson.NewObjectId()
	collection.Insert(experiment)
	log.WithFields(log.Fields{
		"id":                      experiment.Id,
		"randomId":                experiment.RandomId,
		"type":                    experiment.Type,
		"clusterSize":             experiment.ClusterSize,
		"randomSeed":              experiment.RandomSeed,
		"fitnessFunctionName":     experiment.FitnessFunctionName,
		"populationSize":          experiment.PopulationSize,
		"generationsNumber":       experiment.GenerationsNumber,
		"chromosomeSize":          experiment.ChromosomeSize,
		"tournamentSelectionSize": experiment.TournamentSelectionSize,
		"crossoverRate":           experiment.CrossoverRate,
		"mutationRate":            experiment.MutationRate,
		"peaksNumber":             experiment.PeaksNumber,
	}).Info("Experiment registered")
	return experiment.Id
}

func ReportTime(time *Time, collection *mgo.Collection) {
	if collection != nil {
		time.Id = bson.NewObjectId()
		collection.Insert(time)
		log.WithFields(log.Fields{
			"id":         time.Id,
			"experiment": time.Experiment,
			"generation": time.Generation,
			"type":       time.Type,
			"time":       time.Time,
		}).Info("Time registered")
	}
}

func ReportIndividual(individual *Individual, collection *mgo.Collection) {
	if collection != nil {
		individual.Id = bson.NewObjectId()
		collection.Insert(individual)
		log.WithFields(log.Fields{
			"id":           individual.Id,
			"experiment":   individual.Experiment,
			"generation":   individual.Generation,
			"type":         individual.Type,
			"chromosome":   individual.Chromosome,
			"fitnessValue": individual.FitnessValue,
		}).Info("Individual registered")
	}
}

func ReportLatency(latency *Latency, collection *mgo.Collection) {
	if collection != nil {
		latency.Id = bson.NewObjectId()
		collection.Insert(latency)
		log.WithFields(log.Fields{
			"id":                 latency.Id,
			"nodeId":             latency.NodeId,
			"experimentRandomId": latency.ExperimentRandomId,
			"generation":         latency.Generation,
			"individualId":       latency.IndividualId,
			"type":               latency.Type,
			"time":               latency.Time,
		}).Info("Latency registered")
	}
}
