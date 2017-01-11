// Tests for the text package
package text

import (
	"testing"
)

type t struct {
	in  string
	out string
}

var newlineTests = []t{
	{
		in: `mypara
		`,
		out: `<p>mypara</p>`,
	},
}

// TestConvertNewlines tests newlines -> <p></p>
func TestConvertNewlines(t *testing.T) {
	for _, v := range newlineTests {
		r := ConvertNewlines(v.in)
		if r != v.out {
			t.Fatalf("🔥 Failed to transform newlines\n\twanted:%s\n\tgot:%s\n", v.out, r)
		}
	}
}

var activateLinksTests = []t{
	{
		in:  `<a href="https://google.com?foo=bar#123">https://google.com</a>`,
		out: `<a href="https://google.com?foo=bar#123">https://google.com</a>`,
	},
	{
		in:  `https://google.com`,
		out: `<a href="https://google.com">https://google.com</a>`,
	},
	{
		in: `🗑🔥
		
		https://google.com`,
		out: `🗑🔥
		
		<a href="https://google.com">https://google.com</a>`,
	},
	{
		in:  ` https://google.com      `,
		out: ` <a href="https://google.com">https://google.com</a>      `,
	},
	{
		in: `  https://news.ycombinator.com/item?id=13213902?foo=bar#fragment
    `,
		out: `  <a href="https://news.ycombinator.com/item?id=13213902?foo=bar#fragment">https://news.ycombinator.com/item?id=13213902?foo=bar#fragment</a>
    `,
	},
	{
		in:  ` @tester!`,
		out: ` <a href="/u/tester">@tester</a>!`,
	},
	{ // Test medium-style urls with @ usernames
		in:  `https://medium.com/@taylorotwell/measuring-code-complexity-64356da605f9#.wayfi5mch`,
		out: `<a href="https://medium.com/@taylorotwell/measuring-code-complexity-64356da605f9#.wayfi5mch">https://medium.com/@taylorotwell/measuring-code-complexity-64356da605f9#.wayfi5mch</a>`,
	},
	{ // Note this will be escaped when put into the template
		in:  `https://d>medium.com/#.wayfi5mch`,
		out: `<a href="https://d">https://d</a>>medium.com/#.wayfi5mch`,
	},
}

// TestConvertLinks tests links and usernames are converted
func TestConvertLinks(t *testing.T) {
	for _, v := range activateLinksTests {
		r := ConvertLinks(v.in)
		if r != v.out {
			t.Fatalf("🔥 Failed to transform links\n\twanted:%s\n\tgot:%s\n", v.out, r)
		}
	}
}
