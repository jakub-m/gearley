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
	isMatchingTerminal(rune) bool
	// s and input are slices of the full state set and the full input.
	//modifyStateSet(s []*stateSet, input []rune)
}

type terminal struct {
	value rune
}

func Terminal(r rune) terminal {
	return terminal{value: r}
}

func (t terminal) isTerminal() bool {
	return true
}

func (t terminal) String() string {
	return fmt.Sprintf("'%c'", t.value)
}

func (t terminal) isMatchingTerminal(r rune) bool {
	return r == t.value
}

type nonTerminal struct {
	name string
}

func NonTerminal(name string) nonTerminal {
	return nonTerminal{name: name}
}

func (t nonTerminal) isTerminal() bool {
	return false
}

func (t nonTerminal) String() string {
	return t.name
}

func (t nonTerminal) isMatchingTerminal(r rune) bool {
	return false
}

type stateSet struct {
	items []*eitem
	itemSet map[eitem]bool
}

func newStateSet() *stateSet {
	return &stateSet{
		items: []*eitem{},
		itemSet: make(map[eitem]bool),
		}
}

func (s *stateSet) length() int {
	return len(s.items)
}

func (s *stateSet) putItem(item *eitem) {
	// Add items only if they are not already in the item set
	if _, ok := s.itemSet[*item]; ok {
		return
	}
	s.itemSet[*item] = true
	s.items = append(s.items, item)
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
	rightStrings := make([]string, len(t.rule.right))
	for i, r := range t.rule.right {
		rightStrings[i] = r.String()
	}
	return fmt.Sprintf("%v -> %v%v%v (%d)",
		t.rule.left.String(),
		strings.Join(rightStrings[0:t.dot], " "),
		BLACK_CIRCLE,
		strings.Join(rightStrings[t.dot:], " "),
		t.index,
		)
}

func (t *eitem) isCompleted() bool {
	return t.dot == t.rule.length()
}

// check if is terminal and if is matching
func (t *eitem) isNextMatchingTerminal(nextRune rune) bool {
	s := t.getSymbolAt(t.dot)
	return s.isMatchingTerminal(nextRune)
}

func (t *eitem) getSymbolAt(i int) symbol {
	return t.rule.right[i]
}

func (t *eitem) getNext() symbol {
	// TODO edge case when at the end
	return t.rule.right[t.dot]
}

func (t *eitem) createNext() *eitem {
	// TODO edge case when dot is last
	return &eitem{
		dot: t.dot + 1,
		rule: t.rule,
		index: t.index}
}

// state is the highest-level state of the parser.
type state []*stateSet;

func (s *stateSet) String() string {
	return fmt.Sprint(s.items)
}

type rule struct {
	left nonTerminal
	right []symbol
}

func (r *rule) length() int {
	return len(r.right)
}

func (r *rule) String() string {
	rightStrings := make([]string, len(r.right))
	for i, s := range r.right {
		rightStrings[i] = s.String()
	}
	return fmt.Sprintf("%v -> %v", r.left.String(), strings.Join(rightStrings, " "))
}

func Rule(t nonTerminal, symbols ...symbol) *rule {
	return &rule{
		left: t,
		right: symbols}
}

type grammar struct {
	rules []*rule
}

func Grammar(rules ...*rule) *grammar {
	return &grammar{rules: rules}
}

func (g *grammar) Parse(input string) {
	inputRunes := stringToRunes(input)
	st := initializeState(g, inputRunes)
	// the current index in the state 'st' that is being processed - S(stateIndex)
	stateIndex := 0
	// outter loop
	for stateIndex <= len(input) {
		set := (*st)[stateIndex]
		fmt.Println("NOW SET S(", stateIndex, ")", set)
		i := 0
		// inner loop
		for i < set.length() {
			item := set.items[i]
			fmt.Println(i, item)
			i++

			fmt.Println("NOW item", item)

			if item.isCompleted() {
				fmt.Println("completed!")
				// HERE TODO implement what if completed
				continue
			}
			if item.isNextMatchingTerminal(inputRunes[stateIndex]) {
				nextItem := item.createNext()
				// create next item
				// add it to the next stateSet
				fmt.Println("next item", nextItem)
				// TODO edge case when last stateIndex
				nextSet := (*st)[stateIndex+1]
				nextSet.putItem(nextItem)
			}
			if ok := item.getNext().isTerminal(); !ok {
				nextSymbol := item.getNext().(nonTerminal)
				// Find all the rules for the symbol put those rules to the current set
				fmt.Println("NON TERMINAL")
				for _, r := range g.getRulesForSymbol(nextSymbol) {
					nextItem := &eitem{
						rule: r,
						dot: 0,
						index: stateIndex,
						}
					set.putItem(nextItem)
				}
			}
		}

		fmt.Println("S", *st, "\n")
		stateIndex++
	}
}

func (g *grammar) getRulesForSymbol(s symbol) []*rule {
	found := []*rule{}
	for _, r := range g.rules {
		if r.left == s {
		    found = append(found, r)
		}
	}
	return found
}

func initializeState(g *grammar, runes []rune) *state {
	sets := make([]*stateSet, len(runes) + 1)
	for i := range sets {
		sets[i] = newStateSet()
	}
	sets[0] = newStateSetFromRules(g.rules)
	s := state(sets)
	return &s
}

func newStateSetFromRules(rules []*rule) *stateSet {
	items := make([]*eitem, len(rules))
	for i, r := range rules {
		items[i] = &eitem{rule: r, dot: 0, index: 0}
	}
	ss := newStateSet()
	ss.items = items
	return ss
}

func stringToRunes(input string) []rune {
	runes := []rune{}
	for _, r := range input {
		runes = append(runes, r)
	}
	return runes
}
