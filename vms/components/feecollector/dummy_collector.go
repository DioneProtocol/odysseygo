package feecollector

import "math/big"

var _ FeeCollector = &dummyFeeCollector{}

// dummyFeeCollector is used instead of the collector in subnets in order
// not to change fees in the main subnet
type dummyFeeCollector struct{}

func NewDummyCollector() FeeCollector {
	return &dummyFeeCollector{}
}

func (*dummyFeeCollector) AddDChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}

func (*dummyFeeCollector) AddOChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}

func (*dummyFeeCollector) AddAChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}

func (*dummyFeeCollector) GetDChainValue() *big.Int {
	return new(big.Int)
}

func (*dummyFeeCollector) GetOChainValue() *big.Int {
	return new(big.Int)
}

func (*dummyFeeCollector) GetAChainValue() *big.Int {
	return new(big.Int)
}

func (*dummyFeeCollector) SubDChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}

func (*dummyFeeCollector) SubOChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}

func (*dummyFeeCollector) SubAChainValue(amount *big.Int) (*big.Int, error) {
	return new(big.Int), nil
}
