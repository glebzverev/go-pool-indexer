package arb

import (
	"math"

	"github.com/ethereum/go-ethereum/common"
)

type Road struct {
	Route  []common.Address
	Amount float64
}

type Way struct {
	Roads        []Road
	AmountOut    float64
	Price        float64
	InversePrice float64
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
	// one := computeOut(reserveOneHop, amountIn)
	sum := 0.
	alpha := 1.
	for i, dot := range dots {
		out := 0.
		if dot.a1 > 0 && dot.a1 < 1 && alpha-dot.a1 > 0 {
			outMid := computeOut(twoHopReserves[i].first, amountIn*dot.a1)
			out = computeOut(twoHopReserves[i].second, outMid)
			alpha -= dot.a1
		}
		sum += out
	}
	one := computeOut(reserveOneHop, amountIn*alpha)

	sum += one
	// fmt.Println("Opt", sum, amountIn/sum)
	way := arb.buildWay(dots, start, finish, middles, amountIn)
	way.AmountOut = sum
	way.Price = sum / amountIn
	way.InversePrice = amountIn / sum
	return way, nil
}

func (arb *Arb) buildWay(dots []Dot, start, finish common.Address, middles []common.Address, amountIn float64) *Way {
	alpha := 1.
	way := new(Way)
	for i, dot := range dots {
		if dot.a1 > 0 && dot.a1 < 1 && alpha-dot.a1 > 0 {
			divideFirst := arb.linksMerge[start][middles[i]].DividedRoute()
			divideSecond := arb.linksMerge[middles[i]][finish].DividedRoute()
			for _, div1 := range divideFirst {
				for _, div2 := range divideSecond {
					road := Road{
						Route:  []common.Address{div1.Pool.address, div2.Pool.address},
						Amount: dot.a1 * amountIn * div1.Alpha * div2.Alpha,
					}
					way.Roads = append(way.Roads, road)
				}
			}

			alpha -= dot.a1
		}
	}
	divideOne := arb.linksMerge[start][finish].DividedRoute()
	for _, div := range divideOne {
		road := Road{
			Route:  []common.Address{div.Pool.address},
			Amount: amountIn * div.Alpha * alpha,
		}
		way.Roads = append(way.Roads, road)
	}
	return way
}

func findDot(res OneHopResrves, dX float64) float64 {
	z1 := res.Reserve0 / res.Reserve1
	return math.Sqrt(z1) * res.Reserve0 / (dX)
}

func findDots(res TwoHopReserves, dX float64) (float64, float64) {
	if res.first.Reserve0 == 0 || res.first.Reserve1 == 0 ||
		res.second.Reserve0 == 0 || res.second.Reserve1 == 0 {
		return 0, 0
	}
	x1 := res.first.Reserve0
	x2 := res.second.Reserve0
	z1 := res.first.Reserve1 / res.first.Reserve0
	z2 := res.second.Reserve0 / res.second.Reserve1
	b := (2. / 3.) * (math.Sqrt(z1) * x1) / (dX * (1 - res.first.Fee))
	c := -(1. / 3.) * (x1 * x2 * math.Sqrt(z2)) / (dX * dX * (1 - res.first.Fee) * (1 - res.second.Fee) * math.Sqrt(z1))
	a1 := -b/2 + math.Sqrt(b*b/4-c)
	a2 := -b/2 - math.Sqrt(b*b/4-c)

	return a1, a2
}

func FromZeroToOne(a float64) bool {
	return a > 0 && a < 1
}

func computeOut(res OneHopResrves, dX float64) float64 {
	return res.Reserve1 * (1 - res.Fee) * dX / (res.Reserve0 + (1-res.Fee)*dX)
}

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
