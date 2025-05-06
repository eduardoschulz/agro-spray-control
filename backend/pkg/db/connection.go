package db

import (
    "database/sql"
   _ "github.com/lib/pq"
    "log"
)

func Connect(connstr string) *sql.DB {

    db, err := sql.Open("postgres", connstr)
    if err != nil {
        log.Fatalf("falha ao conectar ao banco: %v", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatalf("banco inacessível: %v", err)
    }

    return db
}


