package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Report struct {
	Node  int
	Count int
}

type Reports []*Report

func (reports Reports) Print() {
	var sum int
	l := len(reports)
	fmt.Printf("\nRecursive method started %d times\n", l)
	fmt.Println("--- Bypass order ---------------------------")
	fmt.Println("    N) Node.Id | Route's length on this step")
	fmt.Println("--------------------------------------------")
	for i := l; i > 0; i-- {
		report := reports[i-1]
		sum += report.Count
		fmt.Printf("%5d) %7d | %d\n", l-i+1, report.Node, report.Count)
	}
	fmt.Println("Total number of checked nodes: ", sum)
}

type Node struct {
	Id    int
	Links Nodes
}

type Nodes []*Node

func (nodes Nodes) Print() {
	n := len(nodes)
	fmt.Println("\nNumber of nodes: ", n)
	fmt.Println("--- Nodes: --------------------")
	fmt.Println("    N) Node.Id | Number of arcs")
	fmt.Println("    -> Node.Id (Connection to)")
	fmt.Println("-------------------------------")
	for i := 0; i < n; i++ {
		node := nodes[i]
		l := len(node.Links)
		fmt.Printf("%5d) %7d | %d\n", i+1, node.Id, l)
		for j := 0; j < l; j++ {
			fmt.Printf("    -> %7d\n", node.Links[j].Id)
		}
	}
}

func (node *Node) Calculate() Reports {
	route := make([]int, 0)
	reports := make(Reports, 0)
	node.rcalculate(route, &reports)
	return reports
}

func (node *Node) rcalculate(route []int, reports *Reports) {
	route = append(route, node.Id)
	report := &Report{Node: node.Id}
	for _, n := range node.Links {
		flag := true
		l := len(route) - 1
		report.Count = l
		for i := 0; i < l; i++ {
			if route[i] == n.Id {
				flag = false
				break
			}
		}
		if flag {
			n.rcalculate(route, reports)
		}
	}
	*reports = append(*reports, report)
}

func main() {

	var graph Nodes
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			input := strings.ToUpper(strings.TrimSpace(scanner.Text()))
			if input == "EXIT" {
				break
			}
			if input == "USAGE" || input == "HELP" {
				usage()
			} else if strings.HasPrefix(input, "CREATE") {
				graph = createGraph(strings.TrimPrefix(input, "CREATE"))
				if graph != nil {
					fmt.Println("Graph created")
				} else {
					usage()
				}
			} else if input == "PRINT" {
				if graph != nil {
					graph.Print()
				} else {
					fmt.Println("Graph doesn't created")
				}
			} else if strings.HasPrefix(input, "REPORT") {
				reports := graphReports(graph, strings.TrimPrefix(input, "REPORT"))
				if reports != nil {
					reports.Print()
				} else {
					usage()
				}
			} else {
				usage()
			}
		}
	}
}

func usage() {
	fmt.Println()
	fmt.Println("CREATE [ Number of nodes (digit) ] [ Connectivity (digit) ] - (Re-)Create graph (graph connectivity must be less than number of its nodes)")
	fmt.Println("PRINT                                                       - Print created graph")
	fmt.Println("REPORT [ Node.Id (digit) ]                                  - Calculate graph from this node and print calculation report")
	fmt.Println("USAGE or HELP                                               - Show this text")
	fmt.Println("EXIT                                                        - Exit from the program")
}

func createGraph(input string) Nodes {
	parts := strings.Split(strings.TrimSpace(input), " ")
	if len(parts) != 2 {
		return nil
	}
	nodesNumber, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil
	}
	nodesConnectivity, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil
	}
	nodes := generateNodes(nodesNumber, nodesConnectivity)
	if nodes == nil {
		return nil
	}
	return nodes
}

func graphReports(graph Nodes, input string) Reports {
	id, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return nil
	}
	if id > len(graph) {
		return nil
	}
	return graph[id-1].Calculate()
}

func generateNodes(nodesNumber, nodesConnectivity int) Nodes {
	if nodesConnectivity >= nodesNumber {
		return nil
	}
	nodes := make(Nodes, nodesNumber)
	for i := 0; i < nodesNumber; i++ {
		nodes[i] = &Node{Id: i + 1}
	}
	for i := 0; i < nodesNumber; i++ {
		node := nodes[i]
		node.Links = make(Nodes, 0, nodesConnectivity)
		var l int
		var j = (i + 1) % nodesNumber
		for l < nodesConnectivity {
			if j != i {
				l++
				node.Links = append(node.Links, nodes[j])
			}
			j = (j + 1) % nodesNumber
		}
		nodes[i] = node
	}
	return nodes
}
