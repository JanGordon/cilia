package addons

import "fmt"

type Addon struct {
	OpeningToken string
	ClosingToken string
	Name         string

	ContentModifier func(string) string
	PureSSR         bool
	JS              string
	Open            bool
	StartIndex      int
	Content         string
}

var Addons []*Addon

func initNewAddon(a *Addon) {
	Addons = append(Addons, a)
}

func ReplaceAddons(page string) string {
	var openTokens []*Addon
	for c := 1; c < len(page); c += 1 {
		// check if the sequence fits any addons
		for _, addon := range Addons {

			if string(page[c-1])+string(page[c]) == addon.OpeningToken {
				// found a match
				fmt.Println(addon)
				newToken := *addon
				newToken.StartIndex = c + 1
				openTokens = append(openTokens, &newToken)
			} else if string(page[c-1])+string(page[c]) == addon.ClosingToken {
				// found a closing token
				fmt.Println("FOund a closing token!!!!!")
				for i := len(openTokens) - 1; i >= 0; i-- {
					if openTokens[i].OpeningToken == addon.OpeningToken && addon.Open {
						t := openTokens[i]
						t.Open = false
						t.Content = page[openTokens[i].StartIndex : c-1]
						fmt.Println(openTokens[i].Content)

						modifiedContent := t.ContentModifier(t.Content)
						page = page[:t.StartIndex-2] + fmt.Sprintf("<%v>%v</%v>", t.Name, modifiedContent, t.Name) + page[c+1:]
						// need to know if it should be rerun or only ssr
						// here, convert each token to js but dont run it.
						// then work out how they fit together so can be run in same context
						// once found we need to loop over forward and store the contents in a dom tree not on the page.
					}
				}
			}
		}
	}
	return page
}
