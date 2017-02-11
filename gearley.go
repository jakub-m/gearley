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

func (t *terminal) isMatchingTerminal(r rune) bool {
	return r == t.value
}

//func (t *terminal) modifyStateSet(s []*stateSet, input []rune) {
//	fmt.Printf("%c == %c\n", input[0], t.value)
//	if input[0] == t.value {
//		// There is a match!
//		s[0].putState()
//	}
//}

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

func (t *nonTerminal) isMatchingTerminal(r rune) bool {
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

// check if is terminal and if is matching
func (t *eitem) isNextMatchingTerminal(nextRune rune) bool {
	s := t.getSymbolAt(t.dot)
	return s.isMatchingTerminal(nextRune)
}

func (t *eitem) getSymbolAt(i int) symbol {
	return t.rule.right[i]
}

func (t *eitem) createNext() *eitem {
	// TODO edge case when dot is last
	return &eitem{
		dot: t.dot + 1,
		rule: t.rule,
		index: t.index}
}

// Produce a new item. The produced item depends if the next symbol is Terminal or Non Terminal.
// newItem - the new item
// offset - offset where the new item should be placed
// ok - if the item was created
//func (t *eitem) produceNewItem() (newItem *eitem, offset int, ok bool) {
//	r := t.rule.right
//	if r[t.dot].isTerminal() {
//		// TODO produce new terminal
//	} else {
//		// TOOD produce from non-terminal
//		// !! need the whole grammar to produce the new item !!
//	}
//}


//// isTerminal checks if the next symbol in the item is a terminal symbol
//func (t *eitem) isNextTerminal() bool {
//	r := t.rule.right
//	if t.dot >= len(r) {
//		return false
//	}
//	return r[t.dot].isTerminal()
//}
//
//// s - slice of the state set. The slice starts with the state NEXT after the one currently interpreted.
//// input - slice of the input. The slice corresponds to `s`.
//func (t *eitem) modifyStateSet(s []*stateSet, input []rune) {
//	// TOOD edge case - out of bounds
//	t.rule.right[t.dot].modifyStateSet(s, input)
//}

// state is the highest-level state of the parser.
type state []*stateSet;

func (s *stateSet) String() string {
	return fmt.Sprint(s.items)
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
	inputRunes := stringToRunes(input)
	st := initializeState(g, inputRunes)
	// the current index in the state 'st' that is being processed - S(stateIndex)
	stateIndex := 0

	set := (*st)[stateIndex]
	i := 0
	for i < set.length() {
		item := set.items[i]
		fmt.Println(i, item)
		i++

		// TODO check completed
		if item.isNextMatchingTerminal(inputRunes[stateIndex]) {
			nextItem := item.createNext()
			// create next item
			// add it to the next stateSet
			fmt.Println("next item", nextItem)
			// TODO edge case when last stateIndex
			(*st)[stateIndex+1].putItem(nextItem)
		}
	}

	//fmt.Println(input)
	//st.processStateSet(0, inputRunes)
	fmt.Println(*st)
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

// k - index of the state set to process
// inputRunes - the input as runes
func (s *state) processStateSet(k int, input []rune) {
	fmt.Printf("==== %d ====\n", k)
	if k >= len(*s) {
		panic(fmt.Sprintf("out of bound: %d, len is %d", k, len(*s)))
	}
	set := (*s)[k]
	// This operation mutates set, so set.length() can increase in each loop.
	i := 0
	for i < set.length() {
		item := set.items[i]
		fmt.Println(i, item)
		i++
		// TODO check if is completed
		// For terminal, if scan matches - produce new item.
		// For non-terminal - produce new item.
	//	if newItem, offset, ok := item.produceNewItem(); ok {
	//	}
		//if item.isNextTerminal() {
		//	// Scan.
		//	t.modifyStateSet((*s)[k+1:], input[k:])
		//	// TODO scan
		//} else {
		//	// nonterminal
		//	// TODO predict
		//}
	}
}
