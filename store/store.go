package store

import (
	"github.com/spenczar/tdigest/v2"
)

type Stat struct {
	Symbol     string  `json:"symbol"`
	Occurrence uint64  `json:"occurrence"`
	Median     float64 `json:"median"`
	Ltp        float64 `json:"ltp"`

	td tdigest.TDigest `json:"-"`
}

type Store struct {
	dictionary *MapWithPHF
}

func NewStore(allSymbols []string) Store {
	return Store{
		dictionary: NewMapWithPHF(allSymbols),
	}
}

func (s Store) Add(symbol string, price float64) {
	stat, present := s.dictionary.Get(symbol)
	if !present {
		stat = &Stat{
			Symbol: symbol,
			td:     *tdigest.NewWithCompression(500),
		}
	}
	stat.td.Add(price, 1)
	stat.Occurrence = stat.Occurrence + 1
	stat.Median = stat.td.Quantile(0.5)
	stat.Ltp = price
	s.dictionary.Set(symbol, stat)
}

func (s Store) Get(symbol string) (*Stat, bool) {
	stat, present := s.dictionary.Get(symbol)
	return stat, present
}

func (s Store) Symbols() []string {
	symbols := make([]string, 0)
	for _, stat := range s.dictionary.Values {
		if stat != nil {
			symbols = append(symbols, stat.Symbol)
		}
	}

	return symbols
}
