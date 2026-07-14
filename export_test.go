package timui

import "gitlab.com/bytewright/gmath/mathi"

func (t *Timui) RunAfterForTest(f func())             { t.runAfter(f) }
func (t *Timui) MoveCursorForTest(dir mathi.Vec2)     { t.moveCursor(dir) }
func (t *Timui) FrontCharForTest(pos mathi.Vec2) rune { return t.front.Get(pos).Char }
func (s *SplitOptions) NumSplitsForTest() int         { return len(s.splits) }
