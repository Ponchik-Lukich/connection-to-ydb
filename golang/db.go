package main

import (
	"context"
	"fmt"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	yc "github.com/ydb-platform/ydb-go-yc"
	"time"
)

type Department struct {
	Id   string
	Name string
}

type Config struct {
	Endpoint string
	Database string
}

func main() {
	var cfg Config
	cfg.Endpoint = "ydb.serverless... your endpoint without grpcs://"
	cfg.Database = "/ru... your database"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	db, err := ydb.Open(ctx,
		sugar.DSN(cfg.Endpoint, cfg.Database, true),
		yc.WithInternalCA(),
		yc.WithServiceAccountKeyFileCredentials("./authorized_key.json"),
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		query := `SELECT * FROM departments`
		_, res, err := s.Execute(ctx, table.DefaultTxControl(), query, table.NewQueryParameters())
		if err != nil {
			return err
		}
		fmt.Println(res.ResultSetCount())
		for res.NextResultSet(ctx) {
			for res.NextRow() {
				department := &Department{}
				err := res.ScanWithDefaults(
					&department.Id,
					&department.Name,
				)
				if err != nil {
					return err
				}
				fmt.Println(*department)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer func() {
		_ = db.Close(ctx)
	}()
}
