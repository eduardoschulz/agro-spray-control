package main

import (

    "github.com/eduardoschulz/agro-spray-control/backend/pkg/db"
    "fmt"
)

type Banco struct {
    Addr string
    User string
    Passwd string
    Port int16
    DBName string
    SSL    string
}


func main(){


    /* Isso tem q ser movido para outro lugar posteriormente */
     b := Banco {
        Addr: "localhost",
        User: "postgres",
        Passwd: "postgres", 
        Port: 5432,
        DBName: "go",
        SSL: "disable",
    }


    connstr := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        b.Addr, b.Port, b.User, b.Passwd, b.DBName, b.SSL,
    )
    db.Connect(connstr)

}
