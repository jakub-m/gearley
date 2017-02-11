package gearley

// Eerley parser.
// http://loup-vaillant.fr/tutorials/earley-parsing/recogniser

import (
    "fmt"
    "strings"
)

const BLACK_CIRCLE = "\u25CF"

type symbol interface{
	// isTerminal indicates if the symbol is Terminal Symbol or Non Terminal Symbol.
	isTerminal() bool
	String() string
}

type grammar struct {
	rules []*rule
}

type terminal struct {
	value rune
}

func Terminal(r rune) *terminal {
	return &terminal{value: r}
}

func (t *terminal) isTerminal() bool {
	return true
}

func (t *terminal) String() string {
	return fmt.Sprintf("'%c'", t.value)
}

type nonTerminal struct {
	name string
}

func NonTerminal(name string) *nonTerminal {
	return &nonTerminal{name: name}
}

func (t *nonTerminal) isTerminal() bool {
	return false
}

func (t *nonTerminal) String() string {
	return t.name
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
	rightStrings := make([]string, len(t.rule.right) + 1) // +1 for dot
	for i := 0; i < t.dot; i++ {
		rightStrings[i] = t.rule.right[i].String()
	}
	rightStrings[t.dot] = BLACK_CIRCLE
	for i := t.dot; i < len(t.rule.right); i++ {
		rightStrings[i+1] = t.rule.right[i].String()
	}
	return fmt.Sprintf("%v -> %v (%d)", t.rule.left.String(), strings.Join(rightStrings, " "), t.index)
}

// isTerminal checks if the next symbol in the item is a terminal symbol
func (t *eitem) isNextTerminal() bool {
	r := t.rule.right
	if t.dot >= len(r) {
		return false
	}
	return r[t.dot].isTerminal()
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
	rightStrings := make([]string, len(r.right))
	for i, s := range r.right {
		rightStrings[i] = s.String()
	}
	return fmt.Sprintf("%v -> %v", r.left.String(), strings.Join(rightStrings, " "))
}

func Grammar(rules ...*rule) *grammar {
	return &grammar{rules: rules}
}

func Rule(t *nonTerminal, symbols ...symbol) *rule {
	return &rule{
		left: t,
		right: symbols}
}

func (g *grammar) Parse(input string) {
	s := initializeState(g)
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
			// TODO scan
		} else {
			// nonterminal
			// TODO predict
		}
	}
}

// T -> 'a' 'b'
// T -> 'a' T 'b'
// test on ab, aabb aaabbb

// Define grammar
// Parse input using the grammar
