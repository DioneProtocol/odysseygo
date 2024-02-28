package feecollector

import "github.com/DioneProtocol/odysseygo/ids"

var _ FeeCollector = &dummyFeeCollector{}

// dummyFeeCollector is used instead of the collector in subnets in order
// not to change fees in the main subnet
type dummyFeeCollector struct{}

func NewDummyCollector() FeeCollector {
	return &dummyFeeCollector{}
}

func (*dummyFeeCollector) AddAChainValue(amount uint64) error {
	return nil
}

func (*dummyFeeCollector) AddDChainValue(amount uint64) error {
	return nil
}

func (*dummyFeeCollector) AddOrionsValue(orions []ids.NodeID, amount uint64) error {
	return nil
}

func (*dummyFeeCollector) AddURewardValue(amount uint64) error {
	return nil
}

func (*dummyFeeCollector) GetAChainValue() uint64 {
	return 0
}

func (*dummyFeeCollector) GetDChainValue() uint64 {
	return 0
}

func (*dummyFeeCollector) GetOrionValue(ids.NodeID) uint64 {
	return 0
}

func (*dummyFeeCollector) GetURewardValue() uint64 {
	return 0
}

func (*dummyFeeCollector) SubAChainValue(amount uint64) error {
	return nil
}

func (*dummyFeeCollector) SubDChainValue(amount uint64) error {
	return nil
}

func (*dummyFeeCollector) SubOrionsValue(orions []ids.NodeID, amount uint64) error {
	return nil
}

func (*dummyFeeCollector) SubURewardValue(amount uint64) error {
	return nil
}
