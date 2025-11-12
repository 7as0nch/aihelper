package main

import (
	"os"
	"github.com/example/aichat/backend/models/generator/model"
	myStrings "github.com/example/aichat/backend/tools/strings"
	"strings"
	"text/template"
)

type TemplateData struct {
	Package   string
	ModelPath string
	RepoName  string
	ModelName string
	//ModelVariablePlural string
	//ModelVariable       string
}

func main() {
	// 从命令行参数中读取模型和数据访问接口的名称
	md := model.SysParameter{}
	split := strings.Split(md.TableName(), "_")
	sb := myStrings.NewStrBuilder()
	for _, s := range split {
		sb.StrAppend(strings.ToUpper(s[:1]) + s[1:])
	}
	modelName := sb.ToString()
	repoName := strings.ToLower(modelName[:1]) + modelName[1:]

	// 填充模板变量
	data := TemplateData{
		Package:   "data",
		ModelPath: "phm/models/generator/model",
		ModelName: modelName,
		RepoName:  repoName,
		//ModelVariable:       modelName[:1],
		//ModelVariablePlural: modelName[:1] + "s",
	}

	var fn = make(chan func(), 2)

	var fn1 = func() {
		// 加载模板文件
		tmpl, err := template.ParseFiles("D:\\JavaProjectLearning\\goland\\generatedProject\\cdu-phm\\back\\phm\\cmd\\tmpl\\template.tmpl")
		if err != nil {
			panic(err)
		}
		// 将生成的代码写入文件中
		abPath := "D:\\JavaProjectLearning\\goland\\generatedProject\\cdu-phm\\back\\phm\\internal\\data\\"
		f, err := os.Create(abPath + repoName + ".go")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(f, data)

	}
	var fn2 = func() {
		// biz层
		tmpl, err := template.ParseFiles("D:\\JavaProjectLearning\\goland\\generatedProject\\cdu-phm\\back\\phm\\cmd\\tmpl\\biz.tmpl")
		if err != nil {
			panic(err)
		}
		// 将生成的代码写入文件中
		abPath := "D:\\JavaProjectLearning\\goland\\generatedProject\\cdu-phm\\back\\phm\\internal\\biz\\"
		f, err := os.Create(abPath + repoName + ".go")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = tmpl.Execute(f, data)
		if err != nil {
			panic(err)
		}
	}
	fn <- fn1
	fn <- fn2
	//
	var f = <-fn
	var f2 = <-fn
	f()
	f2()
}
