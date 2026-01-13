// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package vt

import "testing"

func TestCursorStyle(t *testing.T) {
	verifyF(t, SteadyBlock.IsVisible(), "steady block should be visible")
	verifyF(t, SteadyBlock.Show() == SteadyBlock, "steady block show should be non-op")
	verifyF(t, !SteadyBlock.Hide().IsVisible(), "steady block hide should be hidden")
	verifyF(t, SteadyBlock.Hide().Show() == SteadyBlock, "steady block hide and show should be inverse")
	verifyF(t, !SteadyBlock.IsBlinking(), "steady block should be steady")
	verifyF(t, SteadyBlock.Blink().IsBlinking(), "steady block blink should work")
	verifyF(t, SteadyBlock.Blink() == BlinkingBlock, "steady block blink should be blinking block")
	verifyF(t, SteadyBlock.Steady() == SteadyBlock, "steady block should be steady")
	verifyF(t, SteadyBlock.Blink().Steady() == SteadyBlock, "steady block blink and steady should be inverse")

	verifyF(t, SteadyBar.IsVisible(), "steady bar should be visible")
	verifyF(t, SteadyBar.Show() == SteadyBar, "steady bar show should be non-op")
	verifyF(t, !SteadyBar.Hide().IsVisible(), "steady bar hide should be hidden")
	verifyF(t, SteadyBar.Hide().Show() == SteadyBar, "steady bar hide and show should be inverse")
	verifyF(t, !SteadyBar.IsBlinking(), "steady bar should be steady")
	verifyF(t, SteadyBar.Blink().IsBlinking(), "steady bar blink should work")
	verifyF(t, SteadyBar.Blink() == BlinkingBar, "steady bar blink should be blinking bar")
	verifyF(t, SteadyBar.Steady() == SteadyBar, "steady bar should be steady")
	verifyF(t, SteadyBar.Blink().Steady() == SteadyBar, "steady bar blink and steady should be inverse")

	verifyF(t, SteadyUnderline.IsVisible(), "steady underline should be visible")
	verifyF(t, SteadyUnderline.Show() == SteadyUnderline, "steady underline show should be non-op")
	verifyF(t, !SteadyUnderline.Hide().IsVisible(), "steady underline hide should be hidden")
	verifyF(t, SteadyUnderline.Hide().Show() == SteadyUnderline, "steady underline hide and show should be inverse")
	verifyF(t, !SteadyUnderline.IsBlinking(), "steady underline should be steady")
	verifyF(t, SteadyUnderline.Blink().IsBlinking(), "steady underline blink should work")
	verifyF(t, SteadyUnderline.Blink() == BlinkingUnderline, "steady underline blink should be blinking underline")
	verifyF(t, SteadyUnderline.Steady() == SteadyUnderline, "steady underline should be steady")
	verifyF(t, SteadyUnderline.Blink().Steady() == SteadyUnderline, "steady underline blink and steady should be inverse")

	verifyF(t, BlinkingBlock.IsVisible(), "blinking block should be visible")
	verifyF(t, BlinkingBlock.Show() == BlinkingBlock, "blinking block show should be non-op")
	verifyF(t, !BlinkingBlock.Hide().IsVisible(), "blinking block hide should be hidden")
	verifyF(t, BlinkingBlock.Hide().Show() == BlinkingBlock, "blinking block hide and show should be inverse")
	verifyF(t, BlinkingBlock.IsBlinking(), "blinking block should be blinking")
	verifyF(t, BlinkingBlock.Blink().IsBlinking(), "blinking block blink should work")
	verifyF(t, BlinkingBlock.Steady() == SteadyBlock, "blinking block steady should be steady block")
	verifyF(t, BlinkingBlock.Blink() == BlinkingBlock, "blinking block should be blinking block")
	verifyF(t, BlinkingBlock.Steady().Blink() == BlinkingBlock, "blinking block steady and blink should be inverse")

	verifyF(t, BlinkingBar.IsVisible(), "blinking bar should be visible")
	verifyF(t, BlinkingBar.Show() == BlinkingBar, "blinking bar show should be non-op")
	verifyF(t, !BlinkingBar.Hide().IsVisible(), "blinking bar hide should be hidden")
	verifyF(t, BlinkingBar.Hide().Show() == BlinkingBar, "blinking bar hide and show should be inverse")
	verifyF(t, BlinkingBar.IsBlinking(), "blinking bar should be blinking")
	verifyF(t, BlinkingBar.Blink().IsBlinking(), "blinking bar blink should work")
	verifyF(t, BlinkingBar.Steady() == SteadyBar, "blinking bar steady should be steady bar")
	verifyF(t, BlinkingBar.Blink() == BlinkingBar, "blinking bar should be blinking bar")
	verifyF(t, BlinkingBar.Steady().Blink() == BlinkingBar, "blinking bar steady and blink should be inverse")

	verifyF(t, BlinkingUnderline.IsVisible(), "blinking underline should be visible")
	verifyF(t, BlinkingUnderline.Show() == BlinkingUnderline, "blinking underline show should be non-op")
	verifyF(t, !BlinkingUnderline.Hide().IsVisible(), "blinking underline hide should be hidden")
	verifyF(t, BlinkingUnderline.Hide().Show() == BlinkingUnderline, "blinking underline hide and show should be inverse")
	verifyF(t, BlinkingUnderline.IsBlinking(), "blinking underline should be blinking")
	verifyF(t, BlinkingUnderline.Blink().IsBlinking(), "blinking underline blink should work")
	verifyF(t, BlinkingUnderline.Steady() == SteadyUnderline, "blinking underline steady should be steady underline")
	verifyF(t, BlinkingUnderline.Blink() == BlinkingUnderline, "blinking underline should be blinking underline")
	verifyF(t, BlinkingUnderline.Steady().Blink() == BlinkingUnderline, "blinking underline steady and blink should be inverse")
}
