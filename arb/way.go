package arb

import (
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/common"
)

type Road struct {
	route []common.Address
	amoun float64
}

type Way struct {
	roads []Road
}

type OneHopResrves struct {
	Reserve0 float64
	Reserve1 float64
	Fee      float64
}

type TwoHopReserves struct {
	first  OneHopResrves
	second OneHopResrves
}

type Dot struct {
	a1 float64
	a2 float64
}

func (arb *Arb) FindOptimal(start common.Address, finish common.Address, amountIn float64) (*Way, error) {
	reserveOneHop := arb.GetOneHopReserves(start, finish)
	middles := arb.GetMiddles(start, finish)
	twoHopReserves := make([]TwoHopReserves, len(middles))
	dots := make([]Dot, len(middles))
	for i, middle := range middles {
		reserves := TwoHopReserves{
			first:  arb.GetOneHopReserves(start, middle),
			second: arb.GetOneHopReserves(middle, finish),
		}
		twoHopReserves[i] = reserves
		a1, a2 := findDots(reserves, amountIn)
		dots[i] = Dot{a1, a2}
	}
	fmt.Println(findDot(reserveOneHop, amountIn))
	fmt.Println(dots)
	return nil, nil
}

func findDot(res OneHopResrves, dX float64) float64 {
	z1 := res.Reserve0 / res.Reserve1
	return math.Sqrt(z1) * res.Reserve0 / (2 * dX)
}

func findDots(res TwoHopReserves, dX float64) (float64, float64) {
	if res.first.Reserve0 == 0 || res.first.Reserve1 == 0 ||
		res.second.Reserve0 == 0 || res.second.Reserve1 == 0 {
		return 0, 0
	}
	x1 := res.first.Reserve0
	x2 := res.second.Reserve0
	z1 := res.first.Reserve0 / res.first.Reserve1
	z2 := res.second.Reserve0 / res.second.Reserve1
	b := (2. / 3.) * (math.Sqrt(z1) * x1) / dX
	c := (1. / 3.) * (x1 * x2 * math.Sqrt(z2)) / (dX * dX * math.Sqrt(z1))
	a1 := -b/2 + math.Sqrt(b*b/4+c)
	a2 := -b/2 - math.Sqrt(b*b/4+c)
	if !FromZeroToOne(a1) {
		a1 = 0
	}
	if !FromZeroToOne(a2) {
		a2 = 0
	}
	return a1, a2
}

func FromZeroToOne(a float64) bool {
	return a > 0 && a < 1
}

func computeOut(reserve)

func (arb *Arb) GetOneHopReserves(start common.Address, finish common.Address) OneHopResrves {
	oneHop, ok := arb.linksMerge[start][finish]
	var reserveOneHop OneHopResrves
	reserveOneHop.Fee = oneHop.cammulativeFee
	if !ok {
		reserveOneHop.Reserve0 = 0.
		reserveOneHop.Reserve1 = 0.
	} else {
		reserveOneHop.Reserve0 = oneHop.cammulativeReserve1
		reserveOneHop.Reserve1 = oneHop.cammulativeReserve2
		if finish == oneHop.token0.Address {
			reserveOneHop.Reserve0 = oneHop.cammulativeReserve2
			reserveOneHop.Reserve1 = oneHop.cammulativeReserve1
		}
	}
	return reserveOneHop
}
