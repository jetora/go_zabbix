package controllers

import (
	"bapi/models"
    "encoding/json"
	"github.com/astaxie/beego"
)

// Operations about Users
type InsviewController struct {
	beego.Controller
}

// @Title GetGraphResult 
// @Description get graph info
// @Param   body        body    models.getgraphresult true        "body for user content"
// @Success 200 {int} models.GraphResult.Ip
// @Failure 403 body is empty
// @router / [post]
func (u *InsviewController) Post() {
    var graph models.GraphResult
    json.Unmarshal(u.Ctx.Input.RequestBody, &graph)
    ugraph ,err:= models.GetGraphResult(graph)
    //handleError(err)    
    //u.Data["json"] = map[string]string{"uip": uinfo.Ip}
    if err != nil {
        u.Data["json"] = err.Error()
    } else {
        u.Data["json"] = ugraph
    }
    u.ServeJSON()
}
