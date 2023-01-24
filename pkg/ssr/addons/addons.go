package addons

type Addon struct {
	OpeningToken    string
	ClosingToken    string
	Name            string
	ContentModifier func(string) string
}

var Addons []*Addon

func initNewAddon(a *Addon) {
	Addons = append(Addons, a)
}

func ReplaceAddons(page string) string {
	return page
}
