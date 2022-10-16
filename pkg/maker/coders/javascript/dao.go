package javascript

import (
	"go.uber.org/zap"
	"strings"
	"teamide/pkg/maker/coders/common"
	"teamide/pkg/maker/modelers"
	"teamide/pkg/util"
)

type daoCoder struct {
	*appCoder
}

func (this_ *daoCoder) Gen(code *common.Code, model *modelers.DaoModel) (err error) {

	return
}

func GetDaoMethodName(name string) (methodName string) {
	methodName = GetFormatMethodName(name)
	methodName = "dao" + util.Capitalize(methodName)
	return
}

func GetDaoJavascript(app *modelers.Application, dao *modelers.DaoModel) (javascript string, err error) {
	serviceMethodName := GetDaoMethodName(dao.Name)
	javascript += ""
	javascript += "function " + serviceMethodName + "("
	for _, arg := range dao.Args {
		javascript += arg.Name + ", "
	}
	javascript = strings.TrimSuffix(javascript, ", ")

	javascript += ") {"

	javascript += "\n"

	var stepsJavascript string
	stepsJavascript, err = GetJavascriptBySteps(app, dao.Steps, 1)
	if err != nil {
		util.Logger.Error("GetDaoJavascript GetJavascriptBySteps error", zap.Any("dao", dao), zap.Error(err))
		return
	}
	if util.IsNotEmpty(stepsJavascript) {
		javascript += stepsJavascript
	}

	if util.IsNotEmpty(dao.Return) {
		if dao.Return != "-" {
			util.AppendLine(&javascript, "return "+dao.Return, 1)
		} else {
			util.AppendLine(&javascript, "return", 1)
		}
	}

	javascript += "}"
	// fmt.Println(javascript)
	return
}
