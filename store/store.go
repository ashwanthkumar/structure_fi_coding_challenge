package store

import (
	"github.com/spenczar/tdigest/v2"
)

type Stat struct {
	Symbol     string  `json:"symbol"`
	Occurrence uint64  `json:"occurrence"`
	Median     float64 `json:"median"`
	Ltp        float64 `json:"ltp"`
}

type Store struct {
	// TODO: Replace the default map with a custom map later on
	tdigestMap map[string]tdigest.TDigest
	internal   map[string]Stat
}

func NewStore() Store {
	return Store{
		internal: make(map[string]Stat),
	}
}

func (s Store) Add(symbol string, price float64) {
	td, present := s.tdigestMap[symbol]
	if !present {
		td = *tdigest.New()
	}
	td.Add(price, 1)
	s.tdigestMap[symbol] = td

	stat, present := s.internal[symbol]
	if !present {
		stat = Stat{
			Symbol: symbol,
		}
	}
	stat.Occurrence = stat.Occurrence + 1
	stat.Median = td.Quantile(0.5)
	stat.Ltp = price
	s.internal[symbol] = stat
}

func (s Store) Get(symbol string) (Stat, bool) {
	stat, present := s.internal[symbol]
	return stat, present
}

func (s Store) Symbols() []string {
	symbols := make([]string, 0, len(s.internal))
	for k := range s.internal {
		symbols = append(symbols, k)
	}

	return symbols
}
