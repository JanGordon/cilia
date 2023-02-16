package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/JanGordon/cilia-framework/server/pkg/global"
	"github.com/spf13/cobra"
)

type FileEntry struct {
	Name     string
	IsDir    bool
	Contents string
	Children []FileEntry
}

type SkeletonStrucutre struct {
	Name     string
	Contents []FileEntry
}

var Demos = []SkeletonStrucutre{
	{
		Name: "scaffold",
		Contents: []FileEntry{
			{
				"routes",
				true,
				"",
				[]FileEntry{
					{
						"routes/index.html",
						false,
						`<!DOCTYPE html>
						<html lang="en">
						<head>
							<meta charset="UTF-8">
							<meta http-equiv="X-UA-Compatible" content="IE=edge">
							<meta name="viewport" content="width=device-width, initial-scale=1.0">
							<title>Cilia Project</title>
						</head>
						<body>
							<h1>Welcome to Cilia!</h1>
							<a href="https://github.com/JanGordon/cilia">github</a>
							<script src="/test.js"></script>
						</body>`,
						[]FileEntry{},
					},
					{
						"routes/index.html.out",
						false,
						`<!DOCTYPE html>
						<html lang="en"><head>
						<meta charset="UTF-8"/>
						<meta http-equiv="X-UA-Compatible" content="IE=edge"/>
						<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
						<title>Cilia Project</title>
					</head>
					<body>
						<h1>Welcome to Cilia!</h1>
						<a href="https://github.com/JanGordon/cilia">github</a>
						<script src="/test.js"></script>
					</body></html>`,
						[]FileEntry{},
					},
					{
						"routes/test.js",
						false,
						`const devSocket = new WebSocket("ws:localhost:8080/ws");

						devSocket.onopen = (event) => {
							console.log("connection opened")
						
						}
						
						devSocket.onmessage = async (event) => {
							if (event.data == "reload") {
								devSocket.send("reload successful");
								window.location.reload();
							} else if (event.data == "reloadhtml") {
								devSocket.send("reload successful");
								fetch(window.location.pathname)
								.then(function (response) {
									return response.text()
								}).then(function (data) {
									document.body.innerHTML = data
									console.log("reloaded html", data)
								})
							}
						}
						`,
						[]FileEntry{},
					},
				},
			},
			{
				"components",
				true,
				"",
				[]FileEntry{},
			},
			{
				"public",
				true,
				"",
				[]FileEntry{},
			},
			{
				"stem.toml",
				false,
				`Version = 0.1
				Name = 'cilia-project'`,
				[]FileEntry{},
			},
		},
	},
}

var createCmd = &cobra.Command{
	Use:     "create",
	Args:    cobra.ExactArgs(2),
	Aliases: []string{"c"},
	Short:   "creates a base cilia project",
	Run: func(cmd *cobra.Command, args []string) {
		var selectedStructure SkeletonStrucutre
		for _, s := range Demos {
			if s.Name == args[1] {
				selectedStructure = s
			}
		}
		if err := os.Mkdir(filepath.Join(global.ProjectRoot, args[0]), os.ModePerm); err != nil {
			log.Fatal(err)
		}
		global.ProjectRoot = filepath.Join(global.ProjectRoot, args[0])
		loopDir(selectedStructure.Contents)

	},
}

func loopDir(dir []FileEntry) {
	for _, f := range dir {
		if f.IsDir {
			if err := os.Mkdir(filepath.Join(global.ProjectRoot, f.Name), os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("Creating %v", filepath.Join(global.ProjectRoot, f.Name))

			file, err := os.Create(filepath.Join(global.ProjectRoot, f.Name))
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			file.Write([]byte(f.Contents))
		}
		loopDir(f.Children)

	}
}

func init() {
	// createCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "add addon globally")
	rootCmd.AddCommand(createCmd)
}
