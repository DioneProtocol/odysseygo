package feecollector

import (
	"math/big"
	"sync"

	"github.com/ava-labs/avalanchego/database"
)

var (
	_ FeeCollector = &collector{}

	signPrefix = []byte("sign")
	pFeeKey    = []byte("pfee")
	xFeeKey    = []byte("xfee")
	cFeeKey    = []byte("cfee")
)

type FeeCollector interface {
	AddPChainValue(amount *big.Int) (*big.Int, error)
	AddCChainValue(amount *big.Int) (*big.Int, error)
	AddXChainValue(amount *big.Int) (*big.Int, error)

	SubPChainValue(amount *big.Int) (*big.Int, error)
	SubCChainValue(amount *big.Int) (*big.Int, error)
	SubXChainValue(amount *big.Int) (*big.Int, error)

	GetPChainValue() *big.Int
	GetCChainValue() *big.Int
	GetXChainValue() *big.Int
}

type collector struct {
	pChainLock sync.Mutex
	cChainLock sync.Mutex
	xChainLock sync.Mutex

	pChainValue *big.Int
	cChainValue *big.Int
	xChainValue *big.Int

	db database.Database
}

func New(db database.Database) (FeeCollector, error) {
	pChainValue, err := database.GetBigInt(db, pFeeKey)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}
	xChainValue, err := database.GetBigInt(db, xFeeKey)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}
	cChainValue, err := database.GetBigInt(db, cFeeKey)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}
	return &collector{
		db:          db,
		pChainValue: pChainValue,
		xChainValue: xChainValue,
		cChainValue: cChainValue,
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

func (c *collector) GetPChainValue() *big.Int {
	c.pChainLock.Lock()
	defer c.pChainLock.Unlock()
	return new(big.Int).Set(c.pChainValue)
}

func (c *collector) AddPChainValue(amount *big.Int) (*big.Int, error) {
	c.pChainLock.Lock()
	defer c.pChainLock.Unlock()
	return c.updateChainValue(c.pChainValue, amount, pFeeKey)
}

func (c *collector) SubPChainValue(amount *big.Int) (*big.Int, error) {
	c.pChainLock.Lock()
	defer c.pChainLock.Unlock()
	negAmount := new(big.Int).Neg(amount)
	return c.updateChainValue(c.pChainValue, negAmount, pFeeKey)
}

func (c *collector) GetXChainValue() *big.Int {
	c.xChainLock.Lock()
	defer c.xChainLock.Unlock()
	return new(big.Int).Set(c.xChainValue)
}

func (c *collector) AddXChainValue(amount *big.Int) (*big.Int, error) {
	c.xChainLock.Lock()
	defer c.xChainLock.Unlock()
	return c.updateChainValue(c.xChainValue, amount, xFeeKey)
}

func (c *collector) SubXChainValue(amount *big.Int) (*big.Int, error) {
	c.xChainLock.Lock()
	defer c.xChainLock.Unlock()
	negAmount := new(big.Int).Neg(amount)
	return c.updateChainValue(c.xChainValue, negAmount, xFeeKey)
}

func (c *collector) GetCChainValue() *big.Int {
	c.cChainLock.Lock()
	defer c.cChainLock.Unlock()
	return new(big.Int).Set(c.cChainValue)
}

func (c *collector) AddCChainValue(amount *big.Int) (*big.Int, error) {
	c.cChainLock.Lock()
	defer c.cChainLock.Unlock()
	return c.updateChainValue(c.cChainValue, amount, cFeeKey)
}

func (c *collector) SubCChainValue(amount *big.Int) (*big.Int, error) {
	c.cChainLock.Lock()
	defer c.cChainLock.Unlock()
	negAmount := new(big.Int).Neg(amount)
	return c.updateChainValue(c.cChainValue, negAmount, cFeeKey)
}
