package template4app

var (
	SettingInvalid = `
[server]
root_url = %(protocol)s://%(domain)s:%(port)s/innovation/
alt_url = https://innovation.com/

`
	SettingOverride = `
[paths]
data = /tmp/override


`
	SettingOverrideWindows = `
[paths]
data = c:\tmp\override


`
	SettingSession = `
[session]
provider = file

`
)
