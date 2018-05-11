package model

import (
	"time"
)

type Node struct {
	ID     int
	Weight int
	PE, PI float64
	Links  Nodes
}

type Nodes []*Node

func (nodes Nodes) Calculate(timeout time.Duration) Routes {
	routes := make(Routes, 0)
	for _, n := range nodes {
		if n.PE != 0 {
			routes = append(routes, n.Calculate(timeout)...)
		}
	}
	return routes
}

func (nodes Nodes) Get(id int) *Node {
	for _, node := range nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}

type Route struct {
	Path        []int
	Nodes       Nodes
	Weight      int
	Probability float64
}

func (route *Route) Copy() Route {
	l := len(route.Nodes)
	r := Route{Nodes: make(Nodes, l)}
	copy(r.Nodes, route.Nodes)
	return r
}

func (route *Route) Contain(node *Node) bool {
	for _, n := range route.Nodes {
		if n.ID == node.ID {
			return true
		}
	}
	return false
}

func (route *Route) Calculate() {
	node := route.Nodes[0]
	l := len(route.Nodes)
	route.Weight = node.Weight
	route.Probability = node.PE
	route.Path = make([]int, l)
	route.Path[0] = node.ID
	for i := 1; i < l; i++ {
		node = route.Nodes[i]
		route.Path[i] = node.ID
		route.Weight += node.Weight
		route.Probability *= node.PI
	}
}

type Routes []Route

func (node *Node) Calculate(timeout time.Duration) Routes {
	route := Route{Nodes: make(Nodes, 0)}
	routes := make(Routes, 0)
	stopCh := make(chan struct{})
	if timeout != 0 {
		go func() {
			time.Sleep(timeout)
			close(stopCh)
		}()
	}
	node.rcalculate(route, &routes, stopCh)
	select {
	case <-stopCh:
		return nil
	default:
		if timeout == 0 {
			close(stopCh)
		}
		return routes
	}
}

func (node *Node) rcalculate(route Route, routes *Routes, stopCh <-chan struct{}) {
	route.Nodes = append(route.Nodes, node)
	for _, n := range node.Links {
		select {
		case <-stopCh:
			return
		default:
			if !route.Contain(n) {
				n.rcalculate(route.Copy(), routes, stopCh)
			}
		}
	}
	select {
	case <-stopCh:
		return
	default:
		route.Calculate()
		*routes = append(*routes, route)
	}
}

//func (node *Node) CCalculate() Routes {
//	route := Route{Nodes:make(Nodes, 0)}
//	routes := make(Routes, 0)
//	rc := make(chan Route)
//	wg := new(sync.WaitGroup)
//	wg.Add(1)
//	go func() {
//		wg.Wait()
//		close(rc)
//	}()
//	go node.rcalculate(route, nil, rc, wg)
//	for route := range rc {
//		routes = append(routes, route)
//	}
//	return routes
//}
//
//func (node *Node) rcalculate(route Route, routes *Routes, rc chan<- Route, wg *sync.WaitGroup) {
//	route.Nodes = append(route.Nodes, node)
//	for _, n := range node.Links {
//		if !route.Contain(n) {
//			if routes == nil {
//				wg.Add(1)
//				go n.rcalculate(route.Copy(), routes, rc, wg)
//			} else {
//				n.rcalculate(route.Copy(), routes, nil, nil)
//			}
//		}
//	}
//	route.Calculate()
//	if routes == nil {
//		rc <- route
//		wg.Done()
//	} else {
//		*routes = append(*routes, route)
//	}
//}
