package ltl

import (
	"errors"
	"fmt"
	"strings"
)

// ErrNotPropositional is an error for nodes not supporting propositional logic
var ErrNotPropositional = errors.New("not propositional logic")

// ErrCompile is an error for when compilations has failed
var ErrCompile = errors.New("compile error")

// Node is any node of an LTL formula
type Node interface {
	SameAs(Node) bool
	Normalize() Node
	Compile(*RefTable) Node
	Len() int
	String() string
}

// BinaryNode is an ltl node which has two child nodes
type BinaryNode interface {
	Node
	LHSNode() Node
	RHSNode() Node
}

// UnaryNode is an ltl node which only has one child
type UnaryNode interface {
	Node
	ChildNode() Node
}

// Satisfiable is an interface that determines if a node is satisfiable
type Satisfiable interface {
	Satisfied(Set) bool
}

// RefTable references other propositional logic
type RefTable map[Ref]Node

// NewRef adds a new reference to the reference table
func (r *RefTable) NewRef(n Node) Ref {
	ref := Ref{len(r)}
	r[ref] = n
	return ref
}

func binaryNodeString(b BinaryNode, op string) string {
	var sb strings.Builder

	if _, ok := b.LHSNode().(BinaryNode); ok {
		fmt.Fprintf(&sb, "(%v)", b.LHSNode())
	} else {
		fmt.Fprint(&sb, b.LHSNode())
	}
	fmt.Fprintf(&sb, " %s ", op)
	if _, ok := b.RHSNode().(BinaryNode); ok {
		fmt.Fprintf(&sb, "(%v)", b.RHSNode())
	} else {
		fmt.Fprint(&sb, b.RHSNode())
	}
	return sb.String()
}

// Negate negates the ltl formula and removes double negations
func Negate(node Node) Node {
	if n, ok := node.(Not); ok {
		return n.ChildNode()
	}

	return Not{node}
}

// FindAtomicPropositions returns a list of all atomic propositions
// in the LTL formula starting from the node
func FindAtomicPropositions(node Node) Set {
	ap := make(Set, 0)
	return auxFindAtomicPropositions(node, ap)
}

func auxFindAtomicPropositions(node Node, acc Set) Set {
	if ap, ok := node.(AP); ok {
		return acc.Add(ap)
	} else if _, ok := node.(True); ok {
		return acc
	} else if unary, ok := node.(UnaryNode); ok {
		return auxFindAtomicPropositions(unary.ChildNode(), acc)
	} else if binary, ok := node.(BinaryNode); ok {
		aps := auxFindAtomicPropositions(binary.LHSNode(), acc)
		return auxFindAtomicPropositions(binary.RHSNode(), aps)
	}
	panic(fmt.Errorf("unknown ltl node %v", node))
}

// Closure returns a list of all sub-nodes of a given node and the
// node itself.
func Closure(node Node) Set {
	sub := Subformulae(node)

	for s := range sub {
		sub.Add(Negate(s))
	}

	return sub
}

// Subformulae return all subformulaes of a LTL formula
func Subformulae(node Node) Set {
	closureTemp := make(Set, 0)
	return auxSubformulae(node, closureTemp)
}

func auxSubformulae(node Node, acc Set) Set {
	if ap, ok := node.(AP); ok {
		return acc.Add(ap)
	} else if t, ok := node.(True); ok {
		return acc.Add(t)
	} else if unary, ok := node.(UnaryNode); ok {
		acc.Add(unary)
		return auxSubformulae(unary.ChildNode(), acc)
	} else if binary, ok := node.(BinaryNode); ok {
		acc.Add(binary)
		acc = auxSubformulae(binary.LHSNode(), acc)
		return auxSubformulae(binary.RHSNode(), acc)
	}
	panic(fmt.Sprintf("unknown ltl node %v", node))
}

// Satisfied returns true if the set satisfies phi
func Satisfied(phi Node, set Set) (sat bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			if r == ErrNotPropositional {
				err = ErrNotPropositional
			}
		}
	}()
	if p, ok := phi.(Satisfiable); ok {
		return p.Satisfied(set), nil
	}
	return false, ErrNotPropositional
}
