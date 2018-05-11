package model_test

import (
	"testing"
	"apt-model/model"
	"time"
)

var id int
func RandomNode() *model.Node {
	id++
	return &model.Node{ID: id,Links:make(model.Nodes, 0)}
}

func RandomNodeChain(pn *model.Node, c, d int) *model.Node {
	d--
	if pn == nil {
		pn = RandomNode()
	}
	if d > 0 {
		for i := 0; i < c; i++ {
			n := RandomNode()
			pn.Links = append(pn.Links, n)
			RandomNodeChain(n, c, d)
		}
	}
	return pn
}

var node = RandomNodeChain(nil, 50, 5)

func BenchmarkCalculate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := time.Now()
		if routes := node.Calculate(100 * time.Millisecond); routes == nil {
			b.Error("Timeouted")
		}
		b.Error(time.Now().Sub(t).Seconds())
	}
}

//func BenchmarkConcurrencyCalculate(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		node.CCalculate()
//	}
//}