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
	"errors"
	"strings"

	"github.com/opennota/byteutil"
)

// ErrMissingScheme error is returned by Parse if the passed URL starts with a colon.
var ErrMissingScheme = errors.New("missing protocol scheme")

var (
	cs1 [256]bool
	cs2 [256]bool
	cs3 [256]bool
)

func init() {
	for _, b := range "+-.0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		cs1[b] = true
	}
	for _, b := range "#/?" {
		cs2[b] = true
	}
	for _, b := range "\t\r\n \"#%'/;<>?\\^`{|}" {
		cs3[b] = true
	}
}

var slashedProtocol = map[string]bool{
	"http":   true,
	"https":  true,
	"ftp":    true,
	"gopher": true,
	"file":   true,
}

func split(s string, c byte) (string, string) {
	i := strings.IndexByte(s, c)
	if i < 0 {
		return s, ""
	}
	return s[:i], s[i+1:]
}

func findScheme(s string) (int, error) {
	if s == "" {
		return 0, nil
	}

	b := s[0]
	if b == ':' {
		return 0, ErrMissingScheme
	}
	if !byteutil.IsLetter(b) {
		return 0, nil
	}

	for i := 1; i < len(s); i++ {
		b := s[i]
		switch {
		case cs1[b]:
			// do nothing
		case b == ':':
			return i, nil
		default:
			return 0, nil
		}
	}

	return 0, nil
}

// Parse parses rawurl into a URL structure.
func Parse(rawurl string) (*URL, error) {
	n, err := findScheme(rawurl)
	if err != nil {
		return nil, err
	}

	var url URL
	rest := rawurl
	hostless := false
	if n > 0 {
		url.Scheme, rest = byteutil.ToLower(rest[:n]), rest[n+1:]
		if url.Scheme == "javascript" {
			hostless = true
		}
	}

	if !hostless && strings.HasPrefix(rest, "//") {
		url.Slashes, rest = true, rest[2:]
	}

	if !hostless && (url.Slashes || (url.Scheme != "" && !slashedProtocol[url.Scheme])) {
		hostEnd := byteutil.IndexAnyTable(rest, &cs2)
		atSign := -1
		i := hostEnd
		if i == -1 {
			i = len(rest) - 1
		}
		for i >= 0 {
			if rest[i] == '@' {
				atSign = i
				break
			}
			i--
		}

		if atSign != -1 {
			url.Auth, rest = rest[:atSign], rest[atSign+1:]
		}

		hostEnd = byteutil.IndexAnyTable(rest, &cs3)
		if hostEnd == -1 {
			hostEnd = len(rest)
		}
		if hostEnd > 0 && hostEnd < len(rest) && rest[hostEnd-1] == ':' {
			hostEnd--
		}
		host := rest[:hostEnd]

		if len(host) > 1 {
			b := host[hostEnd-1]
			if byteutil.IsDigit(b) {
				for i := len(host) - 2; i >= 0; i-- {
					b := host[i]
					if b == ':' {
						url.Host, url.Port = host[:i], host[i+1:]
						break
					}
					if !byteutil.IsDigit(b) {
						break
					}
				}
			} else if b == ':' {
				host = host[:hostEnd-1]
				hostEnd--
			}
		}
		if url.Port == "" {
			url.Host = host
		}
		rest = rest[hostEnd:]

		if ipv6 := len(url.Host) > 2 &&
			url.Host[0] == '[' &&
			url.Host[len(url.Host)-1] == ']'; ipv6 {
			url.Host = url.Host[1 : len(url.Host)-1]
			url.IPv6 = true
		} else if i := strings.IndexByte(url.Host, ':'); i >= 0 {
			url.Host, rest = url.Host[:i], url.Host[i:]+rest
		}
	}

	rest, url.Fragment = split(rest, '#')
	url.Path, url.RawQuery = split(rest, '?')

	return &url, nil
}
