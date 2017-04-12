package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"sort"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"

	"github.com/pasqualesalza/amqpga/communication"
	"github.com/pasqualesalza/amqpga/config"
	"github.com/pasqualesalza/amqpga/ga"
	"github.com/pasqualesalza/amqpga/report"
	"github.com/pasqualesalza/amqpga/util"
)

// Sends the individuals to the slaves.
func sendIndividualsToSlaves(individuals []*ga.Individual, channel *amqp.Channel, requestQueue *amqp.Queue) {
	for _, individual := range individuals {
		communication.PublishMessage(individual.Encode(), channel, requestQueue)

		log.WithFields(log.Fields{
			"individual": individual,
			"queue":      requestQueue.Name,
		}).Debugf("Published individual on %v queue", requestQueue.Name)
	}

	log.WithFields(log.Fields{
		"queue": requestQueue.Name,
	}).Debugf("Published individuals on %v queue", requestQueue.Name)
}

// Receives individuals from the slaves.
func receiveIndividualsFromSlaves(messages <-chan amqp.Delivery, expected int, channel *amqp.Channel, responseQueue *amqp.Queue) []*ga.Individual {
	individuals := make([]*ga.Individual, expected)

	done := make(chan bool)
	go func() {
		i := 0
		for message := range messages {
			individual := new(ga.Individual)
			individual.Decode(message.Body)
			individuals[i] = individual
			message.Ack(false)

			log.WithFields(log.Fields{
				"individual": individual,
				"queue":      responseQueue.Name,
			}).Debugf("Consumed individual from %v queue", responseQueue.Name)

			i++
			if i >= expected {
				break
			}
		}
		done <- true
	}()
	<-done

	// Sort the individuals by id.
	individualsCopy := make(ga.SortByIdIndividuals, len(individuals))
	copy(individualsCopy, individuals)
	sort.Sort(individualsCopy)
	individuals = individualsCopy

	log.WithFields(log.Fields{
		"queue": responseQueue.Name,
	}).Debugf("Consumed individuals from %v queue", responseQueue.Name)

	return individuals
}

// Receive individuals from the master.
func receiveIndividualsFromMaster(fitnessFunctionName string, fitnessFunctionArguments interface{}, messages <-chan amqp.Delivery, channel *amqp.Channel, requestQueue *amqp.Queue, responseQueue *amqp.Queue) {
	for message := range messages {
		individual := new(ga.Individual)
		individual.Decode(message.Body)

		log.WithFields(log.Fields{
			"individual": individual,
			"queue":      requestQueue.Name,
		}).Debugf("Consumed individual from %v queue", requestQueue.Name)

		individual.FitnessValue = executeFitnessFunction(individual, fitnessFunctionName, fitnessFunctionArguments)

		sendIndividualToMaster(individual, channel, responseQueue)

		log.WithFields(log.Fields{
			"individual": individual,
			"queue":      responseQueue.Name,
		}).Debugf("Published individual to %v queue", responseQueue.Name)

		// Notifies the correct processing to the broker.
		message.Ack(false)
	}
}

// Sends an individual to the master.
func sendIndividualToMaster(individual *ga.Individual, channel *amqp.Channel, responseQueue *amqp.Queue) {
	communication.PublishMessage(individual.Encode(), channel, responseQueue)

	log.WithFields(log.Fields{
		"individual": individual,
		"queue":      responseQueue.Name,
	}).Debugf("Published individual on %v queue", responseQueue.Name)
}

func executeFitnessFunction(individual *ga.Individual, fitnessFunctionName string, fitnessFunctionArguments interface{}) ga.FitnessValue {

	var fitnessValue ga.FitnessValue

	switch fitnessFunctionName {
	case "sleep":
		ga.SleepFitnessFunction(time.Duration(fitnessFunctionArguments.(int64)) * time.Nanosecond)
		fitnessValue = ga.Float64FitnessValue(randomFitnessValue.Float64())
	case "sphere":
		fitnessValue = ga.SphereFunctionFitnessEvaluation(individual.Chromosome.(ga.Float64VectorChromosome))
	case "rastrigin":
		fitnessValue = ga.RastriginFunctionFitnessEvaluation(individual.Chromosome.(ga.Float64VectorChromosome))
	case "ackley":
		fitnessValue = ga.AckleyFunctionFitnessEvaluation(individual.Chromosome.(ga.Float64VectorChromosome))
	case "schwefel":
		fitnessValue = ga.SchwefelFunctionFitnessEvaluation(individual.Chromosome.(ga.Float64VectorChromosome))
	case "rosenbrock":
		fitnessValue = ga.RosenbrockFunctionFitnessEvaluation(individual.Chromosome.(ga.Float64VectorChromosome))
	case "ppeaks":
		fitnessValue = ga.PPeaksFitnessFunction(individual.Chromosome.(ga.ByteVectorChromosome), fitnessFunctionArguments.([]ga.ByteVectorChromosome))
	}

	return fitnessValue
}

// Sends the latency requests to the queue.
func sendLatencyRequests(individuals []*ga.Individual, channel *amqp.Channel, requestQueue *amqp.Queue, mongoLatenciesCollection *mgo.Collection) {
	startTimes := make([]int64, len(individuals))
	for i, individual := range individuals {
		startTimes[i] = time.Now().UnixNano()
		communication.PublishMessage(individual.Encode(), channel, requestQueue)
	}

	for i, individual := range individuals {
		report.ReportLatency(&report.Latency{
			NodeId:             nodeId,
			ExperimentRandomId: randomId,
			Generation:         individual.Generation,
			IndividualId:       individual.Id,
			Type:               "start",
			Time:               startTimes[i],
		}, mongoLatenciesCollection)
	}
}

// Processes latency requests from the queue.
func processLatencyRequests(messages <-chan amqp.Delivery, channel *amqp.Channel, requestQueue *amqp.Queue, responseQueue *amqp.Queue, mongoLatenciesCollection *mgo.Collection) {

	expected := populationSize / int(clusterSize)
	individuals := make([]*ga.Individual, expected)
	finishTimes := make([]int64, expected)

	individualsNumber := 0
	for message := range messages {
		individual := new(ga.Individual)
		individual.Decode(message.Body)
		individuals[individualsNumber] = individual
		message.Ack(false)

		finishTimes[individualsNumber] = time.Now().UnixNano()

		individualsNumber++
	}

	for i := 0; i < individualsNumber; i++ {
		individual := individuals[i]
		report.ReportLatency(&report.Latency{
			NodeId:             nodeId,
			ExperimentRandomId: randomId,
			Generation:         individual.Generation,
			IndividualId:       individual.Id,
			Type:               "finish",
			Time:               finishTimes[i],
		}, mongoLatenciesCollection)
	}
}

type ExperimentConfiguration struct {
	RandomId                string  "id"
	MongoDBDatabase         string  "mongoDBDatabase"
	ClusterSize             int64   "clusterSize"
	RandomSeed              int64   "randomSeed"
	FitnessFunctionName     string  "fitnessFunctionName"
	PopulationSize          int     "populationSize"
	GenerationsNumber       int64   "generationsNumber"
	ChromosomeSize          int     "chromosomeSize"
	TournamentSelectionSize int     "tournamentSelectionSize"
	CrossoverRate           float64 "crossoverRate"
	MutationRate            float64 "mutationRate"
	PeaksNumber             int64   "peaksNumber"
	SleepTime               int64   "sleepTime"
}

var etcdHost string
var role string
var rabbitMQHost string
var mongoDBHost string
var mongoDBDatabase string
var randomId string
var clusterSize int64
var randomSeed int64
var fitnessFunctionName string
var populationSize int
var generationsNumber int64
var chromosomeSize int
var tournamentSelectionSize int
var crossoverRate float64
var mutationRate float64
var peaksNumber int64
var verbose bool
var testSetup bool
var testLatency bool
var nodeId string
var sleepTime int64
var randomFitnessValue *rand.Rand

func init() {
	// Sets the flags for command line.
	flag.StringVar(&etcdHost, "etcd", "", "etcd host to configure")
	flag.StringVar(&role, "role", "sequential", "Role in computation [sequential, master, slave]")
	flag.StringVar(&rabbitMQHost, "rabbitmq", "amqp://guest:guest@localhost:5672", "RabbitMQ host")
	flag.StringVar(&mongoDBHost, "mongodb", "", "MongoDB host")
	flag.StringVar(&mongoDBDatabase, "database", "", "MongoDB database name")
	flag.Int64Var(&clusterSize, "cluster", int64(0), "Cluster size")
	flag.Int64Var(&randomSeed, "seed", int64(42), "Random seed")
	flag.StringVar(&fitnessFunctionName, "fitness", "sphere", "Fitness function name [sphere, rastrigin, ackley, schwefel, rosenbrock, ppeaks]")
	flag.IntVar(&populationSize, "population", 10, "Number of individuals in the population")
	flag.Int64Var(&generationsNumber, "generations", int64(10), "Number of generations")
	flag.IntVar(&chromosomeSize, "chromosome", 10, "Chromosome size")
	flag.IntVar(&tournamentSelectionSize, "selection", 2, "Tournament selection size")
	flag.Float64Var(&crossoverRate, "crossover", float64(1.0), "Crossover rate")
	flag.Float64Var(&mutationRate, "mutation", float64(0.001), "Mutation rate")
	flag.Int64Var(&peaksNumber, "peaks", 512, "Peaks number for P-Peaks function")
	flag.BoolVar(&verbose, "verbose", false, "Log verbosely")
	flag.BoolVar(&testSetup, "test-setup", false, "Send a probe to test the cluster setup")
	flag.BoolVar(&testLatency, "test-latency", false, "Send a probe to test the latency")
	flag.Int64Var(&sleepTime, "sleep-time", 1000000, "Sleep time per sleep function")

	// Sets log options.
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(os.Stderr)
}

func main() {
	// Parses the flags.
	flag.Parse()

	// Generates a random node id.
	nodeId = util.RandomId(32)

	// Sets the log level.
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	// etcd configuration.
	if etcdHost != "" {
		configuration := client.Config{
			Endpoints:               []string{etcdHost},
			Transport:               client.DefaultTransport,
			HeaderTimeoutPerRequest: time.Second,
		}

		connection, err := client.New(configuration)
		util.FailOnError(err, "Failed to connect to etcd")

		etcd := client.NewKeysAPI(connection)

		experimentConfigurationResponse, err := etcd.Get(context.Background(), config.EtcdExperimentConfigurationKey, nil)
		util.FailOnError(err, "Failed to get the experiment configuration key")
		var experimentConfiguration ExperimentConfiguration
		json.Unmarshal([]byte(experimentConfigurationResponse.Node.Value), &experimentConfiguration)

		if role != "sequential" {
			rabbitMQConfigurationResponse, err := etcd.Get(context.Background(), config.EtcdRabbitMQConfigurationKey, nil)
			util.FailOnError(err, "Failed to get the RabbitMQ configuration key")
			rabbitMQHost = rabbitMQConfigurationResponse.Node.Value
		}

		if !testSetup {
			mongoDBConfigurationResponse, err := etcd.Get(context.Background(), config.EtcdMongoDBConfigurationKey, nil)
			util.FailOnError(err, "Failed to get the MongoDB configuration key")
			mongoDBHost = mongoDBConfigurationResponse.Node.Value
			mongoDBDatabase = experimentConfiguration.MongoDBDatabase
		}

		randomId = experimentConfiguration.RandomId
		clusterSize = experimentConfiguration.ClusterSize
		randomSeed = experimentConfiguration.RandomSeed
		fitnessFunctionName = experimentConfiguration.FitnessFunctionName
		populationSize = experimentConfiguration.PopulationSize
		generationsNumber = experimentConfiguration.GenerationsNumber
		chromosomeSize = experimentConfiguration.ChromosomeSize
		tournamentSelectionSize = experimentConfiguration.TournamentSelectionSize
		crossoverRate = experimentConfiguration.CrossoverRate
		mutationRate = experimentConfiguration.MutationRate
		peaksNumber = experimentConfiguration.PeaksNumber
		sleepTime = experimentConfiguration.SleepTime
	}

	log.WithFields(log.Fields{
		"role":                    role,
		"rabbitMQHost":            rabbitMQHost,
		"mongoDBHost":             mongoDBHost,
		"mongoDBDatabase":         mongoDBDatabase,
		"randomId":                randomId,
		"clusterSize":             clusterSize,
		"randomSeed":              randomSeed,
		"fitnessFunctionName":     fitnessFunctionName,
		"populationSize":          populationSize,
		"generationsNumber":       generationsNumber,
		"chromosomeSize":          chromosomeSize,
		"tournamentSelectionSize": tournamentSelectionSize,
		"crossoverRate":           crossoverRate,
		"mutationRate":            mutationRate,
		"peaksNumber":             peaksNumber,
		"verbose":                 verbose,
		"testSetup":               testSetup,
		"testLatency":             testLatency,
		"sleepTime":               sleepTime,
	}).Info("Settings parsed")

	// MongoDB report initialization.
	var mongoSession *mgo.Session
	var mongoTimesCollection *mgo.Collection
	var mongoIndividualsCollection *mgo.Collection
	var mongoLatenciesCollection *mgo.Collection
	var experiment mgo.DBRef
	if mongoDBHost != "" {
		mongoSession = report.Connect(mongoDBHost)
		defer mongoSession.Close()

		mongoSession.SetMode(mgo.Monotonic, true)

		mongoTimesCollection = mongoSession.DB(mongoDBDatabase).C("times")
		mongoIndividualsCollection = mongoSession.DB(mongoDBDatabase).C("individuals")
		mongoLatenciesCollection = mongoSession.DB(mongoDBDatabase).C("latencies")

		// Registers the experiment.
		mongoExperimentsCollection := mongoSession.DB(mongoDBDatabase).C("experiments")
		switch role {
		case "sequential", "master":
			var experimentType string
			if role == "sequential" {
				experimentType = "sequential"
			} else {
				experimentType = "parallel"
			}

			experimentId := report.ReportExperiment(&report.Experiment{
				RandomId:                randomId,
				Type:                    experimentType,
				ClusterSize:             clusterSize,
				RandomSeed:              randomSeed,
				FitnessFunctionName:     fitnessFunctionName,
				PopulationSize:          populationSize,
				GenerationsNumber:       generationsNumber,
				ChromosomeSize:          chromosomeSize,
				TournamentSelectionSize: tournamentSelectionSize,
				CrossoverRate:           crossoverRate,
				MutationRate:            mutationRate,
				PeaksNumber:             peaksNumber,
				SleepTime:               sleepTime,
			}, mongoExperimentsCollection)

			experiment = mgo.DBRef{
				Collection: "experiments",
				Id:         experimentId,
			}
		}
	}
	experimentStartTime := time.Now()

	// Set the random seed.
	rand.Seed(randomSeed)
	randomFitnessValue = rand.New(rand.NewSource(randomSeed))

	var connection *amqp.Connection
	var channel *amqp.Channel
	var requestQueue *amqp.Queue
	var responseQueue *amqp.Queue

	switch role {
	case "master", "slave":
		// Connects to the server.
		connection = communication.Connect(rabbitMQHost)
		defer connection.Close()

		// Opens a channel.
		channel = communication.OpenChannel(connection)
		defer channel.Close()

		// Declares the request queue.
		requestQueue = communication.CreateRequestQueue(channel)

		// Declares the response queue.
		responseQueue = communication.CreateResponseQueue(channel)
	}

	// Setup test variables adjustment.
	if testSetup {
		if clusterSize <= 0 {
			populationSize = 1
		} else {
			populationSize = int(clusterSize)
		}
		generationsNumber = 0
		chromosomeSize = 1
		peaksNumber = 1
	}

	// Latency test variables adjustment.
	if testLatency {
		generationsNumber = 1
	}

	// Sets the bounds for initialization and the fitness function.
	var minBound interface{}
	var maxBound interface{}
	var fitnessFunctionArguments interface{}

	switch fitnessFunctionName {
	case "sphere":
		minBound = ga.SphereFunctionMinBound
		maxBound = ga.SphereFunctionMaxBound
	case "rastrigin":
		minBound = ga.RastriginFunctionMinBound
		maxBound = ga.RastriginFunctionMaxBound
	case "ackley":
		minBound = ga.AckleyFunctionMinBound
		maxBound = ga.AckleyFunctionMaxBound
	case "schwefel":
		minBound = ga.SchwefelFunctionMinBound
		maxBound = ga.SchwefelFunctionMaxBound
	case "rosenbrock":
		minBound = ga.RosenbrockFunctionMinBound
		maxBound = ga.RosenbrockFunctionMaxBound
	case "ppeaks":
		minBound = byte(0)
		maxBound = byte(1)
		peaks := make([]ga.ByteVectorChromosome, peaksNumber)
		for i := int64(0); i < peaksNumber; i++ {
			peaks[i] = ga.ByteVectorChromosomeInitialization(chromosomeSize, 0, 1)
		}
		fitnessFunctionArguments = peaks
	case "sleep":
		minBound = byte(0)
		maxBound = byte(1)
		fitnessFunctionArguments = sleepTime
	}

	// Executes the routines for the selected role.
	switch role {
	case "sequential", "master":
		var responses <-chan amqp.Delivery
		if role == "master" {
			// Consumes the response queue.
			responses = communication.ConsumeQueue(channel, responseQueue)
		}

		// >> Initialization.
		log.Info("Initialization started")
		initializationStartTime := time.Now()

		population := make([]*ga.Individual, populationSize)
		switch fitnessFunctionName {
		case "sphere", "rastrigin", "ackley", "schwefel", "rosenbrock":
			for j := int64(0); j < int64(populationSize); j++ {
				population[j] = new(ga.Individual)
				population[j].Id = j
				population[j].Chromosome = ga.Float64VectorChromosomeInitialization(chromosomeSize, minBound.(float64), maxBound.(float64))
			}
		case "ppeaks", "sleep":
			for j := int64(0); j < int64(populationSize); j++ {
				population[j] = new(ga.Individual)
				population[j].Id = j
				population[j].Chromosome = ga.ByteVectorChromosomeInitialization(chromosomeSize, minBound.(byte), maxBound.(byte))
			}
		}

		log.Info("Initialization finished")
		report.ReportTime(&report.Time{
			Experiment: experiment,
			Type:       "initialization",
			Generation: 0,
			Time:       report.MillisecondsSince(initializationStartTime),
		}, mongoTimesCollection)

		i := int64(0)
		for ; i < generationsNumber; i++ {
			// Frees memory.
			go debug.FreeOSMemory()

			log.Infof("Started generation %v", i)
			generationStartTime := time.Now()

			// Sets the generation number.
			for j := 0; j < populationSize; j++ {
				population[j].Generation = int64(i)
			}

			// >> Fitness.
			log.Info("Fitness evaluation started")
			fitnessEvaluationStartTime := time.Now()

			switch role {
			case "sequential":
				for j := 0; j < populationSize; j++ {
					population[j].FitnessValue = executeFitnessFunction(population[j], fitnessFunctionName, fitnessFunctionArguments)
				}
			case "master":
				if !testLatency {
					sendIndividualsToSlaves(population, channel, requestQueue)
					population = receiveIndividualsFromSlaves(responses, len(population), channel, responseQueue)
				} else {
					sendLatencyRequests(population, channel, requestQueue, mongoLatenciesCollection)
					for {
						queue, _ := channel.QueueInspect(requestQueue.Name)
						if queue.Messages == 0 {
							break
						}
						time.Sleep(1 * time.Second)
					}
					channel.QueueDelete(requestQueue.Name, false, true, false)
				}
			}

			log.Info("Fitness evaluation finished")
			report.ReportTime(&report.Time{
				Experiment: experiment,
				Type:       "fitnessEvaluation",
				Generation: i,
				Time:       report.MillisecondsSince(fitnessEvaluationStartTime),
			}, mongoTimesCollection)

			if testLatency {
				log.Infof("Finished generation %v", i)
				report.ReportTime(&report.Time{
					Experiment: experiment,
					Type:       "generation",
					Generation: i,
					Time:       report.MillisecondsSince(generationStartTime),
				}, mongoTimesCollection)

				break
			}

			// Frees memory.
			go debug.FreeOSMemory()

			// Sort the individuals in a copy population to print best, worst and average individuals.
			populationCopy := make(ga.SortByMinFitnessValueIndividuals, len(population))
			copy(populationCopy, population)
			sort.Sort(populationCopy)
			var bestIndividual *ga.Individual
			var worstIndividual *ga.Individual

			switch fitnessFunctionName {
			case "sphere", "rastrigin", "ackley", "schwefel", "rosenbrock", "sleep":
				bestIndividual = populationCopy[0]
				worstIndividual = populationCopy[len(populationCopy)-1]
			case "ppeaks":
				bestIndividual = populationCopy[len(populationCopy)-1]
				worstIndividual = populationCopy[0]
			}
			fitnessValueSum := float64(0.0)
			for _, x := range populationCopy {
				fitnessValueSum += float64(x.FitnessValue.(ga.Float64FitnessValue))
			}
			averageFitnessValue := ga.Float64FitnessValue(fitnessValueSum / float64(len(populationCopy)))

			log.Infof("Best individual fitness: %v", bestIndividual.FitnessValue)
			report.ReportIndividual(&report.Individual{
				Experiment:   experiment,
				Generation:   i,
				Type:         "bestIndividual",
				Chromosome:   fmt.Sprintf("%v", bestIndividual.Chromosome),
				FitnessValue: fmt.Sprintf("%v", bestIndividual.FitnessValue),
			}, mongoIndividualsCollection)

			log.Infof("Worst individual fitness: %v", worstIndividual.FitnessValue)
			report.ReportIndividual(&report.Individual{
				Experiment:   experiment,
				Generation:   i,
				Type:         "worstIndividual",
				Chromosome:   fmt.Sprintf("%v", worstIndividual.Chromosome),
				FitnessValue: fmt.Sprintf("%v", worstIndividual.FitnessValue),
			}, mongoIndividualsCollection)

			log.Infof("Average fitness: %v", averageFitnessValue)
			report.ReportIndividual(&report.Individual{
				Experiment:   experiment,
				Generation:   i,
				Type:         "averageFitnessValue",
				FitnessValue: fmt.Sprintf("%v", averageFitnessValue),
			}, mongoIndividualsCollection)

			// Frees memory.
			populationCopy = nil
			go debug.FreeOSMemory()

			// >> Selection.
			log.Info("Selection started")
			selectionStartTime := time.Now()

			parents := make([]*ga.Individual, populationSize)
			switch fitnessFunctionName {
			case "sphere", "rastrigin", "ackley", "schwefel", "rosenbrock":
				for j := 0; j < populationSize; j++ {
					parents[j] = ga.TournamentSelection(population, tournamentSelectionSize, true)
				}
			case "ppeaks", "sleep":
				for j := 0; j < populationSize; j++ {
					parents[j] = ga.TournamentSelection(population, tournamentSelectionSize, false)
				}
			}

			log.Info("Selection finished")
			report.ReportTime(&report.Time{
				Experiment: experiment,
				Type:       "selection",
				Generation: i,
				Time:       report.MillisecondsSince(selectionStartTime),
			}, mongoTimesCollection)

			// Frees memory.
			population = nil // Wrong with elitism
			go debug.FreeOSMemory()

			// >> Crossover.
			log.Info("Crossover started")
			crossoverStartTime := time.Now()

			offspring := make([]*ga.Individual, populationSize)
			for j := 0; j < populationSize; j += 2 {
				parent1 := parents[j]
				parent2 := parents[j+1]
				child1, child2 := ga.TwoPointsCrossover(parent1, parent2, crossoverRate)
				offspring[j] = &child1
				offspring[j+1] = &child2
			}

			log.Info("Crossover finished")
			report.ReportTime(&report.Time{
				Experiment: experiment,
				Type:       "crossover",
				Generation: i,
				Time:       report.MillisecondsSince(crossoverStartTime),
			}, mongoTimesCollection)

			// Frees memory.
			parents = nil // Wrong with elitism
			go debug.FreeOSMemory()

			// >> Mutation.
			log.Info("Mutation started")
			mutationStartTime := time.Now()

			switch fitnessFunctionName {
			case "sphere", "rastrigin", "ackley", "schwefel", "rosenbrock":
				for j := 0; j < populationSize; j++ {
					ga.Float64RandomMutation(offspring[j], minBound.(float64), maxBound.(float64), mutationRate)
				}
			case "ppeaks":
				for j := 0; j < populationSize; j++ {
					ga.ByteRandomMutation(offspring[j], minBound.(byte), maxBound.(byte), mutationRate)
				}
			}

			log.Info("Mutation finished")
			report.ReportTime(&report.Time{
				Experiment: experiment,
				Type:       "mutation",
				Generation: i,
				Time:       report.MillisecondsSince(mutationStartTime),
			}, mongoTimesCollection)

			// Replaces population with offspring.
			population = offspring

			// Sets the id.
			for j := int64(0); j < int64(populationSize); j++ {
				population[j].Id = (i+1)*int64(populationSize) + j
			}

			log.Infof("Finished generation %v", i)
			report.ReportTime(&report.Time{
				Experiment: experiment,
				Type:       "generation",
				Generation: i,
				Time:       report.MillisecondsSince(generationStartTime),
			}, mongoTimesCollection)

			// Frees memory.
			go debug.FreeOSMemory()
		}

		// Frees memory.
		go debug.FreeOSMemory()

		if !testLatency {
			// >> Solution fitness evaluation.
			log.Info("Solution fitness evaluation started")
			solutionFitnessEvaluationStartTime := time.Now()

			switch role {
			case "sequential":
				for j := 0; j < populationSize; j++ {
					population[j].FitnessValue = executeFitnessFunction(population[j], fitnessFunctionName, fitnessFunctionArguments)
				}
			case "master":
				sendIndividualsToSlaves(population, channel, requestQueue)
				population = receiveIndividualsFromSlaves(responses, len(population), channel, responseQueue)
			}

			log.Info("Solution fitness evaluation finished")
			report.ReportTime(&report.Time{
				Experiment: experiment,
				Type:       "solutionFitnessEvaluation",
				Generation: i,
				Time:       report.MillisecondsSince(solutionFitnessEvaluationStartTime),
			}, mongoTimesCollection)

			// Frees memory.
			go debug.FreeOSMemory()

			// Sort the individuals in a copy population to print best, worst and average individuals.
			populationCopy := make(ga.SortByMinFitnessValueIndividuals, len(population))
			copy(populationCopy, population)
			var bestIndividual *ga.Individual
			var worstIndividual *ga.Individual
			switch fitnessFunctionName {
			case "sphere", "rastrigin", "ackley", "schwefel", "rosenbrock":
				bestIndividual = populationCopy[0]
				worstIndividual = populationCopy[len(populationCopy)-1]
			case "ppeaks", "sleep":
				bestIndividual = populationCopy[len(populationCopy)-1]
				worstIndividual = populationCopy[0]
			}
			fitnessValueSum := float64(0.0)
			for _, x := range populationCopy {
				fitnessValueSum += float64(x.FitnessValue.(ga.Float64FitnessValue))
			}
			averageFitnessValue := ga.Float64FitnessValue(fitnessValueSum / float64(len(populationCopy)))

			log.Infof("Solution best individual fitness: %v", bestIndividual.FitnessValue)
			report.ReportIndividual(&report.Individual{
				Experiment:   experiment,
				Generation:   i,
				Type:         "solutionBestIndividual",
				Chromosome:   fmt.Sprintf("%v", bestIndividual.Chromosome),
				FitnessValue: fmt.Sprintf("%v", bestIndividual.FitnessValue),
			}, mongoIndividualsCollection)

			log.Infof("Solution worst individual fitness: %v", worstIndividual.FitnessValue)
			report.ReportIndividual(&report.Individual{
				Experiment:   experiment,
				Generation:   i,
				Type:         "solutionWorstIndividual",
				Chromosome:   fmt.Sprintf("%v", worstIndividual.Chromosome),
				FitnessValue: fmt.Sprintf("%v", worstIndividual.FitnessValue),
			}, mongoIndividualsCollection)

			log.Infof("Solution average fitness: %v", averageFitnessValue)
			report.ReportIndividual(&report.Individual{
				Experiment:   experiment,
				Generation:   i,
				Type:         "solutionAverageFitnessValue",
				FitnessValue: fmt.Sprintf("%v", averageFitnessValue),
			}, mongoIndividualsCollection)
		}

		report.ReportTime(&report.Time{
			Experiment: experiment,
			Type:       "experiment",
			Generation: i,
			Time:       report.MillisecondsSince(experimentStartTime),
		}, mongoTimesCollection)

	case "slave":
		// Consumes the request queue.
		requests := communication.ConsumeQueue(channel, requestQueue)

		forever := make(chan bool)
		if !testLatency {
			go receiveIndividualsFromMaster(fitnessFunctionName, fitnessFunctionArguments, requests, channel, requestQueue, responseQueue)
		} else {
			go processLatencyRequests(requests, channel, requestQueue, responseQueue, mongoLatenciesCollection)
		}
		<-forever
	}
}
