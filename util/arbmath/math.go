// Copyright 2021-2022, Offchain Labs, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package arbmath

import (
	"encoding/binary"
	"math"
	"math/big"
	"math/bits"
)

// the smallest power of two greater than the input
func NextPowerOf2(value uint64) uint64 {
	return 1 << Log2ceil(value)
}

// the log2 of the int, rounded up
func Log2ceil(value uint64) uint64 {
	return uint64(64 - bits.LeadingZeros64(value))
}

// the minimum of two ints
func MinInt(value, ceiling int64) int64 {
	if value > ceiling {
		return ceiling
	}
	return value
}

// the minimum of two uints
func MinUint(value, ceiling uint64) uint64 {
	if value > ceiling {
		return ceiling
	}
	return value
}

// the maximum of two ints
func MaxInt(value, floor int64) int64 {
	if value < floor {
		return floor
	}
	return value
}

// casts an int to a huge
func UintToBig(value uint64) *big.Int {
	return new(big.Int).SetUint64(value)
}

// casts a uint to a big float
func UintToBigFloat(value uint64) *big.Float {
	return new(big.Float).SetPrec(53).SetUint64(value)
}

// casts a huge to a uint, saturating if out of bounds
func BigToUintSaturating(value *big.Int) uint64 {
	if value.Sign() < 0 {
		return 0
	}
	if !value.IsUint64() {
		return math.MaxUint64
	}
	return value.Uint64()
}

// casts a huge to a uint, panicking if out of bounds
func BigToUintOrPanic(value *big.Int) uint64 {
	if value.Sign() < 0 {
		panic("big.Int value is less than 0")
	}
	if !value.IsUint64() {
		panic("big.Int value exceeds the max Uint64")
	}
	return value.Uint64()
}

// casts an rational to a big float
func UfracToBigFloat(numerator, denominator uint64) *big.Float {
	float := new(big.Float)
	float.Quo(UintToBigFloat(numerator), UintToBigFloat(denominator))
	return float
}

// check huge equality
func BigEquals(first, second *big.Int) bool {
	return first.Cmp(second) == 0
}

// check if a huge is less than another
func BigLessThan(first, second *big.Int) bool {
	return first.Cmp(second) < 0
}

// check if a huge is greater than another
func BigGreaterThan(first, second *big.Int) bool {
	return first.Cmp(second) > 0
}

// add a huge to another
func BigAdd(augend *big.Int, addend *big.Int) *big.Int {
	return new(big.Int).Add(augend, addend)
}

// subtract from a huge another
func BigSub(minuend *big.Int, subtrahend *big.Int) *big.Int {
	return new(big.Int).Sub(minuend, subtrahend)
}

// multiply a huge by another
func BigMul(multiplicand *big.Int, multiplier *big.Int) *big.Int {
	return new(big.Int).Mul(multiplicand, multiplier)
}

// divide a huge by another
func BigDiv(dividend *big.Int, divisor *big.Int) *big.Int {
	return new(big.Int).Div(dividend, divisor)
}

// multiply a huge by a rational
func BigMulByFrac(value *big.Int, numerator, denominator int64) *big.Int {
	value = new(big.Int).Set(value)
	value.Mul(value, big.NewInt(numerator))
	value.Div(value, big.NewInt(denominator))
	return value
}

// multiply a huge by a rational whose components are non-negative
func BigMulByUfrac(value *big.Int, numerator, denominator uint64) *big.Int {
	value = new(big.Int).Set(value)
	value.Mul(value, new(big.Int).SetUint64(numerator))
	value.Div(value, new(big.Int).SetUint64(denominator))
	return value
}

// multiply a huge by an integer
func BigMulByInt(multiplicand *big.Int, multiplier int64) *big.Int {
	return new(big.Int).Mul(multiplicand, big.NewInt(multiplier))
}

// multiply a huge by a unsigned integer
func BigMulByUint(multiplicand *big.Int, multiplier uint64) *big.Int {
	return new(big.Int).Mul(multiplicand, new(big.Int).SetUint64(multiplier))
}

// divide a huge by an unsigned integer
func BigDivByUint(dividend *big.Int, divisor uint64) *big.Int {
	return BigDiv(dividend, UintToBig(divisor))
}

// divide a huge by an integer
func BigDivByInt(dividend *big.Int, divisor int64) *big.Int {
	return BigDiv(dividend, big.NewInt(divisor))
}

// add two big floats together
func BigAddFloat(augend, addend *big.Float) *big.Float {
	return new(big.Float).Add(augend, addend)
}

// multiply a big float by another
func BigMulFloat(multiplicand, multiplier *big.Float) *big.Float {
	return new(big.Float).Mul(multiplicand, multiplier)
}

// multiply a big float by an unsigned integer
func BigFloatMulByUint(multiplicand *big.Float, multiplier uint64) *big.Float {
	return new(big.Float).Mul(multiplicand, UintToBigFloat(multiplier))
}

// add two int64's without overflow
func SaturatingAdd(augend, addend int64) int64 {
	sum := augend + addend
	if addend > 0 && sum < augend {
		sum = math.MaxInt64
	}
	if addend < 0 && sum > augend {
		sum = math.MinInt64
	}
	return sum
}

// add two uint64's without overflow
func SaturatingUAdd(augend uint64, addend uint64) uint64 {
	sum := augend + addend
	if sum < augend || sum < addend {
		sum = math.MaxUint64
	}
	return sum
}

// subtract an int64 from another without overflow
func SaturatingSub(minuend, subtrahend int64) int64 {
	return SaturatingAdd(minuend, -subtrahend)
}

// subtract a uint64 from another without underflow
func SaturatingUSub(minuend uint64, subtrahend uint64) uint64 {
	if subtrahend >= minuend {
		return 0
	}
	return minuend - subtrahend
}

// multiply two uint64's without overflow
func SaturatingUMul(multiplicand uint64, multiplier uint64) uint64 {
	product := multiplicand * multiplier
	if multiplier != 0 && product/multiplier != multiplicand {
		product = math.MaxUint64
	}
	return product
}

// multiply two int64's without over/underflow
func SaturatingMul(multiplicand int64, multiplier int64) int64 {
	product := multiplicand * multiplier
	if multiplier != 0 && product/multiplier != multiplicand {
		if (multiplicand > 0 && multiplier > 0) || (multiplicand < 0 && multiplier < 0) {
			product = math.MaxInt64
		} else {
			product = math.MinInt64
		}
	}
	return product
}

// cast a uint64 to an int64, clipping to [0, 2^63-1]
func SaturatingCast(value uint64) int64 {
	if value > math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(value)
}

// cast an int64 to a uint64, clipping to [0, 2^63-1]
func SaturatingUCast(value int64) uint64 {
	if value < 0 {
		return 0
	}
	return uint64(value)
}

// the number of eth-words needed to store n bytes
func WordsForBytes(nbytes uint64) uint64 {
	return (nbytes + 31) / 32
}

// casts a uint64 to its big-endian representation
func UintToBytes(value uint64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, value)
	return result
}

// Return the Maclaurin series approximation of e^x, where x is denominated in basis points.
// This quartic polynomial will underestimate e^x by about 5% as x approaches 20000 bips.
func ApproxExpBasisPoints(value Bips) Bips {
	input := value
	negative := value < 0
	if negative {
		input = -value
	}
	x := uint64(input)

	bips := uint64(OneInBips)
	res := bips + x/4
	res = bips + SaturatingUMul(res, x)/(3*bips)
	res = bips + SaturatingUMul(res, x)/(2*bips)
	res = bips + SaturatingUMul(res, x)/(1*bips)
	if negative {
		return Bips(SaturatingCast(bips * bips / res))
	} else {
		return Bips(SaturatingCast(res))
	}
}

// Return the Newton's method approximation of sqrt(x)
// The error should be no more than 1 for values up to 2^63
func ApproxSquareRoot(value uint64) uint64 {

	if value == 0 {
		return 0
	}

	// ensure our starting approximation's square exceeds the value
	approx := value
	for SaturatingUMul(approx, approx)/2 > value {
		approx /= 2
	}

	for i := 0; i < 4; i++ {
		if approx > value/approx {
			diff := approx - value/approx
			approx = SaturatingUAdd(value/approx, diff/2)
		} else {
			diff := value/approx - approx
			approx = SaturatingUAdd(approx, diff/2)
		}
	}
	return approx
}
