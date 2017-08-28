package controllers

import (
	"bapi/models"
    "encoding/json"
    "log"
	"github.com/astaxie/beego"
)

// Operations about Users
type RinfoController struct {
	beego.Controller
}

// @Title CreateRinfo
// @Description create rinfo
// @Param   body        body    models.rinfo true        "body for user content"
// @Success 200 {int} models.Rinfo.Ip
// @Failure 403 body is empty
// @router / [post]
func (u *RinfoController) Post() {
    var rinfo models.Rinfo
    json.Unmarshal(u.Ctx.Input.RequestBody, &rinfo)
    uinfo ,err:= models.GetRinfo(rinfo)
    //handleError(err)    
    //u.Data["json"] = map[string]string{"uip": uinfo.Ip}
    if err != nil {
        u.Data["json"] = err.Error()
    } else {
        u.Data["json"] = uinfo
    }
    u.ServeJSON()
}

func handleError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

