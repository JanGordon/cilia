package addons

import (
	"fmt"

	"github.com/JanGordon/cilia-framework/pkg/page"
	"rogchap.com/v8go"
)

type Addon struct {
	OpeningToken string
	ClosingToken string
	Name         string

	ContentModifier func(string, v8go.Context, string) (string, string)
	PureSSR         bool
	Open            bool
	StartIndex      int
	Content         string
}

type ScriptNode struct {
	Open       bool
	Content    string
	StartIndex int
}

type ExternalAddonCfg struct {
	OpeningToken string
	ClosingToken string
	Name         string
	PureSSR      bool
}

var Addons []*Addon

func initNewAddon(a *Addon) {
	Addons = append(Addons, a)
}

func ReplaceAddons(document *page.Page, ssr bool, jsFile *page.JsFile) page.Page {
	// svae previos
	var page = document.TextContents
	var ctx = document.Js.Ctx
	var openTokens []*Addon
	var ScriptNodes []*ScriptNode

	fmt.Println(page)
	for c := 1; c < len(page); c += 1 {
		if c+7 <= len(page) && string(page[c-1:c+7]) == "<script>" {
			fmt.Println("starting script")
			ScriptNodes = append(ScriptNodes, &ScriptNode{true, "", c + 7})
		} else if c+8 <= len(page) && string(page[c-1:c+8]) == "</script>" {
			for i := len(ScriptNodes) - 1; i >= 0; i-- {
				if ScriptNodes[i].Open {
					fmt.Println("running script", page[ScriptNodes[i].StartIndex:c-2])
					ScriptNodes[i].Open = false
					ctx.RunScript(page[ScriptNodes[i].StartIndex:c-2], "inlinejsscript.js")
					break
				}
			}
		}
		// check if the sequence fits any addons
		for _, addon := range Addons {

			if string(page[c-1])+string(page[c]) == addon.OpeningToken {
				// found a match
				newToken := *addon
				newToken.StartIndex = c + 1
				openTokens = append(openTokens, &newToken)
				fmt.Println("found match")

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
						fmt.Println("Computing token: ", t.Content)

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

						// make content modifier js so it can be added t opage.Scripts and can be bundled for lcient
						id := fmt.Sprintf("%v%v", t.Name, i)
						modifiedContent, script := t.ContentModifier(t.Content, *ctx, id)
						// document.Script = append(document.Script, script)
						jsFile.Contents += script
						fmt.Println("mod", modifiedContent)
						page = page[:t.StartIndex-2] + fmt.Sprintf("<%v cilia-id='%v'>%v</%v>", t.Name, id, modifiedContent, t.Name) + page[c+1:]
						c += len(modifiedContent)
						break
						// need to know if it should be rerun or only ssr
						// here, convert each token to js but dont run it.
						// then work out how they fit together so can be run in same context
						// once found we need to loop over forward and store the contents in a dom tree not on the page.
					}
				}
			}
		}
	}
	document.TextContents = page
	return *document
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
