package template4app

var (
	ServiceLoginAuth = `
package login

import (
	"errors"
	"{{.Dir}}/pkg/bus"
	"{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/registry"
)

var (
	ErrEmailNotAllowed       = errors.New("Required email domain not fulfilled")
	ErrInvalidCredentials    = errors.New("Invalid Username or Password")
	ErrNoEmail               = errors.New("Login provider didn't return an email address")
	ErrProviderDeniedRequest = errors.New("Login provider denied login request")
	ErrSignUpNotAllowed      = errors.New("Signup is not allowed for this adapter")
	ErrTooManyLoginAttempts  = errors.New("Too many consecutive incorrect login attempts for user. Login for user temporarily blocked")
	ErrPasswordEmpty         = errors.New("No password provided")
	ErrUserDisabled          = errors.New("User is disabled")
)

func init() {
	registry.RegisterService(&UserLoginService{})
}

type UserLoginService struct{}

func (l *UserLoginService) Init() error {
	bus.AddHandler("auth", AuthenticateUser)
	return nil
}

// AuthenticateUser authenticates the user via username & password
func AuthenticateUser(query *models.LoginUserQuery) error {
	if err := validateLoginAttempts(query.Username); err != nil {
		return err
	}

	if err := validatePasswordSet(query.Password); err != nil {
		return err
	}

	err := loginUsingGrafanaDB(query)
	if err == nil || (err != models.ErrUserNotFound && err != ErrInvalidCredentials && err != ErrUserDisabled) {
		return err
	}

	if err == models.ErrUserNotFound {
		return ErrInvalidCredentials
	}

	return err
}

func validatePasswordSet(password string) error {
	if len(password) == 0 || len(password) > 256 {
		return ErrPasswordEmpty
	}

	return nil
}

`
	ServiceLoginBruceForceLoginProtection = `
package login

import (
	"time"

	"{{.Dir}}/pkg/bus"
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
)

var (
	maxInvalidLoginAttempts int64 = 5
	loginAttemptsWindow           = time.Minute * 5
)

var validateLoginAttempts = func(username string) error {
	if setting.DisableBruteForceLoginProtection {
		return nil
	}

	loginAttemptCountQuery := m.GetUserLoginAttemptCountQuery{
		Username: username,
		Since:    time.Now().Add(-loginAttemptsWindow),
	}

	if err := bus.Dispatch(&loginAttemptCountQuery); err != nil {
		return err
	}

	if loginAttemptCountQuery.Result >= maxInvalidLoginAttempts {
		return ErrTooManyLoginAttempts
	}

	return nil
}

var saveInvalidLoginAttempt = func(query *m.LoginUserQuery) {
	if setting.DisableBruteForceLoginProtection {
		return
	}

	loginAttemptCommand := m.CreateLoginAttemptCommand{
		Username:  query.Username,
		IpAddress: query.IpAddress,
	}

	bus.Dispatch(&loginAttemptCommand)
}

`
	ServiceLoginBruceForceLoginProtectionTest = `
package login

import (
	"testing"

	"{{.Dir}}/pkg/bus"
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoginAttemptsValidation(t *testing.T) {
	Convey("Validate login attempts", t, func() {
		Convey("Given brute force login protection enabled", func() {
			setting.DisableBruteForceLoginProtection = false

			Convey("When user login attempt count equals max-1 ", func() {
				withLoginAttempts(maxInvalidLoginAttempts - 1)
				err := validateLoginAttempts("user")

				Convey("it should not result in error", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When user login attempt count equals max ", func() {
				withLoginAttempts(maxInvalidLoginAttempts)
				err := validateLoginAttempts("user")

				Convey("it should result in too many login attempts error", func() {
					So(err, ShouldEqual, ErrTooManyLoginAttempts)
				})
			})

			Convey("When user login attempt count is greater than max ", func() {
				withLoginAttempts(maxInvalidLoginAttempts + 5)
				err := validateLoginAttempts("user")

				Convey("it should result in too many login attempts error", func() {
					So(err, ShouldEqual, ErrTooManyLoginAttempts)
				})
			})

			Convey("When saving invalid login attempt", func() {
				defer bus.ClearBusHandlers()
				createLoginAttemptCmd := &m.CreateLoginAttemptCommand{}

				bus.AddHandler("test", func(cmd *m.CreateLoginAttemptCommand) error {
					createLoginAttemptCmd = cmd
					return nil
				})

				saveInvalidLoginAttempt(&m.LoginUserQuery{
					Username:  "user",
					Password:  "pwd",
					IpAddress: "192.168.1.1:56433",
				})

				Convey("it should dispatch command", func() {
					So(createLoginAttemptCmd, ShouldNotBeNil)
					So(createLoginAttemptCmd.Username, ShouldEqual, "user")
					So(createLoginAttemptCmd.IpAddress, ShouldEqual, "192.168.1.1:56433")
				})
			})
		})

		Convey("Given brute force login protection disabled", func() {
			setting.DisableBruteForceLoginProtection = true

			Convey("When user login attempt count equals max-1 ", func() {
				withLoginAttempts(maxInvalidLoginAttempts - 1)
				err := validateLoginAttempts("user")

				Convey("it should not result in error", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When user login attempt count equals max ", func() {
				withLoginAttempts(maxInvalidLoginAttempts)
				err := validateLoginAttempts("user")

				Convey("it should not result in error", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When user login attempt count is greater than max ", func() {
				withLoginAttempts(maxInvalidLoginAttempts + 5)
				err := validateLoginAttempts("user")

				Convey("it should not result in error", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When saving invalid login attempt", func() {
				defer bus.ClearBusHandlers()
				createLoginAttemptCmd := (*m.CreateLoginAttemptCommand)(nil)

				bus.AddHandler("test", func(cmd *m.CreateLoginAttemptCommand) error {
					createLoginAttemptCmd = cmd
					return nil
				})

				saveInvalidLoginAttempt(&m.LoginUserQuery{
					Username:  "user",
					Password:  "pwd",
					IpAddress: "192.168.1.1:56433",
				})

				Convey("it should not dispatch command", func() {
					So(createLoginAttemptCmd, ShouldBeNil)
				})
			})
		})
	})
}

func withLoginAttempts(loginAttempts int64) {
	bus.AddHandler("test", func(query *m.GetUserLoginAttemptCountQuery) error {
		query.Result = loginAttempts
		return nil
	})
}

`
	ServiceLoginGrafanaLogin = `
package login

import (
	"crypto/subtle"

	"{{.Dir}}/pkg/bus"
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/util"
)

var validatePassword = func(providedPassword string, userPassword string, userSalt string) error {
	passwordHashed := util.EncodePassword(providedPassword, userSalt)
	if subtle.ConstantTimeCompare([]byte(passwordHashed), []byte(userPassword)) != 1 {
		return ErrInvalidCredentials
	}

	return nil
}

var loginUsingGrafanaDB = func(query *m.LoginUserQuery) error {
	userQuery := m.GetUserByLoginQuery{LoginOrEmail: query.Username}

	if err := bus.Dispatch(&userQuery); err != nil {
		return err
	}

	user := userQuery.Result

	if err := validatePassword(query.Password, user.Password, user.Salt); err != nil {
		return err
	}

	query.User = user
	return nil
}

`
	ServiceLoginGrafanaLoginTest = `
package login

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"{{.Dir}}/pkg/bus"
	m "{{.Dir}}/pkg/models"
)

func TestGrafanaLogin(t *testing.T) {
	Convey("Login using Grafana DB", t, func() {
		grafanaLoginScenario("When login with non-existing user", func(sc *grafanaLoginScenarioContext) {
			sc.withNonExistingUser()
			err := loginUsingGrafanaDB(sc.loginUserQuery)

			Convey("it should result in user not found error", func() {
				So(err, ShouldEqual, m.ErrUserNotFound)
			})

			Convey("it should not call password validation", func() {
				So(sc.validatePasswordCalled, ShouldBeFalse)
			})

			Convey("it should not pupulate user object", func() {
				So(sc.loginUserQuery.User, ShouldBeNil)
			})
		})

		grafanaLoginScenario("When login with invalid credentials", func(sc *grafanaLoginScenarioContext) {
			sc.withInvalidPassword()
			err := loginUsingGrafanaDB(sc.loginUserQuery)

			Convey("it should result in invalid credentials error", func() {
				So(err, ShouldEqual, ErrInvalidCredentials)
			})

			Convey("it should call password validation", func() {
				So(sc.validatePasswordCalled, ShouldBeTrue)
			})

			Convey("it should not pupulate user object", func() {
				So(sc.loginUserQuery.User, ShouldBeNil)
			})
		})

		grafanaLoginScenario("When login with valid credentials", func(sc *grafanaLoginScenarioContext) {
			sc.withValidCredentials()
			err := loginUsingGrafanaDB(sc.loginUserQuery)

			Convey("it should not result in error", func() {
				So(err, ShouldBeNil)
			})

			Convey("it should call password validation", func() {
				So(sc.validatePasswordCalled, ShouldBeTrue)
			})

			Convey("it should pupulate user object", func() {
				So(sc.loginUserQuery.User, ShouldNotBeNil)
				So(sc.loginUserQuery.User.Login, ShouldEqual, sc.loginUserQuery.Username)
				So(sc.loginUserQuery.User.Password, ShouldEqual, sc.loginUserQuery.Password)
			})
		})

		grafanaLoginScenario("When login with disabled user", func(sc *grafanaLoginScenarioContext) {
			sc.withDisabledUser()
			err := loginUsingGrafanaDB(sc.loginUserQuery)

			Convey("it should return user is disabled error", func() {
				So(err, ShouldEqual, ErrUserDisabled)
			})

			Convey("it should not call password validation", func() {
				So(sc.validatePasswordCalled, ShouldBeFalse)
			})

			Convey("it should not pupulate user object", func() {
				So(sc.loginUserQuery.User, ShouldBeNil)
			})
		})
	})
}

type grafanaLoginScenarioContext struct {
	loginUserQuery         *m.LoginUserQuery
	validatePasswordCalled bool
}

type grafanaLoginScenarioFunc func(c *grafanaLoginScenarioContext)

func grafanaLoginScenario(desc string, fn grafanaLoginScenarioFunc) {
	Convey(desc, func() {
		origValidatePassword := validatePassword

		sc := &grafanaLoginScenarioContext{
			loginUserQuery: &m.LoginUserQuery{
				Username:  "user",
				Password:  "pwd",
				IpAddress: "192.168.1.1:56433",
			},
			validatePasswordCalled: false,
		}

		defer func() {
			validatePassword = origValidatePassword
		}()

		fn(sc)
	})
}

func mockPasswordValidation(valid bool, sc *grafanaLoginScenarioContext) {
	validatePassword = func(providedPassword string, userPassword string, userSalt string) error {
		sc.validatePasswordCalled = true

		if !valid {
			return ErrInvalidCredentials
		}

		return nil
	}
}

func (sc *grafanaLoginScenarioContext) getUserByLoginQueryReturns(user *m.User) {
	bus.AddHandler("test", func(query *m.GetUserByLoginQuery) error {
		if user == nil {
			return m.ErrUserNotFound
		}

		query.Result = user
		return nil
	})
}

func (sc *grafanaLoginScenarioContext) withValidCredentials() {
	sc.getUserByLoginQueryReturns(&m.User{
		Id:       1,
		Login:    sc.loginUserQuery.Username,
		Password: sc.loginUserQuery.Password,
		Salt:     "salt",
	})
	mockPasswordValidation(true, sc)
}

func (sc *grafanaLoginScenarioContext) withNonExistingUser() {
	sc.getUserByLoginQueryReturns(nil)
}

func (sc *grafanaLoginScenarioContext) withInvalidPassword() {
	sc.getUserByLoginQueryReturns(&m.User{
		Password: sc.loginUserQuery.Password,
		Salt:     "salt",
	})
	mockPasswordValidation(false, sc)
}

func (sc *grafanaLoginScenarioContext) withDisabledUser() {
	sc.getUserByLoginQueryReturns(&m.User{
		IsDisabled: true,
	})
}

`
)
