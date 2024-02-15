package feecollector

import (
	"math/big"
	"sync"

	"github.com/DioneProtocol/odysseygo/database"
)

var (
	_ FeeCollector = &collector{}

	signPrefix = []byte("sign")
	oFeeKey    = []byte("ofee")
	aFeeKey    = []byte("afee")
	dFeeKey    = []byte("dfee")
)

type FeeCollector interface {
	AddOChainValue(amount *big.Int) (*big.Int, error)
	AddDChainValue(amount *big.Int) (*big.Int, error)
	AddAChainValue(amount *big.Int) (*big.Int, error)

	SubOChainValue(amount *big.Int) (*big.Int, error)
	SubDChainValue(amount *big.Int) (*big.Int, error)
	SubAChainValue(amount *big.Int) (*big.Int, error)

	GetOChainValue() *big.Int
	GetDChainValue() *big.Int
	GetAChainValue() *big.Int
}

type collector struct {
	oChainLock sync.Mutex
	dChainLock sync.Mutex
	aChainLock sync.Mutex

	oChainValue *big.Int
	dChainValue *big.Int
	aChainValue *big.Int

	db database.Database
}

func New(db database.Database) (FeeCollector, error) {
	oChainValue, err := database.GetBigInt(db, oFeeKey)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}
	aChainValue, err := database.GetBigInt(db, aFeeKey)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}
	dChainValue, err := database.GetBigInt(db, dFeeKey)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}
	return &collector{
		db:          db,
		oChainValue: oChainValue,
		aChainValue: aChainValue,
		dChainValue: dChainValue,
	}, nil
}

func (c *collector) updateChainValue(value *big.Int, diff *big.Int, key []byte) (*big.Int, error) {
	newValue := new(big.Int).Add(value, diff)
	err := database.PutBigInt(c.db, key, newValue)
	if err == nil {
		value.Set(newValue)
	}
	return value, err
}

func (c *collector) GetOChainValue() *big.Int {
	c.oChainLock.Lock()
	defer c.oChainLock.Unlock()
	return new(big.Int).Set(c.oChainValue)
}

func (c *collector) AddOChainValue(amount *big.Int) (*big.Int, error) {
	c.oChainLock.Lock()
	defer c.oChainLock.Unlock()
	return c.updateChainValue(c.oChainValue, amount, oFeeKey)
}

func (c *collector) SubOChainValue(amount *big.Int) (*big.Int, error) {
	c.oChainLock.Lock()
	defer c.oChainLock.Unlock()
	negAmount := new(big.Int).Neg(amount)
	return c.updateChainValue(c.oChainValue, negAmount, oFeeKey)
}

func (c *collector) GetAChainValue() *big.Int {
	c.aChainLock.Lock()
	defer c.aChainLock.Unlock()
	return new(big.Int).Set(c.aChainValue)
}

func (c *collector) AddAChainValue(amount *big.Int) (*big.Int, error) {
	c.aChainLock.Lock()
	defer c.aChainLock.Unlock()
	return c.updateChainValue(c.aChainValue, amount, aFeeKey)
}

func (c *collector) SubAChainValue(amount *big.Int) (*big.Int, error) {
	c.aChainLock.Lock()
	defer c.aChainLock.Unlock()
	negAmount := new(big.Int).Neg(amount)
	return c.updateChainValue(c.aChainValue, negAmount, aFeeKey)
}

func (c *collector) GetDChainValue() *big.Int {
	c.dChainLock.Lock()
	defer c.dChainLock.Unlock()
	return new(big.Int).Set(c.dChainValue)
}

func (c *collector) AddDChainValue(amount *big.Int) (*big.Int, error) {
	c.dChainLock.Lock()
	defer c.dChainLock.Unlock()
	return c.updateChainValue(c.dChainValue, amount, dFeeKey)
}

func (c *collector) SubDChainValue(amount *big.Int) (*big.Int, error) {
	c.dChainLock.Lock()
	defer c.dChainLock.Unlock()
	negAmount := new(big.Int).Neg(amount)
	return c.updateChainValue(c.dChainValue, negAmount, dFeeKey)
}
