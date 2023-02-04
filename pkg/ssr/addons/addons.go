package addons

import (
	"fmt"

	wasmer "github.com/wasmerio/wasmer-go/wasmer"
	"rogchap.com/v8go"
)

type Addon struct {
	OpeningToken string
	ClosingToken string
	Name         string

	ContentModifier func(string, v8go.Context) string
	PureSSR         bool
	Open            bool
	StartIndex      int
	Content         string
}

type ExternalAddonCfg struct {
	OpeningToken string
	ClosingToken string
	Name         string
	PureSSR      bool
}

var Addons []*Addon
var store *wasmer.Store

func initNewAddon(a *Addon) {
	Addons = append(Addons, a)
}

func ReplaceAddons(page string, ctx v8go.Context, ssr bool) string {
	var openTokens []*Addon
	for c := 1; c < len(page); c += 1 {
		// check if the sequence fits any addons
		for _, addon := range Addons {

			if string(page[c-1])+string(page[c]) == addon.OpeningToken {
				// found a match
				newToken := *addon
				newToken.StartIndex = c + 1
				openTokens = append(openTokens, &newToken)
			} else if string(page[c-1])+string(page[c]) == addon.ClosingToken {
				// found a closing token
				fmt.Println("FOund a closing token!!!!!")
				for i := len(openTokens) - 1; i >= 0; i-- {
					if openTokens[i].PureSSR {
						if !ssr {
							continue
						}
						// this means it shouldnt be built as it is build
					}
					if openTokens[i].OpeningToken == addon.OpeningToken && addon.Open {
						t := openTokens[i]
						t.Open = false
						t.Content = page[openTokens[i].StartIndex : c-1]
						fmt.Println("Computing token")

						//wasm stuff
						// store := wasmer.NewStore(wasmer.NewEngine())
						// module, _ := wasmer.NewModule(store, t.ContentModifier)

						// wasiEnv, _ := wasmer.NewWasiStateBuilder("wasi-program").
						// 	// Choose according to your actual situation
						// 	// Argument("--foo").
						// 	// Environment("ABC", "DEF").
						// 	// MapDirectory("./", ".").
						// 	Finalize()
						// importObject, err := wasiEnv.GenerateImportObject(store, module)
						// check(err)

						// instance, err := wasmer.NewInstance(module, importObject)
						// check(err)

						// start, err := instance.Exports.GetWasiStartFunction()
						// check(err)
						// start()

						// modifier, err := instance.Exports.GetFunction("Modifier")
						// check(err)
						modifiedContent := t.ContentModifier(t.Content, ctx)
						fmt.Println(modifiedContent)
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

type config struct {
	Version        float32
	Name           string
	DirectAddons   []string
	IndirectAddons []string
}

func init() {

	// configFile, err := os.ReadFile(filepath.Join(global.ProjectRoot, "stem.toml"))
	// if err != nil {
	// 	panic(err)
	// }
	// var cfg config
	// err = toml.Unmarshal(configFile, &cfg)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(cfg.Name)
	// for _, directAddon := range cfg.DirectAddons {
	// 	resp1, err := http.Get(directAddon + "/raw/main/addon.wasm")
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	resp2, err := http.Get(directAddon + "/raw/main/config.toml")
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	wasmBytes, err := ioutil.ReadAll(resp1.Body)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	identifier, err := ioutil.ReadAll(resp2.Body)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	var cfg ExternalAddonCfg
	// 	err = toml.Unmarshal(identifier, &cfg)
	// 	initNewAddon(&Addon{
	// 		OpeningToken:    cfg.OpeningToken,
	// 		ClosingToken:    cfg.ClosingToken,
	// 		Name:            cfg.Name,
	// 		ContentModifier: wasmBytes,
	// 		PureSSR:         cfg.PureSSR,
	// 		Open:            true,
	// 		StartIndex:      0,
	// 		Content:         "",
	// 	})
	// 	fmt.Println("initialized new addon:", cfg.Name)
	// }

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func runWasm(wasmFile []byte, stringArg string, funcName string) string {
	return ""
}
