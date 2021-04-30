package exp4

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/rs/zerolog/log"
)

const (
	KeyTmpl   = ".tmpl"
	KeyEntity = "entity"
	tmplName  = "TemplateProject"
)

func Run() {
	ParseTemplates("rock")
	fmt.Println("done")

}

func ParseTemplates(projName string) {
	templ := template.New("")
	srcRootPath := "/your/absolute/path/TemplateProject"
	err := filepath.Walk(srcRootPath, func(path string, info os.FileInfo, err error) error {
		data := struct {
			UpperCamelCaseProjName string
			LowerCamelCaseProjName string
		}{
			UpperCamelCaseProjName: strings.Title(projName),
			LowerCamelCaseProjName: projName,
		}

		var f *os.File
		switch mode := info.Mode(); {
		case mode.IsDir():
			mask := syscall.Umask(0)
			defer syscall.Umask(mask)
			targetFolderPath := strings.Replace(path, tmplName, projName, 1)
			if strings.Contains(path, KeyEntity) {
				targetFolderPath = strings.Replace(targetFolderPath, KeyEntity, projName, 1)

			}
			err = os.Mkdir(targetFolderPath, 0777)
			if err != nil {
				fmt.Println(err)

			}
		case mode.IsRegular():
			ext := filepath.Ext(path)
			targetFilePath := strings.Replace(path, tmplName, projName, 1)
			if strings.Contains(path, KeyEntity) {
				targetFilePath = strings.ReplaceAll(targetFilePath, KeyEntity, projName)

			}
			if ext == KeyTmpl {
				targetFilePath = strings.Replace(targetFilePath, KeyTmpl, "", 1)

			}
			f, err = os.Create(targetFilePath)
			if err != nil {
				fmt.Println(err)

			}
			defer f.Close()
			if ext != KeyTmpl {
				srcFile, err := os.Open(path)
				if err != nil {
					fmt.Println(err)

				}
				defer srcFile.Close()
				_, err = io.Copy(f, srcFile)
				if err != nil {
					fmt.Println(err)

				}

			}

		}
		if strings.Contains(path, ".tmpl") {
			tt := template.Must(templ.ParseFiles(path))
			if err := tt.ExecuteTemplate(f, filepath.Base(path), data); err != nil {
				fmt.Println(err)

			}

		}

		return err

	})
	if err != nil {
		log.Error().Err(err).Msg("generate failed")
	}
}
