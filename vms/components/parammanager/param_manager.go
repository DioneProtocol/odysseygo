package parammanager

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/DioneProtocol/odysseygo/database"
)

var (
	_ ParamManager = &manager{}

	minValidatorStakeKey = []byte("minValidatorStake")
	minValidatorStakeDurationKey = []byte("minValidatorStakeDuration")
)

type ParamManager interface {
	SetMinValidatorStake(amount uint64) error
	SetMinValidatorStakeDuration(amount uint64) error

	GetMinValidatorStake() uint64
	GetMinValidatorStakeDuration() time.Duration
}

type manager struct {
	lock sync.Mutex

	minValidatorStake *atomic.Uint64
	minValidatorStakeDuration *atomic.Uint64
	db                database.Database
}

func New(db database.Database) (ParamManager, error) {

	minValidatorStakeUint, err := database.GetUInt64(db, minValidatorStakeKey)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}

	minValidatorStakeDurationUint, err := database.GetUInt64(db, minValidatorStakeDurationKey)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}

	minValidatorStake := atomic.Uint64{}
	minValidatorStake.Store(minValidatorStakeUint)

	minValidatorStakeDuration := atomic.Uint64{}
	minValidatorStakeDuration.Store(minValidatorStakeDurationUint)

	return &manager{
		db:                db,
		minValidatorStake: &minValidatorStake,
		minValidatorStakeDuration: &minValidatorStakeDuration,
	}, nil
}

func (m *manager) SetMinValidatorStake(amount uint64) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.minValidatorStake.Store(amount)

	if err := database.PutUInt64(m.db, minValidatorStakeKey, amount); err != nil {
		return err
	}

	return nil
}

func (m *manager) GetMinValidatorStake() uint64 {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.minValidatorStake.Load()
}

func (m *manager) SetMinValidatorStakeDuration(amount uint64) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.minValidatorStakeDuration.Store(amount)

	if err := database.PutUInt64(m.db, minValidatorStakeDurationKey, amount); err != nil {
		return err
	}

	return nil
}

func (m *manager) GetMinValidatorStakeDuration() time.Duration {
	m.lock.Lock()
	defer m.lock.Unlock()

	return time.Duration(m.minValidatorStakeDuration.Load())
}