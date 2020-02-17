package template4app

var (
	ServiceAuthToken = `
package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"{{.Dir}}/pkg/services/serverlock"

	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/registry"
	"{{.Dir}}/pkg/services/sqlstore"
	"{{.Dir}}/pkg/setting"
	"{{.Dir}}/pkg/util"
)

func init() {
	registry.RegisterService(&UserAuthTokenService{})
}

var getTime = time.Now

const urgentRotateTime = 1 * time.Minute

type UserAuthTokenService struct {
	SQLStore          *sqlstore.SqlStore            ` + "`inject:\"\"`"+`
	ServerLockService *serverlock.ServerLockService ` +"`inject:\"\"`"+`
	Cfg               *setting.Cfg                  ` +"`inject:\"\"`"+`
	log               log.Logger
}

func (s *UserAuthTokenService) Init() error {
	s.log = log.New("auth")
	return nil
}

func (s *UserAuthTokenService) ActiveTokenCount(ctx context.Context) (int64, error) {

	var count int64
	var err error
	err = s.SQLStore.WithDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		var model userAuthToken
		count, err = dbSession.Where(` +"`created_at > ? AND rotated_at > ?`"+`,
			s.createdAfterParam(),
			s.rotatedAfterParam()).
			Count(&model)

		return err
	})

	return count, err
}

func (s *UserAuthTokenService) CreateToken(ctx context.Context, userId int64, clientIP, userAgent string) (*models.UserToken, error) {
	clientIP = util.ParseIPAddress(clientIP)
	token, err := util.RandomHex(16)
	if err != nil {
		return nil, err
	}

	hashedToken := hashToken(token)

	now := getTime().Unix()

	userAuthToken := userAuthToken{
		UserId:        userId,
		AuthToken:     hashedToken,
		PrevAuthToken: hashedToken,
		ClientIp:      clientIP,
		UserAgent:     userAgent,
		RotatedAt:     now,
		CreatedAt:     now,
		UpdatedAt:     now,
		SeenAt:        0,
		AuthTokenSeen: false,
	}

	err = s.SQLStore.WithDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		_, err = dbSession.Insert(&userAuthToken)
		return err
	})

	if err != nil {
		return nil, err
	}

	userAuthToken.UnhashedToken = token

	s.log.Debug("user auth token created", "tokenId", userAuthToken.Id, "userId", userAuthToken.UserId, "clientIP", userAuthToken.ClientIp, "userAgent", userAuthToken.UserAgent, "authToken", userAuthToken.AuthToken)

	var userToken models.UserToken
	err = userAuthToken.toUserToken(&userToken)

	return &userToken, err
}

func (s *UserAuthTokenService) LookupToken(ctx context.Context, unhashedToken string) (*models.UserToken, error) {
	hashedToken := hashToken(unhashedToken)
	if setting.Env == setting.DEV {
		s.log.Debug("looking up token", "unhashed", unhashedToken, "hashed", hashedToken)
	}

	var model userAuthToken
	var exists bool
	var err error
	err = s.SQLStore.WithDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		exists, err = dbSession.Where("(auth_token = ? OR prev_auth_token = ?) AND created_at > ? AND rotated_at > ?",
			hashedToken,
			hashedToken,
			s.createdAfterParam(),
			s.rotatedAfterParam()).
			Get(&model)

		return err

	})

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, models.ErrUserTokenNotFound
	}

	if model.AuthToken != hashedToken && model.PrevAuthToken == hashedToken && model.AuthTokenSeen {
		modelCopy := model
		modelCopy.AuthTokenSeen = false
		expireBefore := getTime().Add(-urgentRotateTime).Unix()

		var affectedRows int64
		err = s.SQLStore.WithTransactionalDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
			affectedRows, err = dbSession.Where("id = ? AND prev_auth_token = ? AND rotated_at < ?",
				modelCopy.Id,
				modelCopy.PrevAuthToken,
				expireBefore).
				AllCols().Update(&modelCopy)

			return err
		})

		if err != nil {
			return nil, err
		}

		if affectedRows == 0 {
			s.log.Debug("prev seen token unchanged", "tokenId", model.Id, "userId", model.UserId, "clientIP", model.ClientIp, "userAgent", model.UserAgent, "authToken", model.AuthToken)
		} else {
			s.log.Debug("prev seen token", "tokenId", model.Id, "userId", model.UserId, "clientIP", model.ClientIp, "userAgent", model.UserAgent, "authToken", model.AuthToken)
		}
	}

	if !model.AuthTokenSeen && model.AuthToken == hashedToken {
		modelCopy := model
		modelCopy.AuthTokenSeen = true
		modelCopy.SeenAt = getTime().Unix()

		var affectedRows int64
		err = s.SQLStore.WithTransactionalDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
			affectedRows, err = dbSession.Where("id = ? AND auth_token = ?",
				modelCopy.Id,
				modelCopy.AuthToken).
				AllCols().Update(&modelCopy)

			return err
		})

		if err != nil {
			return nil, err
		}

		if affectedRows == 1 {
			model = modelCopy
		}

		if affectedRows == 0 {
			s.log.Debug("seen wrong token", "tokenId", model.Id, "userId", model.UserId, "clientIP", model.ClientIp, "userAgent", model.UserAgent, "authToken", model.AuthToken)
		} else {
			s.log.Debug("seen token", "tokenId", model.Id, "userId", model.UserId, "clientIP", model.ClientIp, "userAgent", model.UserAgent, "authToken", model.AuthToken)
		}
	}

	model.UnhashedToken = unhashedToken

	var userToken models.UserToken
	err = model.toUserToken(&userToken)

	return &userToken, err
}

func (s *UserAuthTokenService) TryRotateToken(ctx context.Context, token *models.UserToken, clientIP, userAgent string) (bool, error) {
	if token == nil {
		return false, nil
	}

	model := userAuthTokenFromUserToken(token)

	now := getTime()

	var needsRotation bool
	rotatedAt := time.Unix(model.RotatedAt, 0)
	if model.AuthTokenSeen {
		needsRotation = rotatedAt.Before(now.Add(-time.Duration(s.Cfg.TokenRotationIntervalMinutes) * time.Minute))
	} else {
		needsRotation = rotatedAt.Before(now.Add(-urgentRotateTime))
	}

	if !needsRotation {
		return false, nil
	}

	s.log.Debug("token needs rotation", "tokenId", model.Id, "authTokenSeen", model.AuthTokenSeen, "rotatedAt", rotatedAt)

	clientIP = util.ParseIPAddress(clientIP)
	newToken, err := util.RandomHex(16)
	if err != nil {
		return false, err
	}
	hashedToken := hashToken(newToken)

	// very important that auth_token_seen is set after the prev_auth_token = case when ... for mysql to function correctly
	sql := `+ "`UPDATE user_auth_token " + `
SET
seen_at = 0,
user_agent = ?,
client_ip = ?,
prev_auth_token = case when auth_token_seen = ? then auth_token else prev_auth_token end,
auth_token = ?,
auth_token_seen = ?,
rotated_at = ?
WHERE id = ? AND (auth_token_seen = ? OR rotated_at < ?)` + "` " + `

	var affected int64
	err = s.SQLStore.WithTransactionalDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		res, err := dbSession.Exec(sql, userAgent, clientIP, s.SQLStore.Dialect.BooleanStr(true), hashedToken, s.SQLStore.Dialect.BooleanStr(false), now.Unix(), model.Id, s.SQLStore.Dialect.BooleanStr(true), now.Add(-30*time.Second).Unix())
		if err != nil {
			return err
		}

		affected, err = res.RowsAffected()
		return err
	})

	if err != nil {
		return false, err
	}

	s.log.Debug("auth token rotated", "affected", affected, "auth_token_id", model.Id, "userId", model.UserId)
	if affected > 0 {
		model.UnhashedToken = newToken
		model.toUserToken(token)
		return true, nil
	}

	return false, nil
}

func (s *UserAuthTokenService) RevokeToken(ctx context.Context, token *models.UserToken) error {
	if token == nil {
		return models.ErrUserTokenNotFound
	}

	model := userAuthTokenFromUserToken(token)

	var rowsAffected int64
	var err error
	err = s.SQLStore.WithDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		rowsAffected, err = dbSession.Delete(model)
		return err
	})

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		s.log.Debug("user auth token not found/revoked", "tokenId", model.Id, "userId", model.UserId, "clientIP", model.ClientIp, "userAgent", model.UserAgent)
		return models.ErrUserTokenNotFound
	}

	s.log.Debug("user auth token revoked", "tokenId", model.Id, "userId", model.UserId, "clientIP", model.ClientIp, "userAgent", model.UserAgent)

	return nil
}

func (s *UserAuthTokenService) RevokeAllUserTokens(ctx context.Context, userId int64) error {
	return s.SQLStore.WithDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		sql := ` +"`DELETE from user_auth_token WHERE user_id = ?`"+`
		res, err := dbSession.Exec(sql, userId)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		s.log.Debug("all user tokens for user revoked", "userId", userId, "count", affected)

		return err
	})
}

func (s *UserAuthTokenService) BatchRevokeAllUserTokens(ctx context.Context, userIds []int64) error {
	return s.SQLStore.WithTransactionalDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		if len(userIds) == 0 {
			return nil
		}

		user_id_params := strings.Repeat(",?", len(userIds)-1)
		sql := "DELETE from user_auth_token WHERE user_id IN (?" + user_id_params + ")"

		params := []interface{}{sql}
		for _, v := range userIds {
			params = append(params, v)
		}

		res, err := dbSession.Exec(params...)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		s.log.Debug("all user tokens for given users revoked", "usersCount", len(userIds), "count", affected)

		return err
	})
}

func (s *UserAuthTokenService) GetUserToken(ctx context.Context, userId, userTokenId int64) (*models.UserToken, error) {

	var result models.UserToken
	err := s.SQLStore.WithDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		var token userAuthToken
		exists, err := dbSession.Where("id = ? AND user_id = ?", userTokenId, userId).Get(&token)
		if err != nil {
			return err
		}

		if !exists {
			return models.ErrUserTokenNotFound
		}

		token.toUserToken(&result)
		return nil
	})

	return &result, err
}

func (s *UserAuthTokenService) GetUserTokens(ctx context.Context, userId int64) ([]*models.UserToken, error) {

	result := []*models.UserToken{}
	err := s.SQLStore.WithDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		var tokens []*userAuthToken
		err := dbSession.Where("user_id = ? AND created_at > ? AND rotated_at > ?",
			userId,
			s.createdAfterParam(),
			s.rotatedAfterParam()).
			Find(&tokens)

		if err != nil {
			return err
		}

		for _, token := range tokens {
			var userToken models.UserToken
			token.toUserToken(&userToken)
			result = append(result, &userToken)
		}

		return nil
	})

	return result, err
}

func (s *UserAuthTokenService) createdAfterParam() int64 {
	tokenMaxLifetime := time.Duration(s.Cfg.LoginMaxLifetimeDays) * 24 * time.Hour
	return getTime().Add(-tokenMaxLifetime).Unix()
}

func (s *UserAuthTokenService) rotatedAfterParam() int64 {
	tokenMaxInactiveLifetime := time.Duration(s.Cfg.LoginMaxInactiveLifetimeDays) * 24 * time.Hour
	return getTime().Add(-tokenMaxInactiveLifetime).Unix()
}

func hashToken(token string) string {
	hashBytes := sha256.Sum256([]byte(token + setting.SecretKey))
	return hex.EncodeToString(hashBytes[:])
}

`
	ServiceAuthTokenTest = `
package auth

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"{{.Dir}}/pkg/components/simplejson"
	"{{.Dir}}/pkg/setting"

	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/services/sqlstore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserAuthToken(t *testing.T) {
	Convey("Test user auth token", t, func() {
		ctx := createTestContext(t)
		userAuthTokenService := ctx.tokenService
		userID := int64(10)

		t := time.Date(2018, 12, 13, 13, 45, 0, 0, time.UTC)
		getTime = func() time.Time {
			return t
		}

		Convey("When creating token", func() {
			userToken, err := userAuthTokenService.CreateToken(context.Background(), userID, "192.168.10.11:1234", "some user agent")
			So(err, ShouldBeNil)
			So(userToken, ShouldNotBeNil)
			So(userToken.AuthTokenSeen, ShouldBeFalse)

			Convey("Can count active tokens", func() {
				count, err := userAuthTokenService.ActiveTokenCount(context.Background())
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 1)
			})

			Convey("When lookup unhashed token should return user auth token", func() {
				userToken, err := userAuthTokenService.LookupToken(context.Background(), userToken.UnhashedToken)
				So(err, ShouldBeNil)
				So(userToken, ShouldNotBeNil)
				So(userToken.UserId, ShouldEqual, userID)
				So(userToken.AuthTokenSeen, ShouldBeTrue)

				storedAuthToken, err := ctx.getAuthTokenByID(userToken.Id)
				So(err, ShouldBeNil)
				So(storedAuthToken, ShouldNotBeNil)
				So(storedAuthToken.AuthTokenSeen, ShouldBeTrue)
			})

			Convey("When lookup hashed token should return user auth token not found error", func() {
				userToken, err := userAuthTokenService.LookupToken(context.Background(), userToken.AuthToken)
				So(err, ShouldEqual, models.ErrUserTokenNotFound)
				So(userToken, ShouldBeNil)
			})

			Convey("revoking existing token should delete token", func() {
				err = userAuthTokenService.RevokeToken(context.Background(), userToken)
				So(err, ShouldBeNil)

				model, err := ctx.getAuthTokenByID(userToken.Id)
				So(err, ShouldBeNil)
				So(model, ShouldBeNil)
			})

			Convey("revoking nil token should return error", func() {
				err = userAuthTokenService.RevokeToken(context.Background(), nil)
				So(err, ShouldEqual, models.ErrUserTokenNotFound)
			})

			Convey("revoking non-existing token should return error", func() {
				userToken.Id = 1000
				err = userAuthTokenService.RevokeToken(context.Background(), userToken)
				So(err, ShouldEqual, models.ErrUserTokenNotFound)
			})

			Convey("When creating an additional token", func() {
				userToken2, err := userAuthTokenService.CreateToken(context.Background(), userID, "192.168.10.11:1234", "some user agent")
				So(err, ShouldBeNil)
				So(userToken2, ShouldNotBeNil)

				Convey("Can get first user token", func() {
					token, err := userAuthTokenService.GetUserToken(context.Background(), userID, userToken.Id)
					So(err, ShouldBeNil)
					So(token, ShouldNotBeNil)
					So(token.Id, ShouldEqual, userToken.Id)
				})

				Convey("Can get second user token", func() {
					token, err := userAuthTokenService.GetUserToken(context.Background(), userID, userToken2.Id)
					So(err, ShouldBeNil)
					So(token, ShouldNotBeNil)
					So(token.Id, ShouldEqual, userToken2.Id)
				})

				Convey("Can get user tokens", func() {
					tokens, err := userAuthTokenService.GetUserTokens(context.Background(), userID)
					So(err, ShouldBeNil)
					So(tokens, ShouldHaveLength, 2)
					So(tokens[0].Id, ShouldEqual, userToken.Id)
					So(tokens[1].Id, ShouldEqual, userToken2.Id)
				})

				Convey("Can revoke all user tokens", func() {
					err := userAuthTokenService.RevokeAllUserTokens(context.Background(), userID)
					So(err, ShouldBeNil)

					model, err := ctx.getAuthTokenByID(userToken.Id)
					So(err, ShouldBeNil)
					So(model, ShouldBeNil)

					model2, err := ctx.getAuthTokenByID(userToken2.Id)
					So(err, ShouldBeNil)
					So(model2, ShouldBeNil)
				})
			})

			Convey("When revoking users tokens in a batch", func() {
				Convey("Can revoke all users tokens", func() {
					userIds := []int64{}
					for i := 0; i < 3; i++ {
						userId := userID + int64(i+1)
						userIds = append(userIds, userId)
						userAuthTokenService.CreateToken(context.Background(), userId, "192.168.10.11:1234", "some user agent")
					}

					err := userAuthTokenService.BatchRevokeAllUserTokens(context.Background(), userIds)
					So(err, ShouldBeNil)

					for _, v := range userIds {
						tokens, err := userAuthTokenService.GetUserTokens(context.Background(), v)
						So(err, ShouldBeNil)
						So(len(tokens), ShouldEqual, 0)
					}
				})
			})
		})

		Convey("expires correctly", func() {
			userToken, err := userAuthTokenService.CreateToken(context.Background(), userID, "192.168.10.11:1234", "some user agent")
			So(err, ShouldBeNil)

			userToken, err = userAuthTokenService.LookupToken(context.Background(), userToken.UnhashedToken)
			So(err, ShouldBeNil)

			getTime = func() time.Time {
				return t.Add(time.Hour)
			}

			rotated, err := userAuthTokenService.TryRotateToken(context.Background(), userToken, "192.168.10.11:1234", "some user agent")
			So(err, ShouldBeNil)
			So(rotated, ShouldBeTrue)

			userToken, err = userAuthTokenService.LookupToken(context.Background(), userToken.UnhashedToken)
			So(err, ShouldBeNil)

			stillGood, err := userAuthTokenService.LookupToken(context.Background(), userToken.UnhashedToken)
			So(err, ShouldBeNil)
			So(stillGood, ShouldNotBeNil)

			model, err := ctx.getAuthTokenByID(userToken.Id)
			So(err, ShouldBeNil)

			Convey("when rotated_at is 6:59:59 ago should find token", func() {
				getTime = func() time.Time {
					return time.Unix(model.RotatedAt, 0).Add(24 * 7 * time.Hour).Add(-time.Second)
				}

				stillGood, err = userAuthTokenService.LookupToken(context.Background(), stillGood.UnhashedToken)
				So(err, ShouldBeNil)
				So(stillGood, ShouldNotBeNil)
			})

			Convey("when rotated_at is 7:00:00 ago should not find token", func() {
				getTime = func() time.Time {
					return time.Unix(model.RotatedAt, 0).Add(24 * 7 * time.Hour)
				}

				notGood, err := userAuthTokenService.LookupToken(context.Background(), userToken.UnhashedToken)
				So(err, ShouldEqual, models.ErrUserTokenNotFound)
				So(notGood, ShouldBeNil)

				Convey("should not find active token when expired", func() {
					count, err := userAuthTokenService.ActiveTokenCount(context.Background())
					So(err, ShouldBeNil)
					So(count, ShouldEqual, 0)
				})
			})

			Convey("when rotated_at is 5 days ago and created_at is 29 days and 23:59:59 ago should not find token", func() {
				updated, err := ctx.updateRotatedAt(model.Id, time.Unix(model.CreatedAt, 0).Add(24*25*time.Hour).Unix())
				So(err, ShouldBeNil)
				So(updated, ShouldBeTrue)

				getTime = func() time.Time {
					return time.Unix(model.CreatedAt, 0).Add(24 * 30 * time.Hour).Add(-time.Second)
				}

				stillGood, err = userAuthTokenService.LookupToken(context.Background(), stillGood.UnhashedToken)
				So(err, ShouldBeNil)
				So(stillGood, ShouldNotBeNil)
			})

			Convey("when rotated_at is 5 days ago and created_at is 30 days ago should not find token", func() {
				updated, err := ctx.updateRotatedAt(model.Id, time.Unix(model.CreatedAt, 0).Add(24*25*time.Hour).Unix())
				So(err, ShouldBeNil)
				So(updated, ShouldBeTrue)

				getTime = func() time.Time {
					return time.Unix(model.CreatedAt, 0).Add(24 * 30 * time.Hour)
				}

				notGood, err := userAuthTokenService.LookupToken(context.Background(), userToken.UnhashedToken)
				So(err, ShouldEqual, models.ErrUserTokenNotFound)
				So(notGood, ShouldBeNil)
			})
		})

		Convey("can properly rotate tokens", func() {
			userToken, err := userAuthTokenService.CreateToken(context.Background(), userID, "192.168.10.11:1234", "some user agent")
			So(err, ShouldBeNil)

			prevToken := userToken.AuthToken
			unhashedPrev := userToken.UnhashedToken

			rotated, err := userAuthTokenService.TryRotateToken(context.Background(), userToken, "192.168.10.12:1234", "a new user agent")
			So(err, ShouldBeNil)
			So(rotated, ShouldBeFalse)

			updated, err := ctx.markAuthTokenAsSeen(userToken.Id)
			So(err, ShouldBeNil)
			So(updated, ShouldBeTrue)

			model, err := ctx.getAuthTokenByID(userToken.Id)
			So(err, ShouldBeNil)

			var tok models.UserToken
			err = model.toUserToken(&tok)
			So(err, ShouldBeNil)

			getTime = func() time.Time {
				return t.Add(time.Hour)
			}

			rotated, err = userAuthTokenService.TryRotateToken(context.Background(), &tok, "192.168.10.12:1234", "a new user agent")
			So(err, ShouldBeNil)
			So(rotated, ShouldBeTrue)

			unhashedToken := tok.UnhashedToken

			model, err = ctx.getAuthTokenByID(tok.Id)
			So(err, ShouldBeNil)
			model.UnhashedToken = unhashedToken

			So(model.RotatedAt, ShouldEqual, getTime().Unix())
			So(model.ClientIp, ShouldEqual, "192.168.10.12")
			So(model.UserAgent, ShouldEqual, "a new user agent")
			So(model.AuthTokenSeen, ShouldBeFalse)
			So(model.SeenAt, ShouldEqual, 0)
			So(model.PrevAuthToken, ShouldEqual, prevToken)

			// ability to auth using an old token

			lookedUpUserToken, err := userAuthTokenService.LookupToken(context.Background(), model.UnhashedToken)
			So(err, ShouldBeNil)
			So(lookedUpUserToken, ShouldNotBeNil)
			So(lookedUpUserToken.AuthTokenSeen, ShouldBeTrue)
			So(lookedUpUserToken.SeenAt, ShouldEqual, getTime().Unix())

			lookedUpUserToken, err = userAuthTokenService.LookupToken(context.Background(), unhashedPrev)
			So(err, ShouldBeNil)
			So(lookedUpUserToken, ShouldNotBeNil)
			So(lookedUpUserToken.Id, ShouldEqual, model.Id)
			So(lookedUpUserToken.AuthTokenSeen, ShouldBeTrue)

			getTime = func() time.Time {
				return t.Add(time.Hour + (2 * time.Minute))
			}

			lookedUpUserToken, err = userAuthTokenService.LookupToken(context.Background(), unhashedPrev)
			So(err, ShouldBeNil)
			So(lookedUpUserToken, ShouldNotBeNil)
			So(lookedUpUserToken.AuthTokenSeen, ShouldBeTrue)

			lookedUpModel, err := ctx.getAuthTokenByID(lookedUpUserToken.Id)
			So(err, ShouldBeNil)
			So(lookedUpModel, ShouldNotBeNil)
			So(lookedUpModel.AuthTokenSeen, ShouldBeFalse)

			rotated, err = userAuthTokenService.TryRotateToken(context.Background(), userToken, "192.168.10.12:1234", "a new user agent")
			So(err, ShouldBeNil)
			So(rotated, ShouldBeTrue)

			model, err = ctx.getAuthTokenByID(userToken.Id)
			So(err, ShouldBeNil)
			So(model, ShouldNotBeNil)
			So(model.SeenAt, ShouldEqual, 0)
		})

		Convey("keeps prev token valid for 1 minute after it is confirmed", func() {
			userToken, err := userAuthTokenService.CreateToken(context.Background(), userID, "192.168.10.11:1234", "some user agent")
			So(err, ShouldBeNil)
			So(userToken, ShouldNotBeNil)

			lookedUpUserToken, err := userAuthTokenService.LookupToken(context.Background(), userToken.UnhashedToken)
			So(err, ShouldBeNil)
			So(lookedUpUserToken, ShouldNotBeNil)

			getTime = func() time.Time {
				return t.Add(10 * time.Minute)
			}

			prevToken := userToken.UnhashedToken
			rotated, err := userAuthTokenService.TryRotateToken(context.Background(), userToken, "1.1.1.1", "firefox")
			So(err, ShouldBeNil)
			So(rotated, ShouldBeTrue)

			getTime = func() time.Time {
				return t.Add(20 * time.Minute)
			}

			currentUserToken, err := userAuthTokenService.LookupToken(context.Background(), userToken.UnhashedToken)
			So(err, ShouldBeNil)
			So(currentUserToken, ShouldNotBeNil)

			prevUserToken, err := userAuthTokenService.LookupToken(context.Background(), prevToken)
			So(err, ShouldBeNil)
			So(prevUserToken, ShouldNotBeNil)
		})

		Convey("will not mark token unseen when prev and current are the same", func() {
			userToken, err := userAuthTokenService.CreateToken(context.Background(), userID, "192.168.10.11:1234", "some user agent")
			So(err, ShouldBeNil)
			So(userToken, ShouldNotBeNil)

			lookedUpUserToken, err := userAuthTokenService.LookupToken(context.Background(), userToken.UnhashedToken)
			So(err, ShouldBeNil)
			So(lookedUpUserToken, ShouldNotBeNil)

			lookedUpUserToken, err = userAuthTokenService.LookupToken(context.Background(), userToken.UnhashedToken)
			So(err, ShouldBeNil)
			So(lookedUpUserToken, ShouldNotBeNil)

			lookedUpModel, err := ctx.getAuthTokenByID(lookedUpUserToken.Id)
			So(err, ShouldBeNil)
			So(lookedUpModel, ShouldNotBeNil)
			So(lookedUpModel.AuthTokenSeen, ShouldBeTrue)
		})

		Convey("Rotate token", func() {
			userToken, err := userAuthTokenService.CreateToken(context.Background(), userID, "192.168.10.11:1234", "some user agent")
			So(err, ShouldBeNil)
			So(userToken, ShouldNotBeNil)

			prevToken := userToken.AuthToken

			Convey("Should rotate current token and previous token when auth token seen", func() {
				updated, err := ctx.markAuthTokenAsSeen(userToken.Id)
				So(err, ShouldBeNil)
				So(updated, ShouldBeTrue)

				getTime = func() time.Time {
					return t.Add(10 * time.Minute)
				}

				rotated, err := userAuthTokenService.TryRotateToken(context.Background(), userToken, "1.1.1.1", "firefox")
				So(err, ShouldBeNil)
				So(rotated, ShouldBeTrue)

				storedToken, err := ctx.getAuthTokenByID(userToken.Id)
				So(err, ShouldBeNil)
				So(storedToken, ShouldNotBeNil)
				So(storedToken.AuthTokenSeen, ShouldBeFalse)
				So(storedToken.PrevAuthToken, ShouldEqual, prevToken)
				So(storedToken.AuthToken, ShouldNotEqual, prevToken)

				prevToken = storedToken.AuthToken

				updated, err = ctx.markAuthTokenAsSeen(userToken.Id)
				So(err, ShouldBeNil)
				So(updated, ShouldBeTrue)

				getTime = func() time.Time {
					return t.Add(20 * time.Minute)
				}

				rotated, err = userAuthTokenService.TryRotateToken(context.Background(), userToken, "1.1.1.1", "firefox")
				So(err, ShouldBeNil)
				So(rotated, ShouldBeTrue)

				storedToken, err = ctx.getAuthTokenByID(userToken.Id)
				So(err, ShouldBeNil)
				So(storedToken, ShouldNotBeNil)
				So(storedToken.AuthTokenSeen, ShouldBeFalse)
				So(storedToken.PrevAuthToken, ShouldEqual, prevToken)
				So(storedToken.AuthToken, ShouldNotEqual, prevToken)
			})

			Convey("Should rotate current token, but keep previous token when auth token not seen", func() {
				userToken.RotatedAt = getTime().Add(-2 * time.Minute).Unix()

				getTime = func() time.Time {
					return t.Add(2 * time.Minute)
				}

				rotated, err := userAuthTokenService.TryRotateToken(context.Background(), userToken, "1.1.1.1", "firefox")
				So(err, ShouldBeNil)
				So(rotated, ShouldBeTrue)

				storedToken, err := ctx.getAuthTokenByID(userToken.Id)
				So(err, ShouldBeNil)
				So(storedToken, ShouldNotBeNil)
				So(storedToken.AuthTokenSeen, ShouldBeFalse)
				So(storedToken.PrevAuthToken, ShouldEqual, prevToken)
				So(storedToken.AuthToken, ShouldNotEqual, prevToken)
			})
		})

		Convey("When populating userAuthToken from UserToken should copy all properties", func() {
			ut := models.UserToken{
				Id:            1,
				UserId:        2,
				AuthToken:     "a",
				PrevAuthToken: "b",
				UserAgent:     "c",
				ClientIp:      "d",
				AuthTokenSeen: true,
				SeenAt:        3,
				RotatedAt:     4,
				CreatedAt:     5,
				UpdatedAt:     6,
				UnhashedToken: "e",
			}
			utBytes, err := json.Marshal(ut)
			So(err, ShouldBeNil)
			utJSON, err := simplejson.NewJson(utBytes)
			So(err, ShouldBeNil)
			utMap := utJSON.MustMap()

			var uat userAuthToken
			uat.fromUserToken(&ut)
			uatBytes, err := json.Marshal(uat)
			So(err, ShouldBeNil)
			uatJSON, err := simplejson.NewJson(uatBytes)
			So(err, ShouldBeNil)
			uatMap := uatJSON.MustMap()

			So(uatMap, ShouldResemble, utMap)
		})

		Convey("When populating userToken from userAuthToken should copy all properties", func() {
			uat := userAuthToken{
				Id:            1,
				UserId:        2,
				AuthToken:     "a",
				PrevAuthToken: "b",
				UserAgent:     "c",
				ClientIp:      "d",
				AuthTokenSeen: true,
				SeenAt:        3,
				RotatedAt:     4,
				CreatedAt:     5,
				UpdatedAt:     6,
				UnhashedToken: "e",
			}
			uatBytes, err := json.Marshal(uat)
			So(err, ShouldBeNil)
			uatJSON, err := simplejson.NewJson(uatBytes)
			So(err, ShouldBeNil)
			uatMap := uatJSON.MustMap()

			var ut models.UserToken
			err = uat.toUserToken(&ut)
			So(err, ShouldBeNil)
			utBytes, err := json.Marshal(ut)
			So(err, ShouldBeNil)
			utJSON, err := simplejson.NewJson(utBytes)
			So(err, ShouldBeNil)
			utMap := utJSON.MustMap()

			So(utMap, ShouldResemble, uatMap)
		})

		Reset(func() {
			getTime = time.Now
		})
	})
}

func createTestContext(t *testing.T) *testContext {
	t.Helper()

	sqlstore := sqlstore.InitTestDB(t)
	tokenService := &UserAuthTokenService{
		SQLStore: sqlstore,
		Cfg: &setting.Cfg{
			LoginMaxInactiveLifetimeDays: 7,
			LoginMaxLifetimeDays:         30,
			TokenRotationIntervalMinutes: 10,
		},
		log: log.New("test-logger"),
	}

	return &testContext{
		sqlstore:     sqlstore,
		tokenService: tokenService,
	}
}

type testContext struct {
	sqlstore     *sqlstore.SqlStore
	tokenService *UserAuthTokenService
}

func (c *testContext) getAuthTokenByID(id int64) (*userAuthToken, error) {
	sess := c.sqlstore.NewSession()
	var t userAuthToken
	found, err := sess.ID(id).Get(&t)
	if err != nil || !found {
		return nil, err
	}

	return &t, nil
}

func (c *testContext) markAuthTokenAsSeen(id int64) (bool, error) {
	sess := c.sqlstore.NewSession()
	res, err := sess.Exec("UPDATE user_auth_token SET auth_token_seen = ? WHERE id = ?", c.sqlstore.Dialect.BooleanStr(true), id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}

func (c *testContext) updateRotatedAt(id, rotatedAt int64) (bool, error) {
	sess := c.sqlstore.NewSession()
	res, err := sess.Exec("UPDATE user_auth_token SET rotated_at = ? WHERE id = ?", rotatedAt, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}
`
	ServiceModel = `
package auth

import (
	"fmt"

	"{{.Dir}}/pkg/models"
)

type userAuthToken struct {
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
	UnhashedToken string ` +"`xorm:\"-\"`"+`
}

func userAuthTokenFromUserToken(ut *models.UserToken) *userAuthToken {
	var uat userAuthToken
	uat.fromUserToken(ut)
	return &uat
}

func (uat *userAuthToken) fromUserToken(ut *models.UserToken) error {
	if uat == nil {
		return fmt.Errorf("needs pointer to userAuthToken struct")
	}

	uat.Id = ut.Id
	uat.UserId = ut.UserId
	uat.AuthToken = ut.AuthToken
	uat.PrevAuthToken = ut.PrevAuthToken
	uat.UserAgent = ut.UserAgent
	uat.ClientIp = ut.ClientIp
	uat.AuthTokenSeen = ut.AuthTokenSeen
	uat.SeenAt = ut.SeenAt
	uat.RotatedAt = ut.RotatedAt
	uat.CreatedAt = ut.CreatedAt
	uat.UpdatedAt = ut.UpdatedAt
	uat.UnhashedToken = ut.UnhashedToken

	return nil
}

func (uat *userAuthToken) toUserToken(ut *models.UserToken) error {
	if uat == nil {
		return fmt.Errorf("needs pointer to userAuthToken struct")
	}

	ut.Id = uat.Id
	ut.UserId = uat.UserId
	ut.AuthToken = uat.AuthToken
	ut.PrevAuthToken = uat.PrevAuthToken
	ut.UserAgent = uat.UserAgent
	ut.ClientIp = uat.ClientIp
	ut.AuthTokenSeen = uat.AuthTokenSeen
	ut.SeenAt = uat.SeenAt
	ut.RotatedAt = uat.RotatedAt
	ut.CreatedAt = uat.CreatedAt
	ut.UpdatedAt = uat.UpdatedAt
	ut.UnhashedToken = uat.UnhashedToken

	return nil
}

`
	ServiceTesting = `
package auth

import (
	"context"

	"{{.Dir}}/pkg/models"
)

type FakeUserAuthTokenService struct {
	CreateTokenProvider         func(ctx context.Context, userId int64, clientIP, userAgent string) (*models.UserToken, error)
	TryRotateTokenProvider      func(ctx context.Context, token *models.UserToken, clientIP, userAgent string) (bool, error)
	LookupTokenProvider         func(ctx context.Context, unhashedToken string) (*models.UserToken, error)
	RevokeTokenProvider         func(ctx context.Context, token *models.UserToken) error
	RevokeAllUserTokensProvider func(ctx context.Context, userId int64) error
	ActiveAuthTokenCount        func(ctx context.Context) (int64, error)
	GetUserTokenProvider        func(ctx context.Context, userId, userTokenId int64) (*models.UserToken, error)
	GetUserTokensProvider       func(ctx context.Context, userId int64) ([]*models.UserToken, error)
	BatchRevokedTokenProvider   func(ctx context.Context, userIds []int64) error
}

func NewFakeUserAuthTokenService() *FakeUserAuthTokenService {
	return &FakeUserAuthTokenService{
		CreateTokenProvider: func(ctx context.Context, userId int64, clientIP, userAgent string) (*models.UserToken, error) {
			return &models.UserToken{
				UserId:        0,
				UnhashedToken: "",
			}, nil
		},
		TryRotateTokenProvider: func(ctx context.Context, token *models.UserToken, clientIP, userAgent string) (bool, error) {
			return false, nil
		},
		LookupTokenProvider: func(ctx context.Context, unhashedToken string) (*models.UserToken, error) {
			return &models.UserToken{
				UserId:        0,
				UnhashedToken: "",
			}, nil
		},
		RevokeTokenProvider: func(ctx context.Context, token *models.UserToken) error {
			return nil
		},
		RevokeAllUserTokensProvider: func(ctx context.Context, userId int64) error {
			return nil
		},
		BatchRevokedTokenProvider: func(ctx context.Context, userIds []int64) error {
			return nil
		},
		ActiveAuthTokenCount: func(ctx context.Context) (int64, error) {
			return 10, nil
		},
		GetUserTokenProvider: func(ctx context.Context, userId, userTokenId int64) (*models.UserToken, error) {
			return nil, nil
		},
		GetUserTokensProvider: func(ctx context.Context, userId int64) ([]*models.UserToken, error) {
			return nil, nil
		},
	}
}

func (s *FakeUserAuthTokenService) CreateToken(ctx context.Context, userId int64, clientIP, userAgent string) (*models.UserToken, error) {
	return s.CreateTokenProvider(context.Background(), userId, clientIP, userAgent)
}

func (s *FakeUserAuthTokenService) LookupToken(ctx context.Context, unhashedToken string) (*models.UserToken, error) {
	return s.LookupTokenProvider(context.Background(), unhashedToken)
}

func (s *FakeUserAuthTokenService) TryRotateToken(ctx context.Context, token *models.UserToken, clientIP, userAgent string) (bool, error) {
	return s.TryRotateTokenProvider(context.Background(), token, clientIP, userAgent)
}

func (s *FakeUserAuthTokenService) RevokeToken(ctx context.Context, token *models.UserToken) error {
	return s.RevokeTokenProvider(context.Background(), token)
}

func (s *FakeUserAuthTokenService) RevokeAllUserTokens(ctx context.Context, userId int64) error {
	return s.RevokeAllUserTokensProvider(context.Background(), userId)
}

func (s *FakeUserAuthTokenService) ActiveTokenCount(ctx context.Context) (int64, error) {
	return s.ActiveAuthTokenCount(context.Background())
}

func (s *FakeUserAuthTokenService) GetUserToken(ctx context.Context, userId, userTokenId int64) (*models.UserToken, error) {
	return s.GetUserTokenProvider(context.Background(), userId, userTokenId)
}

func (s *FakeUserAuthTokenService) GetUserTokens(ctx context.Context, userId int64) ([]*models.UserToken, error) {
	return s.GetUserTokensProvider(context.Background(), userId)
}

func (s *FakeUserAuthTokenService) BatchRevokeAllUserTokens(ctx context.Context, userIds []int64) error {
	return s.BatchRevokedTokenProvider(ctx, userIds)
}

`
	ServiceTokenCleanup = `
package auth

import (
	"context"
	"time"

	"{{.Dir}}/pkg/services/sqlstore"
)

func (srv *UserAuthTokenService) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Hour)
	maxInactiveLifetime := time.Duration(srv.Cfg.LoginMaxInactiveLifetimeDays) * 24 * time.Hour
	maxLifetime := time.Duration(srv.Cfg.LoginMaxLifetimeDays) * 24 * time.Hour

	err := srv.ServerLockService.LockAndExecute(ctx, "cleanup expired auth tokens", time.Hour*12, func() {
		srv.deleteExpiredTokens(ctx, maxInactiveLifetime, maxLifetime)
	})

	if err != nil {
		srv.log.Error("failed to lock and execute cleanup of expired auth token", "error", err)
	}

	for {
		select {
		case <-ticker.C:
			err := srv.ServerLockService.LockAndExecute(ctx, "cleanup expired auth tokens", time.Hour*12, func() {
				srv.deleteExpiredTokens(ctx, maxInactiveLifetime, maxLifetime)
			})

			if err != nil {
				srv.log.Error("failed to lock and execute cleanup of expired auth token", "error", err)
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (srv *UserAuthTokenService) deleteExpiredTokens(ctx context.Context, maxInactiveLifetime, maxLifetime time.Duration) (int64, error) {
	createdBefore := getTime().Add(-maxLifetime)
	rotatedBefore := getTime().Add(-maxInactiveLifetime)

	srv.log.Debug("starting cleanup of expired auth tokens", "createdBefore", createdBefore, "rotatedBefore", rotatedBefore)

	var affected int64
	err := srv.SQLStore.WithDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		sql := ` + "`DELETE from user_auth_token WHERE created_at <= ? OR rotated_at <= ?`"+`
		res, err := dbSession.Exec(sql, createdBefore.Unix(), rotatedBefore.Unix())
		if err != nil {
			return err
		}

		affected, err = res.RowsAffected()
		if err != nil {
			srv.log.Error("failed to cleanup expired auth tokens", "error", err)
			return nil
		}

		srv.log.Debug("cleanup of expired auth tokens done", "count", affected)

		return nil
	})

	return affected, err
}

`
	ServiceTokenCleanupTest = `
package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUserAuthTokenCleanup(t *testing.T) {

	Convey("Test user auth token cleanup", t, func() {
		ctx := createTestContext(t)
		ctx.tokenService.Cfg.LoginMaxInactiveLifetimeDays = 7
		ctx.tokenService.Cfg.LoginMaxLifetimeDays = 30

		insertToken := func(token string, prev string, createdAt, rotatedAt int64) {
			ut := userAuthToken{AuthToken: token, PrevAuthToken: prev, CreatedAt: createdAt, RotatedAt: rotatedAt, UserAgent: "", ClientIp: ""}
			_, err := ctx.sqlstore.NewSession().Insert(&ut)
			So(err, ShouldBeNil)
		}

		t := time.Date(2018, 12, 13, 13, 45, 0, 0, time.UTC)
		getTime = func() time.Time {
			return t
		}

		Convey("should delete tokens where token rotation age is older than or equal 7 days", func() {
			from := t.Add(-7 * 24 * time.Hour)

			// insert three old tokens that should be deleted
			for i := 0; i < 3; i++ {
				insertToken(fmt.Sprintf("oldA%d", i), fmt.Sprintf("oldB%d", i), from.Unix(), from.Unix())
			}

			// insert three active tokens that should not be deleted
			for i := 0; i < 3; i++ {
				from = from.Add(time.Second)
				insertToken(fmt.Sprintf("newA%d", i), fmt.Sprintf("newB%d", i), from.Unix(), from.Unix())
			}

			affected, err := ctx.tokenService.deleteExpiredTokens(context.Background(), 7*24*time.Hour, 30*24*time.Hour)
			So(err, ShouldBeNil)
			So(affected, ShouldEqual, 3)
		})

		Convey("should delete tokens where token age is older than or equal 30 days", func() {
			from := t.Add(-30 * 24 * time.Hour)
			fromRotate := t.Add(-time.Second)

			// insert three old tokens that should be deleted
			for i := 0; i < 3; i++ {
				insertToken(fmt.Sprintf("oldA%d", i), fmt.Sprintf("oldB%d", i), from.Unix(), fromRotate.Unix())
			}

			// insert three active tokens that should not be deleted
			for i := 0; i < 3; i++ {
				from = from.Add(time.Second)
				insertToken(fmt.Sprintf("newA%d", i), fmt.Sprintf("newB%d", i), from.Unix(), fromRotate.Unix())
			}

			affected, err := ctx.tokenService.deleteExpiredTokens(context.Background(), 7*24*time.Hour, 30*24*time.Hour)
			So(err, ShouldBeNil)
			So(affected, ShouldEqual, 3)
		})
	})
}

`

)
