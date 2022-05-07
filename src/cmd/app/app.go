package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy"
	messageservice "github.com/dmalykh/axeloy/axeloy/message/service"
	"github.com/dmalykh/axeloy/axeloy/router"
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
	ErrLoadDrivers  = errors.New(`can't load ways driver`)
	ErrOpenDrivers  = errors.New(`can't open ways driver`)
)

func NewApp() *App {
	return &App{}
}

type App struct {
}

func (a *App) Open(configPath string) (*configuration.Config, error) {
	//Parse config
	config, err := configuration.LoadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf(`%w %s`, ErrParseConfig, err.Error())
	}
	return config, nil
}

// Load application creates repositories and run ways
func (a *App) Load(ctx context.Context, config *configuration.Config) (*axeloy.Axeloy, error) {
	//Connect to database
	conn, err := db.Connect(ctx, config.Database.Driver, config.Database.Dsn)
	if err != nil {
		return nil, fmt.Errorf(`%w %s`, ErrDbConnection, err.Error())
	}

	var reform = db.Reform(config.Database.Driver, conn)

	// Create repositories
	var wayRepository = wayrepo.NewWayRepository(reform)
	var routeRepository = routerrepo.NewRouteRepository(reform)
	var trackRepository = routerrepo.NewTrackRepository(reform)
	var messageRepository = message.NewMessageRepository(reform)

	// LoadFile way's drivers
	driverService := wayservice.NewDriverService()
	var drivers = make(map[string]driver.Driver, len(config.Drivers))
	for _, driverConfig := range config.Drivers {
		drv, err := driver.Open(ctx, driverConfig.Path, func(v interface{}) error {
			return configuration.Unmarshal(driverConfig, v)
		})
		if err != nil {
			return nil, fmt.Errorf(`%w %s`, ErrOpenDrivers, err.Error())
		}
		drivers[driverConfig.Name] = drv
	}
	if err := driverService.Load(drivers); err != nil {
		return nil, fmt.Errorf(`%w %s`, ErrLoadDrivers, err.Error())
	}

	// LoadFile ways
	var wayService = wayservice.NewService(wayRepository, driverService)

	//LoadFile services
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
func (a *App) WithShutdown(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()
	return ctx
}
