package main

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/wmsx/xconf/config-srv/broadcast"
	"github.com/wmsx/xconf/config-srv/broadcast/broker"
	"github.com/wmsx/xconf/config-srv/broadcast/database"
	"github.com/wmsx/xconf/config-srv/conf"
	"github.com/wmsx/xconf/config-srv/dao"
	"github.com/wmsx/xconf/config-srv/handler"
	protoConfig "github.com/wmsx/xconf/proto/config"
)

const name  = "wm.sx.svc.config"

var config conf.Config

func main() {
	log.Name("xconf")

	service := micro.NewService(
		micro.Name(name),
		micro.Transport(grpc.NewTransport()),
		micro.Flags(
			cli.StringFlag{
				Name:   "database_driver",
				Usage:  "database driver",
				EnvVar: "DATABASE_DRIVER",
				Value:  "mysql",
			},
			cli.StringFlag{
				Name:   "database_url",
				Usage:  "database url",
				EnvVar: "DATABASE_URL",
				Value:  "root:12345@(127.0.0.1:3306)/xconf?charset=utf8&parseTime=true&loc=Local",
			},
			cli.StringFlag{
				Name:   "broadcast",
				Usage:  "broadcast (db/broker)",
				EnvVar: "BROADCAST",
				Value:  "db",
			}),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			config.DB.DriverName = c.String("database_driver")
			config.DB.URL = c.String("database_url")
			config.BroadcastType = c.String("broadcast")
			log.Infof("database_driver: %s , database_url: %s\n", config.DB.DriverName, config.DB.URL)
		}),
		micro.BeforeStart(func() (err error) {
			if err = dao.Init(&config); err != nil {
				return
			}
			if err = dao.GetDao().Ping(); err != nil {
				return
			}

			var bc broadcast.Broadcast
			switch config.BroadcastType {
			case "db":
				bc, err = database.New()
				if err != nil {
					return err
				}
			case "broker":
				bc, err = broker.New(service)
				if err != nil {
					return err
				}
			default:
				return errors.New("broadcast： Invalid option")
			}
			broadcast.Init(bc)
			return
		}),
		micro.BeforeStop(func() error {
			return dao.GetDao().Disconnect()
		}),
	)

	if err := protoConfig.RegisterConfigHandler(service.Server(), new(handler.Config)); err != nil {
		panic(err)
	}

	if err := service.Run(); err != nil {
		panic(err)
	}
}
