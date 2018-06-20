package avest

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"os"
	"github.com/peyman-abdi/avalanche/app/modules/core/config"
	"github.com/peyman-abdi/avalanche/app/modules/core/logger"
	"github.com/peyman-abdi/avalanche/app/modules/core/database"
	"github.com/peyman-abdi/avalanche/app/modules/core/router"
	"github.com/peyman-abdi/avalanche/app/modules/core/modules"
	application "github.com/peyman-abdi/avalanche/app/modules/core/app"
	"github.com/peyman-abdi/avalanche/app/modules/core/template"
)

var app core.Application
var conf core.Config
var log core.Logger
var repo core.Repository
var mig core.Migrator
var mm core.ModuleManager
var r core.Router
var t core.TemplateEngine

var s = new(ServicesMock)

type ServicesMock struct {
}

func (s *ServicesMock) Repository() core.Repository     { return repo }
func (s *ServicesMock) Migrator() core.Migrator         { return mig }
func (s *ServicesMock) Localization() core.Localization { return nil }
func (s *ServicesMock) Config() core.Config             { return conf }
func (s *ServicesMock) Logger() core.Logger             { return log }
func (s *ServicesMock) Modules() core.ModuleManager     { return mm }
func (s *ServicesMock) App() core.Application           { return app }
func (s *ServicesMock) Router() core.Router             { return r }
func (s *ServicesMock) TemplateEngine() core.TemplateEngine             { return t }

func MockServices(configs map[string]interface{}, envs map[string]string) core.Services {
	app = application.Initialize(0, "test")
	os.MkdirAll(app.StoragePath(""), 0700)

	CreateConfigFiles(app, configs)

	conf = config.Initialize(app)
	log = logger.Initialize(conf)
	log.LoadConsole()

	repo, mig = database.Initialize(conf, log)

	t = template.Initialize(app, log)
	r = router.Initialize(conf, log, t)

	mm = modules.Initialize(conf, mig)
	mm.LoadModules(s)

	return s
}
