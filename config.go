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
	{"misc", "Misc", "/crafting_table_0.png"},
}

var ExcludeFiles []string = []string{"snippet.md", "snippet.json", "snippet_icon.png"}

type SnippetData struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

const REPOSITORY_ROOT string = "https://github.com/Hatchibombotar/bedrock-snippets/"
const ROOT_DIRECTORY string = "/bedrock-snippets"
