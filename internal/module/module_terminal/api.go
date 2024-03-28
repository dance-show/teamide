package module_terminal

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	goSSH "golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"teamide/internal/module/module_node"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
	"teamide/pkg/ssh"
	"teamide/pkg/terminal"
)

type api struct {
	*WorkerFactory
	terminalCommandService *TerminalCommandService
}

func NewApi(toolboxService_ *module_toolbox.ToolboxService, nodeService_ *module_node.NodeService, res *TerminalCommandService) *api {
	return &api{
		WorkerFactory:          NewWorkerFactory(toolboxService_, nodeService_),
		terminalCommandService: NewTerminalCommandService(toolboxService_.ServerContext),
	}
}

var (
	// Terminal 权限

	// Power 文件管理器 基本 权限
	Power           = base.AppendPower(&base.PowerAction{Action: "terminal", Text: "终端", ShouldLogin: true, StandAlone: true})
	websocketPower  = base.AppendPower(&base.PowerAction{Action: "websocket", Text: "终端WebSocket", ShouldLogin: true, StandAlone: true, Parent: Power})
	check           = base.AppendPower(&base.PowerAction{Action: "check", Text: "终端测试", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower      = base.AppendPower(&base.PowerAction{Action: "close", Text: "终端关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
	keyPower        = base.AppendPower(&base.PowerAction{Action: "key", Text: "终端Key", ShouldLogin: true, StandAlone: true, Parent: Power})
	changeSizePower = base.AppendPower(&base.PowerAction{Action: "changeSize", Text: "终端窗口大小变更", ShouldLogin: true, StandAlone: true, Parent: Power})
	getLogs         = base.AppendPower(&base.PowerAction{Action: "getLogs", Text: "getLogs", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteLog       = base.AppendPower(&base.PowerAction{Action: "deleteLog", Text: "deleteLog", ShouldLogin: true, StandAlone: true, Parent: Power})
	cleanLog        = base.AppendPower(&base.PowerAction{Action: "cleanLog", Text: "cleanLog", ShouldLogin: true, StandAlone: true, Parent: Power})
	downloadLog     = base.AppendPower(&base.PowerAction{Action: "downloadLog", Text: "downloadLog", ShouldLogin: true, StandAlone: true, Parent: Power})
	upload          = base.AppendPower(&base.PowerAction{Action: "upload", Text: "upload", ShouldLogin: true, StandAlone: true, Parent: Power})
	systemInfo      = base.AppendPower(&base.PowerAction{Action: "system/info", Text: "system", ShouldLogin: true, StandAlone: true, Parent: Power})
	systemMonitor   = base.AppendPower(&base.PowerAction{Action: "system/monitor", Text: "system", ShouldLogin: true, StandAlone: true, Parent: Power})

	command       = base.AppendPower(&base.PowerAction{Action: "command", Text: "命令行", ShouldLogin: true, StandAlone: true, Parent: Power})
	commandSave   = base.AppendPower(&base.PowerAction{Action: "save", Text: "插入", ShouldLogin: true, StandAlone: true, Parent: command})
	commandQuery  = base.AppendPower(&base.PowerAction{Action: "query", Text: "查询", ShouldLogin: true, StandAlone: true, Parent: command})
	commandCount  = base.AppendPower(&base.PowerAction{Action: "count", Text: "查询", ShouldLogin: true, StandAlone: true, Parent: command})
	commandDelete = base.AppendPower(&base.PowerAction{Action: "delete", Text: "删除", ShouldLogin: true, StandAlone: true, Parent: command})
	commandClean  = base.AppendPower(&base.PowerAction{Action: "clean", Text: "清理", ShouldLogin: true, StandAlone: true, Parent: command})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: keyPower, Do: this_.key})
	apis = append(apis, &base.ApiWorker{Power: websocketPower, Do: this_.websocket, IsWebSocket: true})
	apis = append(apis, &base.ApiWorker{Power: changeSizePower, Do: this_.changeSize})
	apis = append(apis, &base.ApiWorker{Power: check, Do: this_.check})
	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})
	apis = append(apis, &base.ApiWorker{Power: getLogs, Do: this_.getLogs})
	apis = append(apis, &base.ApiWorker{Power: deleteLog, Do: this_.deleteLog})
	apis = append(apis, &base.ApiWorker{Power: cleanLog, Do: this_.cleanLog})
	apis = append(apis, &base.ApiWorker{Power: downloadLog, Do: this_.downloadLog})
	apis = append(apis, &base.ApiWorker{Power: upload, Do: this_.upload, IsUpload: true, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: systemInfo, Do: this_.systemInfo})
	apis = append(apis, &base.ApiWorker{Power: systemMonitor, Do: this_.systemMonitor, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: commandSave, Do: this_.commandSave, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: commandQuery, Do: this_.commandQuery, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: commandCount, Do: this_.commandCount, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: commandClean, Do: this_.commandClean, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: commandDelete, Do: this_.commandDelete, NotRecodeLog: true})

	return
}

func (this_ *api) key(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	service, _, err := this_.createService(&CreateParam{
		place:    request.Place,
		placeId:  request.PlaceId,
		workerId: request.WorkerId,
		lastUser: request.LastUser,
		lastDir:  request.LastDir,
	})
	if err != nil {
		return
	}

	data := make(map[string]interface{})

	data["isWindows"], err = service.service.IsWindows()
	if err != nil {
		return
	}
	data["key"] = util.GetUUID()
	res = data
	return
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  32 * 1024,
	WriteBufferSize: 32 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (this_ *api) websocket(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	if request.JWT == nil || request.JWT.UserId == 0 {
		err = errors.New("登录用户获取失败")
		return
	}
	key := c.Query("key")
	if key == "" {
		err = errors.New("key获取失败")
		return
	}
	place := c.Query("place")
	if place == "" {
		err = errors.New("place获取失败")
		return
	}
	placeId := c.Query("placeId")
	workerId := c.Query("workerId")
	if workerId == "" {
		err = errors.New("workerId获取失败")
		return
	}

	cols, _ := strconv.Atoi(c.Query("cols"))
	rows, _ := strconv.Atoi(c.Query("rows"))
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	service := this_.GetService(key)
	if service != nil {
		err = errors.New("会话[" + key + "]已存在")

		_ = ws.WriteMessage(websocket.BinaryMessage, []byte("service create error:"+err.Error()))
		this_.Logger.Error("websocket start error", zap.Error(err))
		_ = ws.Close()
		return
	}

	err = this_.Start(key,
		&CreateParam{
			place:    place,
			placeId:  placeId,
			workerId: workerId,
			lastUser: c.Query("lastUser"),
			lastDir:  c.Query("lastDir"),
		},
		&terminal.Size{
			Cols: cols,
			Rows: rows,
		}, ws)
	if err != nil {
		_ = ws.WriteMessage(websocket.BinaryMessage, []byte("start error:"+err.Error()))
		this_.Logger.Error("websocket start error", zap.Error(err))
		_ = ws.Close()
		return
	}

	res = base.HttpNotResponse
	return
}

func (this_ *api) upload(r *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	key := c.PostForm("key")
	if key == "" {
		err = errors.New("key获取失败")
		return
	}
	value := c.PostForm("value")
	var bs []byte
	if value != "" {
		bs, err = base64.StdEncoding.DecodeString(value)
		if err != nil {
			return
		}
	}
	//fmt.Println("bs:", bs)
	service := this_.GetService(key)
	if service == nil || service.service == nil {
		err = errors.New("会话[" + key + "]不存在")
		return
	}
	_, err = service.service.Write(bs)
	return
}

func (this_ *api) systemInfo(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	service := this_.GetService(request.Key)
	if service == nil || service.service == nil {
		return
	}

	res, err = service.service.SystemInfo()

	return
}

func (this_ *api) systemMonitor(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	service := this_.GetService(request.Key)
	if service == nil || service.service == nil {
		return
	}

	res, err = service.service.SystemMonitorData()
	return
}

func (this_ *api) commandSave(r *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &TerminalCommandModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	userAgentStr := c.Request.UserAgent()
	request.Ip = c.ClientIP()
	request.UserAgent = userAgentStr
	request.UserId = r.JWT.UserId
	request.UserName = r.JWT.Name
	request.UserAccount = r.JWT.Account
	request.LoginId = r.JWT.LoginId

	err = this_.terminalCommandService.Save(request)
	if err != nil {
		return
	}
	res = request
	return
}

func (this_ *api) commandQuery(r *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &TerminalCommandModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	request.UserId = r.JWT.UserId
	res, err = this_.terminalCommandService.Query(request)
	return
}

func (this_ *api) commandCount(r *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &TerminalCommandModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	request.UserId = r.JWT.UserId
	res, err = this_.terminalCommandService.Count(request)
	return
}

func (this_ *api) commandClean(r *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &TerminalCommandModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	request.UserId = r.JWT.UserId
	err = this_.terminalCommandService.Clean(request)
	return
}

func (this_ *api) commandDelete(r *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &TerminalCommandModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	request.UserId = r.JWT.UserId
	err = this_.terminalCommandService.Delete(request.TerminalCommandId)
	return
}

type Request struct {
	Place    string `json:"place,omitempty"`
	PlaceId  string `json:"placeId,omitempty"`
	Key      string `json:"key,omitempty"`
	WorkerId string `json:"workerId"`
	LastUser string `json:"lastUser,omitempty"`
	LastDir  string `json:"lastDir,omitempty"`
	*terminal.Size
}

func (this_ *api) check(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &module_toolbox.ToolboxModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.Option == "" {
		err = errors.New("SSH 配置不存在")
		return
	}

	var config *ssh.Config
	var sshConfig *ssh.Config
	config, sshConfig, err = this_.toolboxService.GetSSHConfig(request.Option)
	if err != nil {
		return
	}
	if sshConfig != nil {
		var sshClient *goSSH.Client
		sshClient, err = ssh.NewClient(*sshConfig)
		if err != nil {
			util.Logger.Error("getSSHService ssh NewClient error", zap.Any("address", sshConfig.Address), zap.Error(err))
			return
		}
		config.SSHClient = sshClient
	}
	service := ssh.NewTerminalService(config, "", "")

	err = service.TestClient()
	return
}

func (this_ *api) close(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}
	this_.stopService(request.Key)
	return
}

func (this_ *api) changeSize(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}
	service := this_.GetService(request.Key)
	if service == nil || service.service == nil {
		return
	}
	err = service.service.ChangeSize(request.Size)
	return
}

func (this_ *api) getLogs(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = this_.WorkerFactory.getLogs(request.Place, request.PlaceId)
	return
}

func (this_ *api) deleteLog(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	path := this_.WorkerFactory.getLogPath(request.Place, request.PlaceId, request.WorkerId)
	if ex, _ := util.PathExists(path); ex {
		err = os.Remove(path)
	}
	return
}

func (this_ *api) cleanLog(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	path := this_.WorkerFactory.getLogPath(request.Place, request.PlaceId, request.WorkerId)
	if ex, _ := util.PathExists(path); !ex {
		return
	}
	f, err := os.Create(path)
	defer func() { _ = f.Close() }()
	_, err = f.WriteString("")
	return
}

func (this_ *api) downloadLog(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	this_.Logger.Info("下载日志记录 start")
	res = base.HttpNotResponse
	defer func() {
		if err != nil {
			_, _ = c.Writer.WriteString(err.Error())
		}
	}()

	request := map[string]string{}

	err = c.Bind(&request)
	if err != nil {
		return
	}

	fileName := "" + request["fileName"] + ".log"
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(fileName)))

	// 此处不设置 文件大小，如果设置文件大小，将无法终止下载
	//c.Header("Content-Length", fmt.Sprint(fileInfo.Size))
	c.Header("download-file-name", fileName)

	path := this_.WorkerFactory.getLogPath(request["place"], request["placeId"], request["workerId"])

	if ex, _ := util.PathExists(path); ex {
		var f *os.File
		f, err = os.Open(path)
		if err != nil {
			return
		}
		defer func() { _ = f.Close() }()
		_, err = io.Copy(c.Writer, f)
	} else {
		_, err = c.Writer.WriteString("暂无日志记录")
	}
	c.Status(http.StatusOK)
	return
}
