package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hatchibombotar/bedrock-snippets/website/helper"
	"html"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/chasefleming/elem-go/styles"

	"github.com/alecthomas/chroma/v2"
	chromaHtmlFormatter "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
)

type ContentPageData struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

func generatePagesForSnippet(snippetId string) {
	var snippetData SnippetData
	snippetDataJson, err := os.ReadFile(SNIPPET_DIRECTORY + "/" + snippetId + "/meta.json")

	if err == nil {
		err := json.Unmarshal([]byte(snippetDataJson), &snippetData)
		if err != nil {
			panic("Error decoding JSON.")
		}
	} else {
		fmt.Println("Snippet", snippetId, "is missing a meta.json file.")
	}

	sourceHref := SNIPPET_REPO_ROOT + "tree/main/resources/" + snippetId
	dirName := "./build/snippets/" + snippetId

	err = os.Mkdir(dirName, os.ModePerm)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.ParseFiles("./html/layout.html", "./html/snippet.html")
	if err != nil {
		fmt.Println("Error parsing template files")
		panic(err)
	}

	sidebar := generateSidebarElement(snippetId, "", 0)

	filepath.WalkDir(SNIPPET_DIRECTORY+"/"+snippetId+"/", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		rawPath, _ := filepath.Rel(SNIPPET_DIRECTORY, path)
		ext := filepath.Ext(path)
		outputPath := "build/snippets/" + rawPath

		if d.IsDir() {
			err := os.Mkdir(outputPath, os.ModePerm)
			if err != nil && !os.IsExist(err) {
				panic(err)
			}
			return nil
		}
		fileName := filepath.Base(path)

		if slices.Contains(ExcludeFiles, fileName) {
			return nil
		}

		outputFilePath := outputPath + ".html"
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		data := struct {
			Title          string
			SidebarHTML    template.HTML
			Breadcrumbs    string
			TextContent    template.HTML
			OverviewPage   bool
			SourceHref     string
			ROOT_DIRECTORY string
		}{
			Title:          snippetData.Name,
			SidebarHTML:    template.HTML(sidebar.Render()),
			Breadcrumbs:    strings.Join(strings.Split(filepath.ToSlash(rawPath), "/"), " > "),
			OverviewPage:   false,
			SourceHref:     sourceHref,
			ROOT_DIRECTORY: ROOT_DIRECTORY,
		}

		switch ext {
		case ".json", ".material":
			data.TextContent = CreateJSONPreview(path)
		case ".png":
			data.TextContent = CreatePNGPreview(path)
		case ".md":
			data.TextContent = CreateMDPreview(path)
		case ".js", ".ts":
			data.TextContent = CreateJSPreview(path)
		case ".mcfunction":
			data.TextContent = CreateMcFunctionPreview(path)
		case ".lang":
			data.TextContent = CreateLangPreview(path)
		default:
			data.TextContent = CreateTextPreview(path)
			log.Printf("No preview avaliable for %s", ext)
		}

		err = tmpl.ExecuteTemplate(outputFile, "layout.html", data)
		if err != nil {
			fmt.Println("Error executing template")
			panic(err)
		}

		return nil
	})

	indexFilePath := dirName + "/index.html"
	indexFile, err := os.Create(indexFilePath)
	if err != nil {
		panic(err)
	}
	defer indexFile.Close()

	data := struct {
		Title          string
		SidebarHTML    template.HTML
		Breadcrumbs    string
		TextContent    template.HTML
		OverviewPage   bool
		SourceHref     string
		ROOT_DIRECTORY string
	}{
		Title:          snippetData.Name,
		SidebarHTML:    template.HTML(sidebar.Render()),
		Breadcrumbs:    snippetId + " > overview",
		OverviewPage:   true,
		SourceHref:     sourceHref,
		ROOT_DIRECTORY: ROOT_DIRECTORY,
	}

	err = tmpl.ExecuteTemplate(indexFile, "layout.html", data)
	if err != nil {
		fmt.Println("Error executing template")
		panic(err)
	}
}

func generateSidebarElement(snippetName string, base string, level int) *elem.Element {
	contents, err := os.ReadDir(SNIPPET_DIRECTORY + "/" + snippetName + "/" + base)
	if err != nil {
		panic(err)
	}
	content := elem.Div(attrs.Props{
		attrs.Class: "flex flex-col",
	})

	if level == 0 {
		anchorElement := elem.Button(
			attrs.Props{
				attrs.Class: "w-full hover:bg-neutral-200 focus:bg-neutral-200 dark:hover:bg-neutral-700 dark:focus:bg-neutral-700 px-1 py-0.5 h-7 text-left truncate",
				"onclick":   fmt.Sprintf("window.location.href='%s'", ROOT_DIRECTORY+"/snippets/"+snippetName+"/"),
			},

			elem.Text("📄README"),
		)

		content.Children = append(content.Children, anchorElement)
	}

	for _, e := range contents {
		if slices.Contains(ExcludeFiles, e.Name()) {
			continue
		}

		anchorElementIcon := helper.Ternary(e.IsDir(), "🗂️", "📄")
		anchorElementLabel := anchorElementIcon + e.Name()
		if e.Name() == "rp" {
			anchorElementLabel = "🗂️Resource Pack"
		} else if e.Name() == "bp" {
			anchorElementLabel = "🗂️Behavior Pack"
		}

		anchorElement := elem.Button(
			attrs.Props{
				attrs.Class: "hover:bg-neutral-200 focus:bg-neutral-200 dark:hover:bg-neutral-700 dark:focus:bg-neutral-700 px-1 py-0.5 h-7",
			},

			elem.P(attrs.Props{
				attrs.Class: "w-full text-left",
				attrs.Style: styles.Props{
					styles.PaddingLeft: helper.Ternary(level > 0, fmt.Sprint("calc(var(--spacing) * ", (2*level)+1, ")"), "0"),
				}.ToInline(),
			}, elem.Text(anchorElementLabel)),
		)

		if !e.IsDir() {
			anchorElement.Attrs["onclick"] = fmt.Sprintf("window.location.href='%s'", ROOT_DIRECTORY+"/snippets/"+snippetName+"/"+base+e.Name())
			content.Children = append(content.Children, anchorElement)
			continue
		}

		// If directory: create collapsible container

		containerID := fmt.Sprintf("dir-%s-%d", base+e.Name(), level)

		childContainer := elem.Div(attrs.Props{
			attrs.Class: "flex flex-col",
			attrs.ID:    containerID,
		}, generateSidebarElement(snippetName, base+e.Name()+"/", level+1))

		// Add onclick handler to toggle visibility
		anchorElement.Attrs["onclick"] = fmt.Sprintf("toggleDir('%s')", containerID)

		// Append the anchor and its collapsible child container
		content.Children = append(content.Children, anchorElement, childContainer)
	}

	return content
}

func CreateJSONPreview(filePath string) template.HTML {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Unable to read file.", err)
	}

	htmlFormatter := chromaHtmlFormatter.New(
		chromaHtmlFormatter.LineNumbersInTable(true),
		chromaHtmlFormatter.WithClasses(true),
		chromaHtmlFormatter.ClassPrefix("chroma-"),
	)

	lexer := lexers.Get("json")

	iterator, err := lexer.Tokenise(nil, string(content))
	if err != nil {
		panic(err)
	}

	var result bytes.Buffer
	err = htmlFormatter.Format(&result, &chroma.Style{}, iterator)
	if err != nil {
		panic(err)
	}

	container := elem.Div(attrs.Props{
		attrs.ID:            "snippet-content",
		"data-content-type": "text",
		"data-content-text": html.EscapeString(string(content)),
	},
		elem.Raw(result.String()),
	)

	return template.HTML(container.Render())
}

func CreateJSPreview(filePath string) template.HTML {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Unable to read file.", err)
	}

	htmlFormatter := chromaHtmlFormatter.New(
		chromaHtmlFormatter.LineNumbersInTable(true),
		chromaHtmlFormatter.WithClasses(true),
		chromaHtmlFormatter.ClassPrefix("chroma-"),
	)

	lexer := lexers.Get("ts")

	iterator, err := lexer.Tokenise(nil, string(content))
	if err != nil {
		panic(err)
	}

	var result bytes.Buffer
	err = htmlFormatter.Format(&result, &chroma.Style{}, iterator)
	if err != nil {
		panic(err)
	}

	container := elem.Div(attrs.Props{
		attrs.ID:            "snippet-content",
		"data-content-type": "text",
		"data-content-text": html.EscapeString(string(content)),
	},
		elem.Raw(result.String()),
	)

	return template.HTML(container.Render())
}

func CreateMcFunctionPreview(filePath string) template.HTML {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Unable to read file.", err)
	}

	htmlFormatter := chromaHtmlFormatter.New(
		chromaHtmlFormatter.LineNumbersInTable(true),
		chromaHtmlFormatter.WithClasses(true),
		chromaHtmlFormatter.ClassPrefix("chroma-"),
	)

	lexer := lexers.Get("mcfunction")

	iterator, err := lexer.Tokenise(nil, string(content))
	if err != nil {
		panic(err)
	}

	var result bytes.Buffer
	err = htmlFormatter.Format(&result, &chroma.Style{}, iterator)
	if err != nil {
		panic(err)
	}

	container := elem.Div(attrs.Props{
		attrs.ID:            "snippet-content",
		"data-content-type": "text",
		"data-content-text": html.EscapeString(string(content)),
	},
		elem.Raw(result.String()),
	)

	return template.HTML(container.Render())
}
func CreateLangPreview(filePath string) template.HTML {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Unable to read file.", err)
	}

	htmlFormatter := chromaHtmlFormatter.New(
		chromaHtmlFormatter.LineNumbersInTable(true),
		chromaHtmlFormatter.WithClasses(true),
		chromaHtmlFormatter.ClassPrefix("chroma-"),
	)

	lexer := lexers.Get("ini")

	iterator, err := lexer.Tokenise(nil, string(content))
	if err != nil {
		panic(err)
	}

	var result bytes.Buffer
	err = htmlFormatter.Format(&result, &chroma.Style{}, iterator)
	if err != nil {
		panic(err)
	}

	container := elem.Div(attrs.Props{
		attrs.ID:            "snippet-content",
		"data-content-type": "text",
		"data-content-text": html.EscapeString(string(content)),
	},
		elem.Raw(result.String()),
	)

	return template.HTML(container.Render())
}

func CreatePNGPreview(filePath string) template.HTML {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Unable to read file.", err)
	}

	imgBase64Str := base64.StdEncoding.EncodeToString(content)

	result := "<img id=\"snippet-content\" data-content-type=\"image\" class=\"preview-image\" src=\"data:image/png;base64," + imgBase64Str + "\" >"

	return template.HTML(result)
}

func CreateMDPreview(filePath string) template.HTML {
	indexFileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Unable to read file.", err)
	}

	renderedIndexFileContent := mdToHTML([]byte(indexFileContent))
	return template.HTML(renderedIndexFileContent)
}

func CreateTextPreview(filePath string) template.HTML {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Unable to read file.", err)
	}

	container := elem.Div(attrs.Props{
		attrs.ID:            "snippet-content",
		"data-content-type": "text",
		"data-content-text": html.EscapeString(string(content)),
	},
		elem.Raw(string(content)),
	)

	return template.HTML(container.Render())
}
