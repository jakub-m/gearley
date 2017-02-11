package gearley
import (
	"testing"
)

var T = NonTerminal("T")
var A = Terminal('a')
var B = Terminal('b')

func TestAbba(t *testing.T) {
	g := Grammar(
	    Rule(T, A, B),	    // T -> 'a' 'b'
	    Rule(T, A, T, B),	    // T -> 'a' T 'b'
	)

	g.Parse("abba")
}
