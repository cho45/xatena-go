package xatena

import (
	"testing"
)

const superPre2TestData = `
### test
::: input
>|perl|
use strict;
use warnings;
warn "helloworld";
||<
::: expected
<pre class="code lang-perl">use strict;
use warnings;
warn &quot;helloworld&quot;;</pre>
</pre>
### test
::: input
>|c|
#include <stdio.h>

int main() {
  puts("helloworld");
}
||<
::: expected
<pre class="code lang-c">#include &lt;stdio.h&gt;

int main() {
  puts(&quot;helloworld&quot;);
}</pre>
</pre>

`

func TestFormat_SuperPre2(t *testing.T) {
	blocks := parseTestBlocksWithDelim(superPre2TestData, "###", ":::")
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}
