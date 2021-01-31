package nscript

import (
	"github.com/smartwalle/nscript/internal"
	"strings"
)

type Function struct {
	name     string
	segments []*Segment
}

func NewFunction(name string) *Function {
	var f = &Function{}
	f.name = name
	return f
}

func (this *Function) parse(lines []string) error {
	var sLines []string
	for idx, line := range lines {
		if idx != 0 && strings.HasPrefix(line, internal.KeyIf) {
			// 如果不是第一行，并且又发现了 #IF
			// 则表示一个代码片断结束了，需要开启新的代码片断
			var nSegment = NewSegment()
			if err := nSegment.parse(sLines); err != nil {
				return err
			}
			this.segments = append(this.segments, nSegment)
			sLines = nil
		}
		sLines = append(sLines, line)
	}
	if len(sLines) > 0 {
		var nSegment = NewSegment()
		if err := nSegment.parse(sLines); err != nil {
			return err
		}
		this.segments = append(this.segments, nSegment)
		sLines = nil
	}
	return nil
}

func (this *Function) exec(ctx Context) ([]string, string, error) {
	for _, seg := range this.segments {
		var ok, err = seg.check(ctx)
		if err != nil {
			// 若有错误，返回错误
			return nil, "", err
		}
		if ok {
			// 执行 Action
			return seg.execAction(ctx)
		} else if seg.hasElse() {
			// 执行 ElseAction
			return seg.execElseAction(ctx)
		}
	}
	return nil, "", nil
}
