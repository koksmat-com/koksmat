package kitchen

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

type Pair struct {
	Key   string
	Value string
}

type Environment struct {
	Type   string
	Tenant string
	Name   string
	Title  string
	Pairs  []Pair
}

type Status struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	About       string `json:"about"`
	Markdown    string `json:"markdown"`
	Description string `json:"description"`
	Url         string `json:"url"`

	Environments []Environment `json:"environments"`
}

func GetStatus(kitchen string, parseMD bool) (Status, error) {
	root := viper.GetString("KITCHENROOT")
	status := Status{}
	kitchenPath := path.Join(root, kitchen)
	about, meta, err := ReadMarkdown(false, kitchenPath, "readme.md")
	if parseMD {
		html, _, err := ParseMarkdown(false, kitchenPath, about)
		if err != nil {
			return status, err
		}
		status.About = html
	}

	status.Markdown = about
	status.Title = GetMetadataProperty(meta, "title", kitchen)
	status.Description = GetMetadataProperty(meta, "description", "")

	sharePointPath := path.Join(kitchenPath, ".koksmat", "sharepoint")

	sharePointEnvironments, err := os.ReadDir(sharePointPath)
	if err != nil {
		return status, nil
	}
	for _, c := range sharePointEnvironments {
		if c.IsDir() {
			env := Environment{}
			env.Name = c.Name()
			env.Title = c.Name()
			env.Type = "sharepoint"
			env.Tenant = c.Name()
			env.Pairs = []Pair{}
			status.Environments = append(status.Environments, env)

		}
	}

	return status, nil
}
