package main

import (
	"encoding/json"
	"html/template"
	"log"
	"os"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
)

func generateHomepage() {
	homepageLinks := generateHomepageLinks()

	tmpl, err := template.ParseFiles("./html/layout.html", "./html/home.html")
	if err != nil {
		log.Fatalf("Error parsing template files: %v", err)
	}
	outputFile, err := os.Create("build/index.html")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	data := struct {
		Content        template.HTML
		Title          string
		ROOT_DIRECTORY string
	}{
		Content:        homepageLinks,
		Title:          "Home",
		ROOT_DIRECTORY: ROOT_DIRECTORY,
	}

	err = tmpl.ExecuteTemplate(
		outputFile, "layout.html",
		data,
	)

	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

type snippetIdAndData struct {
	id   string
	data SnippetData
}

func generateHomepageLinks() template.HTML {
	content := elem.Div(attrs.Props{
		attrs.Class: "flex flex-row flex-wrap gap-3 p-2",
	})

	snippets, err := os.ReadDir("./snippets")
	if err != nil {
		panic(err)
	}

	snippetsInCategories := make(map[string][]snippetIdAndData)
	for _, category := range Categories {
		snippetsInCategories[category.Id] = make([]snippetIdAndData, 0)
	}

	for _, e := range snippets {
		if !e.IsDir() {
			continue
		}
		snippetName := e.Name()

		var snippetData SnippetData
		snippetDataJson, err := os.ReadFile("snippets/" + snippetName + "/snippet.json")
		if err != nil {
			panic(`Snippet is missing a snippet.json file.`)
		}

		err = json.Unmarshal([]byte(snippetDataJson), &snippetData)
		if err != nil {
			panic("Error decoding JSON.")
		}

		snippetsInCategories[snippetData.Category] = append(snippetsInCategories[snippetData.Category], snippetIdAndData{snippetName, snippetData})
	}

	for _, category := range Categories {
		categoryDiv := elem.Div(
			attrs.Props{
				attrs.Class: "flex flex-col",
			},
			elem.Div(
				attrs.Props{
					attrs.Class: "flex flex-row",
				},
				elem.Img(
					attrs.Props{
						attrs.Src:   ROOT_DIRECTORY + category.Icon,
						attrs.Class: "w-6 h-6 inline mr-1",
						attrs.Alt:   "",
					},
				),
				elem.H2(attrs.Props{
					attrs.Class: "font-semibold",
				},
					elem.Text(category.DisplayName),
				),
			),
		)

		for _, snippet := range snippetsInCategories[category.Id] {
			categoryDiv.Children = append(categoryDiv.Children,
				elem.A(
					attrs.Props{
						attrs.Href: ROOT_DIRECTORY + "/snippets/" + snippet.id,
					},
					elem.Span(attrs.Props{
						attrs.Class: "text-blue-600 dark:text-blue-500",
					},
						elem.Text(snippet.data.Name),
					),
				),
			)
		}

		content.Children = append(content.Children, categoryDiv)
	}

	return template.HTML(content.Render())
}
