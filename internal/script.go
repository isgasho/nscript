package internal

import "strings"

type Script struct {
	Functions map[string]*Function
}

func NewScript() *Script {
	var s = &Script{}
	s.Functions = make(map[string]*Function)
	return s
}

func (this *Script) Add(p *Function) {
	if p == nil {
		return
	}
	this.Functions[strings.ToUpper(p.Name)] = p
}

func (this *Script) Take(name string) *Function {
	name = strings.ToUpper(name)
	var f = this.Functions[name]
	delete(this.Functions, name)
	return f
}
