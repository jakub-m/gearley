package gearley

// Eerley parser.
// http://loup-vaillant.fr/tutorials/earley-parsing/recogniser

import (
    "fmt"
)

type symbol interface{}

type grammar struct {
	rules []*rule
}

type terminal struct {
}

type nonTerminal struct {
	name string
}

type stateSet struct {
	items []*eitem
}

func (s *stateSet) length() int {
	return len(s.items)
}

// eitem is a single Earley item
// rule - the corresponding rule
// dot - position in the rule. Effectively it indicates the NEXT symbol in the rule that will be processed.
// index - position of the item in the parsed string
type eitem struct {
	rule *rule
	dot int
	index int
}

func (t *eitem) String() string {
	return fmt.Sprint(t.rule)
}

// isTerminal checks if the next symbol in the item is a terminal symbol
func (t *eitem) isNextTerminal() bool {
	r := t.rule.right
	if len(r) >= t.dot {
		return false
	}
	switch r[t.dot].(type) {
		case terminal:
			return true
		default:
			return false
	}
}

// state is the highest-level state of the parser.
type state []*stateSet;

func newStateSet(rules []*rule) *stateSet {
	items := make([]*eitem, len(rules))
	for i, r := range rules {
		items[i] = &eitem{rule: r, dot: 0, index: 0}
	}
	return &stateSet{items: items}
}

func (s *stateSet) String() string {
	return fmt.Sprint(s.items)
}

func initializeState(g *grammar) *state {
	sets :=  []*stateSet{newStateSet(g.rules)}
	s := state(sets)
	return &s
}

type rule struct {
	left *nonTerminal
	right []symbol
}

func (r *rule) String() string {
	return "rule"
}

func Grammar(rules ...*rule) *grammar {
	return &grammar{rules: rules}
}

func Terminal(r rune) *terminal {
	return &terminal{}
}

func NonTerminal(name string) *nonTerminal {
	return &nonTerminal{name: name}
}

func Rule(t *nonTerminal, symbols ...symbol) *rule {
	return &rule{
		left: t,
		right: symbols}
}

func (g *grammar) Parse(input string) {
	s := initializeState(g)
	fmt.Println(s)
	s.processStateSet(0)
}

// k - index of the state set to process
func (s *state) processStateSet(k int) {
	if k >= len(*s) {
		panic(fmt.Sprintf("out of bound: %d, len is %d", k, len(*s)))
	}
	set := (*s)[k]
	// This operation mutates set, so set.length() can change in each loop.
	i := 0
	for i < set.length() {
		t := set.items[i]
		fmt.Println(i, t)
		i++
		if t.isNextTerminal() {
			// terminal symbol
		} else {
			// nonterminal
		}
	}
}

// T -> 'a' 'b'
// T -> 'a' T 'b'
// test on ab, aabb aaabbb

// Define grammar
// Parse input using the grammar
