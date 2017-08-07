package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
)

var (
	DB     *sql.DB
	Driver = "mysql"
)

func Init() error {
	var err error
	if DB == nil {
		//This is so local dev will still work
		if os.Getenv("DATABASE_INSTANCE") == "" {
			DB, err = sql.Open(Driver, connectionString())
		} else {
			cfg := mysql.Cfg(os.Getenv("DATABASE_INSTANCE"), os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"))
			cfg.DBName = os.Getenv("CURT_DEV_NAME")
			cfg.ParseTime = true
			DB, err = mysql.DialCfg(cfg)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func connectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		databaseName := os.Getenv("CURT_DEV_NAME")
		if databaseName == "" {
			databaseName = "CurtData"
		}
		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, databaseName, "America%2FChicago")
	}

	return "root:@tcp(127.0.0.1:3306)/CurtData?parseTime=true&loc=America%2FChicago"
}
