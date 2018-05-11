package jo

import "apt-model/model"

type Node struct {
	ID     int
	Weight int
	PE, PI float64
	Links  []int
}

type Nodes []*Node

type Request struct {
	Nodes   Nodes
	Timeout int
}

type Route struct {
	Path []int
	Weight int
	Probability float64
}

type Routes []Route

func ConvertRoutes(mroutes model.Routes) Routes {
	l := len(mroutes)
	mr := new(model.Route)
	routes := make(Routes, 0, l)
	for i := 0; i < l; i++ {
		mr = &mroutes[i]
		routes = append(routes, Route{Path:mr.Path, Weight:mr.Weight,Probability:mr.Probability})
	}
	return routes
}

type Response struct {
	Routes Routes
}

func (nodes Nodes) Parse() model.Nodes {
	var n *Node
	var mn *model.Node
	l := len(nodes)
	mnodes := make(model.Nodes, l)
	for i := 0; i < l; i++ {
		n = nodes[i]
		mnodes[i] = &model.Node{ID: n.ID, Weight: n.Weight, PE: n.PE, PI: n.PI}
	}
	for i := 0; i < l; i++ {
		n = nodes[i]
		l := len(n.Links)
		mn = mnodes[i]
		mn.Links = make(model.Nodes, 0, l)
		for i := 0; i < l; i++ {
			mn.Links = append(mn.Links, mnodes.Get(n.Links[i]))
		}
	}
	return mnodes
}
