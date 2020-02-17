package template4app

var (
	ModelApikey = `
package models

import (
	"errors"
	"time"
)

var ErrInvalidApiKey = errors.New("Invalid API Key")
var ErrInvalidApiKeyExpiration = errors.New("Negative value for SecondsToLive")

type ApiKey struct {
	Id      int64
	OrgId   int64
	Name    string
	Key     string
	Role    RoleType
	Created time.Time
	Updated time.Time
	Expires *int64
}

// ---------------------
// COMMANDS
type AddApiKeyCommand struct {
	Name          string   ` +"`json:\"name\" binding:\"Required\"`"+`
	Role          RoleType ` +"`json:\"role\" binding:\"Required\"`"+ `
	OrgId         int64    ` +"`json:\"-\"`"+`
	Key           string   ` +"`json:\"-\"`"+`
	SecondsToLive int64    ` + "`json:\"secondsToLive\"`"+`

	Result *ApiKey ` +"`json:\"-\"`"+`
}

type UpdateApiKeyCommand struct {
	Id   int64    `+"`json:\"id\"`"+`
	Name string   ` + "`json:\"name\"`"+`
	Role RoleType ` +"`json:\"role\"`"+`

	OrgId int64 ` +"`json:\"-\"`"+`
}

type DeleteApiKeyCommand struct {
	Id    int64 ` +"`json:\"id\"`"+`
	OrgId int64 ` +"`json:\"-\"`"+`
}

// ----------------------
// QUERIES

type GetApiKeysQuery struct {
	OrgId          int64
	IncludeInvalid bool
	Result         []*ApiKey
}

type GetApiKeyByNameQuery struct {
	KeyName string
	OrgId   int64
	Result  *ApiKey
}

type GetApiKeyByIdQuery struct {
	ApiKeyId int64
	Result   *ApiKey
}

// ------------------------
// DTO & Projections

type ApiKeyDTO struct {
	Id         int64      `+"`json:\"id\"`"+`
	Name       string     ` +"`json:\"name\"`"+`
	Role       RoleType   ` +"`json:\"role\"`"+`
	Expiration *time.Time ` +"`json:\"expiration,omitempty\"`"+`
}

`
	ModelContext = `
package models

import (
	"strings"

	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/setting"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/macaron.v1"
)

type ReqContext struct {
	*macaron.Context
	*SignedInUser
	UserToken *UserToken

	IsSignedIn     bool
	IsRenderCall   bool
	AllowAnonymous bool
	SkipCache      bool
	Logger         log.Logger
}

// Handle handles and logs error by given status.
func (ctx *ReqContext) Handle(status int, title string, err error) {
	if err != nil {
		ctx.Logger.Error(title, "error", err)
		if setting.Env != setting.PROD {
			ctx.Data["ErrorMsg"] = err
		}
	}

	ctx.Data["Title"] = title
	ctx.Data["AppSubUrl"] = setting.AppSubUrl
	ctx.Data["Theme"] = "dark"

	ctx.HTML(status, setting.ERR_TEMPLATE_NAME)
}

func (ctx *ReqContext) JsonOK(message string) {
	resp := make(map[string]interface{})
	resp["message"] = message
	ctx.JSON(200, resp)
}

func (ctx *ReqContext) IsApiRequest() bool {
	return strings.HasPrefix(ctx.Req.URL.Path, "/api")
}

func (ctx *ReqContext) JsonApiErr(status int, message string, err error) {
	resp := make(map[string]interface{})

	if err != nil {
		ctx.Logger.Error(message, "error", err)
		if setting.Env != setting.PROD {
			resp["error"] = err.Error()
		}
	}

	switch status {
	case 404:
		resp["message"] = "Not Found"
	case 500:
		resp["message"] = "Internal Server Error"
	}

	if message != "" {
		resp["message"] = message
	}

	ctx.JSON(status, resp)
}

func (ctx *ReqContext) HasUserRole(role RoleType) bool {
	return ctx.OrgRole.Includes(role)
}

func (ctx *ReqContext) HasHelpFlag(flag HelpFlags1) bool {
	return ctx.HelpFlags1.HasFlag(flag)
}

func (ctx *ReqContext) TimeRequest(timer prometheus.Summary) {
	ctx.Data["perfmon.timer"] = timer
}

`
	ModelHealth = `
package models

type GetDBHealthQuery struct{}

`
	ModelHelpFlags = `
package models

type HelpFlags1 uint64

const (
	HelpFlagGettingStartedPanelDismissed HelpFlags1 = 1 << iota
	HelpFlagDashboardHelp1
)

func (f HelpFlags1) HasFlag(flag HelpFlags1) bool { return f&flag != 0 }
func (f *HelpFlags1) AddFlag(flag HelpFlags1)     { *f |= flag }
func (f *HelpFlags1) ClearFlag(flag HelpFlags1)   { *f &= ^flag }
func (f *HelpFlags1) ToggleFlag(flag HelpFlags1)  { *f ^= flag }

type SetUserHelpFlagCommand struct {
	HelpFlags1 HelpFlags1
	UserId     int64
}

`
	ModelLoginAttempt = `
package models

import (
	"time"
)

type LoginAttempt struct {
	Id        int64
	Username  string
	IpAddress string
	Created   int64
}

// ---------------------
// COMMANDS

type CreateLoginAttemptCommand struct {
	Username  string
	IpAddress string

	Result LoginAttempt
}

type DeleteOldLoginAttemptsCommand struct {
	OlderThan   time.Time
	DeletedRows int64
}

// ---------------------
// QUERIES

type GetUserLoginAttemptCountQuery struct {
	Username string
	Since    time.Time
	Result   int64
}

`
	ModelNotifications = `
package models

import "errors"

var ErrInvalidEmailCode = errors.New("Invalid or expired email code")
var ErrSmtpNotEnabled = errors.New("SMTP not configured, check your grafana.ini config file's [smtp] section")

type SendEmailCommand struct {
	To           []string
	Template     string
	Subject      string
	Data         map[string]interface{}
	Info         string
	EmbededFiles []string
}

type SendEmailCommandSync struct {
	SendEmailCommand
}

type SendWebhookSync struct {
	Url         string
	User        string
	Password    string
	Body        string
	HttpMethod  string
	HttpHeader  map[string]string
	ContentType string
}

type SendResetPasswordEmailCommand struct {
	User *User
}

type ValidateResetPasswordCodeQuery struct {
	Code   string
	Result *User
}

`
	ModelOrg = `
package models

import (
	"errors"
	"time"
)

// Typed errors
var (
	ErrOrgNotFound  = errors.New("Organization not found")
	ErrOrgNameTaken = errors.New("Organization name is taken")
)

type Org struct {
	Id      int64
	Version int
	Name    string

	Address1 string
	Address2 string
	City     string
	ZipCode  string
	State    string
	Country  string

	Created time.Time
	Updated time.Time
}

// ---------------------
// COMMANDS

type CreateOrgCommand struct {
	Name string ` +"`json:\"name\" binding:\"Required\"`"+`

	// initial admin user for account
	UserId int64 ` +"`json:\"-\"`"+`
	Result Org   ` +"`json:\"-\"`"+`
}

type DeleteOrgCommand struct {
	Id int64
}

type UpdateOrgCommand struct {
	Name  string
	OrgId int64
}

type UpdateOrgAddressCommand struct {
	OrgId int64
	//Address
}

type GetOrgByIdQuery struct {
	Id     int64
	Result *Org
}

type GetOrgByNameQuery struct {
	Name   string
	Result *Org
}

type SearchOrgsQuery struct {
	Query string
	Name  string
	Limit int
	Page  int

	Result []*OrgDTO
}

type OrgDTO struct {
	Id   int64  `+"`json:\"id\"`"+`
	Name string `+ "`json:\"name\"`"+`
}

type OrgDetailsDTO struct {
	Id   int64  ` +"`json:\"id\"`"+`
	Name string ` +"`json:\"name\"`"+`
	//Address Address ` +"`json:\"address\"`"+`
}

type UserOrgDTO struct {
	OrgId int64    `+"`json:\"orgId\"`"+`
	Name  string   ` +"`json:\"name\"`"+`
	Role  RoleType ` +"`json:\"role\"`"+`
}

`
	ModelOrgUser = `
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Typed errors
var (
	ErrInvalidRoleType     = errors.New("Invalid role type")
	ErrLastOrgAdmin        = errors.New("Cannot remove last organization admin")
	ErrOrgUserNotFound     = errors.New("Cannot find the organization user")
	ErrOrgUserAlreadyAdded = errors.New("User is already added to organization")
)

type RoleType string

const (
	ROLE_VIEWER RoleType = "Viewer"
	ROLE_EDITOR RoleType = "Editor"
	ROLE_ADMIN  RoleType = "Admin"
)

func (r RoleType) IsValid() bool {
	return r == ROLE_VIEWER || r == ROLE_ADMIN || r == ROLE_EDITOR
}

func (r RoleType) Includes(other RoleType) bool {
	if r == ROLE_ADMIN {
		return true
	}

	if r == ROLE_EDITOR {
		return other != ROLE_ADMIN
	}

	if r == ROLE_VIEWER {
		return other == ROLE_VIEWER
	}

	return false
}

func (r *RoleType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	*r = RoleType(str)

	if !(*r).IsValid() {
		if (*r) != "" {
			return fmt.Errorf("JSON validation error: invalid role value: %s", *r)
		}

		*r = ROLE_VIEWER
	}

	return nil
}

type OrgUser struct {
	Id      int64
	OrgId   int64
	UserId  int64
	Role    RoleType
	Created time.Time
	Updated time.Time
}

// ---------------------
// COMMANDS

type RemoveOrgUserCommand struct {
	UserId                   int64
	OrgId                    int64
	ShouldDeleteOrphanedUser bool
	UserWasDeleted           bool
}

type AddOrgUserCommand struct {
	LoginOrEmail string   `+"`json:\"loginOrEmail\" binding:\"Required\"`"+`
	Role         RoleType ` +"`json:\"role\" binding:\"Required\"`"+`

	OrgId  int64 `+"`json:\"-\"`"+`
	UserId int64 ` +"`json:\"-\"`"+`
}

type UpdateOrgUserCommand struct {
	Role RoleType ` +"`json:\"role\" binding:\"Required\"`"+`

	OrgId  int64 ` +"`json:\"-\"`"+`
	UserId int64 ` +"`json:\"-\"`"+`
}

// ----------------------
// QUERIES

type GetOrgUsersQuery struct {
	OrgId int64
	Query string
	Limit int

	Result []*OrgUserDTO
}

// ----------------------
// Projections and DTOs

type OrgUserDTO struct {
	OrgId         int64     `+"`json:\"orgId\"`"+`
	UserId        int64     ` +"`json:\"userId\"`"+`
	Email         string    ` +"`json:\"email\"`"+`
	AvatarUrl     string    ` +"`json:\"avatarUrl\"`"+`
	Login         string    ` +"`json:\"login\"`"+`
	Role          string    ` +"`json:\"role\"`"+`
	LastSeenAt    time.Time ` +"`json:\"lastSeenAt\"`"+`
	LastSeenAtAge string    ` +"`json:\"lastSeenAtAge\"`"+`
}

`
	ModelQuotas = `
package models

import (
	"errors"
	"{{.Dir}}/pkg/setting"
	"time"
)

var ErrInvalidQuotaTarget = errors.New("Invalid quota target")

type Quota struct {
	Id      int64
	OrgId   int64
	UserId  int64
	Target  string
	Limit   int64
	Created time.Time
	Updated time.Time
}

type QuotaScope struct {
	Name         string
	Target       string
	DefaultLimit int64
}

type OrgQuotaDTO struct {
	OrgId  int64  `+"`json:\"org_id\"`"+`
	Target string `+"`json:\"target\"`"+`
	Limit  int64  ` +"`json:\"limit\"`"+`
	Used   int64  ` +"`json:\"used\"`"+`
}

type UserQuotaDTO struct {
	UserId int64  `+"`json:\"user_id\"`"+`
	Target string ` +"`json:\"target\"`"+`
	Limit  int64  ` +"`json:\"limit\"`"+`
	Used   int64  ` +"`json:\"used\"`"+`
}

type GlobalQuotaDTO struct {
	Target string ` +"`json:\"target\"`"+`
	Limit  int64  ` +"`json:\"limit\"`"+`
	Used   int64  ` +"`json:\"used\"`"+`
}

type GetOrgQuotaByTargetQuery struct {
	Target  string
	OrgId   int64
	Default int64
	Result  *OrgQuotaDTO
}

type GetOrgQuotasQuery struct {
	OrgId  int64
	Result []*OrgQuotaDTO
}

type GetUserQuotaByTargetQuery struct {
	Target  string
	UserId  int64
	Default int64
	Result  *UserQuotaDTO
}

type GetUserQuotasQuery struct {
	UserId int64
	Result []*UserQuotaDTO
}

type GetGlobalQuotaByTargetQuery struct {
	Target  string
	Default int64
	Result  *GlobalQuotaDTO
}

type UpdateOrgQuotaCmd struct {
	Target string `+"`json:\"target\"`"+`
	Limit  int64  ` +"`json:\"limit\"`"+`
	OrgId  int64  ` +"`json:\"-\"`"+`
}

type UpdateUserQuotaCmd struct {
	Target string ` +"`json:\"target\"`"+`
	Limit  int64  ` +"`json:\"limit\"`"+`
	UserId int64  ` +"`json:\"-\"`"+`
}

func GetQuotaScopes(target string) ([]QuotaScope, error) {
	scopes := make([]QuotaScope, 0)
	switch target {
	case "user":
		scopes = append(scopes,
			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.User},
			QuotaScope{Name: "org", Target: "org_user", DefaultLimit: setting.Quota.Org.User},
		)
		return scopes, nil
	case "org":
		scopes = append(scopes,
			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.Org},
			QuotaScope{Name: "user", Target: "org_user", DefaultLimit: setting.Quota.User.Org},
		)
		return scopes, nil
	case "dashboard":
		scopes = append(scopes,
			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.Dashboard},
			QuotaScope{Name: "org", Target: target, DefaultLimit: setting.Quota.Org.Dashboard},
		)
		return scopes, nil
	case "data_source":
		scopes = append(scopes,
			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.DataSource},
			QuotaScope{Name: "org", Target: target, DefaultLimit: setting.Quota.Org.DataSource},
		)
		return scopes, nil
	case "api_key":
		scopes = append(scopes,
			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.ApiKey},
			QuotaScope{Name: "org", Target: target, DefaultLimit: setting.Quota.Org.ApiKey},
		)
		return scopes, nil
	case "session":
		scopes = append(scopes,
			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.Session},
		)
		return scopes, nil
	default:
		return scopes, ErrInvalidQuotaTarget
	}
}

`
	ModelTags = `
package models

import (
	"strings"
)

type Tag struct {
	Id    int64
	Key   string
	Value string
}

func ParseTagPairs(tagPairs []string) (tags []*Tag) {
	if tagPairs == nil {
		return []*Tag{}
	}

	for _, tagPair := range tagPairs {
		var tag Tag

		if strings.Contains(tagPair, ":") {
			keyValue := strings.Split(tagPair, ":")
			tag.Key = strings.Trim(keyValue[0], " ")
			tag.Value = strings.Trim(keyValue[1], " ")
		} else {
			tag.Key = strings.Trim(tagPair, " ")
		}

		if tag.Key == "" || ContainsTag(tags, &tag) {
			continue
		}

		tags = append(tags, &tag)
	}

	return tags
}

func ContainsTag(existingTags []*Tag, tag *Tag) bool {
	for _, t := range existingTags {
		if t.Key == tag.Key && t.Value == tag.Value {
			return true
		}
	}
	return false
}

func JoinTagPairs(tags []*Tag) []string {
	tagPairs := []string{}

	for _, tag := range tags {
		if tag.Value != "" {
			tagPairs = append(tagPairs, tag.Key+":"+tag.Value)
		} else {
			tagPairs = append(tagPairs, tag.Key)
		}
	}

	return tagPairs
}

`
	ModelTagsTest = `
package models

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsingTags(t *testing.T) {
	Convey("Testing parsing a tag pairs into tags", t, func() {
		Convey("Can parse one empty tag", func() {
			tags := ParseTagPairs([]string{""})
			So(len(tags), ShouldEqual, 0)
		})

		Convey("Can parse valid tags", func() {
			tags := ParseTagPairs([]string{"outage", "type:outage", "error"})
			So(len(tags), ShouldEqual, 3)
			So(tags[0].Key, ShouldEqual, "outage")
			So(tags[0].Value, ShouldEqual, "")
			So(tags[1].Key, ShouldEqual, "type")
			So(tags[1].Value, ShouldEqual, "outage")
			So(tags[2].Key, ShouldEqual, "error")
			So(tags[2].Value, ShouldEqual, "")
		})

		Convey("Can parse tags with spaces", func() {
			tags := ParseTagPairs([]string{" outage ", " type : outage ", "error "})
			So(len(tags), ShouldEqual, 3)
			So(tags[0].Key, ShouldEqual, "outage")
			So(tags[0].Value, ShouldEqual, "")
			So(tags[1].Key, ShouldEqual, "type")
			So(tags[1].Value, ShouldEqual, "outage")
			So(tags[2].Key, ShouldEqual, "error")
			So(tags[2].Value, ShouldEqual, "")
		})

		Convey("Can parse empty tags", func() {
			tags := ParseTagPairs([]string{" outage ", "", "", ":", "type : outage ", "error ", "", ""})
			So(len(tags), ShouldEqual, 3)
			So(tags[0].Key, ShouldEqual, "outage")
			So(tags[0].Value, ShouldEqual, "")
			So(tags[1].Key, ShouldEqual, "type")
			So(tags[1].Value, ShouldEqual, "outage")
			So(tags[2].Key, ShouldEqual, "error")
			So(tags[2].Value, ShouldEqual, "")
		})

		Convey("Can parse tags with extra colons", func() {
			tags := ParseTagPairs([]string{" outage", "type : outage:outage2 :outage3 ", "error :"})
			So(len(tags), ShouldEqual, 3)
			So(tags[0].Key, ShouldEqual, "outage")
			So(tags[0].Value, ShouldEqual, "")
			So(tags[1].Key, ShouldEqual, "type")
			So(tags[1].Value, ShouldEqual, "outage")
			So(tags[2].Key, ShouldEqual, "error")
			So(tags[2].Value, ShouldEqual, "")
		})

		Convey("Can parse tags that contains key and values with spaces", func() {
			tags := ParseTagPairs([]string{" outage 1", "type 1: outage 1 ", "has error "})
			So(len(tags), ShouldEqual, 3)
			So(tags[0].Key, ShouldEqual, "outage 1")
			So(tags[0].Value, ShouldEqual, "")
			So(tags[1].Key, ShouldEqual, "type 1")
			So(tags[1].Value, ShouldEqual, "outage 1")
			So(tags[2].Key, ShouldEqual, "has error")
			So(tags[2].Value, ShouldEqual, "")
		})

		Convey("Can filter out duplicate tags", func() {
			tags := ParseTagPairs([]string{"test", "test", "key:val1", "key:val2"})
			So(len(tags), ShouldEqual, 3)
			So(tags[0].Key, ShouldEqual, "test")
			So(tags[0].Value, ShouldEqual, "")
			So(tags[1].Key, ShouldEqual, "key")
			So(tags[1].Value, ShouldEqual, "val1")
			So(tags[2].Key, ShouldEqual, "key")
			So(tags[2].Value, ShouldEqual, "val2")
		})

		Convey("Can join tag pairs", func() {
			tagPairs := []*Tag{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: ""},
				{Key: "key3"},
			}
			tags := JoinTagPairs(tagPairs)
			So(len(tags), ShouldEqual, 3)
			So(tags[0], ShouldEqual, "key1:val1")
			So(tags[1], ShouldEqual, "key2")
			So(tags[2], ShouldEqual, "key3")
		})
	})
}

`
	ModelTeam = `
package models

import (
	"errors"
	"time"
)

// Typed errors
var (
	ErrTeamNotFound                         = errors.New("Team not found")
	ErrTeamNameTaken                        = errors.New("Team name is taken")
	ErrTeamMemberNotFound                   = errors.New("Team member not found")
	ErrLastTeamAdmin                        = errors.New("Not allowed to remove last admin")
	ErrNotAllowedToUpdateTeam               = errors.New("User not allowed to update team")
	ErrNotAllowedToUpdateTeamInDifferentOrg = errors.New("User not allowed to update team in another org")
)

// Team model
type Team struct {
	Id    int64  `+"`json:\"id\"`"+`
	OrgId int64  `+"`json:\"orgId\"`"+`
	Name  string ` +"`json:\"name\"`"+`
	Email string ` +"`json:\"email\"`"+`

	Created time.Time `+"`json:\"created\"`"+`
	Updated time.Time ` +"`json:\"updated\"`"+`
}

// ---------------------
// COMMANDS

type CreateTeamCommand struct {
	Name  string `+"`json:\"name\" binding:\"Required\"`"+`
	Email string ` +"`json:\"email\"`"+`
	OrgId int64  ` +"`json:\"-\"`"+`

	Result Team ` +"`json:\"-\"`"+`
}

type UpdateTeamCommand struct {
	Id    int64
	Name  string
	Email string
	OrgId int64 `+"`json:\"-\"`"+`
}

type DeleteTeamCommand struct {
	OrgId int64
	Id    int64
}

type GetTeamByIdQuery struct {
	OrgId  int64
	Id     int64
	Result *TeamDTO
}

type GetTeamsByUserQuery struct {
	OrgId  int64
	UserId int64      ` +"`json:\"userId\"`"+`
	Result []*TeamDTO ` +"`json:\"teams\"`"+`
}

type SearchTeamsQuery struct {
	Query        string
	Name         string
	Limit        int
	Page         int
	OrgId        int64
	UserIdFilter int64

	Result SearchTeamQueryResult
}

type TeamDTO struct {
	Id          int64  ` +"`json:\"id\"`"+`
	OrgId       int64  ` +"`json:\"orgId\"`"+`
	Name        string ` +"`json:\"name\"`"+`
	Email       string ` +"`json:\"email\"`"+`
	AvatarUrl   string ` +"`json:\"avatarUrl\"`"+`
	MemberCount int64  `+"`json:\"memberCount\"`"+`
	//Permission  PermissionType `+"`json:\"permission\"`"+`
}

type SearchTeamQueryResult struct {
	TotalCount int64      ` +"`json:\"totalCount\"`"+`
	Teams      []*TeamDTO ` +"`json:\"teams\"`"+`
	Page       int        ` +"`json:\"page\"`"+`
	PerPage    int        ` +"`json:\"perPage\"`"+`
}

`

	ModelTeamMember = `
package models

import (
	"errors"
	"time"
)

// Typed errors
var (
	ErrTeamMemberAlreadyAdded = errors.New("User is already added to this team")
)

// TeamMember model
type TeamMember struct {
	Id       int64
	OrgId    int64
	TeamId   int64
	UserId   int64
	External bool // Signals that the membership has been created by an external systems, such as LDAP
	//Permission PermissionType

	Created time.Time
	Updated time.Time
}

// ---------------------
// COMMANDS

type AddTeamMemberCommand struct {
	UserId   int64 `+"`json:\"userId\" binding:\"Required\"`"+`
	OrgId    int64 ` +"`json:\"-\"`"+`
	TeamId   int64 ` +"`json:\"-\"`"+`
	External bool  ` +"`json:\"-\"`"+`
	//Permission PermissionType `+"`json:\"-\"`"+`
}

type UpdateTeamMemberCommand struct {
	UserId int64 ` +"`json:\"-\"`"+`
	OrgId  int64 ` +"`json:\"-\"`"+`
	TeamId int64 ` +"`json:\"-\"`"+`
	//Permission       PermissionType ` +"`json:\"permission\"`"+`
	ProtectLastAdmin bool ` +"`json:\"-\"`"+`
}

type RemoveTeamMemberCommand struct {
	OrgId            int64 ` +"`json:\"-\"`"+`
	UserId           int64
	TeamId           int64
	ProtectLastAdmin bool ` +"`json:\"-\"`"+`
}

// ----------------------
// QUERIES

type GetTeamMembersQuery struct {
	OrgId    int64
	TeamId   int64
	UserId   int64
	External bool
	Result   []*TeamMemberDTO
}

// ----------------------
// Projections and DTOs

type TeamMemberDTO struct {
	OrgId     int64    ` +"`json:\"orgId\"`"+`
	TeamId    int64    ` +"`json:\"teamId\"`"+`
	UserId    int64    ` +"`json:\"userId\"`"+`
	External  bool     ` +"`json:\"-\"`"+`
	Email     string   ` +"`json:\"email\"`"+`
	Login     string   ` +"`json:\"login\"`"+`
	AvatarUrl string   ` +"`json:\"avatarUrl\"`"+`
	Labels    []string ` +"`json:\"labels\"`"+`
	//Permission PermissionType ` +"`json:\"permission\"`"+`
}

`
	ModelTempUser = `
package models

import (
	"errors"
	"time"
)

// Typed errors
var (
	ErrTempUserNotFound = errors.New("User not found")
)

type TempUserStatus string

const (
	TmpUserSignUpStarted TempUserStatus = "SignUpStarted"
	TmpUserInvitePending TempUserStatus = "InvitePending"
	TmpUserCompleted     TempUserStatus = "Completed"
	TmpUserRevoked       TempUserStatus = "Revoked"
)

// TempUser holds data for org invites and unconfirmed sign ups
type TempUser struct {
	Id              int64
	OrgId           int64
	Version         int
	Email           string
	Name            string
	InvitedByUserId int64
	Status          TempUserStatus

	EmailSent   bool
	EmailSentOn time.Time
	Code        string
	RemoteAddr  string

	Created time.Time
	Updated time.Time
}

// ---------------------
// COMMANDS

type CreateTempUserCommand struct {
	Email           string
	Name            string
	OrgId           int64
	InvitedByUserId int64
	Status          TempUserStatus
	Code            string
	RemoteAddr      string

	Result *TempUser
}

type UpdateTempUserStatusCommand struct {
	Code   string
	Status TempUserStatus
}

type UpdateTempUserWithEmailSentCommand struct {
	Code string
}

type GetTempUsersQuery struct {
	OrgId  int64
	Email  string
	Status TempUserStatus

	Result []*TempUserDTO
}

type GetTempUserByCodeQuery struct {
	Code string

	Result *TempUserDTO
}

type TempUserDTO struct {
	Id             int64          ` +"`json:\"id\"`"+`
	OrgId          int64          ` +"`json:\"orgId\"`"+`
	Name           string         ` +"`json:\"name\"`"+`
	Email          string         ` +"`json:\"email\"`"+`
	InvitedByLogin string         ` +"`json:\"invitedByLogin\"`"+`
	InvitedByEmail string         ` +"`json:\"invitedByEmail\"`"+`
	InvitedByName  string         ` +"`json:\"invitedByName\"`"+`
	Code           string         ` +"`json:\"code\"`"+`
	Status         TempUserStatus ` +"`json:\"status\"`"+`
	Url            string         ` +"`json:\"url\"`"+`
	EmailSent      bool           ` +"`json:\"emailSent\"`"+`
	EmailSentOn    time.Time      ` +"`json:\"emailSentOn\"`"+`
	Created        time.Time      ` +"`json:\"createdOn\"`"+`
}

`
	ModelUser = `
package models

import (
	"errors"
	"time"
)

// Typed errors
var (
	ErrUserNotFound     = errors.New("User not found")
	ErrLastGrafanaAdmin = errors.New("Cannot remove last grafana admin")
)

type Password string

func (p Password) IsWeak() bool {
	return len(p) <= 4
}

type User struct {
	Id            int64
	Version       int
	Email         string
	Name          string
	Login         string
	Password      string
	Salt          string
	Rands         string
	Company       string
	EmailVerified bool
	Theme         string
	HelpFlags1    HelpFlags1

	IsAdmin bool
	OrgId   int64

	Created    time.Time
	Updated    time.Time
	LastSeenAt time.Time
}

func (u *User) NameOrFallback() string {
	if u.Name != "" {
		return u.Name
	} else if u.Login != "" {
		return u.Login
	} else {
		return u.Email
	}
}

// ---------------------
// COMMANDS

type CreateUserCommand struct {
	Email          string
	Login          string
	Name           string
	Company        string
	OrgName        string
	Password       string
	EmailVerified  bool
	IsAdmin        bool
	SkipOrgSetup   bool
	DefaultOrgRole string

	Result User
}

type UpdateUserCommand struct {
	Name  string `+"`json:\"name\"`"+`
	Email string ` +"`json:\"email\"`"+`
	Login string ` +"`json:\"login\"`"+`
	Theme string ` +"`json:\"theme\"`"+`

	UserId int64 ` +"`json:\"-\"`"+ `
}

type ChangeUserPasswordCommand struct {
	OldPassword string `+"`json:\"oldPassword\"`"+`
	NewPassword string ` +"`json:\"newPassword\"`"+ `

	UserId int64 ` +"`json:\"-\"`"+`
}

type UpdateUserPermissionsCommand struct {
	IsGrafanaAdmin bool
	UserId         int64 ` +"`json:\"-\"`"+`
}

type DisableUserCommand struct {
	UserId     int64
	IsDisabled bool
}

type BatchDisableUsersCommand struct {
	UserIds    []int64
	IsDisabled bool
}

type DeleteUserCommand struct {
	UserId int64
}

type SetUsingOrgCommand struct {
	UserId int64
	OrgId  int64
}

// ----------------------
// QUERIES

type GetUserByLoginQuery struct {
	LoginOrEmail string
	Result       *User
}

type GetUserByEmailQuery struct {
	Email  string
	Result *User
}

type GetUserByIdQuery struct {
	Id     int64
	Result *User
}

type GetSignedInUserQuery struct {
	UserId int64
	Login  string
	Email  string
	OrgId  int64
	Result *SignedInUser
}

type GetUserProfileQuery struct {
	UserId int64
	Result UserProfileDTO
}

type SearchUsersQuery struct {
	OrgId      int64
	Query      string
	Page       int
	Limit      int
	AuthModule string

	Result SearchUserQueryResult
}

type SearchUserQueryResult struct {
	TotalCount int64               ` +"`json:\"totalCount\"`"+`
	Users      []*UserSearchHitDTO ` +"`json:\"users\"`"+`
	Page       int                 ` +"`json:\"page\"`"+`
	PerPage    int                 ` +"`json:\"perPage\"`"+`
}

type GetUserOrgListQuery struct {
	UserId int64
	Result []*UserOrgDTO
}

// ------------------------
// DTO & Projections

type SignedInUser struct {
	UserId         int64
	OrgId          int64
	OrgName        string
	OrgRole        RoleType
	Login          string
	Name           string
	Email          string
	ApiKeyId       int64
	OrgCount       int
	IsGrafanaAdmin bool
	IsAnonymous    bool
	HelpFlags1     HelpFlags1
	LastSeenAt     time.Time
	Teams          []int64
}

func (u *SignedInUser) ShouldUpdateLastSeenAt() bool {
	return u.UserId > 0 && time.Since(u.LastSeenAt) > time.Minute*5
}

func (u *SignedInUser) NameOrFallback() string {
	if u.Name != "" {
		return u.Name
	} else if u.Login != "" {
		return u.Login
	} else {
		return u.Email
	}
}

type UpdateUserLastSeenAtCommand struct {
	UserId int64
}

func (user *SignedInUser) HasRole(role RoleType) bool {
	if user.IsGrafanaAdmin {
		return true
	}

	return user.OrgRole.Includes(role)
}

type UserProfileDTO struct {
	Id             int64    `+"`json:\"id\"`"+`
	Email          string   `+"`json:\"email\"`"+`
	Name           string   `+"`json:\"name\"`"+`
	Login          string   ` +"`json:\"login\"`"+`
	Theme          string   ` +"`json:\"theme\"`"+`
	OrgId          int64    `+"`json:\"orgId\"`"+`
	IsGrafanaAdmin bool     ` +"`json:\"isGrafanaAdmin\"`"+`
	IsDisabled     bool     ` +"`json:\"isDisabled\"`"+`
	AuthModule     []string ` +"`json:\"authModule\"`"+`
}

type UserSearchHitDTO struct {
	Id            int64                `+"`json:\"id\"`"+`
	Name          string               `+"`json:\"name\"`"+`
	Login         string               `+"`json:\"login\"`"+`
	Email         string               ` +"`json:\"email\"`"+`
	AvatarUrl     string               ` +"`json:\"avatarUrl\"`"+`
	IsAdmin       bool                 ` +"`json:\"isAdmin\"`"+`
	IsDisabled    bool                 ` +"`json:\"isDisabled\"`"+`
	LastSeenAt    time.Time            `+"`json:\"lastSeenAt\"`"+`
	LastSeenAtAge string               ` +"`json:\"lastSeenAtAge\"`"+`
	AuthModule    AuthModuleConversion ` +"`json:\"authModule\"`"+`
}

type UserIdDTO struct {
	Id      int64  ` +"`json:\"id\"`"+`
	Message string ` +"`json:\"message\"`"+`
}

// implement Conversion interface to define custom field mapping (xorm feature)
type AuthModuleConversion []string

func (auth *AuthModuleConversion) FromDB(data []byte) error {
	auth_module := string(data)
	*auth = []string{auth_module}
	return nil
}

// Just a stub, we don't wanna write to database
func (auth *AuthModuleConversion) ToDB() ([]byte, error) {
	return []byte{}, nil
}

`
	ModelUserAuth = `
package models

import (
	"time"

	"golang.org/x/oauth2"
)

const (
	AuthModuleLDAP = "ldap"
)

type UserAuth struct {
	Id                int64
	UserId            int64
	AuthModule        string
	AuthId            string
	Created           time.Time
	OAuthAccessToken  string
	OAuthRefreshToken string
	OAuthTokenType    string
	OAuthExpiry       time.Time
}

type ExternalUserInfo struct {
	OAuthToken     *oauth2.Token
	AuthModule     string
	AuthId         string
	UserId         int64
	Email          string
	Login          string
	Name           string
	Groups         []string
	OrgRoles       map[int64]RoleType
	IsGrafanaAdmin *bool // This is a pointer to know if we should sync this or not (nil = ignore sync)
	IsDisabled     bool
}

// ---------------------
// COMMANDS

type UpsertUserCommand struct {
	ReqContext    *ReqContext
	ExternalUser  *ExternalUserInfo
	SignupAllowed bool

	Result *User
}

type SetAuthInfoCommand struct {
	AuthModule string
	AuthId     string
	UserId     int64
	OAuthToken *oauth2.Token
}

type UpdateAuthInfoCommand struct {
	AuthModule string
	AuthId     string
	UserId     int64
	OAuthToken *oauth2.Token
}

type DeleteAuthInfoCommand struct {
	UserAuth *UserAuth
}

// ----------------------
// QUERIES

type LoginUserQuery struct {
	ReqContext *ReqContext
	Username   string
	Password   string
	User       *User
	IpAddress  string
}

type GetUserByAuthInfoQuery struct {
	AuthModule string
	AuthId     string
	UserId     int64
	Email      string
	Login      string

	Result *User
}

type GetExternalUserInfoByLoginQuery struct {
	LoginOrEmail string

	Result *ExternalUserInfo
}

type GetAuthInfoQuery struct {
	UserId     int64
	AuthModule string
	AuthId     string

	Result *UserAuth
}

type SyncTeamsCommand struct {
	ExternalUser *ExternalUserInfo
	User         *User
}

`
	ModelUserToken = `
package models

import (
	"context"
	"errors"
)

// Typed errors
var (
	ErrUserTokenNotFound = errors.New("user token not found")
)

// UserToken represents a user token
type UserToken struct {
	Id            int64
	UserId        int64
	AuthToken     string
	PrevAuthToken string
	UserAgent     string
	ClientIp      string
	AuthTokenSeen bool
	SeenAt        int64
	RotatedAt     int64
	CreatedAt     int64
	UpdatedAt     int64
	UnhashedToken string
}

type RevokeAuthTokenCmd struct {
	AuthTokenId int64 `+"`json:\"authTokenId\"`"+`
}

// UserTokenService are used for generating and validating user tokens
type UserTokenService interface {
	CreateToken(ctx context.Context, userId int64, clientIP, userAgent string) (*UserToken, error)
	LookupToken(ctx context.Context, unhashedToken string) (*UserToken, error)
	TryRotateToken(ctx context.Context, token *UserToken, clientIP, userAgent string) (bool, error)
	RevokeToken(ctx context.Context, token *UserToken) error
	RevokeAllUserTokens(ctx context.Context, userId int64) error
	ActiveTokenCount(ctx context.Context) (int64, error)
	GetUserToken(ctx context.Context, userId, userTokenId int64) (*UserToken, error)
	GetUserTokens(ctx context.Context, userId int64) ([]*UserToken, error)
}

`

)
