package internal

import (
	"fmt"
	"reflect"
	"sync"
)

// DependencyGraph 依赖图
type DependencyGraph struct {
	nodes map[RegistrationKey]*GraphNode
	mu    sync.RWMutex
}

// GraphNode 依赖图中的节点
type GraphNode struct {
	Key          RegistrationKey
	Dependencies []RegistrationKey
	Visited      bool
	InStack      bool
}

func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		nodes: make(map[RegistrationKey]*GraphNode),
	}
}

func (g *DependencyGraph) AddNode(key RegistrationKey, dependencies []reflect.Type) {
	g.mu.Lock()
	defer g.mu.Unlock()

	node := &GraphNode{
		Key:          key,
		Dependencies: make([]RegistrationKey, len(dependencies)),
	}

	for i, dep := range dependencies {
		node.Dependencies[i] = RegistrationKey{Type: dep, Name: ""}
	}

	g.nodes[key] = node
}

func (g *DependencyGraph) TopologicalSort() ([]RegistrationKey, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	var sorted []RegistrationKey

	var visit func(*GraphNode) error
	visit = func(node *GraphNode) error {
		if node.InStack {
			return fmt.Errorf("circular dependency detected: %v", node.Key.Type)
		}
		if node.Visited {
			return nil
		}

		node.Visited = true
		node.InStack = true

		for _, depKey := range node.Dependencies {
			if depNode, exists := g.nodes[depKey]; exists {
				if err := visit(depNode); err != nil {
					return err
				}
			}
		}

		node.InStack = false
		sorted = append(sorted, node.Key)
		return nil
	}

	for _, node := range g.nodes {
		if !node.Visited {
			if err := visit(node); err != nil {
				return nil, err
			}
		}
	}

	return sorted, nil
}
