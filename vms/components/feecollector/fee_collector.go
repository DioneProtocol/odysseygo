package feecollector

import (
	"sync"
	"sync/atomic"

	"github.com/DioneProtocol/odysseygo/database"
)

var (
	_ FeeCollector = &collector{}

	aFeeKey = []byte("afee")
	dFeeKey = []byte("dfee")
)

type FeeCollector interface {
	AddDChainValue(amount uint64) error
	AddAChainValue(amount uint64) error

	SubDChainValue(amount uint64) error
	SubAChainValue(amount uint64) error

	GetDChainValue() uint64
	GetAChainValue() uint64
}

type collector struct {
	lock        sync.Mutex
	dChainValue *atomic.Uint64
	aChainValue *atomic.Uint64

	db database.Database
}

func New(db database.Database) (FeeCollector, error) {
	aChainValueUint, err := database.GetUInt64(db, aFeeKey)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}
	dChainValueUint, err := database.GetUInt64(db, dFeeKey)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}

	aChainValue := atomic.Uint64{}
	dChainValue := atomic.Uint64{}

	aChainValue.Store(aChainValueUint)
	dChainValue.Store(dChainValueUint)

	return &collector{
		db:          db,
		aChainValue: &aChainValue,
		dChainValue: &dChainValue,
	}, nil
}

func (c *collector) updateChainValue(newValue uint64, key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return database.PutUInt64(c.db, key, newValue)
}

func (c *collector) GetAChainValue() uint64 {
	return c.aChainValue.Load()
}

func (c *collector) AddAChainValue(amount uint64) error {
	newValue := c.aChainValue.Add(amount)
	return c.updateChainValue(newValue, aFeeKey)
}

func (c *collector) SubAChainValue(amount uint64) error {
	newValue := c.aChainValue.Add(^(amount - 1))
	return c.updateChainValue(newValue, aFeeKey)
}

func (c *collector) GetDChainValue() uint64 {
	return c.dChainValue.Load()
}

func (c *collector) AddDChainValue(amount uint64) error {
	newValue := c.dChainValue.Add(amount)
	return c.updateChainValue(newValue, dFeeKey)
}

func (c *collector) SubDChainValue(amount uint64) error {
	newValue := c.dChainValue.Add(^(amount - 1))
	return c.updateChainValue(newValue, dFeeKey)
}
