package util

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Computes the int between [min, max]
func RandomByteInRange(min, max byte) byte {
	return byte(RandomIntInRange(int(min), int(max)))
}

// Computes the int between [min, max]
func RandomIntInRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// Computes the int64 between [min, max]
func RandomInt64InRange(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

// Computes the float32 between [min, max)
func RandomFloat32InRange(min, max float32) float32 {
	return float32(RandomFloat64InRange(float64(min), float64(max)))
}

// Computes the float64 between [min, max)
func RandomFloat64InRange(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func Hamming(x, y []byte) int {
	distance := 0
	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			distance += 1
		}
	}

	return distance
}

func EuclideanDistance(a, b [2]int) int {
	xd := float64(a[0] - b[0])
	yd := float64(a[1] - b[1])
	return Round(math.Sqrt(xd*xd + yd*yd))
}

func Round(x float64) int {
	return int(math.Floor(x + 0.5))
}

// Checks the error and prints a message.
func FailOnError(err error, message string) {
	if err != nil {
		log.WithFields(log.Fields{
			"message": message,
			"error":   err,
		}).Fatalf("%s: %s", message, err)
	}
}

func RandomId(size int64) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	randomBytes := make([]byte, size/2)
	random.Read(randomBytes)
	return fmt.Sprintf("%x", randomBytes)
}
