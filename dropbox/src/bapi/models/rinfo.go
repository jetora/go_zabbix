package models

import (
    "errors"
    "bapi/db"
    "bapi/ssh"
    "log"
    "strconv"
    "strings"
)

var (
    Rinfos map[string]*Rinfo
)

type Rinfo struct {
    Ip string
    Ismaster    string
    Sessmode    string
    Filemode    string
    Sysmode string
}

func GetRinfo(u Rinfo) (rinfo *Rinfo, err error) {
    t_sql1 := "show global variables like 'read_only'"
    t_sql2 := "show slave status"
    cmd1 := "cat /export/servers/mysql/etc/my.cnf |grep read_only|awk -F '=' '{print $2}'"
    cmd2 := "touch /export/data/mysql/dumps/aa"
    
    tsysmode, err1 := ssh.Ssh(u.Ip, cmd1)
    handleError(err1)
    if err1 != nil {
        return nil,errors.New("Server Cannot Connect...")
    } else {
        t_arr := db.Get_data_arr(u.Ip, t_sql1)
        u.Sessmode=t_arr[1]
        tmp_sysmode, err2 := strconv.Atoi(strings.Replace(strings.Replace(tsysmode, " ", "", -1), "\n", "", -1))
        handleError(err2)
        if tmp_sysmode == 1 {
            u.Sysmode = "ON"
        } else {
            u.Sysmode = "OFF"
        }
        _, err3 := ssh.Ssh(u.Ip, cmd2)
        if err3 != nil {
            u.Filemode = "ON"
        } else {
            u.Filemode = "OFF"
        }
        arr_tmp := db.Get_data_map(u.Ip,t_sql2)
        if len(arr_tmp) == 0 {
            u.Ismaster = "Y"
        } else if arr_tmp[0]["Master_Host"] == "1.1.1.1"{
            u.Ismaster = "Y"   
        } else {
            u.Ismaster = "N"
        }
        return &u,nil
    }
}

func handleError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

