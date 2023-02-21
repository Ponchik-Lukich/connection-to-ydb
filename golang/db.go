package main

import (
	"context"
	"fmt"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	yc "github.com/ydb-platform/ydb-go-yc"
	"reflect"
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

func createTableExample(ctx context.Context, db ydb.Connection, cfg Config) {
	err := db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			return s.CreateTable(ctx, cfg.Database+"/series",
				options.WithColumn("series_id", types.TypeUint64), // not null column
				options.WithColumn("title", types.Optional(types.TypeUTF8)),
				options.WithColumn("release_date", types.Optional(types.TypeDate)),
				options.WithPrimaryKeyColumn("series_id"),
			)
		},
	)
	if err != nil {
		panic(err)
	}
}

func SelectExample(ctx context.Context, db ydb.Connection, cfg Config) {
	err := db.Table().Do(ctx, func(ctx context.Context, s table.Session) (err error) {
		query := `SELECT * FROM series`
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
		panic(err)
	}
}

func main() {
	var cfg Config
	cfg.Endpoint = "ydb.serverless... your endpoint without grpcs://"
	cfg.Database = "/ru... your database"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	fmt.Println()
	db, err := ydb.Open(ctx,
		sugar.DSN(cfg.Endpoint, cfg.Database, true),
		yc.WithInternalCA(),
		yc.WithServiceAccountKeyFileCredentials("./authorized_key.json"),
	)
	fmt.Println(reflect.TypeOf(db))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// create table
	//createTableExample(ctx, db, cfg)
	SelectExample(ctx, db, cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close(ctx)
	}()
}
