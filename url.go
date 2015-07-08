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

// Package url provides functions for parsing, decoding and encoding URLs.
package url

// A URL represents a parsed URL.
type URL struct {
	Scheme   string
	Slashes  bool
	Auth     string
	Host     string
	Port     string
	Path     string
	RawQuery string
	Fragment string
	IPv6     bool
}

// String reassembles the URL into a URL string.
func (u *URL) String() string {
	size := len(u.Path)
	if u.Scheme != "" {
		size += len(u.Scheme) + 1
	}
	if u.Slashes {
		size += 2
	}
	if u.Auth != "" {
		size += len(u.Auth) + 1
	}
	if u.Host != "" {
		size += len(u.Host)
		if u.IPv6 {
			size += 2
		}
	}
	if u.Port != "" {
		size += len(u.Port) + 1
	}
	if u.RawQuery != "" {
		size += len(u.RawQuery) + 1
	}
	if u.Fragment != "" {
		size += len(u.Fragment) + 1
	}
	if size == 0 {
		return ""
	}

	buf := make([]byte, size)
	i := 0
	if u.Scheme != "" {
		i += copy(buf, u.Scheme)
		buf[i] = ':'
		i++
	}
	if u.Slashes {
		buf[i] = '/'
		i++
		buf[i] = '/'
		i++
	}
	if u.Auth != "" {
		i += copy(buf[i:], u.Auth)
		buf[i] = '@'
		i++
	}
	if u.Host != "" {
		if u.IPv6 {
			buf[i] = '['
			i++
			i += copy(buf[i:], u.Host)
			buf[i] = ']'
			i++
		} else {
			i += copy(buf[i:], u.Host)
		}
	}
	if u.Port != "" {
		buf[i] = ':'
		i++
		i += copy(buf[i:], u.Port)
	}
	i += copy(buf[i:], u.Path)
	if u.RawQuery != "" {
		buf[i] = '?'
		i++
		i += copy(buf[i:], u.RawQuery)
	}
	if u.Fragment != "" {
		buf[i] = '#'
		i++
		i += copy(buf[i:], u.Fragment)
	}
	return string(buf)
}
