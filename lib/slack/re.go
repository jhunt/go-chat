package slack

import (
	"regexp"
)

// this should really be in the regexp library
func replaceb(re *regexp.Regexp, b []byte, fn func([][]byte) [][]byte) []byte {
	out := make([]byte, 0, len(b))
	last := 0
	mm := re.FindAllSubmatchIndex(b, -1)

	for _, m := range mm {
		// append the NON-matching segment
		out = append(out, b[last:m[0]]...)
		last = m[1]

		// make subgroups and track non-matching interim segments
		gg := [][]byte{}
		gi := [][2]int{}

		for i := 2; i < len(m); i += 2 {
			gg = append(gg, b[m[i]:m[i+1]])
			gi = append(gi, [2]int{m[i], m[i+1]})
		}

		end := m[0]
		for i, repl := range fn(gg) {
			// append uncaptured segment of match
			out = append(out, b[end:gi[i][0]]...)
			end = gi[i][1]

			// append the replaced segment
			out = append(out, repl...)
		}
		// append uncaptured tail segment of match
		out = append(out, b[end:m[1]]...)
	}

	// append the final NON-matching segment
	out = append(out, b[last:]...)
	return out
}

func replace(re *regexp.Regexp, s string, fn func([]string) []string) string {
	return string(replaceb(re, []byte(s), func(gg [][]byte) [][]byte {
		in := make([]string, len(gg))
		for i := range gg {
			in[i] = string(gg[i])
		}

		out := make([][]byte, len(gg))
		for i, s := range fn(in) {
			out[i] = []byte(s)
		}

		return out
	}))
}
