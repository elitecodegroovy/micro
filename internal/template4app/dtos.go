package template4app

var (
	DtosModels = `
package dtos

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"

	"{{.Dir}}/pkg/components/simplejson"
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
)

type AnyId struct {
	Id int64 ` +"`json:\"id\"`"+`
}

type LoginCommand struct {
	User     string ` +"`json:\"user\" binding:\"Required\"`"+`
	Password string ` +"`json:\"password\" binding:\"Required\"`"+`
	Remember bool   ` +"`json:\"remember\"`"+`
}

type CurrentUser struct {
	IsSignedIn                 bool         ` +"`json:\"isSignedIn\"`"+`
	Id                         int64        ` +"`json:\"id\"`"+`
	Login                      string       ` +"`json:\"login\"`"+`
	Email                      string       ` +"`json:\"email\"`"+`
	Name                       string       ` +"`json:\"name\"`"+`
	LightTheme                 bool         ` +"`json:\"lightTheme\"`"+`
	OrgCount                   int          ` +"`json:\"orgCount\"`"+`
	OrgId                      int64        ` +"`json:\"orgId\"`"+`
	OrgName                    string       ` +"`json:\"orgName\"`"+`
	OrgRole                    m.RoleType   ` +"`json:\"orgRole\"`"+`
	IsGrafanaAdmin             bool         ` +"`json:\"isGrafanaAdmin\"`"+`
	GravatarUrl                string       ` +"`json:\"gravatarUrl\"`"+`
	Timezone                   string       ` +"`json:\"timezone\"`"+`
	Locale                     string       ` +"`json:\"locale\"`"+`
	HelpFlags1                 m.HelpFlags1 ` +"`json:\"helpFlags1\"`"+`
	HasEditPermissionInFolders bool         ` +"`json:\"hasEditPermissionInFolders\"`"+`
}

type MetricRequest struct {
	From    string             ` +"`json:\"from\"`"+`
	To      string             ` +"`json:\"to\"`"+`
	Queries []*simplejson.Json ` +"`json:\"queries\"`"+`
	Debug   bool               ` +"`json:\"debug\"`"+`
}

type UserStars struct {
	DashboardIds map[string]bool `+"`json:\"dashboardIds\"`"+`
}

func GetGravatarUrl(text string) string {
	if setting.DisableGravatar {
		return setting.AppSubUrl + "/public/img/user_profile.png"
	}

	if text == "" {
		return ""
	}

	hasher := md5.New()
	hasher.Write([]byte(strings.ToLower(text)))
	return fmt.Sprintf(setting.AppSubUrl+"/avatar/%x", hasher.Sum(nil))
}

func GetGravatarUrlWithDefault(text string, defaultText string) string {
	if text != "" {
		return GetGravatarUrl(text)
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	if err != nil {
		return ""
	}

	text = reg.ReplaceAllString(defaultText, "") + "@localhost"

	return GetGravatarUrl(text)
}

`
	DtosUser = `package dtos

type SignUpForm struct {
	Email string ` +"`json:\"email\" binding:\"Required\"`"+`
}

type SignUpStep2Form struct {
	Email    string ` +"`json:\"email\"`"+`
	Name     string ` +"`json:\"name\"`"+`
	Username string ` +"`json:\"username\"`"+`
	Password string ` +"`json:\"password\"`"+`
	Code     string ` +"`json:\"code\"`"+`
	OrgName  string ` +"`json:\"orgName\"`"+`
}

type AdminCreateUserForm struct {
	Email    string ` +"`json:\"email\"`"+`
	Login    string ` +"`json:\"login\"`"+`
	Name     string ` +"`json:\"name\"`"+`
	Password string ` +"`json:\"password\" binding:\"Required\"`"+`
}

type AdminUpdateUserForm struct {
	Email string ` +"`json:\"email\"`"+`
	Login string ` +"`json:\"login\"`"+`
	Name  string ` +"`json:\"name\"`"+`
}

type AdminUpdateUserPasswordForm struct {
	Password string ` +"`json:\"password\" binding:\"Required\"`"+`
}

type AdminUpdateUserPermissionsForm struct {
	IsGrafanaAdmin bool ` +"`json:\"isGrafanaAdmin\"`"+`
}

type AdminUserListItem struct {
	Email          string ` +"`json:\"email\"`"+`
	Name           string ` +"`json:\"name\"`"+`
	Login          string ` +"`json:\"login\"`"+`
	IsGrafanaAdmin bool   ` +"`json:\"isGrafanaAdmin\"`"+`
}

type SendResetPasswordEmailForm struct {
	UserOrEmail string ` +"`json:\"userOrEmail\" binding:\"Required\"`"+`
}

type ResetUserPasswordForm struct {
	Code            string ` +"`json:\"code\"`"+`
	NewPassword     string ` +"`json:\"newPassword\"`"+`
	ConfirmPassword string ` +"`json:\"confirmPassword\"`"+`
}

`

)
