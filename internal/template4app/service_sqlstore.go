package template4app

var (
	ServiceSqlstoreAnnotations = `
package sqlstore

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/services/annotations"
)

type SqlAnnotationRepo struct {
}

func (r *SqlAnnotationRepo) Save(item *annotations.Item) error {
	return inTransaction(func(sess *DBSession) error {
		tags := models.ParseTagPairs(item.Tags)
		item.Tags = models.JoinTagPairs(tags)
		item.Created = time.Now().UnixNano() / int64(time.Millisecond)
		item.Updated = item.Created
		if item.Epoch == 0 {
			item.Epoch = item.Created
		}

		if _, err := sess.Table("annotation").Insert(item); err != nil {
			return err
		}

		if item.Tags != nil {
			tags, err := EnsureTagsExist(sess, tags)
			if err != nil {
				return err
			}
			for _, tag := range tags {
				if _, err := sess.Exec("INSERT INTO annotation_tag (annotation_id, tag_id) VALUES(?,?)", item.Id, tag.Id); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *SqlAnnotationRepo) Update(item *annotations.Item) error {
	return inTransaction(func(sess *DBSession) error {
		var (
			isExist bool
			err     error
		)
		existing := new(annotations.Item)

		if item.Id == 0 && item.RegionId != 0 {
			// Update region end time
			isExist, err = sess.Table("annotation").Where("region_id=? AND id!=? AND org_id=?", item.RegionId, item.RegionId, item.OrgId).Get(existing)
		} else {
			isExist, err = sess.Table("annotation").Where("id=? AND org_id=?", item.Id, item.OrgId).Get(existing)
		}

		if err != nil {
			return err
		}
		if !isExist {
			return errors.New("Annotation not found")
		}

		existing.Updated = time.Now().UnixNano() / int64(time.Millisecond)
		existing.Epoch = item.Epoch
		existing.Text = item.Text
		if item.RegionId != 0 {
			existing.RegionId = item.RegionId
		}

		if item.Tags != nil {
			tags, err := EnsureTagsExist(sess, models.ParseTagPairs(item.Tags))
			if err != nil {
				return err
			}
			if _, err := sess.Exec("DELETE FROM annotation_tag WHERE annotation_id = ?", existing.Id); err != nil {
				return err
			}
			for _, tag := range tags {
				if _, err := sess.Exec("INSERT INTO annotation_tag (annotation_id, tag_id) VALUES(?,?)", existing.Id, tag.Id); err != nil {
					return err
				}
			}
		}

		existing.Tags = item.Tags

		_, err = sess.Table("annotation").ID(existing.Id).Cols("epoch", "text", "region_id", "updated", "tags").Update(existing)
		return err
	})
}

func (r *SqlAnnotationRepo) Find(query *annotations.ItemQuery) ([]*annotations.ItemDTO, error) {
	var sql bytes.Buffer
	params := make([]interface{}, 0)

	sql.WriteString(` + "` " +`
SELECT
	annotation.id,
annotation.epoch as time,
annotation.dashboard_id,
annotation.panel_id,
annotation.new_state,
annotation.prev_state,
annotation.alert_id,
annotation.region_id,
annotation.text,
annotation.tags,
annotation.data,
annotation.created,
annotation.updated,
usr.email,
usr.login,
alert.name as alert_name
FROM annotation
LEFT OUTER JOIN ` + "`+"+ "dialect.Quote(\"user\")+`" + ` as usr on usr.id = annotation.user_id
LEFT OUTER JOIN alert on alert.id = annotation.alert_id
` +"`"+ `)

	sql.WriteString(` +"`WHERE annotation.org_id = ?`"+`)
	params = append(params, query.OrgId)

	if query.AnnotationId != 0 {
		// fmt.Print("annotation query")
		sql.WriteString(`  + "`AND annotation.id = ?`"+`)
		params = append(params, query.AnnotationId)
	}

	if query.RegionId != 0 {
		sql.WriteString(`  +"`AND annotation.region_id = ?`"+`)
		params = append(params, query.RegionId)
	}

	if query.AlertId != 0 {
		sql.WriteString(`  +"`AND annotation.alert_id = ?`"+`)
		params = append(params, query.AlertId)
	}

	if query.DashboardId != 0 {
		sql.WriteString(`  +"`AND annotation.dashboard_id = ?`"+`)
		params = append(params, query.DashboardId)
	}

	if query.PanelId != 0 {
		sql.WriteString(`  +"`AND annotation.panel_id = ?`"+`)
		params = append(params, query.PanelId)
	}

	if query.UserId != 0 {
		sql.WriteString(` +"`AND annotation.user_id = ?`"+`)
		params = append(params, query.UserId)
	}

	if query.From > 0 && query.To > 0 {
		sql.WriteString(` +"`AND annotation.epoch BETWEEN ? AND ?`"+`)
		params = append(params, query.From, query.To)
	}

	if query.Type == "alert" {
		sql.WriteString(` +"`AND annotation.alert_id > 0`"+`)
	} else if query.Type == "annotation" {
		sql.WriteString(` +"`AND annotation.alert_id = 0`"+`)
	}

	if len(query.Tags) > 0 {
		keyValueFilters := []string{}

		tags := models.ParseTagPairs(query.Tags)
		for _, tag := range tags {
			if tag.Value == "" {
				keyValueFilters = append(keyValueFilters, "(tag."+dialect.Quote("key")+" = ?)")
				params = append(params, tag.Key)
			} else {
				keyValueFilters = append(keyValueFilters, "(tag."+dialect.Quote("key")+" = ? AND tag."+dialect.Quote("value")+" = ?)")
				params = append(params, tag.Key, tag.Value)
			}
		}

		if len(tags) > 0 {
			tagsSubQuery := fmt.Sprintf(`+"`SELECT SUM(1) FROM annotation_tag at INNER JOIN tag on tag.id = at.tag_id WHERE at.annotation_id = annotation.id AND (%s)`"+`, strings.Join(keyValueFilters, " OR "))

			if query.MatchAny {
				sql.WriteString(fmt.Sprintf(" AND (%s) > 0 ", tagsSubQuery))
			} else {
				sql.WriteString(fmt.Sprintf(" AND (%s) = %d ", tagsSubQuery, len(tags)))
			}

		}
	}

	if query.Limit == 0 {
		query.Limit = 100
	}

	sql.WriteString(" ORDER BY epoch DESC" + dialect.Limit(query.Limit))

	items := make([]*annotations.ItemDTO, 0)

	if err := x.SQL(sql.String(), params...).Find(&items); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *SqlAnnotationRepo) Delete(params *annotations.DeleteParams) error {
	return inTransaction(func(sess *DBSession) error {
		var (
			sql         string
			annoTagSql  string
			queryParams []interface{}
		)

		sqlog.Info("delete", "orgId", params.OrgId)
		if params.RegionId != 0 {
			annoTagSql = "DELETE FROM annotation_tag WHERE annotation_id IN (SELECT id FROM annotation WHERE region_id = ? AND org_id = ?)"
			sql = "DELETE FROM annotation WHERE region_id = ? AND org_id = ?"
			queryParams = []interface{}{params.RegionId, params.OrgId}
		} else if params.Id != 0 {
			annoTagSql = "DELETE FROM annotation_tag WHERE annotation_id IN (SELECT id FROM annotation WHERE id = ? AND org_id = ?)"
			sql = "DELETE FROM annotation WHERE id = ? AND org_id = ?"
			queryParams = []interface{}{params.Id, params.OrgId}
		} else {
			annoTagSql = "DELETE FROM annotation_tag WHERE annotation_id IN (SELECT id FROM annotation WHERE dashboard_id = ? AND panel_id = ? AND org_id = ?)"
			sql = "DELETE FROM annotation WHERE dashboard_id = ? AND panel_id = ? AND org_id = ?"
			queryParams = []interface{}{params.DashboardId, params.PanelId, params.OrgId}
		}

		sqlOrArgs := append([]interface{}{annoTagSql}, queryParams...)

		if _, err := sess.Exec(sqlOrArgs...); err != nil {
			return err
		}

		sqlOrArgs = append([]interface{}{sql}, queryParams...)

		if _, err := sess.Exec(sqlOrArgs...); err != nil {
			return err
		}

		return nil
	})
}

`
	ServiceSqlstoreAnnotationsTest = `
package sqlstore

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"{{.Dir}}/pkg/services/annotations"
)

func TestAnnotations(t *testing.T) {
	InitTestDB(t)

	Convey("Testing annotation saving/loading", t, func() {
		repo := SqlAnnotationRepo{}

		Convey("Can save annotation", func() {
			Reset(func() {
				_, err := x.Exec("DELETE FROM annotation WHERE 1=1")
				So(err, ShouldBeNil)
				_, err = x.Exec("DELETE FROM annotation_tag WHERE 1=1")
				So(err, ShouldBeNil)
			})

			annotation := &annotations.Item{
				OrgId:       1,
				UserId:      1,
				DashboardId: 1,
				Text:        "hello",
				Type:        "alert",
				Epoch:       10,
				Tags:        []string{"outage", "error", "type:outage", "server:server-1"},
			}
			err := repo.Save(annotation)

			So(err, ShouldBeNil)
			So(annotation.Id, ShouldBeGreaterThan, 0)

			annotation2 := &annotations.Item{
				OrgId:       1,
				UserId:      1,
				DashboardId: 2,
				Text:        "hello",
				Type:        "alert",
				Epoch:       20,
				Tags:        []string{"outage", "error", "type:outage", "server:server-1"},
				RegionId:    1,
			}
			err = repo.Save(annotation2)
			So(err, ShouldBeNil)
			So(annotation2.Id, ShouldBeGreaterThan, 0)

			globalAnnotation1 := &annotations.Item{
				OrgId:  1,
				UserId: 1,
				Text:   "deploy",
				Type:   "",
				Epoch:  15,
				Tags:   []string{"deploy"},
			}
			err = repo.Save(globalAnnotation1)
			So(err, ShouldBeNil)
			So(globalAnnotation1.Id, ShouldBeGreaterThan, 0)

			globalAnnotation2 := &annotations.Item{
				OrgId:  1,
				UserId: 1,
				Text:   "rollback",
				Type:   "",
				Epoch:  17,
				Tags:   []string{"rollback"},
			}
			err = repo.Save(globalAnnotation2)
			So(err, ShouldBeNil)
			So(globalAnnotation2.Id, ShouldBeGreaterThan, 0)

			Convey("Can query for annotation by dashboard id", func() {
				items, err := repo.Find(&annotations.ItemQuery{
					OrgId:       1,
					DashboardId: 1,
					From:        0,
					To:          15,
				})

				So(err, ShouldBeNil)
				So(items, ShouldHaveLength, 1)

				Convey("Can read tags", func() {
					So(items[0].Tags, ShouldResemble, []string{"outage", "error", "type:outage", "server:server-1"})
				})

				Convey("Has created and updated values", func() {
					So(items[0].Created, ShouldBeGreaterThan, 0)
					So(items[0].Updated, ShouldBeGreaterThan, 0)
					So(items[0].Updated, ShouldEqual, items[0].Created)
				})
			})

			Convey("Can query for annotation by id", func() {
				items, err := repo.Find(&annotations.ItemQuery{
					OrgId:        1,
					AnnotationId: annotation2.Id,
				})

				So(err, ShouldBeNil)
				So(items, ShouldHaveLength, 1)
				So(items[0].Id, ShouldEqual, annotation2.Id)
			})

			Convey("Can query for annotation by region id", func() {
				items, err := repo.Find(&annotations.ItemQuery{
					OrgId:    1,
					RegionId: annotation2.RegionId,
				})

				So(err, ShouldBeNil)
				So(items, ShouldHaveLength, 1)
				So(items[0].Id, ShouldEqual, annotation2.Id)
			})

			Convey("Should not find any when item is outside time range", func() {
				items, err := repo.Find(&annotations.ItemQuery{
					OrgId:       1,
					DashboardId: 1,
					From:        12,
					To:          15,
				})

				So(err, ShouldBeNil)
				So(items, ShouldHaveLength, 0)
			})

			Convey("Should not find one when tag filter does not match", func() {
				items, err := repo.Find(&annotations.ItemQuery{
					OrgId:       1,
					DashboardId: 1,
					From:        1,
					To:          15,
					Tags:        []string{"asd"},
				})

				So(err, ShouldBeNil)
				So(items, ShouldHaveLength, 0)
			})

			Convey("Should not find one when type filter does not match", func() {
				items, err := repo.Find(&annotations.ItemQuery{
					OrgId:       1,
					DashboardId: 1,
					From:        1,
					To:          15,
					Type:        "alert",
				})

				So(err, ShouldBeNil)
				So(items, ShouldHaveLength, 0)
			})

			Convey("Should find one when all tag filters does match", func() {
				items, err := repo.Find(&annotations.ItemQuery{
					OrgId:       1,
					DashboardId: 1,
					From:        1,
					To:          15, //this will exclude the second test annotation
					Tags:        []string{"outage", "error"},
				})

				So(err, ShouldBeNil)
				So(items, ShouldHaveLength, 1)
			})

			Convey("Should find two annotations using partial match", func() {
				items, err := repo.Find(&annotations.ItemQuery{
					OrgId:    1,
					From:     1,
					To:       25,
					MatchAny: true,
					Tags:     []string{"rollback", "deploy"},
				})

				So(err, ShouldBeNil)
				So(items, ShouldHaveLength, 2)
			})

			Convey("Should find one when all key value tag filters does match", func() {
				items, err := repo.Find(&annotations.ItemQuery{
					OrgId:       1,
					DashboardId: 1,
					From:        1,
					To:          15,
					Tags:        []string{"type:outage", "server:server-1"},
				})

				So(err, ShouldBeNil)
				So(items, ShouldHaveLength, 1)
			})

			Convey("Can update annotation and remove all tags", func() {
				query := &annotations.ItemQuery{
					OrgId:       1,
					DashboardId: 1,
					From:        0,
					To:          15,
				}
				items, err := repo.Find(query)

				So(err, ShouldBeNil)

				annotationId := items[0].Id

				err = repo.Update(&annotations.Item{
					Id:    annotationId,
					OrgId: 1,
					Text:  "something new",
					Tags:  []string{},
				})

				So(err, ShouldBeNil)

				items, err = repo.Find(query)

				So(err, ShouldBeNil)

				Convey("Can read tags", func() {
					So(items[0].Id, ShouldEqual, annotationId)
					So(len(items[0].Tags), ShouldEqual, 0)
					So(items[0].Text, ShouldEqual, "something new")
				})
			})

			Convey("Can update annotation with new tags", func() {
				query := &annotations.ItemQuery{
					OrgId:       1,
					DashboardId: 1,
					From:        0,
					To:          15,
				}
				items, err := repo.Find(query)

				So(err, ShouldBeNil)

				annotationId := items[0].Id

				err = repo.Update(&annotations.Item{
					Id:    annotationId,
					OrgId: 1,
					Text:  "something new",
					Tags:  []string{"newtag1", "newtag2"},
				})

				So(err, ShouldBeNil)

				items, err = repo.Find(query)

				So(err, ShouldBeNil)

				Convey("Can read tags", func() {
					So(items[0].Id, ShouldEqual, annotationId)
					So(items[0].Tags, ShouldResemble, []string{"newtag1", "newtag2"})
					So(items[0].Text, ShouldEqual, "something new")
				})

				Convey("Updated time has increased", func() {
					So(items[0].Updated, ShouldBeGreaterThan, items[0].Created)
				})
			})

			Convey("Can delete annotation", func() {
				query := &annotations.ItemQuery{
					OrgId:       1,
					DashboardId: 1,
					From:        0,
					To:          15,
				}
				items, err := repo.Find(query)
				So(err, ShouldBeNil)

				annotationId := items[0].Id

				err = repo.Delete(&annotations.DeleteParams{Id: annotationId, OrgId: 1})
				So(err, ShouldBeNil)

				items, err = repo.Find(query)
				So(err, ShouldBeNil)

				Convey("Should be deleted", func() {
					So(len(items), ShouldEqual, 0)
				})
			})

		})
	})
}

`
	ServiceSqlstoreHealth = `
package sqlstore

import (
	"fmt"
	"{{.Dir}}/pkg/bus"
	m "{{.Dir}}/pkg/models"
)

func init() {
	bus.AddHandler("sql", GetDBHealthQuery)
	fmt.Println("Initialized sqlstore DB health....")
}

func GetDBHealthQuery(query *m.GetDBHealthQuery) error {
	return x.Ping()
}

`
	ServiceSqlstoreLogger = `
package sqlstore

import (
	"fmt"
	glog "{{.Dir}}/pkg/infra/log"

	"github.com/go-xorm/core"
)

type XormLogger struct {
	grafanaLog glog.Logger
	level      glog.Lvl
	showSQL    bool
}

func NewXormLogger(level glog.Lvl, grafanaLog glog.Logger) *XormLogger {
	return &XormLogger{
		grafanaLog: grafanaLog,
		level:      level,
		showSQL:    true,
	}
}

// Error implement core.ILogger
func (s *XormLogger) Error(v ...interface{}) {
	if s.level <= glog.LvlError {
		s.grafanaLog.Error(fmt.Sprint(v...))
	}
}

// Errorf implement core.ILogger
func (s *XormLogger) Errorf(format string, v ...interface{}) {
	if s.level <= glog.LvlError {
		s.grafanaLog.Error(fmt.Sprintf(format, v...))
	}
}

// Debug implement core.ILogger
func (s *XormLogger) Debug(v ...interface{}) {
	if s.level <= glog.LvlDebug {
		s.grafanaLog.Debug(fmt.Sprint(v...))
	}
}

// Debugf implement core.ILogger
func (s *XormLogger) Debugf(format string, v ...interface{}) {
	if s.level <= glog.LvlDebug {
		s.grafanaLog.Debug(fmt.Sprintf(format, v...))
	}
}

// Info implement core.ILogger
func (s *XormLogger) Info(v ...interface{}) {
	if s.level <= glog.LvlInfo {
		s.grafanaLog.Info(fmt.Sprint(v...))
	}
}

// Infof implement core.ILogger
func (s *XormLogger) Infof(format string, v ...interface{}) {
	if s.level <= glog.LvlInfo {
		s.grafanaLog.Info(fmt.Sprintf(format, v...))
	}
}

// Warn implement core.ILogger
func (s *XormLogger) Warn(v ...interface{}) {
	if s.level <= glog.LvlWarn {
		s.grafanaLog.Warn(fmt.Sprint(v...))
	}
}

// Warnf implement core.ILogger
func (s *XormLogger) Warnf(format string, v ...interface{}) {
	if s.level <= glog.LvlWarn {
		s.grafanaLog.Warn(fmt.Sprintf(format, v...))
	}
}

// Level implement core.ILogger
func (s *XormLogger) Level() core.LogLevel {
	switch s.level {
	case glog.LvlError:
		return core.LOG_ERR
	case glog.LvlWarn:
		return core.LOG_WARNING
	case glog.LvlInfo:
		return core.LOG_INFO
	case glog.LvlDebug:
		return core.LOG_DEBUG
	default:
		return core.LOG_ERR
	}
}

// SetLevel implement core.ILogger
func (s *XormLogger) SetLevel(l core.LogLevel) {
}

// ShowSQL implement core.ILogger
func (s *XormLogger) ShowSQL(show ...bool) {
	s.grafanaLog.Error("ShowSQL", "show", "show")
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (s *XormLogger) IsShowSQL() bool {
	return s.showSQL
}

`
	ServiceSqlstoreLoginAttempt = `
package sqlstore

import (
	"strconv"
	"time"

	"{{.Dir}}/pkg/bus"
	m "{{.Dir}}/pkg/models"
)

var getTimeNow = time.Now

func init() {
	bus.AddHandler("sql", CreateLoginAttempt)
	bus.AddHandler("sql", DeleteOldLoginAttempts)
	bus.AddHandler("sql", GetUserLoginAttemptCount)
}

func CreateLoginAttempt(cmd *m.CreateLoginAttemptCommand) error {
	return inTransaction(func(sess *DBSession) error {
		loginAttempt := m.LoginAttempt{
			Username:  cmd.Username,
			IpAddress: cmd.IpAddress,
			Created:   getTimeNow().Unix(),
		}

		if _, err := sess.Insert(&loginAttempt); err != nil {
			return err
		}

		cmd.Result = loginAttempt

		return nil
	})
}

func DeleteOldLoginAttempts(cmd *m.DeleteOldLoginAttemptsCommand) error {
	return inTransaction(func(sess *DBSession) error {
		var maxId int64
		sql := "SELECT max(id) as id FROM login_attempt WHERE created < ?"
		result, err := sess.Query(sql, cmd.OlderThan.Unix())

		if err != nil {
			return err
		}
		// nolint: gosimple
		if result == nil || len(result) == 0 || result[0] == nil {
			return nil
		}

		maxId = toInt64(result[0]["id"])

		if maxId == 0 {
			return nil
		}

		sql = "DELETE FROM login_attempt WHERE id <= ?"

		if result, err := sess.Exec(sql, maxId); err != nil {
			return err
		} else if cmd.DeletedRows, err = result.RowsAffected(); err != nil {
			return err
		}

		return nil
	})
}

func GetUserLoginAttemptCount(query *m.GetUserLoginAttemptCountQuery) error {
	loginAttempt := new(m.LoginAttempt)
	total, err := x.
		Where("username = ?", query.Username).
		And("created >= ?", query.Since.Unix()).
		Count(loginAttempt)

	if err != nil {
		return err
	}

	query.Result = total
	return nil
}

func toInt64(i interface{}) int64 {
	switch i := i.(type) {
	case []byte:
		n, _ := strconv.ParseInt(string(i), 10, 64)
		return n
	case int:
		return int64(i)
	case int64:
		return i
	}
	return 0
}

`
	ServiceSqlstoreLoginAttemptTest = `
package sqlstore

import (
	"testing"
	"time"

	m "{{.Dir}}/pkg/models"
	. "github.com/smartystreets/goconvey/convey"
)

func mockTime(mock time.Time) time.Time {
	getTimeNow = func() time.Time { return mock }
	return mock
}

func TestLoginAttempts(t *testing.T) {
	Convey("Testing Login Attempts DB Access", t, func() {
		InitTestDB(t)

		user := "user"
		beginningOfTime := mockTime(time.Date(2017, 10, 22, 8, 0, 0, 0, time.Local))

		err := CreateLoginAttempt(&m.CreateLoginAttemptCommand{
			Username:  user,
			IpAddress: "192.168.0.1",
		})
		So(err, ShouldBeNil)

		timePlusOneMinute := mockTime(beginningOfTime.Add(time.Minute * 1))

		err = CreateLoginAttempt(&m.CreateLoginAttemptCommand{
			Username:  user,
			IpAddress: "192.168.0.1",
		})
		So(err, ShouldBeNil)

		timePlusTwoMinutes := mockTime(beginningOfTime.Add(time.Minute * 2))

		err = CreateLoginAttempt(&m.CreateLoginAttemptCommand{
			Username:  user,
			IpAddress: "192.168.0.1",
		})
		So(err, ShouldBeNil)

		Convey("Should return a total count of zero login attempts when comparing since beginning of time + 2min and 1s", func() {
			query := m.GetUserLoginAttemptCountQuery{
				Username: user,
				Since:    timePlusTwoMinutes.Add(time.Second * 1),
			}
			err := GetUserLoginAttemptCount(&query)
			So(err, ShouldBeNil)
			So(query.Result, ShouldEqual, 0)
		})

		Convey("Should return the total count of login attempts since beginning of time", func() {
			query := m.GetUserLoginAttemptCountQuery{
				Username: user,
				Since:    beginningOfTime,
			}
			err := GetUserLoginAttemptCount(&query)
			So(err, ShouldBeNil)
			So(query.Result, ShouldEqual, 3)
		})

		Convey("Should return the total count of login attempts since beginning of time + 1min", func() {
			query := m.GetUserLoginAttemptCountQuery{
				Username: user,
				Since:    timePlusOneMinute,
			}
			err := GetUserLoginAttemptCount(&query)
			So(err, ShouldBeNil)
			So(query.Result, ShouldEqual, 2)
		})

		Convey("Should return the total count of login attempts since beginning of time + 2min", func() {
			query := m.GetUserLoginAttemptCountQuery{
				Username: user,
				Since:    timePlusTwoMinutes,
			}
			err := GetUserLoginAttemptCount(&query)
			So(err, ShouldBeNil)
			So(query.Result, ShouldEqual, 1)
		})

		Convey("Should return deleted rows older than beginning of time", func() {
			cmd := m.DeleteOldLoginAttemptsCommand{
				OlderThan: beginningOfTime,
			}
			err := DeleteOldLoginAttempts(&cmd)

			So(err, ShouldBeNil)
			So(cmd.DeletedRows, ShouldEqual, 0)
		})

		Convey("Should return deleted rows older than beginning of time + 1min", func() {
			cmd := m.DeleteOldLoginAttemptsCommand{
				OlderThan: timePlusOneMinute,
			}
			err := DeleteOldLoginAttempts(&cmd)

			So(err, ShouldBeNil)
			So(cmd.DeletedRows, ShouldEqual, 1)
		})

		Convey("Should return deleted rows older than beginning of time + 2min", func() {
			cmd := m.DeleteOldLoginAttemptsCommand{
				OlderThan: timePlusTwoMinutes,
			}
			err := DeleteOldLoginAttempts(&cmd)

			So(err, ShouldBeNil)
			So(cmd.DeletedRows, ShouldEqual, 2)
		})

		Convey("Should return deleted rows older than beginning of time + 2min and 1s", func() {
			cmd := m.DeleteOldLoginAttemptsCommand{
				OlderThan: timePlusTwoMinutes.Add(time.Second * 1),
			}
			err := DeleteOldLoginAttempts(&cmd)

			So(err, ShouldBeNil)
			So(cmd.DeletedRows, ShouldEqual, 3)
		})
	})
}

`
	ServiceSqlstoreSession = `
package sqlstore

import (
	"context"
	"reflect"

	"github.com/go-xorm/xorm"
)

type DBSession struct {
	*xorm.Session
	events []interface{}
}

type dbTransactionFunc func(sess *DBSession) error

func (sess *DBSession) publishAfterCommit(msg interface{}) {
	sess.events = append(sess.events, msg)
}

// NewSession returns a new DBSession
func (ss *SqlStore) NewSession() *DBSession {
	return &DBSession{Session: ss.engine.NewSession()}
}

func newSession() *DBSession {
	return &DBSession{Session: x.NewSession()}
}

func startSession(ctx context.Context, engine *xorm.Engine, beginTran bool) (*DBSession, error) {
	value := ctx.Value(ContextSessionName)
	var sess *DBSession
	sess, ok := value.(*DBSession)

	if ok {
		return sess, nil
	}

	newSess := &DBSession{Session: engine.NewSession()}
	if beginTran {
		err := newSess.Begin()
		if err != nil {
			return nil, err
		}
	}
	return newSess, nil
}

// WithDbSession calls the callback with an session attached to the context.
func (ss *SqlStore) WithDbSession(ctx context.Context, callback dbTransactionFunc) error {
	sess, err := startSession(ctx, ss.engine, false)
	if err != nil {
		return err
	}

	return callback(sess)
}

func withDbSession(ctx context.Context, callback dbTransactionFunc) error {
	sess, err := startSession(ctx, x, false)
	if err != nil {
		return err
	}

	return callback(sess)
}

func (sess *DBSession) InsertId(bean interface{}) (int64, error) {
	//table := sess.DB().Mapper.Obj2Table(getTypeName(bean))

	//dialect.PreInsertId(table, sess.Session)

	id, err := sess.Session.InsertOne(bean)

	//dialect.PostInsertId(table, sess.Session)

	return id, err
}

func getTypeName(bean interface{}) (res string) {
	t := reflect.TypeOf(bean)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

`
	ServiceSqlstore = `
package sqlstore

import (
	"fmt"
	"{{.Dir}}/pkg/bus"
	"{{.Dir}}/pkg/infra/localcache"
	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/registry"
	"{{.Dir}}/pkg/services/annotations"
	"{{.Dir}}/pkg/services/sqlstore/migrator"
	"{{.Dir}}/pkg/services/sqlstore/sqlutil"
	"{{.Dir}}/pkg/setting"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const ContextSessionName = "db-session"

var (
	x       *xorm.Engine
	dialect migrator.Dialect

	sqlog log.Logger = log.New("sqlstore")
)

func init() {
	// This change will make xorm use an empty default schema for postgres and
	// by that mimic the functionality of how it was functioning before
	// xorm's changes above.
	xorm.DefaultPostgresSchema = ""

	registry.Register(&registry.Descriptor{
		Name:         "SqlStore",
		Instance:     &SqlStore{},
		InitPriority: registry.High,
	})
	fmt.Println("Initialized sqlstore ....")
}

type SqlStore struct {
	Cfg          *setting.Cfg             `+"`inject:\"\"`"+`
	Bus          bus.Bus         		  `+"`inject:\"\"`"+`
	CacheService *localcache.CacheService `+"`inject:\"\"`"+`

	dbCfg   DatabaseConfig
	Dialect migrator.Dialect
	engine  *xorm.Engine
	log     log.Logger
}

type DatabaseConfig struct {
	Type             string
	Host             string
	Name             string
	User             string
	Pwd              string
	Path             string
	SslMode          string
	CaCertPath       string
	ClientKeyPath    string
	ClientCertPath   string
	ServerCertName   string
	ConnectionString string
	MaxOpenConn      int
	MaxIdleConn      int
	ConnMaxLifetime  int
	CacheMode        string
	UrlQueryParams   map[string][]string
}

func (ss *SqlStore) Init() error {
	ss.log = log.New("sqlstore")
	ss.readConfig()

	engine, err := ss.getEngine()

	if err != nil {
		return fmt.Errorf("Fail to connect to database: %v", err)
	}

	ss.engine = engine

	// temporarily still set global var
	x = engine
	dialect = ss.Dialect

	// Init repo instances
	annotations.SetRepository(&SqlAnnotationRepo{})
	ss.Bus.SetTransactionManager(ss)

	// Register handlers
	ss.addUserQueryAndCommandHandlers()

	return nil
}

func (ss *SqlStore) readConfig() {
	sec := ss.Cfg.Raw.Section("database")

	cfgURL := sec.Key("url").String()
	if len(cfgURL) != 0 {
		dbURL, _ := url.Parse(cfgURL)
		ss.dbCfg.Type = dbURL.Scheme
		ss.dbCfg.Host = dbURL.Host

		pathSplit := strings.Split(dbURL.Path, "/")
		if len(pathSplit) > 1 {
			ss.dbCfg.Name = pathSplit[1]
		}

		userInfo := dbURL.User
		if userInfo != nil {
			ss.dbCfg.User = userInfo.Username()
			ss.dbCfg.Pwd, _ = userInfo.Password()
		}

		ss.dbCfg.UrlQueryParams = dbURL.Query()
	} else {
		ss.dbCfg.Type = sec.Key("type").String()
		ss.dbCfg.Host = sec.Key("host").String()
		ss.dbCfg.Name = sec.Key("name").String()
		ss.dbCfg.User = sec.Key("user").String()
		ss.dbCfg.ConnectionString = sec.Key("connection_string").String()
		ss.dbCfg.Pwd = sec.Key("password").String()
	}

	ss.dbCfg.MaxOpenConn = sec.Key("max_open_conn").MustInt(0)
	ss.dbCfg.MaxIdleConn = sec.Key("max_idle_conn").MustInt(2)
	ss.dbCfg.ConnMaxLifetime = sec.Key("conn_max_lifetime").MustInt(14400)

	ss.dbCfg.SslMode = sec.Key("ssl_mode").String()
	ss.dbCfg.CaCertPath = sec.Key("ca_cert_path").String()
	ss.dbCfg.ClientKeyPath = sec.Key("client_key_path").String()
	ss.dbCfg.ClientCertPath = sec.Key("client_cert_path").String()
	ss.dbCfg.ServerCertName = sec.Key("server_cert_name").String()
	ss.dbCfg.Path = sec.Key("path").MustString("data/grafana.db")

	ss.dbCfg.CacheMode = sec.Key("cache_mode").MustString("private")
}

func (ss *SqlStore) buildConnectionString() (string, error) {
	cnnstr := ss.dbCfg.ConnectionString

	// special case used by integration tests
	if cnnstr != "" {
		return cnnstr, nil
	}

	switch ss.dbCfg.Type {
	case migrator.MYSQL:
		protocol := "tcp"
		if strings.HasPrefix(ss.dbCfg.Host, "/") {
			protocol = "unix"
		}

		cnnstr = fmt.Sprintf("%s:%s@%s(%s)/%s?collation=utf8mb4_unicode_ci&allowNativePasswords=true",
			ss.dbCfg.User, ss.dbCfg.Pwd, protocol, ss.dbCfg.Host, ss.dbCfg.Name)

		if ss.dbCfg.SslMode == "true" || ss.dbCfg.SslMode == "skip-verify" {
			tlsCert, err := makeCert(ss.dbCfg)
			if err != nil {
				return "", err
			}
			mysql.RegisterTLSConfig("custom", tlsCert)
			cnnstr += "&tls=custom"
		}

		cnnstr += ss.buildExtraConnectionString('&')

	case migrator.SQLITE:
		// special case for tests
		if !filepath.IsAbs(ss.dbCfg.Path) {
			ss.dbCfg.Path = filepath.Join(ss.Cfg.DataPath, ss.dbCfg.Path)
		}
		os.MkdirAll(path.Dir(ss.dbCfg.Path), os.ModePerm)
		cnnstr = fmt.Sprintf("file:%s?cache=%s&mode=rwc", ss.dbCfg.Path, ss.dbCfg.CacheMode)
		cnnstr += ss.buildExtraConnectionString('&')
	default:
		return "", fmt.Errorf("Unknown database type: %s", ss.dbCfg.Type)
	}

	return cnnstr, nil
}

func (ss *SqlStore) getEngine() (*xorm.Engine, error) {
	connectionString, err := ss.buildConnectionString()

	if err != nil {
		return nil, err
	}

	sqlog.Info("Connecting to DB", "dbtype", ss.dbCfg.Type)
	engine, err := xorm.NewEngine(ss.dbCfg.Type, connectionString)
	if err != nil {
		return nil, err
	}

	engine.SetMaxOpenConns(ss.dbCfg.MaxOpenConn)
	engine.SetMaxIdleConns(ss.dbCfg.MaxIdleConn)
	engine.SetConnMaxLifetime(time.Second * time.Duration(ss.dbCfg.ConnMaxLifetime))

	// configure sql logging
	debugSql := ss.Cfg.Raw.Section("database").Key("log_queries").MustBool(false)
	if !debugSql {
		engine.SetLogger(&xorm.DiscardLogger{})
	} else {
		engine.SetLogger(NewXormLogger(log.LvlInfo, log.New("sqlstore.xorm")))
		engine.ShowSQL(true)
		engine.ShowExecTime(true)
	}

	return engine, nil
}

func (ss *SqlStore) buildExtraConnectionString(sep rune) string {
	if ss.dbCfg.UrlQueryParams == nil {
		return ""
	}

	var sb strings.Builder
	for key, values := range ss.dbCfg.UrlQueryParams {
		for _, value := range values {
			sb.WriteRune(sep)
			sb.WriteString(key)
			sb.WriteRune('=')
			sb.WriteString(value)
		}
	}
	return sb.String()
}

// Interface of arguments for testing db
type ITestDB interface {
	Helper()
	Fatalf(format string, args ...interface{})
}

// InitTestDB initiliaze test DB
func InitTestDB(t ITestDB) *SqlStore {
	t.Helper()
	sqlstore := &SqlStore{}
	//sqlstore.skipEnsureAdmin = true
	sqlstore.Bus = bus.New()

	dbType := migrator.SQLITE

	// environment variable present for test db?
	if db, present := os.LookupEnv("GRAFANA_TEST_DB"); present {
		dbType = db
	}

	// set test db config
	sqlstore.Cfg = setting.NewCfg()
	sec, _ := sqlstore.Cfg.Raw.NewSection("database")
	sec.NewKey("type", dbType)

	switch dbType {
	case "mysql":
		sec.NewKey("connection_string", sqlutil.TestDB_Mysql.ConnStr)
	case "postgres":
		sec.NewKey("connection_string", sqlutil.TestDB_Postgres.ConnStr)
	default:
		sec.NewKey("connection_string", sqlutil.TestDB_Sqlite3.ConnStr)
	}

	// need to get engine to clean db before we init
	engine, err := xorm.NewEngine(dbType, sec.Key("connection_string").String())
	if err != nil {
		t.Fatalf("Failed to init test database: %v", err)
	}

	sqlstore.Dialect = migrator.NewDialect(engine)

	// temp global var until we get rid of global vars
	dialect = sqlstore.Dialect

	if err := dialect.CleanDB(); err != nil {
		t.Fatalf("Failed to clean test db %v", err)
	}

	if err := sqlstore.Init(); err != nil {
		t.Fatalf("Failed to init test database: %v", err)
	}

	sqlstore.engine.DatabaseTZ = time.UTC
	sqlstore.engine.TZLocation = time.UTC

	return sqlstore
}

`
	ServiceSqlstoreTest = `
package sqlstore

import (
	"{{.Dir}}/pkg/setting"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type sqlStoreTest struct {
	name          string
	dbType        string
	dbHost        string
	connStrValues []string
}

var sqlStoreTestCases = []sqlStoreTest{
	{
		name:          "MySQL IPv4",
		dbType:        "mysql",
		dbHost:        "1.2.3.4:5678",
		connStrValues: []string{"tcp(1.2.3.4:5678)"},
	},
	{
		name:          "MySQL IPv4 (Default Port)",
		dbType:        "mysql",
		dbHost:        "1.2.3.4",
		connStrValues: []string{"tcp(1.2.3.4)"},
	},
	{
		name:          "MySQL IPv6",
		dbType:        "mysql",
		dbHost:        "[fe80::24e8:31b2:91df:b177]:1234",
		connStrValues: []string{"tcp([fe80::24e8:31b2:91df:b177]:1234)"},
	},
	{
		name:          "MySQL IPv6 (Default Port)",
		dbType:        "mysql",
		dbHost:        "::1",
		connStrValues: []string{"tcp(::1)"},
	},
}

func TestSqlConnectionString(t *testing.T) {
	Convey("Testing SQL Connection Strings", t, func() {
		t.Helper()

		for _, testCase := range sqlStoreTestCases {
			Convey(testCase.name, func() {
				sqlstore := &SqlStore{}
				sqlstore.Cfg = makeSqlStoreTestConfig(testCase.dbType, testCase.dbHost)
				sqlstore.readConfig()

				connStr, err := sqlstore.buildConnectionString()

				So(err, ShouldBeNil)
				for _, connSubStr := range testCase.connStrValues {
					So(connStr, ShouldContainSubstring, connSubStr)
				}
			})
		}
	})
}

func makeSqlStoreTestConfig(dbType string, host string) *setting.Cfg {
	cfg := setting.NewCfg()

	sec, _ := cfg.Raw.NewSection("database")
	sec.NewKey("type", dbType)
	sec.NewKey("host", host)
	sec.NewKey("user", "user")
	sec.NewKey("name", "test_db")
	sec.NewKey("password", "pass")

	return cfg
}

`
	ServiceSqlstoreTags = `
package sqlstore

import "{{.Dir}}/pkg/models"

// Will insert if needed any new key/value pars and return ids
func EnsureTagsExist(sess *DBSession, tags []*models.Tag) ([]*models.Tag, error) {
	for _, tag := range tags {
		var existingTag models.Tag

		// check if it exists
		if exists, err := sess.Table("tag").Where("`+"`key`"+`=? AND ` +"`value`"+`=?", tag.Key, tag.Value).Get(&existingTag); err != nil {
			return nil, err
		} else if exists {
			tag.Id = existingTag.Id
		} else {
			if _, err := sess.Table("tag").Insert(tag); err != nil {
				return nil, err
			}
		}
	}

	return tags, nil
}

`
	ServiceSqlstoreTeam = `
package sqlstore

import (
	"bytes"
	"fmt"
	"time"

	"{{.Dir}}/pkg/bus"
	m "{{.Dir}}/pkg/models"
)

func init() {
	bus.AddHandler("sql", CreateTeam)
	bus.AddHandler("sql", UpdateTeam)
	bus.AddHandler("sql", DeleteTeam)
	bus.AddHandler("sql", SearchTeams)
	bus.AddHandler("sql", GetTeamById)
	bus.AddHandler("sql", GetTeamsByUser)

	bus.AddHandler("sql", AddTeamMember)
	bus.AddHandler("sql", UpdateTeamMember)
	bus.AddHandler("sql", RemoveTeamMember)
	bus.AddHandler("sql", GetTeamMembers)
}

func getTeamSearchSqlBase() string {
	return ` +"`"+ `SELECT
team.id as id,
team.org_id,
team.name as name,
team.email as email,
(SELECT COUNT(*) from team_member where team_member.team_id = team.id) as member_count,
team_member.permission
FROM team as team
INNER JOIN team_member on team.id = team_member.team_id AND team_member.user_id = ? ` +"`"+`
}

func getTeamSelectSqlBase() string {
	return ` +"`"+`SELECT
team.id as id,
team.org_id,
team.name as name,
team.email as email,
(SELECT COUNT(*) from team_member where team_member.team_id = team.id) as member_count
FROM team as team ` +"`"+` 
}

func CreateTeam(cmd *m.CreateTeamCommand) error {
	return inTransaction(func(sess *DBSession) error {

		if isNameTaken, err := isTeamNameTaken(cmd.OrgId, cmd.Name, 0, sess); err != nil {
			return err
		} else if isNameTaken {
			return m.ErrTeamNameTaken
		}

		team := m.Team{
			Name:    cmd.Name,
			Email:   cmd.Email,
			OrgId:   cmd.OrgId,
			Created: time.Now(),
			Updated: time.Now(),
		}

		_, err := sess.Insert(&team)

		cmd.Result = team

		return err
	})
}

func UpdateTeam(cmd *m.UpdateTeamCommand) error {
	return inTransaction(func(sess *DBSession) error {

		if isNameTaken, err := isTeamNameTaken(cmd.OrgId, cmd.Name, cmd.Id, sess); err != nil {
			return err
		} else if isNameTaken {
			return m.ErrTeamNameTaken
		}

		team := m.Team{
			Name:    cmd.Name,
			Email:   cmd.Email,
			Updated: time.Now(),
		}

		sess.MustCols("email")

		affectedRows, err := sess.ID(cmd.Id).Update(&team)

		if err != nil {
			return err
		}

		if affectedRows == 0 {
			return m.ErrTeamNotFound
		}

		return nil
	})
}

// DeleteTeam will delete a team, its member and any permissions connected to the team
func DeleteTeam(cmd *m.DeleteTeamCommand) error {
	return inTransaction(func(sess *DBSession) error {
		if _, err := teamExists(cmd.OrgId, cmd.Id, sess); err != nil {
			return err
		}

		deletes := []string{
			"DELETE FROM team_member WHERE org_id=? and team_id = ?",
			"DELETE FROM team WHERE org_id=? and id = ?",
			"DELETE FROM dashboard_acl WHERE org_id=? and team_id = ?",
		}

		for _, sql := range deletes {
			_, err := sess.Exec(sql, cmd.OrgId, cmd.Id)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func teamExists(orgId int64, teamId int64, sess *DBSession) (bool, error) {
	if res, err := sess.Query("SELECT 1 from team WHERE org_id=? and id=?", orgId, teamId); err != nil {
		return false, err
	} else if len(res) != 1 {
		return false, m.ErrTeamNotFound
	}

	return true, nil
}

func isTeamNameTaken(orgId int64, name string, existingId int64, sess *DBSession) (bool, error) {
	var team m.Team
	exists, err := sess.Where("org_id=? and name=?", orgId, name).Get(&team)

	if err != nil {
		return false, nil
	}

	if exists && existingId != team.Id {
		return true, nil
	}

	return false, nil
}

func SearchTeams(query *m.SearchTeamsQuery) error {
	query.Result = m.SearchTeamQueryResult{
		Teams: make([]*m.TeamDTO, 0),
	}
	queryWithWildcards := "%" + query.Query + "%"

	var sql bytes.Buffer
	params := make([]interface{}, 0)

	if query.UserIdFilter > 0 {
		sql.WriteString(getTeamSearchSqlBase())
		params = append(params, query.UserIdFilter)
	} else {
		sql.WriteString(getTeamSelectSqlBase())
	}
	sql.WriteString(` +"`WHERE team.org_id = ?`"+`)

	params = append(params, query.OrgId)

	if query.Query != "" {
		sql.WriteString(` +"`and team.name`"+ ` + dialect.LikeStr() + ` +"`?`"+`)
		params = append(params, queryWithWildcards)
	}

	if query.Name != "" {
		sql.WriteString(` +"`and team.name = ?`"+`)
		params = append(params, query.Name)
	}

	sql.WriteString(` +"`order by team.name asc`"+`)

	if query.Limit != 0 {
		offset := query.Limit * (query.Page - 1)
		sql.WriteString(dialect.LimitOffset(int64(query.Limit), int64(offset)))
	}

	if err := x.SQL(sql.String(), params...).Find(&query.Result.Teams); err != nil {
		return err
	}

	team := m.Team{}
	countSess := x.Table("team")
	if query.Query != "" {
		countSess.Where(` +"`name`"+ `+dialect.LikeStr()+` +"`?`"+`, queryWithWildcards)
	}

	if query.Name != "" {
		countSess.Where("name=?", query.Name)
	}

	count, err := countSess.Count(&team)
	query.Result.TotalCount = count

	return err
}

func GetTeamById(query *m.GetTeamByIdQuery) error {
	var sql bytes.Buffer

	sql.WriteString(getTeamSelectSqlBase())
	sql.WriteString(` +"`WHERE team.org_id = ? and team.id = ?`"+`)

	var team m.TeamDTO
	exists, err := x.SQL(sql.String(), query.OrgId, query.Id).Get(&team)

	if err != nil {
		return err
	}

	if !exists {
		return m.ErrTeamNotFound
	}

	query.Result = &team
	return nil
}

// GetTeamsByUser is used by the Guardian when checking a users' permissions
func GetTeamsByUser(query *m.GetTeamsByUserQuery) error {
	query.Result = make([]*m.TeamDTO, 0)

	var sql bytes.Buffer

	sql.WriteString(getTeamSelectSqlBase())
	sql.WriteString(` +"`INNER JOIN team_member on team.id = team_member.team_id`"+`)
	sql.WriteString(` +"`WHERE team.org_id = ? and team_member.user_id = ?`"+`)

	err := x.SQL(sql.String(), query.OrgId, query.UserId).Find(&query.Result)
	return err
}

// AddTeamMember adds a user to a team
func AddTeamMember(cmd *m.AddTeamMemberCommand) error {
	return inTransaction(func(sess *DBSession) error {
		if res, err := sess.Query("SELECT 1 from team_member WHERE org_id=? and team_id=? and user_id=?", cmd.OrgId, cmd.TeamId, cmd.UserId); err != nil {
			return err
		} else if len(res) == 1 {
			return m.ErrTeamMemberAlreadyAdded
		}

		if _, err := teamExists(cmd.OrgId, cmd.TeamId, sess); err != nil {
			return err
		}

		entity := m.TeamMember{
			OrgId:    cmd.OrgId,
			TeamId:   cmd.TeamId,
			UserId:   cmd.UserId,
			External: cmd.External,
			Created:  time.Now(),
			Updated:  time.Now(),
			//Permission: cmd.Permission,
		}

		_, err := sess.Insert(&entity)
		return err
	})
}

func getTeamMember(sess *DBSession, orgId int64, teamId int64, userId int64) (m.TeamMember, error) {
	rawSql := ` +"`SELECT * FROM team_member WHERE org_id=? and team_id=? and user_id=?`"+`
	var member m.TeamMember
	exists, err := sess.SQL(rawSql, orgId, teamId, userId).Get(&member)

	if err != nil {
		return member, err
	}
	if !exists {
		return member, m.ErrTeamMemberNotFound
	}

	return member, nil
}

// UpdateTeamMember updates a team member
func UpdateTeamMember(cmd *m.UpdateTeamMemberCommand) error {
	return inTransaction(func(sess *DBSession) error {
		member, err := getTeamMember(sess, cmd.OrgId, cmd.TeamId, cmd.UserId)
		if err != nil {
			return err
		}

		if cmd.ProtectLastAdmin {
			_, err := isLastAdmin(sess, cmd.OrgId, cmd.TeamId, cmd.UserId)
			if err != nil {
				return err
			}
		}

		//if cmd.Permission != m.PERMISSION_ADMIN { // make sure we don't get invalid permission levels in store
		//	cmd.Permission = 0
		//}

		//member.Permission = cmd.Permission
		_, err = sess.Cols("permission").Where("org_id=? and team_id=? and user_id=?", cmd.OrgId, cmd.TeamId, cmd.UserId).Update(member)

		return err
	})
}

// RemoveTeamMember removes a member from a team
func RemoveTeamMember(cmd *m.RemoveTeamMemberCommand) error {
	return inTransaction(func(sess *DBSession) error {
		if _, err := teamExists(cmd.OrgId, cmd.TeamId, sess); err != nil {
			return err
		}

		if cmd.ProtectLastAdmin {
			_, err := isLastAdmin(sess, cmd.OrgId, cmd.TeamId, cmd.UserId)
			if err != nil {
				return err
			}
		}

		var rawSql = "DELETE FROM team_member WHERE org_id=? and team_id=? and user_id=?"
		res, err := sess.Exec(rawSql, cmd.OrgId, cmd.TeamId, cmd.UserId)
		if err != nil {
			return err
		}
		rows, err := res.RowsAffected()
		if rows == 0 {
			return m.ErrTeamMemberNotFound
		}

		return err
	})
}

func isLastAdmin(sess *DBSession, orgId int64, teamId int64, userId int64) (bool, error) {
	rawSql := "SELECT user_id FROM team_member WHERE org_id=? and team_id=? and permission=?"
	userIds := []*int64{}

	//m.PERMISSION_ADMIN
	err := sess.SQL(rawSql, orgId, teamId, 2).Find(&userIds)
	if err != nil {
		return false, err
	}

	isAdmin := false
	for _, adminId := range userIds {
		if userId == *adminId {
			isAdmin = true
			break
		}
	}

	if isAdmin && len(userIds) == 1 {
		return true, m.ErrLastTeamAdmin
	}

	return false, err
}

// GetTeamMembers return a list of members for the specified team
func GetTeamMembers(query *m.GetTeamMembersQuery) error {
	query.Result = make([]*m.TeamMemberDTO, 0)
	sess := x.Table("team_member")
	sess.Join("INNER", x.Dialect().Quote("user"), fmt.Sprintf("team_member.user_id=%s.id", x.Dialect().Quote("user")))
	if query.OrgId != 0 {
		sess.Where("team_member.org_id=?", query.OrgId)
	}
	if query.TeamId != 0 {
		sess.Where("team_member.team_id=?", query.TeamId)
	}
	if query.UserId != 0 {
		sess.Where("team_member.user_id=?", query.UserId)
	}
	if query.External {
		sess.Where("team_member.external=?", dialect.BooleanStr(true))
	}
	sess.Cols("team_member.org_id", "team_member.team_id", "team_member.user_id", "user.email", "user.login", "team_member.external", "team_member.permission")
	sess.Asc("user.login", "user.email")

	err := sess.Find(&query.Result)
	return err
}

`
	ServiceSqlstoreTlsMySQL = `
package sqlstore

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"{{.Dir}}/pkg/infra/log"
	"io/ioutil"
)

var tlslog = log.New("tls_mysql")

func makeCert(config DatabaseConfig) (*tls.Config, error) {
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(config.CaCertPath)
	if err != nil {
		return nil, fmt.Errorf("Could not read DB CA Cert path: %v", config.CaCertPath)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs: rootCertPool,
	}
	if config.ClientCertPath != "" && config.ClientKeyPath != "" {
		tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			tlslog.Debug("Loading client certificate")
			cert, err := tls.LoadX509KeyPair(config.ClientCertPath, config.ClientKeyPath)
			return &cert, err
		}
	}
	tlsConfig.ServerName = config.ServerCertName
	if config.SslMode == "skip-verify" {
		tlsConfig.InsecureSkipVerify = true
	}
	// Return more meaningful error before it is too late
	if config.ServerCertName == "" && !tlsConfig.InsecureSkipVerify {
		return nil, fmt.Errorf("server_cert_name is missing. Consider using ssl_mode = skip-verify")
	}
	return tlsConfig, nil
}

`
	ServiceSqlstoreTransaction = `
package sqlstore

import (
	"context"
	"time"

	"{{.Dir}}/pkg/bus"
	"{{.Dir}}/pkg/infra/log"
	"github.com/go-xorm/xorm"
	sqlite3 "github.com/mattn/go-sqlite3"
)

// WithTransactionalDbSession calls the callback with an session within a transaction
func (ss *SqlStore) WithTransactionalDbSession(ctx context.Context, callback dbTransactionFunc) error {
	return inTransactionWithRetryCtx(ctx, ss.engine, callback, 0)
}

func (ss *SqlStore) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return ss.inTransactionWithRetry(ctx, fn, 0)
}

func (ss *SqlStore) inTransactionWithRetry(ctx context.Context, fn func(ctx context.Context) error, retry int) error {
	return inTransactionWithRetryCtx(ctx, ss.engine, func(sess *DBSession) error {
		withValue := context.WithValue(ctx, ContextSessionName, sess)
		return fn(withValue)
	}, retry)
}

func inTransactionWithRetry(callback dbTransactionFunc, retry int) error {
	return inTransactionWithRetryCtx(context.Background(), x, callback, retry)
}

func inTransactionWithRetryCtx(ctx context.Context, engine *xorm.Engine, callback dbTransactionFunc, retry int) error {
	sess, err := startSession(ctx, engine, true)
	if err != nil {
		return err
	}

	defer sess.Close()

	err = callback(sess)

	// special handling of database locked errors for sqlite, then we can retry 5 times
	if sqlError, ok := err.(sqlite3.Error); ok && retry < 5 {
		if sqlError.Code == sqlite3.ErrLocked || sqlError.Code == sqlite3.ErrBusy {
			sess.Rollback()
			time.Sleep(time.Millisecond * time.Duration(10))
			sqlog.Info("Database locked, sleeping then retrying", "error", err, "retry", retry)
			return inTransactionWithRetry(callback, retry+1)
		}
	}

	if err != nil {
		sess.Rollback()
		return err
	} else if err = sess.Commit(); err != nil {
		return err
	}

	if len(sess.events) > 0 {
		for _, e := range sess.events {
			if err = bus.Publish(e); err != nil {
				log.Error(3, "Failed to publish event after commit. error: %v", err)
			}
		}
	}

	return nil
}

func inTransaction(callback dbTransactionFunc) error {
	return inTransactionWithRetry(callback, 0)
}

func inTransactionCtx(ctx context.Context, callback dbTransactionFunc) error {
	return inTransactionWithRetryCtx(ctx, x, callback, 0)
}

`
	ServiceSqlstoreTransactionTest = `
package sqlstore

import (
	"errors"
	"testing"
)

var ErrProvokedError = errors.New("testing error")

func TestTransaction(t *testing.T) {
	//ss := InitTestDB(t)

	//Convey("InTransaction asdf asdf", t, func() {
	//	cmd := &models.AddApiKeyCommand{Key: "secret-key", Name: "key", OrgId: 1}
	//
	//	err := AddApiKey(cmd)
	//	So(err, ShouldBeNil)
	//
	//	deleteApiKeyCmd := &models.DeleteApiKeyCommand{Id: cmd.Result.Id, OrgId: 1}
	//
	//	Convey("can update key", func() {
	//		err := ss.InTransaction(context.Background(), func(ctx context.Context) error {
	//			return DeleteApiKeyCtx(ctx, deleteApiKeyCmd)
	//		})
	//
	//		So(err, ShouldBeNil)
	//
	//		query := &models.GetApiKeyByIdQuery{ApiKeyId: cmd.Result.Id}
	//		err = GetApiKeyById(query)
	//		So(err, ShouldEqual, models.ErrInvalidApiKey)
	//	})
	//
	//	Convey("won't update if one handler fails", func() {
	//		err := ss.InTransaction(context.Background(), func(ctx context.Context) error {
	//			err := DeleteApiKeyCtx(ctx, deleteApiKeyCmd)
	//			if err != nil {
	//				return err
	//			}
	//
	//			return ErrProvokedError
	//		})
	//
	//		So(err, ShouldEqual, ErrProvokedError)
	//
	//		query := &models.GetApiKeyByIdQuery{ApiKeyId: cmd.Result.Id}
	//		err = GetApiKeyById(query)
	//		So(err, ShouldBeNil)
	//		So(query.Result.Id, ShouldEqual, cmd.Result.Id)
	//	})
	//})
}

`
	ServiceSqlstoreUser = `
package sqlstore

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"{{.Dir}}/pkg/bus"
	"{{.Dir}}/pkg/events"
	"{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	"{{.Dir}}/pkg/util"
)

func (ss *SqlStore) addUserQueryAndCommandHandlers() {
	ss.Bus.AddHandler(ss.GetSignedInUserWithCache)

	bus.AddHandler("sql", GetUserById)
	bus.AddHandler("sql", UpdateUser)
	bus.AddHandler("sql", ChangeUserPassword)
	bus.AddHandler("sql", GetUserByLogin)
	bus.AddHandler("sql", GetUserByEmail)
	bus.AddHandler("sql", SetUsingOrg)
	bus.AddHandler("sql", UpdateUserLastSeenAt)
	bus.AddHandler("sql", GetUserProfile)
	bus.AddHandler("sql", SearchUsers)
	bus.AddHandler("sql", GetUserOrgList)
	bus.AddHandler("sql", DisableUser)
	bus.AddHandler("sql", BatchDisableUsers)
	bus.AddHandler("sql", DeleteUser)
	bus.AddHandler("sql", UpdateUserPermissions)
	bus.AddHandler("sql", SetUserHelpFlag)
	bus.AddHandlerCtx("sql", CreateUser)
}

func getOrgIdForNewUser(cmd *models.CreateUserCommand, sess *DBSession) (int64, error) {
	if cmd.SkipOrgSetup {
		return -1, nil
	}

	var org models.Org

	if setting.AutoAssignOrg {
		has, err := sess.Where("id=?", setting.AutoAssignOrgId).Get(&org)
		if err != nil {
			return 0, err
		}
		if has {
			return org.Id, nil
		}
		if setting.AutoAssignOrgId == 1 {
			org.Name = "Main Org."
			org.Id = int64(setting.AutoAssignOrgId)
		} else {
			sqlog.Info("Could not create user: organization id %v does not exist",
				setting.AutoAssignOrgId)
			return 0, fmt.Errorf("Could not create user: organization id %v does not exist",
				setting.AutoAssignOrgId)
		}
	} else {
		org.Name = cmd.OrgName
		if len(org.Name) == 0 {
			org.Name = util.StringsFallback2(cmd.Email, cmd.Login)
		}
	}

	org.Created = time.Now()
	org.Updated = time.Now()

	if org.Id != 0 {
		if _, err := sess.InsertId(&org); err != nil {
			return 0, err
		}
	} else {
		if _, err := sess.InsertOne(&org); err != nil {
			return 0, err
		}
	}

	sess.publishAfterCommit(&events.OrgCreated{
		Timestamp: org.Created,
		Id:        org.Id,
		Name:      org.Name,
	})

	return org.Id, nil
}

func CreateUser(ctx context.Context, cmd *models.CreateUserCommand) error {
	return inTransactionCtx(ctx, func(sess *DBSession) error {
		orgId, err := getOrgIdForNewUser(cmd, sess)
		if err != nil {
			return err
		}

		if cmd.Email == "" {
			cmd.Email = cmd.Login
		}

		// create user
		user := models.User{
			Email:         cmd.Email,
			Name:          cmd.Name,
			Login:         cmd.Login,
			Company:       cmd.Company,
			IsAdmin:       cmd.IsAdmin,
			OrgId:         orgId,
			EmailVerified: cmd.EmailVerified,
			Created:       time.Now(),
			Updated:       time.Now(),
			LastSeenAt:    time.Now().AddDate(-10, 0, 0),
		}

		user.Salt = util.GetRandomString(10)
		user.Rands = util.GetRandomString(10)

		if len(cmd.Password) > 0 {
			user.Password = util.EncodePassword(cmd.Password, user.Salt)
		}

		sess.UseBool("is_admin")

		if _, err := sess.Insert(&user); err != nil {
			return err
		}

		sess.publishAfterCommit(&events.UserCreated{
			Timestamp: user.Created,
			Id:        user.Id,
			Name:      user.Name,
			Login:     user.Login,
			Email:     user.Email,
		})

		cmd.Result = user

		// create org user link
		if !cmd.SkipOrgSetup {
			orgUser := models.OrgUser{
				OrgId:   orgId,
				UserId:  user.Id,
				Role:    models.ROLE_ADMIN,
				Created: time.Now(),
				Updated: time.Now(),
			}

			if setting.AutoAssignOrg && !user.IsAdmin {
				if len(cmd.DefaultOrgRole) > 0 {
					orgUser.Role = models.RoleType(cmd.DefaultOrgRole)
				} else {
					orgUser.Role = models.RoleType(setting.AutoAssignOrgRole)
				}
			}

			if _, err = sess.Insert(&orgUser); err != nil {
				return err
			}
		}

		return nil
	})
}

func GetUserById(query *models.GetUserByIdQuery) error {
	user := new(models.User)
	has, err := x.Id(query.Id).Get(user)

	if err != nil {
		return err
	} else if !has {
		return models.ErrUserNotFound
	}

	query.Result = user

	return nil
}

func GetUserByLogin(query *models.GetUserByLoginQuery) error {
	if query.LoginOrEmail == "" {
		return models.ErrUserNotFound
	}

	// Try and find the user by login first.
	// It's not sufficient to assume that a LoginOrEmail with an "@" is an email.
	user := &models.User{Login: query.LoginOrEmail}
	has, err := x.Get(user)

	if err != nil {
		return err
	}

	if !has && strings.Contains(query.LoginOrEmail, "@") {
		// If the user wasn't found, and it contains an "@" fallback to finding the
		// user by email.
		user = &models.User{Email: query.LoginOrEmail}
		has, err = x.Get(user)
	}

	if err != nil {
		return err
	} else if !has {
		return models.ErrUserNotFound
	}

	query.Result = user

	return nil
}

func GetUserByEmail(query *models.GetUserByEmailQuery) error {
	if query.Email == "" {
		return models.ErrUserNotFound
	}

	user := &models.User{Email: query.Email}
	has, err := x.Get(user)

	if err != nil {
		return err
	} else if !has {
		return models.ErrUserNotFound
	}

	query.Result = user

	return nil
}

func UpdateUser(cmd *models.UpdateUserCommand) error {
	return inTransaction(func(sess *DBSession) error {

		user := models.User{
			Name:    cmd.Name,
			Email:   cmd.Email,
			Login:   cmd.Login,
			Theme:   cmd.Theme,
			Updated: time.Now(),
		}

		if _, err := sess.ID(cmd.UserId).Update(&user); err != nil {
			return err
		}

		sess.publishAfterCommit(&events.UserUpdated{
			Timestamp: user.Created,
			Id:        user.Id,
			Name:      user.Name,
			Login:     user.Login,
			Email:     user.Email,
		})

		return nil
	})
}

func ChangeUserPassword(cmd *models.ChangeUserPasswordCommand) error {
	return inTransaction(func(sess *DBSession) error {

		user := models.User{
			Password: cmd.NewPassword,
			Updated:  time.Now(),
		}

		_, err := sess.ID(cmd.UserId).Update(&user)
		return err
	})
}

func UpdateUserLastSeenAt(cmd *models.UpdateUserLastSeenAtCommand) error {
	return inTransaction(func(sess *DBSession) error {
		user := models.User{
			Id:         cmd.UserId,
			LastSeenAt: time.Now(),
		}

		_, err := sess.ID(cmd.UserId).Update(&user)
		return err
	})
}

func SetUsingOrg(cmd *models.SetUsingOrgCommand) error {
	getOrgsForUserCmd := &models.GetUserOrgListQuery{UserId: cmd.UserId}
	GetUserOrgList(getOrgsForUserCmd)

	valid := false
	for _, other := range getOrgsForUserCmd.Result {
		if other.OrgId == cmd.OrgId {
			valid = true
		}
	}

	if !valid {
		return fmt.Errorf("user does not belong to org")
	}

	return inTransaction(func(sess *DBSession) error {
		return setUsingOrgInTransaction(sess, cmd.UserId, cmd.OrgId)
	})
}

func setUsingOrgInTransaction(sess *DBSession, userID int64, orgID int64) error {
	user := models.User{
		Id:    userID,
		OrgId: orgID,
	}

	_, err := sess.ID(userID).Update(&user)
	return err
}

func GetUserProfile(query *models.GetUserProfileQuery) error {
	var user models.User
	has, err := x.Id(query.UserId).Get(&user)

	if err != nil {
		return err
	} else if !has {
		return models.ErrUserNotFound
	}

	query.Result = models.UserProfileDTO{
		Id:             user.Id,
		Name:           user.Name,
		Email:          user.Email,
		Login:          user.Login,
		Theme:          user.Theme,
		IsGrafanaAdmin: user.IsAdmin,
		OrgId:          user.OrgId,
	}

	return err
}

func GetUserOrgList(query *models.GetUserOrgListQuery) error {
	query.Result = make([]*models.UserOrgDTO, 0)
	sess := x.Table("org_user")
	sess.Join("INNER", "org", "org_user.org_id=org.id")
	sess.Where("org_user.user_id=?", query.UserId)
	sess.Cols("org.name", "org_user.role", "org_user.org_id")
	sess.OrderBy("org.name")
	err := sess.Find(&query.Result)
	return err
}

func newSignedInUserCacheKey(orgID, userID int64) string {
	return fmt.Sprintf("signed-in-user-%d-%d", userID, orgID)
}

func (ss *SqlStore) GetSignedInUserWithCache(query *models.GetSignedInUserQuery) error {
	cacheKey := newSignedInUserCacheKey(query.OrgId, query.UserId)
	if cached, found := ss.CacheService.Get(cacheKey); found {
		query.Result = cached.(*models.SignedInUser)
		return nil
	}

	err := GetSignedInUser(query)
	if err != nil {
		return err
	}

	cacheKey = newSignedInUserCacheKey(query.Result.OrgId, query.UserId)
	ss.CacheService.Set(cacheKey, query.Result, time.Second*5)
	return nil
}

func GetSignedInUser(query *models.GetSignedInUserQuery) error {
	orgId := "u.org_id"
	if query.OrgId > 0 {
		orgId = strconv.FormatInt(query.OrgId, 10)
	}

	var rawSql = ` +"`"+ `SELECT
u.id             as user_id,
u.is_admin       as is_grafana_admin,
u.email          as email,
u.login          as login,
u.name           as name,
u.help_flags1    as help_flags1,
u.last_seen_at   as last_seen_at,
(SELECT COUNT(*) FROM org_user where org_user.user_id = u.id) as org_count,
org.name         as org_name,
org_user.role    as org_role,
org.id           as org_id
FROM `+"`"+ "+ dialect.Quote(\"user\") +`"+ ` as u
LEFT OUTER JOIN org_user on org_user.org_id = ` +"`"+"+  orgId + `"+` and org_user.user_id = u.id
LEFT OUTER JOIN org on org.id = org_user.org_id ` +"`"+ `

	sess := x.Table("user")
	if query.UserId > 0 {
		sess.SQL(rawSql+"WHERE u.id=?", query.UserId)
	} else if query.Login != "" {
		sess.SQL(rawSql+"WHERE u.login=?", query.Login)
	} else if query.Email != "" {
		sess.SQL(rawSql+"WHERE u.email=?", query.Email)
	}

	var user models.SignedInUser
	has, err := sess.Get(&user)
	if err != nil {
		return err
	} else if !has {
		return models.ErrUserNotFound
	}

	if user.OrgRole == "" {
		user.OrgId = -1
		user.OrgName = "Org missing"
	}

	getTeamsByUserQuery := &models.GetTeamsByUserQuery{OrgId: user.OrgId, UserId: user.UserId}
	err = GetTeamsByUser(getTeamsByUserQuery)
	if err != nil {
		return err
	}

	user.Teams = make([]int64, len(getTeamsByUserQuery.Result))
	for i, t := range getTeamsByUserQuery.Result {
		user.Teams[i] = t.Id
	}

	query.Result = &user
	return err
}

func SearchUsers(query *models.SearchUsersQuery) error {
	query.Result = models.SearchUserQueryResult{
		Users: make([]*models.UserSearchHitDTO, 0),
	}

	queryWithWildcards := "%" + query.Query + "%"

	whereConditions := make([]string, 0)
	whereParams := make([]interface{}, 0)
	sess := x.Table("user").Alias("u")

	// Join with only most recent auth module
	joinCondition := ` +"`"+`(
SELECT id from user_auth
WHERE user_auth.user_id = u.id
ORDER BY user_auth.created DESC ` +"`"+ `
	joinCondition = "user_auth.id=" + joinCondition + dialect.Limit(1) + ")"
	sess.Join("LEFT", "user_auth", joinCondition)

	if query.OrgId > 0 {
		whereConditions = append(whereConditions, "org_id = ?")
		whereParams = append(whereParams, query.OrgId)
	}

	if query.Query != "" {
		whereConditions = append(whereConditions, "(email "+dialect.LikeStr()+" ? OR name "+dialect.LikeStr()+" ? OR login "+dialect.LikeStr()+" ?)")
		whereParams = append(whereParams, queryWithWildcards, queryWithWildcards, queryWithWildcards)
	}

	if query.AuthModule != "" {
		whereConditions = append(
			whereConditions,
			` +"`"+`u.id IN (SELECT user_id
FROM user_auth
WHERE auth_module=?)`+"`"+`,
		)

		whereParams = append(whereParams, query.AuthModule)
	}

	if len(whereConditions) > 0 {
		sess.Where(strings.Join(whereConditions, " AND "), whereParams...)
	}

	offset := query.Limit * (query.Page - 1)
	sess.Limit(query.Limit, offset)
	sess.Cols("u.id", "u.email", "u.name", "u.login", "u.is_admin", "u.is_disabled", "u.last_seen_at", "user_auth.auth_module")
	sess.OrderBy("u.id")
	if err := sess.Find(&query.Result.Users); err != nil {
		return err
	}

	// get total
	user := models.User{}
	countSess := x.Table("user").Alias("u")

	if len(whereConditions) > 0 {
		countSess.Where(strings.Join(whereConditions, " AND "), whereParams...)
	}

	count, err := countSess.Count(&user)
	query.Result.TotalCount = count

	for _, user := range query.Result.Users {
		user.LastSeenAtAge = util.GetAgeString(user.LastSeenAt)
	}

	return err
}

func DisableUser(cmd *models.DisableUserCommand) error {
	user := models.User{}
	sess := x.Table("user")
	sess.ID(cmd.UserId).Get(&user)

	sess.UseBool("is_disabled")

	_, err := sess.ID(cmd.UserId).Update(&user)
	return err
}

func BatchDisableUsers(cmd *models.BatchDisableUsersCommand) error {
	return inTransaction(func(sess *DBSession) error {
		userIds := cmd.UserIds

		if len(userIds) == 0 {
			return nil
		}

		user_id_params := strings.Repeat(",?", len(userIds)-1)
		disableSQL := "UPDATE " + dialect.Quote("user") + " SET is_disabled=? WHERE Id IN (?" + user_id_params + ")"

		disableParams := []interface{}{disableSQL, cmd.IsDisabled}
		for _, v := range userIds {
			disableParams = append(disableParams, v)
		}

		_, err := sess.Exec(disableParams...)
		if err != nil {
			return err
		}

		return nil
	})
}

func DeleteUser(cmd *models.DeleteUserCommand) error {
	return inTransaction(func(sess *DBSession) error {
		return deleteUserInTransaction(sess, cmd)
	})
}

func deleteUserInTransaction(sess *DBSession, cmd *models.DeleteUserCommand) error {
	deletes := []string{
		"DELETE FROM star WHERE user_id = ?",
		"DELETE FROM " + dialect.Quote("user") + " WHERE id = ?",
		"DELETE FROM org_user WHERE user_id = ?",
		"DELETE FROM dashboard_acl WHERE user_id = ?",
		"DELETE FROM preferences WHERE user_id = ?",
		"DELETE FROM team_member WHERE user_id = ?",
		"DELETE FROM user_auth WHERE user_id = ?",
		"DELETE FROM user_auth_token WHERE user_id = ?",
		"DELETE FROM quota WHERE user_id = ?",
	}

	for _, sql := range deletes {
		_, err := sess.Exec(sql, cmd.UserId)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateUserPermissions(cmd *models.UpdateUserPermissionsCommand) error {
	return inTransaction(func(sess *DBSession) error {
		user := models.User{}
		sess.ID(cmd.UserId).Get(&user)

		user.IsAdmin = cmd.IsGrafanaAdmin
		sess.UseBool("is_admin")

		_, err := sess.ID(user.Id).Update(&user)
		if err != nil {
			return err
		}

		// validate that after update there is at least one server admin
		if err := validateOneAdminLeft(sess); err != nil {
			return err
		}

		return nil
	})
}

func SetUserHelpFlag(cmd *models.SetUserHelpFlagCommand) error {
	return inTransaction(func(sess *DBSession) error {

		user := models.User{
			Id:         cmd.UserId,
			HelpFlags1: cmd.HelpFlags1,
			Updated:    time.Now(),
		}

		_, err := sess.ID(cmd.UserId).Cols("help_flags1").Update(&user)
		return err
	})
}

func validateOneAdminLeft(sess *DBSession) error {
	// validate that there is an admin user left
	count, err := sess.Where("is_admin=?", true).Count(&models.User{})
	if err != nil {
		return err
	}

	if count == 0 {
		return models.ErrLastGrafanaAdmin
	}

	return nil
}

`
	ServiceSqlstoreUserAuth = `
package sqlstore

import (
	"encoding/base64"
	"time"

	"{{.Dir}}/pkg/bus"
	"{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	"{{.Dir}}/pkg/util"
)

var getTime = time.Now

func init() {
	bus.AddHandler("sql", GetUserByAuthInfo)
	bus.AddHandler("sql", GetExternalUserInfoByLogin)
	bus.AddHandler("sql", GetAuthInfo)
	bus.AddHandler("sql", SetAuthInfo)
	bus.AddHandler("sql", UpdateAuthInfo)
	bus.AddHandler("sql", DeleteAuthInfo)
}

func GetUserByAuthInfo(query *models.GetUserByAuthInfoQuery) error {
	user := &models.User{}
	has := false
	var err error
	authQuery := &models.GetAuthInfoQuery{}

	// Try to find the user by auth module and id first
	if query.AuthModule != "" && query.AuthId != "" {
		authQuery.AuthModule = query.AuthModule
		authQuery.AuthId = query.AuthId

		err = GetAuthInfo(authQuery)
		if err != models.ErrUserNotFound {
			if err != nil {
				return err
			}

			// if user id was specified and doesn't match the user_auth entry, remove it
			if query.UserId != 0 && query.UserId != authQuery.Result.UserId {
				err = DeleteAuthInfo(&models.DeleteAuthInfoCommand{
					UserAuth: authQuery.Result,
				})
				if err != nil {
					sqlog.Error("Error removing user_auth entry", "error", err)
				}

				authQuery.Result = nil
			} else {
				has, err = x.Id(authQuery.Result.UserId).Get(user)
				if err != nil {
					return err
				}

				if !has {
					// if the user has been deleted then remove the entry
					err = DeleteAuthInfo(&models.DeleteAuthInfoCommand{
						UserAuth: authQuery.Result,
					})
					if err != nil {
						sqlog.Error("Error removing user_auth entry", "error", err)
					}

					authQuery.Result = nil
				}
			}
		}
	}

	// If not found, try to find the user by id
	if !has && query.UserId != 0 {
		has, err = x.Id(query.UserId).Get(user)
		if err != nil {
			return err
		}
	}

	// If not found, try to find the user by email address
	if !has && query.Email != "" {
		user = &models.User{Email: query.Email}
		has, err = x.Get(user)
		if err != nil {
			return err
		}
	}

	// If not found, try to find the user by login
	if !has && query.Login != "" {
		user = &models.User{Login: query.Login}
		has, err = x.Get(user)
		if err != nil {
			return err
		}
	}

	// No user found
	if !has {
		return models.ErrUserNotFound
	}

	// create authInfo record to link accounts
	if authQuery.Result == nil && query.AuthModule != "" {
		cmd2 := &models.SetAuthInfoCommand{
			UserId:     user.Id,
			AuthModule: query.AuthModule,
			AuthId:     query.AuthId,
		}
		if err := SetAuthInfo(cmd2); err != nil {
			return err
		}
	}

	query.Result = user
	return nil
}

func GetExternalUserInfoByLogin(query *models.GetExternalUserInfoByLoginQuery) error {
	userQuery := models.GetUserByLoginQuery{LoginOrEmail: query.LoginOrEmail}
	err := bus.Dispatch(&userQuery)
	if err != nil {
		return err
	}

	authInfoQuery := &models.GetAuthInfoQuery{UserId: userQuery.Result.Id}
	if err := bus.Dispatch(authInfoQuery); err != nil {
		return err
	}

	query.Result = &models.ExternalUserInfo{
		UserId:     userQuery.Result.Id,
		Login:      userQuery.Result.Login,
		Email:      userQuery.Result.Email,
		Name:       userQuery.Result.Name,
		AuthModule: authInfoQuery.Result.AuthModule,
		AuthId:     authInfoQuery.Result.AuthId,
	}
	return nil
}

func GetAuthInfo(query *models.GetAuthInfoQuery) error {
	userAuth := &models.UserAuth{
		UserId:     query.UserId,
		AuthModule: query.AuthModule,
		AuthId:     query.AuthId,
	}
	has, err := x.Desc("created").Get(userAuth)
	if err != nil {
		return err
	}
	if !has {
		return models.ErrUserNotFound
	}

	secretAccessToken, err := decodeAndDecrypt(userAuth.OAuthAccessToken)
	if err != nil {
		return err
	}
	secretRefreshToken, err := decodeAndDecrypt(userAuth.OAuthRefreshToken)
	if err != nil {
		return err
	}
	secretTokenType, err := decodeAndDecrypt(userAuth.OAuthTokenType)
	if err != nil {
		return err
	}
	userAuth.OAuthAccessToken = secretAccessToken
	userAuth.OAuthRefreshToken = secretRefreshToken
	userAuth.OAuthTokenType = secretTokenType

	query.Result = userAuth
	return nil
}

func SetAuthInfo(cmd *models.SetAuthInfoCommand) error {
	return inTransaction(func(sess *DBSession) error {
		authUser := &models.UserAuth{
			UserId:     cmd.UserId,
			AuthModule: cmd.AuthModule,
			AuthId:     cmd.AuthId,
			Created:    getTime(),
		}

		if cmd.OAuthToken != nil {
			secretAccessToken, err := encryptAndEncode(cmd.OAuthToken.AccessToken)
			if err != nil {
				return err
			}
			secretRefreshToken, err := encryptAndEncode(cmd.OAuthToken.RefreshToken)
			if err != nil {
				return err
			}
			secretTokenType, err := encryptAndEncode(cmd.OAuthToken.TokenType)
			if err != nil {
				return err
			}

			authUser.OAuthAccessToken = secretAccessToken
			authUser.OAuthRefreshToken = secretRefreshToken
			authUser.OAuthTokenType = secretTokenType
			authUser.OAuthExpiry = cmd.OAuthToken.Expiry
		}

		_, err := sess.Insert(authUser)
		return err
	})
}

func UpdateAuthInfo(cmd *models.UpdateAuthInfoCommand) error {
	return inTransaction(func(sess *DBSession) error {
		authUser := &models.UserAuth{
			UserId:     cmd.UserId,
			AuthModule: cmd.AuthModule,
			AuthId:     cmd.AuthId,
			Created:    getTime(),
		}

		if cmd.OAuthToken != nil {
			secretAccessToken, err := encryptAndEncode(cmd.OAuthToken.AccessToken)
			if err != nil {
				return err
			}
			secretRefreshToken, err := encryptAndEncode(cmd.OAuthToken.RefreshToken)
			if err != nil {
				return err
			}
			secretTokenType, err := encryptAndEncode(cmd.OAuthToken.TokenType)
			if err != nil {
				return err
			}

			authUser.OAuthAccessToken = secretAccessToken
			authUser.OAuthRefreshToken = secretRefreshToken
			authUser.OAuthTokenType = secretTokenType
			authUser.OAuthExpiry = cmd.OAuthToken.Expiry
		}

		cond := &models.UserAuth{
			UserId:     cmd.UserId,
			AuthModule: cmd.AuthModule,
		}
		upd, err := sess.Update(authUser, cond)
		sqlog.Debug("Updated user_auth", "user_id", cmd.UserId, "auth_module", cmd.AuthModule, "rows", upd)
		return err
	})
}

func DeleteAuthInfo(cmd *models.DeleteAuthInfoCommand) error {
	return inTransaction(func(sess *DBSession) error {
		_, err := sess.Delete(cmd.UserAuth)
		return err
	})
}

// decodeAndDecrypt will decode the string with the standard bas64 decoder
// and then decrypt it with grafana's secretKey
func decodeAndDecrypt(s string) (string, error) {
	// Bail out if empty string since it'll cause a segfault in util.Decrypt
	if s == "" {
		return "", nil
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	decrypted, err := util.Decrypt(decoded, setting.SecretKey)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// encryptAndEncode will encrypt a string with grafana's secretKey, and
// then encode it with the standard bas64 encoder
func encryptAndEncode(s string) (string, error) {
	encrypted, err := util.Encrypt([]byte(s), setting.SecretKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

`
	ServiceSqlstoreUserTest = `
package sqlstore

import (
	"fmt"
	"{{.Dir}}/pkg/models"
)

func GetOrgUsersForTest(query *models.GetOrgUsersQuery) error {
	query.Result = make([]*models.OrgUserDTO, 0)
	sess := x.Table("org_user")
	sess.Join("LEFT ", x.Dialect().Quote("user"), fmt.Sprintf("org_user.user_id=%s.id", x.Dialect().Quote("user")))
	sess.Where("org_user.org_id=?", query.OrgId)
	sess.Cols("org_user.org_id", "org_user.user_id", "user.email", "user.login", "org_user.role")

	err := sess.Find(&query.Result)
	return err
}

`

)
