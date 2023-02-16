package url

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/JanGordon/cilia-framework/server/pkg/global"
)

var reservedPaths = []string{
	"/public",
	"/components",
	"/assets",
}

func ResolveURL(url string) string {
	// takes a url and returns a path to the file relative the project root
	// opposite of ResolvePath
	for _, p := range reservedPaths {
		if strings.HasPrefix(url, p) {
			url = "/.." + url
			break
		}
	}
	url = "/routes" + url
	if filepath.Ext(url) == "" {
		url = filepath.Join(url, "index.html.out")
	}
	return global.ProjectRoot + url
}

func ResolvePath(path string, locationPath string) (string, error) {
	// takes a path and returns a url that points to the file
	// opposite of ResolveURL
	url := getAbs(path, locationPath)
	url = strings.Replace(url, "/routes", "", 1)
	return url, nil
}

func PathRelativeToRoot(path string, locationPath string) (string, error) {
	// takes a path and returns a url that points to the file
	// opposite of ResolveURL
	url := getAbs(path, locationPath)
	return url, nil
}
func getAbs(path string, locationPath string) string {
	new_file_path := path
	new_abs_path := path
	if filepath.IsAbs(new_file_path) {
		new_abs_path = new_file_path
	} else {
		new_abs_path, _ = filepath.Abs(filepath.Join(locationPath, path))
		fmt.Println(new_abs_path, locationPath)
	}
	return new_abs_path

}
