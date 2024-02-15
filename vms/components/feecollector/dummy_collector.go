package feecollector

import "math/big"

var _ FeeCollector = &dummyFeeCollector{}

// dummyFeeCollector is used instead of the collector in subnets in order
// not to change fees in the main subnet
type dummyFeeCollector struct{}

func NewDummyCollector() FeeCollector {
	return &dummyFeeCollector{}
}

func (*dummyFeeCollector) AddCChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}

func (*dummyFeeCollector) AddPChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}

func (*dummyFeeCollector) AddXChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}

func (*dummyFeeCollector) GetCChainValue() *big.Int {
	return new(big.Int)
}

func (*dummyFeeCollector) GetPChainValue() *big.Int {
	return new(big.Int)
}

func (*dummyFeeCollector) GetXChainValue() *big.Int {
	return new(big.Int)
}

func (*dummyFeeCollector) SubCChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}

func (*dummyFeeCollector) SubPChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}

func (*dummyFeeCollector) SubXChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}
