package ltl

// And is the ltl structure for the logical implication
type And struct {
	LHS, RHS Node
}

// SameAs returns true if both nodes are implications and has identical sub-trees
func (c And) SameAs(node Node) bool {
	if c2, ok := node.(And); ok {
		return c.LHS.SameAs(c2.LHS) && c.RHS.SameAs(c2.RHS)
	}
	return false
}

func (c And) LHSNode() Node {
	return c.LHS
}

func (c And) RHSNode() Node {
	return c.RHS
}

func (c And) String() string {
	return binaryNodeString(c, "&")
}

func (c And) Normalize() Node {
	return And{c.LHSNode().Normalize(), c.RHSNode().Normalize()}
}
