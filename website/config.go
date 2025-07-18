package main

type CategoryData struct {
	Id          string
	DisplayName string
	Icon        string
}

var Categories []CategoryData = []CategoryData{
	{"blocks", "Blocks", "/diamond_ore_0.png"},
	{"entities", "Entities", "/spawn_egg_30.png"},
	{"items", "Items", "/iron_pickaxe_0.png"},
	{"world-generation", "World Generation", "/buildplate.png"},
	{"functions", "Functions", "/command_block.png"},
	// {"script-api", "Script API", "/levers.png"},
}

var ExcludeFiles []string = []string{"snippet.md", "meta.json", "snippet_icon.png"}

type SnippetData struct {
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

const SNIPPET_REPO_OWNER string = "bedrock-oss"
const SNIPPET_REPO_NAME string = "bedrock-examples"

const SNIPPET_REPO_ROOT string = "https://github.com/" + SNIPPET_REPO_OWNER + "/" + SNIPPET_REPO_NAME + "/"

const ROOT_DIRECTORY string = "/bedrock-snippets"
const SNIPPET_DIRECTORY string = "./.tmp/snippet_repo/resources/"
