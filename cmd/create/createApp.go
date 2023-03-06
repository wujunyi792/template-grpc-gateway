package create

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"pinnacle-primary-be/pkg/fs"
	"strings"
	"text/template"
)

var (
	appName  string
	dir      string
	force    bool
	StartCmd = &cobra.Command{
		Use:     "create",
		Short:   "create a new app",
		Example: "app create -n users",
		Run: func(cmd *cobra.Command, args []string) {
			err := load()
			if err != nil {
				println(err.Error())
				os.Exit(1)
			}
			println("App " + appName + " generate success under " + dir)
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&appName, "name", "n", "", "create a new app with provided name")
	StartCmd.PersistentFlags().StringVarP(&dir, "path", "p", "internal/app", "new file will generate under provided path")
	StartCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force generate the app")
}

func load() error {
	if appName == "" {
		return errors.New("app name should not be empty, use -n")
	}

	router := path.Join(dir, appName, "router")
	handlerMain := path.Join(dir, appName, "handler", "v1")
	handlerType := path.Join(dir, appName, "handler", "v1")
	dto := path.Join(dir, appName, "dto")
	e := path.Join(dir, appName, "e")
	service := path.Join(dir, appName, "service")
	model := path.Join(dir, appName, "model")
	trigger := path.Join(dir, "routerInitialize")

	_ = fs.IsNotExistMkDir(router)
	_ = fs.IsNotExistMkDir(handlerType)
	_ = fs.IsNotExistMkDir(handlerType)
	_ = fs.IsNotExistMkDir(e)
	_ = fs.IsNotExistMkDir(dto)
	_ = fs.IsNotExistMkDir(service)
	_ = fs.IsNotExistMkDir(model)
	_ = fs.IsNotExistMkDir(trigger)

	m := map[string]string{}
	m["appNameExport"] = strings.ToUpper(appName[:1]) + appName[1:]
	m["appName"] = strings.ToLower(appName[:1]) + appName[1:]

	router += "/" + m["appName"] + ".go"
	service += "/" + m["appName"] + ".go"
	model += "/" + m["appName"] + ".go"
	handlerMain += "/" + m["appName"] + ".go"
	handlerType += "/" + "type.go"
	dto += "/" + m["appName"] + ".go"
	trigger += "/" + m["appName"] + ".go"
	e += "/" + m["appName"] + ".go"

	if !force && (fs.FileExist(router) || fs.FileExist(handlerMain) || fs.FileExist(handlerType) ||
		fs.FileExist(dto) || fs.FileExist(trigger)) || fs.FileExist(service) || fs.FileExist(model) {
		return errors.New("target file already exist, use -f flag to cover")
	}

	if rt, err := template.ParseFiles("template/router.template"); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, m)
		fs.FileCreate(b, router)
	}

	if rt, err := template.ParseFiles("template/handler.template"); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, m)
		fs.FileCreate(b, handlerMain)
	}
	if rt, err := template.ParseFiles("template/type.template"); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, m)
		fs.FileCreate(b, handlerType)
	}

	if rt, err := template.ParseFiles("template/dto.template"); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, m)
		fs.FileCreate(b, dto)
	}
	if rt, err := template.ParseFiles("template/e.template"); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, m)
		fs.FileCreate(b, e)
	}
	if rt, err := template.ParseFiles("template/service.template"); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, m)
		fs.FileCreate(b, service)
	}
	if rt, err := template.ParseFiles("template/model.template"); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, m)
		fs.FileCreate(b, model)
	}

	if rt, err := template.ParseFiles("template/trigger.template"); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, m)
		fs.FileCreate(b, trigger)
	}

	return nil
}
