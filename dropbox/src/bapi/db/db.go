package db 

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

//var myDB *sql.DB

type Hostinfo struct {
    DBUser,
    DBPassword,
    DBname,
    DBHost,
    DBPort,
    DBChar string
}

func connMysql(host *Hostinfo) (*sql.DB, error) {
    if host.DBHost != "" {
        host.DBHost = "tcp(" + host.DBHost + ":" + host.DBPort + ")"
    }
    db, err := sql.Open("mysql", host.DBUser+":"+host.DBPassword+"@"+host.DBHost+"/"+host.DBname+"?charset="+host.DBChar)
    return db, err
}
func SetDB(ip string) (myDB *sql.DB) {
    var server_info Hostinfo
    server_info.DBUser = "xxx"
    server_info.DBPassword = "xxx"
    server_info.DBname = "test"
    server_info.DBHost = ip
    server_info.DBPort = "xxx"
    server_info.DBChar = "utf8"
    myDB, _ = connMysql(&server_info)
    return myDB
}

func Get_data_arr(ip, tmp_sql string) []string {
    db := SetDB(ip)
    defer db.Close()
    rows, err := db.Query(tmp_sql)
    defer rows.Close()
    handleError(err)
    columns, err := rows.Columns()
    handleError(err)
    values := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }
    var tmpstr []string
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        handleError(err)
        var value string
        for _, col := range values {
            if col == nil {
                value = "NULL"
            } else {
                value = string(col)
            }
            tmpstr = append(tmpstr, value)
        }
    }
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }
    return tmpstr
}

func Get_data_map(ip, tmp_sql string) []map[string]string {
    db := SetDB(ip)
    defer db.Close()
    rows, err := db.Query(tmp_sql)
    defer rows.Close()
    handleError(err)
    columns, err := rows.Columns()
    handleError(err)

    values := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }
    var result []map[string]string

    for rows.Next() {
        err = rows.Scan(scanArgs...)
        handleError(err)
        each := make(map[string]string)

        for i, col := range values {
            each[columns[i]] = string(col)
        }

        result = append(result, each)

    }
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }
    return result
}

func handleError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

