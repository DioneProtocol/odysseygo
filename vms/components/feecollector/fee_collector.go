package feecollector

import (
	"sync"
	"sync/atomic"

	"github.com/DioneProtocol/odysseygo/database"
	"github.com/DioneProtocol/odysseygo/database/linkeddb"
	"github.com/DioneProtocol/odysseygo/database/prefixdb"
	"github.com/DioneProtocol/odysseygo/ids"
)

var (
	_ FeeCollector = &collector{}

	aFeeKey  = []byte("afee")
	dFeeKey  = []byte("dfee")
	orionKey = []byte("orion")
)

type FeeCollector interface {
	AddDChainValue(amount uint64) error
	AddAChainValue(amount uint64) error
	AddOrionsValue(orions []ids.NodeID, amount uint64) error

	SubDChainValue(amount uint64) error
	SubAChainValue(amount uint64) error
	SubOrionsValue(orions []ids.NodeID, amount uint64) error

	GetDChainValue() uint64
	GetAChainValue() uint64
	GetOrionValue(ids.NodeID) uint64
}

type collector struct {
	lock        sync.Mutex
	dChainValue *atomic.Uint64
	aChainValue *atomic.Uint64
	orions      map[ids.NodeID]uint64

	orionsDb linkeddb.LinkedDB
	db       database.Database
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

	orions := make(map[ids.NodeID]uint64)
	orionsDb := prefixdb.New(orionKey, db)
	orionsListDb := linkeddb.NewDefault(orionsDb)
	iter := orionsListDb.NewIterator()
	defer iter.Release()
	for iter.Next() {
		orionBytes := iter.Key()
		orion := ids.NodeID(orionBytes)
		value, err := database.ParseUInt64(iter.Value())
		if err != nil {
			return nil, err
		}
		orions[orion] = value
	}

	return &collector{
		db:          db,
		aChainValue: &aChainValue,
		dChainValue: &dChainValue,
		orions:      orions,
		orionsDb:    orionsListDb,
	}, nil
}

func (c *collector) updateChainValue(newValue uint64, key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return database.PutUInt64(c.db, key, newValue)
}

func (c *collector) updateOrions(orions []ids.NodeID, value uint64) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, orion := range orions {
		c.orions[orion] += value
		if err := database.PutUInt64(c.orionsDb, orion.Bytes(), c.orions[orion]); err != nil {
			return err
		}
	}

	return nil
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

func (c *collector) AddOrionsValue(orions []ids.NodeID, amount uint64) error {
	return c.updateOrions(orions, amount)
}

func (c *collector) SubOrionsValue(orions []ids.NodeID, amount uint64) error {
	return c.updateOrions(orions, ^(amount - 1))
}

func (c *collector) GetOrionValue(nodeID ids.NodeID) uint64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.orions[nodeID]
}
