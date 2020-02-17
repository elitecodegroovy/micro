package template4app

var (
	ServiceServerlockModel = `
package serverlock

type serverLock struct {
	Id            int64
	OperationUid  string
	LastExecution int64
	Version       int64
}

`
	ServiceServerlock = `
package serverlock

import (
	"context"
	"time"

	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/registry"
	"{{.Dir}}/pkg/services/sqlstore"
)

func init() {
	registry.RegisterService(&ServerLockService{})
}

// ServerLockService allows servers in HA mode to claim a lock
// and execute an function if the server was granted the lock
type ServerLockService struct {
	SQLStore *sqlstore.SqlStore ` +"`inject:\"\"`"+`
	log      log.Logger
}

// Init this service
func (sl *ServerLockService) Init() error {
	sl.log = log.New("lockservice")
	return nil
}

// LockAndExecute try to create a lock for this server and only executes the
// ` +"`fn`"+` function when successful. This should not be used at low internal. But services
// that needs to be run once every ex 10m.
func (sl *ServerLockService) LockAndExecute(ctx context.Context, actionName string, maxInterval time.Duration, fn func()) error {
	// gets or creates a lockable row
	rowLock, err := sl.getOrCreate(ctx, actionName)
	if err != nil {
		return err
	}

	// avoid execution if last lock happened less than ` +"`maxInterval`"+` ago
	if rowLock.LastExecution != 0 {
		lastExeuctionTime := time.Unix(rowLock.LastExecution, 0)
		if lastExeuctionTime.Unix() > time.Now().Add(-maxInterval).Unix() {
			return nil
		}
	}

	// try to get lock based on rowLow version
	acquiredLock, err := sl.acquireLock(ctx, rowLock)
	if err != nil {
		return err
	}

	if acquiredLock {
		fn()
	}

	return nil
}

func (sl *ServerLockService) acquireLock(ctx context.Context, serverLock *serverLock) (bool, error) {
	var result bool

	err := sl.SQLStore.WithDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		newVersion := serverLock.Version + 1
		sql := `+"`UPDATE server_lock SET version = ?, last_execution = ? WHERE id = ? AND version = ?`"+`

		res, err := dbSession.Exec(sql, newVersion, time.Now().Unix(), serverLock.Id, serverLock.Version)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		result = affected == 1

		return err
	})

	return result, err
}

func (sl *ServerLockService) getOrCreate(ctx context.Context, actionName string) (*serverLock, error) {
	var result *serverLock

	err := sl.SQLStore.WithTransactionalDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		lockRows := []*serverLock{}
		err := dbSession.Where("operation_uid = ?", actionName).Find(&lockRows)
		if err != nil {
			return err
		}

		if len(lockRows) > 0 {
			result = lockRows[0]
			return nil
		}

		lockRow := &serverLock{
			OperationUid:  actionName,
			LastExecution: 0,
		}

		_, err = dbSession.Insert(lockRow)
		if err != nil {
			return err
		}

		result = lockRow

		return nil
	})

	return result, err
}

`
	ServiceServerlockIntegrationTest = `
// +build integration

package serverlock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerLok(t *testing.T) {
	sl := createTestableServerLock(t)

	counter := 0
	fn := func() { counter++ }
	atInterval := time.Second * 1
	ctx := context.Background()

	//this time ` +"`fn`"+` should be executed
	assert.Nil(t, sl.LockAndExecute(ctx, "test-operation", atInterval, fn))

	//this should not execute ` +"`fn`"+`
	assert.Nil(t, sl.LockAndExecute(ctx, "test-operation", atInterval, fn))
	assert.Nil(t, sl.LockAndExecute(ctx, "test-operation", atInterval, fn))

	// wait 2 second.
	<-time.After(time.Second * 2)

	// now ` +"`fn`"+` should be executed again
	err := sl.LockAndExecute(ctx, "test-operation", atInterval, fn)
	assert.Nil(t, err)
	assert.Equal(t, counter, 2)
}

`
	ServiceServerlockTest = `
package serverlock

import (
	"context"
	"testing"

	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/services/sqlstore"
	. "github.com/smartystreets/goconvey/convey"
)

func createTestableServerLock(t *testing.T) *ServerLockService {
	t.Helper()

	sqlstore := sqlstore.InitTestDB(t)

	return &ServerLockService{
		SQLStore: sqlstore,
		log:      log.New("test-logger"),
	}
}

func TestServerLock(t *testing.T) {
	Convey("Server lock", t, func() {
		sl := createTestableServerLock(t)
		operationUID := "test-operation"

		first, err := sl.getOrCreate(context.Background(), operationUID)
		So(err, ShouldBeNil)

		lastExecution := first.LastExecution
		Convey("trying to create three new row locks", func() {
			for i := 0; i < 3; i++ {
				first, err = sl.getOrCreate(context.Background(), operationUID)
				So(err, ShouldBeNil)
				So(first.OperationUid, ShouldEqual, operationUID)
				So(first.Id, ShouldEqual, 1)
			}

			Convey("Should not create new since lock already exist", func() {
				So(lastExecution, ShouldEqual, first.LastExecution)
			})
		})

		Convey("Should be able to create lock on first row", func() {
			gotLock, err := sl.acquireLock(context.Background(), first)
			So(err, ShouldBeNil)
			So(gotLock, ShouldBeTrue)

			gotLock, err = sl.acquireLock(context.Background(), first)
			So(err, ShouldBeNil)
			So(gotLock, ShouldBeFalse)
		})
	})
}

`
)
