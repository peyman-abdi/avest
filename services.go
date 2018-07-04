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
var mig services.Migrator
var mm services.ModuleManager
var r services.Router
var t services.RenderEngine
var cache services.Cache
var redis services.RedisClient

var s = new(ServicesMock)

type ServicesMock struct {
}

func (s *ServicesMock) Repository() services.Repository       { return repo }
func (s *ServicesMock) Migrator() services.Migrator           { return mig }
func (s *ServicesMock) Localization() services.Localization   { return nil }
func (s *ServicesMock) Config() services.Config               { return conf }
func (s *ServicesMock) Logger() services.Logger               { return log }
func (s *ServicesMock) Modules() services.ModuleManager       { return mm }
func (s *ServicesMock) App() services.Application             { return app }
func (s *ServicesMock) Router() services.Router               { return r }
func (s *ServicesMock) Renderer() services.RenderEngine 	  { return t }
func (s *ServicesMock) Cache() services.Cache 	  { return cache }
func (s *ServicesMock) Redis() services.RedisClient 	  { return redis }
func (s *ServicesMock) GetByName(name string) interface{} 	  { return nil }

func MockServices(configs map[string]interface{}, envs map[string]string) services.Services {
	app = application.Initialize(0, "test")
	os.MkdirAll(app.StoragePath(""), 0700)

	CreateConfigFiles(app, configs)

	conf = config.Initialize(app)
	log = logger.Initialize(conf)
	log.LoadConsole()

	repo, mig = database.Initialize(conf, log)

	t = renderer.Initialize(app, log)

	redis = redis2.Initialize(conf)
	cache = cache2.Initialize(app, conf, log, redis)

	r = router.Initialize(app, conf, log, redis, t)

	mm = modules.Initialize(conf, mig)
	mm.LoadModules(s)

	return s
}
