package feecollector

import "github.com/DioneProtocol/odysseygo/ids"

var _ FeeCollector = &dummyFeeCollector{}

// dummyFeeCollector is used instead of the collector in subnets in order
// not to change fees in the main subnet
type dummyFeeCollector struct{}

func NewDummyCollector() FeeCollector {
	return &dummyFeeCollector{}
}

func (*dummyFeeCollector) AddDChainValue(amount uint64) error {
	return nil
}

func (*dummyFeeCollector) AddAChainValue(amount uint64) error {
	return nil
}

func (*dummyFeeCollector) AddOrionsValue(orions []ids.NodeID, amount uint64) error {
	return nil
}

func (*dummyFeeCollector) GetDChainValue() uint64 {
	return 0
}

func (*dummyFeeCollector) GetAChainValue() uint64 {
	return 0
}

func (*dummyFeeCollector) GetOrionValue(ids.NodeID) uint64 {
	return 0
}

func (*dummyFeeCollector) SubDChainValue(amount uint64) error {
	return nil
}

func (*dummyFeeCollector) SubAChainValue(amount uint64) error {
	return nil
}

func (*dummyFeeCollector) SubOrionsValue(orions []ids.NodeID, amount uint64) error {
	return nil
}
