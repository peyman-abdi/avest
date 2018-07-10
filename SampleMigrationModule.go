package avest

import (
	"time"
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
)

type TestMigrationModel struct {
	ID int64
	MyTest string
	CreatedAt *time.Time
}
func (t *TestMigrationModel) PrimaryKey() string {
	return "id"
}
func (t *TestMigrationModel) TableName() string {
	return "tests"
}
type TestMigMigratable struct {
}
func (t *TestMigMigratable) Up(migrator services.Migratory) error {
	var err error
	if err = migrator.AutoMigrate(&TestMigrationModel{}); err != nil {
		return err
	}
	return nil
}
func (t *TestMigMigratable) Down(migrator services.Migratory) error {
	var err error
	if err = migrator.DropTableIfExists(&TestMigrationModel{}); err != nil {
		return err
	}
	return nil
}
type TestMigrationModule struct {
}
var _ services.Module = (*TestMigrationModule)(nil)
func (t *TestMigrationModule) Title() string       { return "TestMigrateModule" }
func (t *TestMigrationModule) Description() string { return "Test module" }
func (t *TestMigrationModule) Version() string     { return "1.0" }
func (t *TestMigrationModule) Activated() bool     { return true }
func (t *TestMigrationModule) Installed() bool     { return true }
func (t *TestMigrationModule) Deactivated()        { }
func (t *TestMigrationModule) Purged()             { }
func (t *TestMigrationModule) Migrations() []services.Migratable {
	return []services.Migratable {
		new(TestMigMigratable),
	}
}
func (t *TestMigrationModule) Routes() []*services.Route {
	return nil
}
func (t *TestMigrationModule) MiddleWares() []*services.MiddleWare {
	return nil
}
func (t *TestMigrationModule) GroupsHandlers() []*services.RouteGroup {
	return nil
}
func (t *TestMigrationModule) Templates() []*services.Template {
	return nil
}
func (t *TestMigrationModule) Services() map[string]func() interface{} {
	return nil
}
