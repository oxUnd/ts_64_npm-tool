package main

import (
    "fmt"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql" // import your used driver
)

type Components struct {
    Id          int64  `orm:"auto"`
    Name        string `orm:"size(256)"`
    Status      int64  `orm:"null"`
    Version     string `orm:"size(128)"`
    User        string `orm:"size(256)"`
    create_date string `orm:"size(32)"`
}

func init() {
    orm.RegisterModel(new(Components))
    orm.RegisterDataBase("default", "mysql", "root@/plg?charset=utf8", 30)
}

func main() {
    orm.Debug = true
    fmt.Println("TEST")
    o := orm.NewOrm()

    c := Components{Id: 640}
    err := o.Read(&c)

    fmt.Println(c)

    if err != nil {
        panic(err)
    }

    var cms []*Components

    qs := o.QueryTable("components")
    num, err := qs.All(&cms)
    if err != nil {
        panic(err)
    }
    fmt.Println(num, cms)

    for _, v := range cms {
        fmt.Println(*v)
    }
}
