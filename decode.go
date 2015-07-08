// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package url

import (
	"bytes"
	"unicode/utf8"

	"github.com/opennota/byteutil"
)

func advance(s string, pos int) (byte, int) {
	if pos >= len(s) {
		return 0, len(s) + 1
	}
	if s[pos] != '%' {
		return s[pos], pos + 1
	}
	if pos+2 < len(s) &&
		byteutil.IsHexDigit(s[pos+1]) &&
		byteutil.IsHexDigit(s[pos+2]) {
		return byteutil.Unhex(s[pos+1])<<4 | byteutil.Unhex(s[pos+2]), pos + 3
	}
	return '%', pos + 1
}

// Decode decodes a percent-encoded URL.
// Invalid percent-encoded sequences are left as is.
// Invalid UTF-8 sequences are replaced with U+FFFD.
func Decode(rawurl string) string {
	var buf bytes.Buffer
	i := 0
	const replacement = "\xEF\xBF\xBD"
outer:
	for i < len(rawurl) {
		r, rlen := utf8.DecodeRuneInString(rawurl[i:])
		if r == '%' && i+2 < len(rawurl) &&
			byteutil.IsHexDigit(rawurl[i+1]) &&
			byteutil.IsHexDigit(rawurl[i+2]) {
			b := byteutil.Unhex(rawurl[i+1])<<4 | byteutil.Unhex(rawurl[i+2])
			if b < 0x80 {
				buf.WriteByte(b)
				i += 3
				continue
			}
			var n int
			if b&0xe0 == 0xc0 {
				n = 1
			} else if b&0xf0 == 0xe0 {
				n = 2
			} else if b&0xf8 == 0xf0 {
				n = 3
			}
			if n == 0 {
				buf.WriteString(replacement)
				i += 3
				continue
			}
			rb := make([]byte, n+1)
			rb[0] = b
			j := i + 3
			for k := 0; k < n; k++ {
				b, j = advance(rawurl, j)
				if j > len(rawurl) || b&0xc0 != 0x80 {
					buf.WriteString(replacement)
					i += 3
					continue outer
				}
				rb[k+1] = b
			}
			r, _ := utf8.DecodeRune(rb)
			buf.WriteRune(r)
			i = j
			continue
		}
		buf.WriteRune(r)
		i += rlen
	}
	return buf.String()
}
