package parammanager

import "time"

var _ ParamManager = &dummyParamManager{}

type dummyParamManager struct{}

func NewDummyManager() ParamManager {
	return &dummyParamManager{}
}

func (*dummyParamManager) GetMinValidatorStake() uint64 {
	return 0
}

func (*dummyParamManager) SetMinValidatorStake(amount uint64) error {
	return nil
}

func (*dummyParamManager) GetMinValidatorStakeDuration() time.Duration {
	return 0
}

func (*dummyParamManager) SetMinValidatorStakeDuration(amount uint64) error {
	return nil
}