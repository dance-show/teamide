package module_toolbox

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
	"teamide/internal/context"
	"teamide/pkg/db"
	"teamide/pkg/toolbox"
)

type ToolboxApi struct {
	*context.ServerContext
	ToolboxService *ToolboxService
}

func NewToolboxApi(ToolboxService *ToolboxService) *ToolboxApi {
	return &ToolboxApi{
		ServerContext:  ToolboxService.ServerContext,
		ToolboxService: ToolboxService,
	}
}

var (
	// 工具 权限

	// PowerToolbox 工具基本 权限
	PowerToolbox          = base.AppendPower(&base.PowerAction{Action: "toolbox", Text: "工具", ShouldLogin: true, StandAlone: true})
	PowerToolboxPage      = base.AppendPower(&base.PowerAction{Action: "toolbox_page", Text: "工具页面", Parent: PowerToolbox, ShouldLogin: true, StandAlone: true})
	PowerToolboxContext   = base.AppendPower(&base.PowerAction{Action: "toolbox_context", Text: "工具上下文", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxInsert    = base.AppendPower(&base.PowerAction{Action: "toolbox_insert", Text: "工具新增", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxUpdate    = base.AppendPower(&base.PowerAction{Action: "toolbox_update", Text: "工具修改", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxRename    = base.AppendPower(&base.PowerAction{Action: "toolbox_rename", Text: "工具重命名", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxDelete    = base.AppendPower(&base.PowerAction{Action: "toolbox_delete", Text: "工具删除", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxMoveGroup = base.AppendPower(&base.PowerAction{Action: "toolbox_move_group", Text: "工具分组", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})

	PowerToolboxGroupInsert = base.AppendPower(&base.PowerAction{Action: "toolbox_group_insert", Text: "工具分组新增", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxGroupUpdate = base.AppendPower(&base.PowerAction{Action: "toolbox_group_update", Text: "工具分组修改", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxGroupDelete = base.AppendPower(&base.PowerAction{Action: "toolbox_group_delete", Text: "工具分组删除", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})

	PowerToolboxOpen                   = base.AppendPower(&base.PowerAction{Action: "toolbox_open", Text: "工具打开", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxUpdateOpenExtend       = base.AppendPower(&base.PowerAction{Action: "toolbox_update_open_extend", Text: "工具打开", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxQueryOpens             = base.AppendPower(&base.PowerAction{Action: "toolbox_query_opens", Text: "工具查询打开", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxClose                  = base.AppendPower(&base.PowerAction{Action: "toolbox_close", Text: "工具关闭", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxOpenTab                = base.AppendPower(&base.PowerAction{Action: "toolbox_open_tab", Text: "工具打开", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxQueryOpenTabs          = base.AppendPower(&base.PowerAction{Action: "toolbox_query_open_tabs", Text: "工具查询打开", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxCloseTab               = base.AppendPower(&base.PowerAction{Action: "toolbox_close_tab", Text: "工具关闭", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxUpdateOpenTabExtend    = base.AppendPower(&base.PowerAction{Action: "toolbox_update_open_tab_extend", Text: "工具打开", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxWork                   = base.AppendPower(&base.PowerAction{Action: "toolbox_work", Text: "工具工作", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxSSHShell               = base.AppendPower(&base.PowerAction{Action: "toolbox_ssh_shell", Text: "工具SSH Shell连接", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxSSHFtp                 = base.AppendPower(&base.PowerAction{Action: "toolbox_ssh_ftp", Text: "工具SSH FTP连接", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxSSHFtpUpload           = base.AppendPower(&base.PowerAction{Action: "toolbox_ssh_ftp_upload", Text: "工具SSH FTP上传", Parent: PowerToolboxSSHFtp, ShouldLogin: true, StandAlone: true})
	PowerToolboxSSHFtpDownload         = base.AppendPower(&base.PowerAction{Action: "toolbox_ssh_ftp_download", Text: "工具SSH FTP下载", Parent: PowerToolboxSSHFtp, ShouldLogin: true, StandAlone: true})
	PowerToolboxDatabaseExportDownload = base.AppendPower(&base.PowerAction{Action: "toolbox_ssh_ftp_download", Text: "工具SSH FTP下载", Parent: PowerToolboxSSHFtp, ShouldLogin: true, StandAlone: true})
)

func (this_ *ToolboxApi) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox"}, Power: PowerToolbox, Do: this_.index})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/page"}, Power: PowerToolboxPage, Do: this_.context})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/context"}, Power: PowerToolboxContext, Do: this_.context})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/insert"}, Power: PowerToolboxInsert, Do: this_.insert})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/update"}, Power: PowerToolboxUpdate, Do: this_.update})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/rename"}, Power: PowerToolboxRename, Do: this_.rename})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/delete"}, Power: PowerToolboxDelete, Do: this_.delete})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/moveGroup"}, Power: PowerToolboxMoveGroup, Do: this_.moveGroup})

	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/group/insert"}, Power: PowerToolboxGroupInsert, Do: this_.insertGroup})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/group/update"}, Power: PowerToolboxGroupUpdate, Do: this_.updateGroup})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/group/delete"}, Power: PowerToolboxGroupDelete, Do: this_.deleteGroup})

	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/work"}, Power: PowerToolboxWork, Do: this_.work})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/open"}, Power: PowerToolboxOpen, Do: this_.open})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/queryOpens"}, Power: PowerToolboxQueryOpens, Do: this_.queryOpens})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/close"}, Power: PowerToolboxClose, Do: this_.close})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/updateOpenExtend"}, Power: PowerToolboxUpdateOpenExtend, Do: this_.updateOpenExtend})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/queryOpenTabs"}, Power: PowerToolboxQueryOpenTabs, Do: this_.queryOpenTabs})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/openTab"}, Power: PowerToolboxOpenTab, Do: this_.openTab})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/closeTab"}, Power: PowerToolboxCloseTab, Do: this_.closeTab})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/updateOpenTabExtend"}, Power: PowerToolboxUpdateOpenTabExtend, Do: this_.updateOpenTabExtend})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/ssh/shell"}, Power: PowerToolboxSSHShell, Do: this_.sshShell, IsWebSocket: true})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/ssh/ftp"}, Power: PowerToolboxSSHFtp, Do: this_.sshFtp, IsWebSocket: true})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/ssh/ftp/upload"}, Power: PowerToolboxSSHFtpUpload, Do: this_.sshFtpUpload})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/ssh/ftp/download"}, Power: PowerToolboxSSHFtpDownload, Do: this_.sshFtpDownload, IsGet: true})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/database/export/download"}, Power: PowerToolboxDatabaseExportDownload, Do: this_.databaseExportDownload, IsGet: true})

	return
}

type IndexResponse struct {
	Types                    []*toolbox.Worker                  `json:"types,omitempty"`
	SqlConditionalOperations []*toolbox.SqlConditionalOperation `json:"sqlConditionalOperations,omitempty"`
	MysqlColumnTypeInfos     []*db.ColumnTypeInfo               `json:"mysqlColumnTypeInfos,omitempty"`
	DatabaseTypes            []*db.DatabaseType                 `json:"databaseTypes,omitempty"`
}

func (this_ *ToolboxApi) index(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	response := &IndexResponse{}

	response.Types = toolbox.GetWorkers()
	response.SqlConditionalOperations = toolbox.SqlConditionalOperations
	response.MysqlColumnTypeInfos = db.MySqlColumnTypeInfos
	response.DatabaseTypes = db.DatabaseTypes

	res = response
	return
}

func (this_ *ToolboxApi) page(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	return
}

type ContextRequest struct {
}

type ContextResponse struct {
	Context map[string][]*ToolboxModel `json:"context,omitempty"`
	Groups  []*ToolboxGroupModel       `json:"groups,omitempty"`
}

func (this_ *ToolboxApi) context(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ContextRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &ContextResponse{}

	response.Groups, err = this_.ToolboxService.QueryGroup(&ToolboxGroupModel{
		UserId: requestBean.JWT.UserId,
	})
	if err != nil {
		return
	}

	list, err := this_.ToolboxService.Query(&ToolboxModel{
		UserId: requestBean.JWT.UserId,
	})
	if err != nil {
		return
	}
	response.Context = make(map[string][]*ToolboxModel)
	for _, one := range list {
		response.Context[one.ToolboxType] = append(response.Context[one.ToolboxType], one)
	}
	response.Context["other"] = append(response.Context["other"], Others...)
	res = response
	return
}
