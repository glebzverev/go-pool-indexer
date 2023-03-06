package indexer

import (
	"encoding/binary"
	"math"
	"math/big"

	"golang.org/x/exp/constraints"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func Max[T constraints.Ordered](args ...T) T {
	if len(args) == 0 {
		return *new(T)
	}
	max := args[0]
	for _, arg := range args[1:] {
		if arg > max {
			max = arg
		}
	}
	return max
}

func Min[T constraints.Ordered](args ...T) T {
	if len(args) == 0 {
		return *new(T)
	}
	min := args[0]
	for _, arg := range args[1:] {
		if arg < min {
			min = arg
		}
	}
	return min
}

func Abs[T constraints.Integer | constraints.Float](arg T) T {
	return T(math.Abs(float64(arg)))
}

func Pow[A, B constraints.Integer | constraints.Float](x A, n B) float64 {
	return math.Pow(float64(x), float64(n))
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func BigIntToBigFloat(i *big.Int) *big.Float {
	return new(big.Float).SetInt(i)
}

func BigIntToFloat(i *big.Int, decimals uint8) float64 {
	quotient := BigIntToBigFloat(i)
	denominator := big.NewFloat(Pow(10, decimals))
	quotient.Quo(quotient, denominator)
	f, _ := quotient.Float64()
	return f
}

func ToBigInt[A constraints.Integer | constraints.Float, B constraints.Integer](f A, decimals B) *big.Int {
	result := new(big.Int)
	bigF := big.NewFloat(float64(f))
	denominator := big.NewFloat(Pow(10, decimals))
	bigF.Mul(bigF, denominator).Int(result)
	return result
}

func Float64ToBytes(buf []byte, f float64) {
	bits := math.Float64bits(f)
	binary.BigEndian.PutUint64(buf, bits)
}

func Float64FromBytes(buf []byte) float64 {
	bits := binary.BigEndian.Uint64(buf)
	return math.Float64frombits(bits)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
