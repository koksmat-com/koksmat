package kitchen

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/koksmat-com/koksmat/connectors"
	"github.com/spf13/viper"
	img64 "github.com/tenkoh/goldmark-img64"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

func GetMetadataProperty(meta map[string]interface{}, key string, defaultValue string) string {

	if meta[key] != nil {
		return meta[key].(string)
	}
	return defaultValue
}

type Metadata map[string]interface{}

func wrapperRenderer(w util.BufWriter, ctx highlighting.CodeBlockContext, entering bool) {
	language, ok := ctx.Language()
	lang := string(language)
	// code block with a language
	if ok && lang != "" {
		if entering {
			w.WriteString(`<div style="border-radius:10px; background-color: #282a36;padding:20px;margin-top:20px;margin-bottom:20px;overflow-x:auto" data-lang=` + lang + ">")
		} else {
			w.WriteString(`</div>`)
		}
		return
	}

	// code block with no language specified
	if language == nil {
		if entering {
			w.WriteString(`<pre style="padding:10px"><code>`)
		} else {
			w.WriteString(`</code></pre>`)
		}
	}
}
func ParseMarkdown(noRender bool, parentPath string, content string) (string, Metadata, error) {
	var buf bytes.Buffer
	context := parser.NewContext()
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, meta.Meta),
		// goldmark.WithExtensions(&mermaid.Extender{
		// 	RenderMode: mermaid.RenderModeServer,
		// 	CLI:        mermaid.MMDC("/usr/local/bin/mmdc"),
		// }),
		goldmark.WithExtensions(img64.Img64),
		goldmark.WithRendererOptions(img64.WithParentPath(parentPath)),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
		// goldmark.WithExtensions(
		// 	highlighting.Highlighting,
		// ),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithWrapperRenderer(wrapperRenderer),
			),
		),
	)
	if noRender {
		md = goldmark.New(
			goldmark.WithExtensions(extension.GFM, meta.Meta),
		)
	}
	if err := md.Convert([]byte(content), &buf, parser.WithContext(context)); err != nil {
		return "", nil, err
	}
	metaData := meta.Get(context)
	return buf.String(), metaData, nil
}
func ParseMarkdownGetMetadata(parentPath string, content string) (Metadata, error) {
	var buf bytes.Buffer
	context := parser.NewContext()
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, meta.Meta),
		// goldmark.WithExtensions(&mermaid.Extender{
		// 	RenderMode: mermaid.RenderModeServer,
		// 	CLI:        mermaid.MMDC("/usr/local/bin/mmdc"),
		// }),
		goldmark.WithExtensions(img64.Img64),
		goldmark.WithRendererOptions(img64.WithParentPath(parentPath)),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
		// goldmark.WithExtensions(
		// 	highlighting.Highlighting,
		// ),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithWrapperRenderer(wrapperRenderer),
			),
		),
	)
	if err := md.Convert([]byte(content), &buf, parser.WithContext(context)); err != nil {
		return nil, err
	}
	metaData := meta.Get(context)
	return metaData, nil
}
func ReadMarkdown(dontrenderMarkdown bool, pathname string, filename string) (string, Metadata, error) {

	filepath := filepath.Join(pathname, filename)
	if !FileExists(filepath) {
		return "", nil, nil
	}
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return "", nil, err
	}
	_, metadata, err := ParseMarkdown(dontrenderMarkdown, pathname, string(fileContent))

	return string(fileContent), metadata, err

}

func List() (*[]Kitchen, error) {
	// userHome, err := os.UserHomeDir()
	// if err != nil {
	// 	return nil, err
	// }

	// root := path.Join(userHome, "kitchens")
	root := viper.GetString("KITCHENROOT")
	dirs, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	kitchens := []Kitchen{}

	for _, dir := range dirs {
		if dir.IsDir() && !strings.HasPrefix(dir.Name(), ".") {

			k, err := GetStatus(dir.Name(), false)
			if err != nil {
				return nil, err
			}
			kitchen := Kitchen{
				Name:        dir.Name(),
				Title:       k.Title,
				Description: k.Description,
				Path:        path.Join(root, dir.Name()),
			}
			kitchens = append(kitchens, kitchen)
			//	}
		}
	}

	return &kitchens, nil
}

func GetScripts(stationPath string, subPath string) ([]Script, error) {

	filePath := path.Join(stationPath, subPath)

	scripts, err := os.ReadDir(filePath)
	result := []Script{}
	if err != nil {
		return nil, err
	}
	for _, script := range scripts {
		if script.IsDir() {
			s, err := GetScripts(stationPath, path.Join(subPath, script.Name()))
			if err != nil {
				return nil, err
			}
			result = append(result, s...)
		}
		if !script.IsDir() && strings.HasSuffix(script.Name(), ".ps1") && !strings.HasPrefix(script.Name(), "test-") {
			markdown, envs, err := ReadMarkdownFromPowerShell(path.Join(filePath, script.Name()))
			if err != nil {
				return nil, err
			}

			_, scriptmeta, _ := ParseMarkdown(false, filePath, markdown)
			defaultTag := strings.ReplaceAll(strings.ToLower(script.Name()), " ", "-")
			script := Script{
				Name:        path.Join(subPath, script.Name()),
				Title:       path.Join(subPath, GetMetadataProperty(scriptmeta, "title", script.Name())),
				Description: GetMetadataProperty(scriptmeta, "description", ""),
				Environment: envs,
				Input:       GetMetadataProperty(scriptmeta, "input", ""),
				Output:      GetMetadataProperty(scriptmeta, "output", ""),
				Cron:        GetMetadataProperty(scriptmeta, "cron", ""),
				Connection:  GetMetadataProperty(scriptmeta, "connection", ""),
				Tag:         GetMetadataProperty(scriptmeta, "tag", defaultTag),
				Trigger:     GetMetadataProperty(scriptmeta, "trigger", ""),
				API:         GetMetadataProperty(scriptmeta, "api", ""),
			}
			result = append(result, script)

		}
		if !script.IsDir() && strings.HasSuffix(script.Name(), ".go") && !strings.HasSuffix(script.Name(), "_test.go") {
			markdown, err := ReadMarkdownFromGo(path.Join(filePath, script.Name()))
			if err != nil {
				return nil, err
			}

			_, scriptmeta, _ := ParseMarkdown(false, filePath, markdown)

			script := Script{
				Name:        path.Join(subPath, script.Name()),
				Title:       path.Join(subPath, GetMetadataProperty(scriptmeta, "title", script.Name())),
				Description: GetMetadataProperty(scriptmeta, "description", ""),
			}
			result = append(result, script)

		}
	}
	return result, nil
}

func GetStations(kitchenName string) (*Kitchen, error) {

	root := viper.GetString("KITCHENROOT")
	filepath := filepath.Join(root, kitchenName)
	dirs, err := os.ReadDir(filepath)
	if err != nil {
		return nil, err
	}

	kitchen := &Kitchen{
		Name: kitchenName,
		Path: filepath,
	}

	readme, meta, _ := ReadMarkdown(false, filepath, "readme.md")

	defaultKitchenTag := strings.ReplaceAll(strings.ToLower(kitchenName), " ", "-")

	kitchen.Readme = readme
	kitchen.Title = GetMetadataProperty(meta, "title", kitchenName)
	kitchen.Description = GetMetadataProperty(meta, "description", "")
	kitchen.Tag = GetMetadataProperty(meta, "tag", defaultKitchenTag)
	for _, dir := range dirs {
		if dir.IsDir() && !strings.HasPrefix(dir.Name(), ".") {
			stationPath := path.Join(filepath, dir.Name())
			stationReadme, stationmeta, _ := ReadMarkdown(false, path.Join(stationPath), "readme.md")
			defaultStationTag := strings.ReplaceAll(strings.ToLower(dir.Name()), " ", "-")

			station := Station{
				Readme: stationReadme,
				Name:   dir.Name(),
				Path:   stationPath,
				Title:  GetMetadataProperty(stationmeta, "title", dir.Name()),
				Tag:    GetMetadataProperty(stationmeta, "tag", defaultStationTag),
			}

			station.Description = GetMetadataProperty(meta, "description", "")
			station.Scripts, err = GetScripts(stationPath, "")
			/*
				scripts, err := os.ReadDir(stationPath)
				if err != nil {
					return nil, err
				}
				for _, script := range scripts {
					if !script.IsDir() && strings.HasSuffix(script.Name(), ".ps1") && !strings.HasPrefix(script.Name(), "test-") {
						markdown, err := ReadMarkdownFromPowerShell(path.Join(stationPath, script.Name()))
						if err != nil {
							return nil, err
						}

						_, scriptmeta, _ := ParseMarkdown(markdown)

						script := Script{
							Name:        script.Name(),
							Title:       getMetadataProperty(scriptmeta, "title", script.Name()),
							Description: getMetadataProperty(scriptmeta, "description", ""),
						}
						station.Scripts = append(station.Scripts, script)

					}
				}
			*/
			kitchen.Stations = append(kitchen.Stations, station)
			//	}
		}
	}

	return kitchen, nil
}

type KitchenJobs struct {
	Name string
	Data string
}

func GetJobs(kitchenName string) (*KitchenJobs, error) {

	root := viper.GetString("KITCHENROOT")
	filepath := filepath.Join(root, kitchenName)

	data, err := connectors.Execute("gh", *&connectors.Options{Dir: filepath}, "jobs", "list")
	if err != nil {
		return nil, err
	}
	kitchen := &KitchenJobs{
		Name: kitchenName,
		Data: string(data),
	}

	return kitchen, nil
}
