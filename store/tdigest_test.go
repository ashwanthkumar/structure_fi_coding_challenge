package store

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/spenczar/tdigest/v2"
	"github.com/stretchr/testify/assert"
)

func TestMedianCalculationUsingTDigestForOddNumberOfElements(t *testing.T) {
	td := tdigest.New()
	values := []float64{1.0, 3.0, 2.0}
	for _, v := range values {
		td.Add(v, 1)
	}

	q := td.Quantile(0.5)
	assert.Equal(t, 2.0, q)
}

func TestMedianCalculationUsingTDigestForEvenNumberOfElements(t *testing.T) {
	td := tdigest.New()
	values := []float64{1.0, 3.0, 2.0, 4.0}
	for _, v := range values {
		td.Add(v, 1)
	}

	q := td.Quantile(0.5)
	assert.Equal(t, 2.5, q)
}

// go test -bench=. -benchmem

func BenchmarkTDigest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Seed(5678)
		values := make(chan float64)

		var (
			// Generate 1M uniform random data between 0 and 100
			n        int     = 1_000_000
			min, max float64 = 0, 100
		)
		go func() {
			for i := 0; i < n; i++ {
				values <- min + rand.Float64()*(max-min)
			}
			close(values)
		}()

		td := tdigest.New()

		for val := range values {
			// Add the value with weight 1
			td.Add(val, 1)
		}

		q := td.Quantile(0.5)
		if q < 48 || q > 52 {
			fmt.Printf("Median is way too less / large")
			b.Failed()
		}
	}
}
