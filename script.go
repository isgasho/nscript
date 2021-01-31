package nscript

import (
	"fmt"
	"github.com/smartwalle/nscript/internal"
	"strings"
)

type Script struct {
	pages map[string]*Page
}

func NewScript(file string) (*Script, error) {
	var iScript, err = internal.LoadFile(file)
	if err != nil {
		return nil, err
	}

	var nScript = &Script{}
	nScript.pages = make(map[string]*Page)

	for _, iPage := range iScript.Pages {
		var nPage = NewPage(iPage.Key)
		if err = nPage.parse(iPage.Lines); err != nil {
			return nil, err
		}
		nScript.pages[nPage.key] = nPage
	}

	return nScript, nil
}

func (this *Script) Exec(key string, ctx Context) ([]string, error) {
	key = strings.ToUpper(key)
	var page = this.pages[key]
	if page == nil {
		return nil, fmt.Errorf("%s not found", key)
	}

	var says, gotoKey, err = page.exec(ctx)
	if err != nil {
		return nil, err
	}

	if gotoKey != "" {
		return this.Exec(gotoKey, ctx)
	}

	return says, nil
}
