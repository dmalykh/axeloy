package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy"
	messageservice "github.com/dmalykh/axeloy/axeloy/message/service"
	"github.com/dmalykh/axeloy/axeloy/router"
	"github.com/dmalykh/axeloy/axeloy/way"
	"github.com/dmalykh/axeloy/axeloy/way/driver"
	wayservice "github.com/dmalykh/axeloy/axeloy/way/service"
	configuration "github.com/dmalykh/axeloy/config"
	"github.com/dmalykh/axeloy/repository/db"
	"github.com/dmalykh/axeloy/repository/db/message"
	routerrepo "github.com/dmalykh/axeloy/repository/db/router"
	wayrepo "github.com/dmalykh/axeloy/repository/db/way"
	"log"
	"os"
	"os/signal"
)

var (
	ErrParseConfig  = errors.New(`can't parse config`)
	ErrDbConnection = errors.New(`can't connect database`)
)

type App struct {
}

func (a *App) Load(ctx context.Context, configPath string) (*axeloy.Axeloy, error) {
	//Parse config
	config, err := configuration.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf(`%w %s`, ErrParseConfig, err.Error())
	}

	//Connect to database
	conn, err := db.Connect(ctx, config.Db.Driver, config.Db.Dsn)
	if err != nil {
		return nil, fmt.Errorf(`%w %s`, ErrDbConnection, err.Error())
	}

	//Create repositories
	var wayRepository = wayrepo.NewWayRepository(conn)
	var routeRepository = routerrepo.NewRouteRepository(conn)
	var trackRepository = routerrepo.NewTrackRepository(conn)
	var messageRepository = message.NewMessageRepository(conn)

	//Load way's drivers
	wayService, err := wayservice.NewService(ctx, &way.Config{
		WayRepository: wayRepository,
		Drivers: func(config *configuration.Config) map[string]driver.Config {
			var drivers = make(map[string]driver.Config)
			for name, d := range config.Ways.Drivers {
				drivers[name] = driver.Config{
					Path:   d.DriverPath,
					Config: d.DriverConfig,
				}
			}
			return drivers
		}(config),
	})
	if err != nil {
		return nil, err
	}

	//Load services
	var messageService = messageservice.NewMessager(messageRepository)
	var routerService = router.NewRouter(routeRepository, wayService)
	var trackService = router.NewTracker(trackRepository, wayService, messageService)
	var ax = axeloy.New(&axeloy.Config{
		Router:   routerService,
		Tracker:  trackService,
		Messager: messageService,
		Wayer:    wayService,
	})
	return ax, nil
}

//Graceful shutdown https://play.golang.org/p/uBMCywO5O0w
func (a *App) NewContext() context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()
	return ctx
}
