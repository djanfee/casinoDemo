package main

import (
	"flag"
	"fmt"

	"casinoDemo/api/casino/internal/config"
	"casinoDemo/api/casino/internal/handler"
	"casinoDemo/api/casino/internal/svc"
	"casinoDemo/api/casino/svc/casino_svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var configFile = flag.String("f", "etc/casino.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	casinoDbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Mysql.Casino.User, c.Mysql.Casino.Pwd, c.Mysql.Casino.Host,
		c.Mysql.Casino.Port, c.Mysql.Casino.DbName)
	casinoDb, err := gorm.Open(mysql.Open(casinoDbDsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	casinoSvc := casino_svc.NewCasinoSvc(casinoDb)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c, casinoSvc, casinoDb)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
