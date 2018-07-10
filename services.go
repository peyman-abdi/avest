package avest

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	"os"
	"github.com/peyman-abdi/avalanche/app/modules/services/config"
	"github.com/peyman-abdi/avalanche/app/modules/services/logger"
	"github.com/peyman-abdi/avalanche/app/modules/services/database"
	"github.com/peyman-abdi/avalanche/app/modules/services/router"
	"github.com/peyman-abdi/avalanche/app/modules/services/modules"
	application "github.com/peyman-abdi/avalanche/app/modules/services/app"
	"github.com/peyman-abdi/avalanche/app/modules/services/renderer"
	redis2 "github.com/peyman-abdi/avalanche/app/modules/services/redis"
	cache2 "github.com/peyman-abdi/avalanche/app/modules/services/cache"
)

var app services.Application
var conf services.Config
var log services.Logger
var repo services.Repository
var mig services.Migratory
var mm services.ModuleManager
var r services.Router
var t services.RenderEngine
var cache services.Cache
var redis services.RedisClient

var s = new(ServicesMock)

type ServicesMock struct {
}

func (s *ServicesMock) Repository() services.Repository     { return repo }
func (s *ServicesMock) Migrator() services.Migratory        { return mig }
func (s *ServicesMock) Localization() services.Localization { return nil }
func (s *ServicesMock) Config() services.Config             { return conf }
func (s *ServicesMock) Logger() services.Logger             { return log }
func (s *ServicesMock) Modules() services.ModuleManager     { return mm }
func (s *ServicesMock) App() services.Application           { return app }
func (s *ServicesMock) Router() services.Router             { return r }
func (s *ServicesMock) Renderer() services.RenderEngine     { return t }
func (s *ServicesMock) Cache() services.Cache 	  			  { return cache }
func (s *ServicesMock) Redis() services.RedisClient 	  	  { return redis }
func (s *ServicesMock) GetByName(name string) interface{} 	  { return nil }

func MockServices(configs map[string]interface{}, envs map[string]string) services.Services {
	app = application.New(0, "test")
	os.MkdirAll(app.StoragePath(""), 0700)

	CreateConfigFiles(app, configs)

	conf = config.New(app)
	log = logger.New(conf)
	log.LoadConsole()

	repo, mig = database.New(conf, log)

	t = renderer.New(app, log)

	redis = redis2.New(conf, log)
	cache = cache2.New(app, conf, log, redis)

	r = router.New(app, conf, log, redis, t)

	mm = modules.New(conf, mig)
	mm.LoadModules(s)

	return s
}
