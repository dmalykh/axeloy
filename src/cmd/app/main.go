package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy"
	messageservice "github.com/dmalykh/axeloy/axeloy/message/service"
	"github.com/dmalykh/axeloy/axeloy/router"
	"github.com/dmalykh/axeloy/axeloy/way"
	configuration "github.com/dmalykh/axeloy/config"
	"github.com/dmalykh/axeloy/repository/db"
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

func load(ctx context.Context, configPath string) (*axeloy.Axeloy, error) {
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

	//Load way's drivers
	wayService, err := way.NewService(ctx, &way.Config{
		WayRepository: wayrepo.NewWayRepository(conn),
		Drivers: func(config *configuration.Config) map[string]way.DriverConfig {
			var drivers = make(map[string]way.DriverConfig)
			for name, driver := range config.Ways.Drivers {
				drivers[name] = way.DriverConfig{
					DriverPath: driver.DriverPath,
					Params:     driver.Params,
				}
			}
			return drivers
		}(config),
	})
	if err != nil {
		return nil, err
	}

	var messageService = messageservice.NewMessager()
	var routerService = router.NewRouter(routerrepo.NewRouteRepository(conn), wayService)
	var trackService = router.NewTracker(routerrepo.NewTrackRepository(conn), wayService, messageService)
	var ax = axeloy.New(&axeloy.Config{
		Router:   routerService,
		Tracker:  trackService,
		Messager: messageService,
		Wayer:    wayService,
	})
	return ax, nil
}

//https://play.golang.org/p/uBMCywO5O0w
func newContext() context.Context {
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

func main() {
	var ctx = newContext()
	if err := ax.Run(ctx); err != nil {
		log.Fatalln(err.Error())
	}
}
