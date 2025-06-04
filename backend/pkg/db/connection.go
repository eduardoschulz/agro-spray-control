package db

import (
   _ "database/sql"
    "github.com/jmoiron/sqlx"

   _ "github.com/lib/pq"
    "log"
)

func Connect(connstr string) *sqlx.DB {

    db, err := sqlx.Open("postgres", connstr)
    if err != nil {
        log.Fatalf("falha ao conectar ao banco: %v", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatalf("banco inacessível: %v", err)
    }

    return db
}


