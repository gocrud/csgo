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

// Duplicate formatDependencyTree here or move it to a shared utils file.
// Since they are in the same package 'internal', we can use the one from engine.go if it's exported or in the same package.
// It is in the same package 'internal', so we can use it directly if it's not private to engine.go
// Wait, formatDependencyTree was defined in engine.go which is package internal.
// So we can use it here.

func (g *DependencyGraph) TopologicalSort() ([]RegistrationKey, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	var sorted []RegistrationKey

	var visit func(*GraphNode, []string) error
	visit = func(node *GraphNode, path []string) error {
		currentName := formatType(node.Key.Type)
		if node.Key.Name != "" {
			currentName += fmt.Sprintf("(%s)", node.Key.Name)
		}

		if node.InStack {
			// Found a cycle
			// Format cycle as a tree-like structure too, or just a simple loop visualization
			// Since it's a loop, a simple arrow chain might be clearer, but let's stick to the tree style requested
			// actually, for a cycle, "A -> B -> A", the tree view would be:
			//   └─ A
			//      └─ B
			//         └─ ❌ A (Circular)

			// We need to construct the chain for the formatter
			// Note: 'path' contains [A, B], and currentName is A.
			// However, formatDependencyTree expects the chain leading TO the error.
			tree := formatDependencyTree(path, currentName+" (Circular)")
			return fmt.Errorf("circular dependency detected:%s", tree)
		}
		if node.Visited {
			return nil
		}

		node.Visited = true
		node.InStack = true

		newPath := append(path, currentName)

		for _, depKey := range node.Dependencies {
			if depNode, exists := g.nodes[depKey]; exists {
				if err := visit(depNode, newPath); err != nil {
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
			if err := visit(node, []string{}); err != nil {
				return nil, err
			}
		}
	}

	return sorted, nil
}
