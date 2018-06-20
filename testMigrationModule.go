package avest

import (
	"time"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
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
func (t *TestMigMigratable) Up(migrator core.Migrator) error {
	var err error
	if err = migrator.AutoMigrate(&TestMigrationModel{}); err != nil {
		return err
	}
	return nil
}
func (t *TestMigMigratable) Down(migrator core.Migrator) error {
	var err error
	if err = migrator.DropTableIfExists(&TestMigrationModel{}); err != nil {
		return err
	}
	return nil
}
type TestMigrationModule struct {
}
var _ core.Module = (*TestMigrationModule)(nil)
func (t *TestMigrationModule) Title() string       { return "TestMigrateModule" }
func (t *TestMigrationModule) Description() string { return "Test module" }
func (t *TestMigrationModule) Version() string     { return "1.0" }
func (t *TestMigrationModule) Activated() bool     { return true }
func (t *TestMigrationModule) Installed() bool     { return true }
func (t *TestMigrationModule) Deactivated()        { }
func (t *TestMigrationModule) Purged()             { }
func (t *TestMigrationModule) Migrations() []core.Migratable {
	return []core.Migratable {
		new(TestMigMigratable),
	}
}
func (t *TestMigrationModule) Routes() []*core.Route {
	return nil
}
func (t *TestMigrationModule) MiddleWares() []*core.MiddleWare {
	return nil
}
func (t *TestMigrationModule) GroupsHandlers() []*core.RouteGroup {
	return nil
}
func (t *TestMigrationModule) Templates() []*core.Template {
	return nil
}