package gearley
import (
	"testing"
)

var T = NonTerminal("T")
var A = Terminal('a')
var B = Terminal('b')

func Test_parse_Aabb(t *testing.T) {
	g := Grammar(
	    Rule(T, A, B),	    // T -> 'a' 'b'
	    Rule(T, A, T, B),	    // T -> 'a' T 'b'
	)

	g.Parse("aabb")
}


func Test_stateSet_putItem(t *testing.T) {
	ruleA := Rule(T, A)
	ruleB := Rule(T, B)

	s := newStateSet()
	if s.length() != 0 {
		t.Errorf("length not 0: %v", s)
	}

	item1a := &eitem{rule: ruleA, dot: 0, index: 0}
	item1b := &eitem{rule: ruleA, dot: 0, index: 0}
	item2a := &eitem{rule: ruleB, dot: 0, index: 0}
	//item2b := &eitem{rule: Rule(T, B), dot: 0, index: 0}

	s.putItem(item1a)
	if s.length() != 1 {
		t.Errorf("length not 1: %v", s)
	}

	s.putItem(item1b)
	if s.length() != 1 {
		t.Errorf("length not 1: %v", s)
	}

	s.putItem(item2a)
	if s.length() != 2 {
		t.Errorf("length not 2: %v", s)
	}

	//s.putItem(item2b)
	//if s.length() != 2 {
	//	t.Errorf("length not 2: %v", s)
	//}
}
