package avest

var CommonEnvs = map[string]string{}
var CommonConfigs = map[string]interface{}{
	"app.hjson": map[string]interface{}{
		"debug": true,
	},
	"database.hjson": map[string]interface{}{
		"app": "sqlite3",
		"runtime": map[string]interface{}{
			"migrations": "migrations",
			"connection": "sqlite3",
		},
		"connections": map[string]interface{}{
			"sqlite3": map[string]interface{}{
				"driver": "sqlite3",
				"file":   "storage(\"test.db\")",
			},
		},
	},
	"server.hjson": map[string]interface{}{
		"address": "127.0.0.1",
		"port":    8181,
	},
}
var SimpleTemplates = map[string]string {
	"home.jet": `
{{ extends "layout.jet" }}

{{ block body() }}
  <main>
    This content will be yielded in the layout above.
  </main>
{{ end }}
	`,
	"layout.jet": `
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    {{yield body()}}
  </body>
</html>
	`,
	"error.jet": `
{{ extends "layout.jet" }}

{{ block body() }}
<main>
    This content will be yielded in the layout above.
  </main>
	`,
}
