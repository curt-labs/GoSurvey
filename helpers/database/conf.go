package database

import (
	"fmt"
	"os"
)

func ConnectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		return fmt.Sprintf("%s:%s@%s(%s)/CurtDev?parseTime=true&loc=America%2FChicago", user, pass, proto, addr)
	}

	return "root:@tcp(127.0.0.1:3306)/CurtDev?parseTime=true&loc=America%2FChicago"
}
