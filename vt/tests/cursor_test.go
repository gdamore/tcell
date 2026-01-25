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
package tests

import (
	"testing"

	"github.com/gdamore/tcell/v3/vt"
)

func TestCursorStyle(t *testing.T) {
	VerifyF(t, vt.SteadyBlock.IsVisible(), "steady block should be visible")
	VerifyF(t, vt.SteadyBlock.Show() == vt.SteadyBlock, "steady block show should be non-op")
	VerifyF(t, !vt.SteadyBlock.Hide().IsVisible(), "steady block hide should be hidden")
	VerifyF(t, vt.SteadyBlock.Hide().Show() == vt.SteadyBlock, "steady block hide and show should be inverse")
	VerifyF(t, !vt.SteadyBlock.IsBlinking(), "steady block should be steady")
	VerifyF(t, vt.SteadyBlock.Blink().IsBlinking(), "steady block blink should work")
	VerifyF(t, vt.SteadyBlock.Blink() == vt.BlinkingBlock, "steady block blink should be blinking block")
	VerifyF(t, vt.SteadyBlock.Steady() == vt.SteadyBlock, "steady block should be steady")
	VerifyF(t, vt.SteadyBlock.Blink().Steady() == vt.SteadyBlock, "steady block blink and steady should be inverse")

	VerifyF(t, vt.SteadyBar.IsVisible(), "steady bar should be visible")
	VerifyF(t, vt.SteadyBar.Show() == vt.SteadyBar, "steady bar show should be non-op")
	VerifyF(t, !vt.SteadyBar.Hide().IsVisible(), "steady bar hide should be hidden")
	VerifyF(t, vt.SteadyBar.Hide().Show() == vt.SteadyBar, "steady bar hide and show should be inverse")
	VerifyF(t, !vt.SteadyBar.IsBlinking(), "steady bar should be steady")
	VerifyF(t, vt.SteadyBar.Blink().IsBlinking(), "steady bar blink should work")
	VerifyF(t, vt.SteadyBar.Blink() == vt.BlinkingBar, "steady bar blink should be blinking bar")
	VerifyF(t, vt.SteadyBar.Steady() == vt.SteadyBar, "steady bar should be steady")
	VerifyF(t, vt.SteadyBar.Blink().Steady() == vt.SteadyBar, "steady bar blink and steady should be inverse")

	VerifyF(t, vt.SteadyUnderline.IsVisible(), "steady underline should be visible")
	VerifyF(t, vt.SteadyUnderline.Show() == vt.SteadyUnderline, "steady underline show should be non-op")
	VerifyF(t, !vt.SteadyUnderline.Hide().IsVisible(), "steady underline hide should be hidden")
	VerifyF(t, vt.SteadyUnderline.Hide().Show() == vt.SteadyUnderline, "steady underline hide and show should be inverse")
	VerifyF(t, !vt.SteadyUnderline.IsBlinking(), "steady underline should be steady")
	VerifyF(t, vt.SteadyUnderline.Blink().IsBlinking(), "steady underline blink should work")
	VerifyF(t, vt.SteadyUnderline.Blink() == vt.BlinkingUnderline, "steady underline blink should be blinking underline")
	VerifyF(t, vt.SteadyUnderline.Steady() == vt.SteadyUnderline, "steady underline should be steady")
	VerifyF(t, vt.SteadyUnderline.Blink().Steady() == vt.SteadyUnderline, "steady underline blink and steady should be inverse")

	VerifyF(t, vt.BlinkingBlock.IsVisible(), "blinking block should be visible")
	VerifyF(t, vt.BlinkingBlock.Show() == vt.BlinkingBlock, "blinking block show should be non-op")
	VerifyF(t, !vt.BlinkingBlock.Hide().IsVisible(), "blinking block hide should be hidden")
	VerifyF(t, vt.BlinkingBlock.Hide().Show() == vt.BlinkingBlock, "blinking block hide and show should be inverse")
	VerifyF(t, vt.BlinkingBlock.IsBlinking(), "blinking block should be blinking")
	VerifyF(t, vt.BlinkingBlock.Blink().IsBlinking(), "blinking block blink should work")
	VerifyF(t, vt.BlinkingBlock.Steady() == vt.SteadyBlock, "blinking block steady should be steady block")
	VerifyF(t, vt.BlinkingBlock.Blink() == vt.BlinkingBlock, "blinking block should be blinking block")
	VerifyF(t, vt.BlinkingBlock.Steady().Blink() == vt.BlinkingBlock, "blinking block steady and blink should be inverse")

	VerifyF(t, vt.BlinkingBar.IsVisible(), "blinking bar should be visible")
	VerifyF(t, vt.BlinkingBar.Show() == vt.BlinkingBar, "blinking bar show should be non-op")
	VerifyF(t, !vt.BlinkingBar.Hide().IsVisible(), "blinking bar hide should be hidden")
	VerifyF(t, vt.BlinkingBar.Hide().Show() == vt.BlinkingBar, "blinking bar hide and show should be inverse")
	VerifyF(t, vt.BlinkingBar.IsBlinking(), "blinking bar should be blinking")
	VerifyF(t, vt.BlinkingBar.Blink().IsBlinking(), "blinking bar blink should work")
	VerifyF(t, vt.BlinkingBar.Steady() == vt.SteadyBar, "blinking bar steady should be steady bar")
	VerifyF(t, vt.BlinkingBar.Blink() == vt.BlinkingBar, "blinking bar should be blinking bar")
	VerifyF(t, vt.BlinkingBar.Steady().Blink() == vt.BlinkingBar, "blinking bar steady and blink should be inverse")

	VerifyF(t, vt.BlinkingUnderline.IsVisible(), "blinking underline should be visible")
	VerifyF(t, vt.BlinkingUnderline.Show() == vt.BlinkingUnderline, "blinking underline show should be non-op")
	VerifyF(t, !vt.BlinkingUnderline.Hide().IsVisible(), "blinking underline hide should be hidden")
	VerifyF(t, vt.BlinkingUnderline.Hide().Show() == vt.BlinkingUnderline, "blinking underline hide and show should be inverse")
	VerifyF(t, vt.BlinkingUnderline.IsBlinking(), "blinking underline should be blinking")
	VerifyF(t, vt.BlinkingUnderline.Blink().IsBlinking(), "blinking underline blink should work")
	VerifyF(t, vt.BlinkingUnderline.Steady() == vt.SteadyUnderline, "blinking underline steady should be steady underline")
	VerifyF(t, vt.BlinkingUnderline.Blink() == vt.BlinkingUnderline, "blinking underline should be blinking underline")
	VerifyF(t, vt.BlinkingUnderline.Steady().Blink() == vt.BlinkingUnderline, "blinking underline steady and blink should be inverse")
}
