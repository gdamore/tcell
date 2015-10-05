// Copyright 2015 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcell

import (
	"sync"

	"golang.org/x/text/encoding"
)

var encodings map[string]encoding.Encoding
var encodingLk sync.Mutex

// RegisterEncoding may be called by the application to register an encoding.
// The presence of additional encodings will facilitate application usage with
// terminal environments where the I/O subsystem does not support Unicode.
// Please see the Go documentation for golang.org/x/text/encoding -- most of
// the common ones exist already as stock variables.  For example, ISO8859-15
// can be registered using the following code:
//
//   import "golang.org/x/text/encoding/charmap"
//
//     ...
//     RegisterEncoding("ISO8859-15", charmap.ISO8859_15)
//
// Aliases can be registered as well, for example "8859-15" could be an alias
// for "ISO8859-15".
//
// For POSIX systems, the tcell pacakge will check the environment variables
// LC_ALL, LC_CTYPE,  and LANG (in that order) to determine the character set.
// These are expected to have the following pattern:
//
//	 $language[.$codeset[@$variant]

// We extract only the $codeset part, which will usually be something like
// UTF-8 or ISO8859-15 or KOI8-R.  Note that if the locale is either "POSIX"
// or "C", then we assume US-ASCII (the POSIX 'portable character set'
// and assume all other characters are somehow invalid.)
//
// On Windows systems, the Console is assumed to be UTF-16LE.  As we
// communicate with the console subsystem using UTF-16LE, no conversions are
// necessary.  So none of this is required for Windows systems.
//
// Modern POSIX systems and terminal emulators may use UTF-8, and for those
// systems, this API is also unnecessary.  For example, Darwin (MacOS X) and
// modern Linux running modern xterm generally will out of the box without
// any of this.  Use of UTF-8 is recommended when possible, as it saves
// quite a lot processing overhead.
//
// Note that some encodings are quite large (for example GB18030 which is a
// superset of Unicode) and so the application size can be expected ot
// increase quite a bit as each encoding is added.  The East Asian encodings
// have been seen to add 100-200K per encoding to the application size.
//
func RegisterEncoding(name string, enc encoding.Encoding) {
	encodingLk.Lock()
	if encodings == nil {
		encodings = make(map[string]encoding.Encoding)
	}
	encodings[name] = enc
	encodingLk.Unlock()
}

// GetEncoding is used by Screen implementors who want to locate an encoding
// for the given character set name.  Note that this will return nil for
// either the Unicode (UTF-8) or ASCII encodings, since we don't use
// encodings for them but instead have our own native methods.
func GetEncoding(name string) encoding.Encoding {
	encodingLk.Lock()
	defer encodingLk.Unlock()
	if enc, ok := encodings[name]; ok {
		return enc
	}
	return nil
}
