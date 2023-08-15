// Copyright 2022 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package tcell

import (
	"fmt"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func ExampleRegisterEncoding() {
	RegisterEncoding("GBK", simplifiedchinese.GBK)
	enc := GetEncoding("GBK")
	glyph, _ := enc.NewDecoder().Bytes([]byte{0x82, 0x74})
	fmt.Println(string(glyph))
	// Output: 倀
}
