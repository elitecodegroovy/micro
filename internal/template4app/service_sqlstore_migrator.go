package template4app

var (
	ServiceSqlstoreMigratorColumn = `
package migrator

// Notice
// code based on parts from from https://github.com/go-xorm/core/blob/3e0fa232ab5c90996406c0cd7ae86ad0e5ecf85f/column.go

type Column struct {
	Name            string
	Type            string
	Length          int
	Length2         int
	Nullable        bool
	IsPrimaryKey    bool
	IsAutoIncrement bool
	Default         string
}

func (col *Column) String(d Dialect) string {
	return d.ColString(col)
}

func (col *Column) StringNoPk(d Dialect) string {
	return d.ColStringNoPk(col)
}

`
	ServiceSqlStoreMigratorConditions = `
package migrator

type MigrationCondition interface {
	Sql(dialect Dialect) (string, []interface{})
	IsFulfilled(results []map[string][]byte) bool
}

type ExistsMigrationCondition struct{}

func (c *ExistsMigrationCondition) IsFulfilled(results []map[string][]byte) bool {
	return len(results) >= 1
}

type NotExistsMigrationCondition struct{}

func (c *NotExistsMigrationCondition) IsFulfilled(results []map[string][]byte) bool {
	return len(results) == 0
}

type IfIndexExistsCondition struct {
	ExistsMigrationCondition
	TableName string
	IndexName string
}

func (c *IfIndexExistsCondition) Sql(dialect Dialect) (string, []interface{}) {
	return dialect.IndexCheckSql(c.TableName, c.IndexName)
}

type IfIndexNotExistsCondition struct {
	NotExistsMigrationCondition
	TableName string
	IndexName string
}

func (c *IfIndexNotExistsCondition) Sql(dialect Dialect) (string, []interface{}) {
	return dialect.IndexCheckSql(c.TableName, c.IndexName)
}

type IfColumnNotExistsCondition struct {
	NotExistsMigrationCondition
	TableName  string
	ColumnName string
}

func (c *IfColumnNotExistsCondition) Sql(dialect Dialect) (string, []interface{}) {
	return dialect.ColumnCheckSql(c.TableName, c.ColumnName)
}

`
	ServiceSqlstoreMigratorDialect = `
package migrator

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
)

type Dialect interface {
	DriverName() string
	Quote(string) string
	AndStr() string
	AutoIncrStr() string
	OrStr() string
	EqStr() string
	ShowCreateNull() bool
	SqlType(col *Column) string
	SupportEngine() bool
	LikeStr() string
	Default(col *Column) string
	BooleanStr(bool) string
	DateTimeFunc(string) string

	CreateIndexSql(tableName string, index *Index) string
	CreateTableSql(table *Table) string
	AddColumnSql(tableName string, col *Column) string
	CopyTableData(sourceTable string, targetTable string, sourceCols []string, targetCols []string) string
	DropTable(tableName string) string
	DropIndexSql(tableName string, index *Index) string

	RenameTable(oldName string, newName string) string
	UpdateTableSql(tableName string, columns []*Column) string

	IndexCheckSql(tableName, indexName string) (string, []interface{})
	ColumnCheckSql(tableName, columnName string) (string, []interface{})

	ColString(*Column) string
	ColStringNoPk(*Column) string

	Limit(limit int64) string
	LimitOffset(limit int64, offset int64) string

	PreInsertId(table string, sess *xorm.Session) error
	PostInsertId(table string, sess *xorm.Session) error

	CleanDB() error
	NoOpSql() string

	IsUniqueConstraintViolation(err error) bool
	IsDeadlock(err error) bool
}

func NewDialect(engine *xorm.Engine) Dialect {
	name := engine.DriverName()
	switch name {
	case MYSQL:
		return NewMysqlDialect(engine)
	case SQLITE:
		return NewSqlite3Dialect(engine)
		//case POSTGRES:
		//	return NewPostgresDialect(engine)
	}

	panic("Unsupported database type: " + name)
}

type BaseDialect struct {
	dialect    Dialect
	engine     *xorm.Engine
	driverName string
}

func (d *BaseDialect) DriverName() string {
	return d.driverName
}

func (b *BaseDialect) ShowCreateNull() bool {
	return true
}

func (b *BaseDialect) AndStr() string {
	return "AND"
}

func (b *BaseDialect) LikeStr() string {
	return "LIKE"
}

func (b *BaseDialect) OrStr() string {
	return "OR"
}

func (b *BaseDialect) EqStr() string {
	return "="
}

func (b *BaseDialect) Default(col *Column) string {
	return col.Default
}

func (db *BaseDialect) DateTimeFunc(value string) string {
	return value
}

func (b *BaseDialect) CreateTableSql(table *Table) string {
	sql := "CREATE TABLE IF NOT EXISTS "
	sql += b.dialect.Quote(table.Name) + " (\n"

	pkList := table.PrimaryKeys

	for _, col := range table.Columns {
		if col.IsPrimaryKey && len(pkList) == 1 {
			sql += col.String(b.dialect)
		} else {
			sql += col.StringNoPk(b.dialect)
		}
		sql = strings.TrimSpace(sql)
		sql += "\n, "
	}

	if len(pkList) > 1 {
		quotedCols := []string{}
		for _, col := range pkList {
			quotedCols = append(quotedCols, b.dialect.Quote(col))
		}

		sql += "PRIMARY KEY ( " + strings.Join(quotedCols, ",") + " ), "
	}

	sql = sql[:len(sql)-2] + ")"
	if b.dialect.SupportEngine() {
		sql += " ENGINE=InnoDB DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci"
	}

	sql += ";"
	return sql
}

func (db *BaseDialect) AddColumnSql(tableName string, col *Column) string {
	return fmt.Sprintf("alter table %s ADD COLUMN %s", db.dialect.Quote(tableName), col.StringNoPk(db.dialect))
}

func (db *BaseDialect) CreateIndexSql(tableName string, index *Index) string {
	quote := db.dialect.Quote
	var unique string
	if index.Type == UniqueIndex {
		unique = " UNIQUE"
	}

	idxName := index.XName(tableName)

	quotedCols := []string{}
	for _, col := range index.Cols {
		quotedCols = append(quotedCols, db.dialect.Quote(col))
	}

	return fmt.Sprintf("CREATE%s INDEX %v ON %v (%v);", unique, quote(idxName), quote(tableName), strings.Join(quotedCols, ","))
}

func (db *BaseDialect) QuoteColList(cols []string) string {
	var sourceColsSql = ""
	for _, col := range cols {
		sourceColsSql += db.dialect.Quote(col)
		sourceColsSql += "\n, "
	}
	return strings.TrimSuffix(sourceColsSql, "\n, ")
}

func (db *BaseDialect) CopyTableData(sourceTable string, targetTable string, sourceCols []string, targetCols []string) string {
	sourceColsSql := db.QuoteColList(sourceCols)
	targetColsSql := db.QuoteColList(targetCols)

	quote := db.dialect.Quote
	return fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s", quote(targetTable), targetColsSql, sourceColsSql, quote(sourceTable))
}

func (db *BaseDialect) DropTable(tableName string) string {
	quote := db.dialect.Quote
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", quote(tableName))
}

func (db *BaseDialect) RenameTable(oldName string, newName string) string {
	quote := db.dialect.Quote
	return fmt.Sprintf("ALTER TABLE %s RENAME TO %s", quote(oldName), quote(newName))
}

func (db *BaseDialect) ColumnCheckSql(tableName, columnName string) (string, []interface{}) {
	return "", nil
}

func (db *BaseDialect) DropIndexSql(tableName string, index *Index) string {
	quote := db.dialect.Quote
	name := index.XName(tableName)
	return fmt.Sprintf("DROP INDEX %v ON %s", quote(name), quote(tableName))
}

func (db *BaseDialect) UpdateTableSql(tableName string, columns []*Column) string {
	return "-- NOT REQUIRED"
}

func (db *BaseDialect) ColString(col *Column) string {
	sql := db.dialect.Quote(col.Name) + " "

	sql += db.dialect.SqlType(col) + " "

	if col.IsPrimaryKey {
		sql += "PRIMARY KEY "
		if col.IsAutoIncrement {
			sql += db.dialect.AutoIncrStr() + " "
		}
	}

	if db.dialect.ShowCreateNull() {
		if col.Nullable {
			sql += "NULL "
		} else {
			sql += "NOT NULL "
		}
	}

	if col.Default != "" {
		sql += "DEFAULT " + db.dialect.Default(col) + " "
	}

	return sql
}

func (db *BaseDialect) ColStringNoPk(col *Column) string {
	sql := db.dialect.Quote(col.Name) + " "

	sql += db.dialect.SqlType(col) + " "

	if db.dialect.ShowCreateNull() {
		if col.Nullable {
			sql += "NULL "
		} else {
			sql += "NOT NULL "
		}
	}

	if col.Default != "" {
		sql += "DEFAULT " + db.dialect.Default(col) + " "
	}

	return sql
}

func (db *BaseDialect) Limit(limit int64) string {
	return fmt.Sprintf(" LIMIT %d", limit)
}

func (db *BaseDialect) LimitOffset(limit int64, offset int64) string {
	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}

func (db *BaseDialect) PreInsertId(table string, sess *xorm.Session) error {
	return nil
}

func (db *BaseDialect) PostInsertId(table string, sess *xorm.Session) error {
	return nil
}

func (db *BaseDialect) CleanDB() error {
	return nil
}

func (db *BaseDialect) NoOpSql() string {
	return "SELECT 0;"
}

`
	ServiceSqlstoreMigrator = `
package migrator

import (
	"time"

	"{{.Dir}}/pkg/infra/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Migrator struct {
	x          *xorm.Engine
	Dialect    Dialect
	migrations []Migration
	Logger     log.Logger
}

type MigrationLog struct {
	Id          int64
	MigrationId string
	Sql         string
	Success     bool
	Error       string
	Timestamp   time.Time
}

func NewMigrator(engine *xorm.Engine) *Migrator {
	mg := &Migrator{}
	mg.x = engine
	mg.Logger = log.New("migrator")
	mg.migrations = make([]Migration, 0)
	mg.Dialect = NewDialect(mg.x)
	return mg
}

func (mg *Migrator) MigrationsCount() int {
	return len(mg.migrations)
}

func (mg *Migrator) AddMigration(id string, m Migration) {
	m.SetId(id)
	mg.migrations = append(mg.migrations, m)
}

func (mg *Migrator) GetMigrationLog() (map[string]MigrationLog, error) {
	logMap := make(map[string]MigrationLog)
	logItems := make([]MigrationLog, 0)

	exists, err := mg.x.IsTableExist(new(MigrationLog))
	if err != nil {
		return nil, err
	}

	if !exists {
		return logMap, nil
	}

	if err = mg.x.Find(&logItems); err != nil {
		return nil, err
	}

	for _, logItem := range logItems {
		if !logItem.Success {
			continue
		}
		logMap[logItem.MigrationId] = logItem
	}

	return logMap, nil
}

func (mg *Migrator) Start() error {
	mg.Logger.Info("Starting DB migration")

	logMap, err := mg.GetMigrationLog()
	if err != nil {
		return err
	}

	for _, m := range mg.migrations {
		_, exists := logMap[m.Id()]
		if exists {
			mg.Logger.Debug("Skipping migration: Already executed", "id", m.Id())
			continue
		}

		sql := m.Sql(mg.Dialect)

		record := MigrationLog{
			MigrationId: m.Id(),
			Sql:         sql,
			Timestamp:   time.Now(),
		}

		err := mg.inTransaction(func(sess *xorm.Session) error {
			err := mg.exec(m, sess)
			if err != nil {
				mg.Logger.Error("Exec failed", "error", err, "sql", sql)
				record.Error = err.Error()
				sess.Insert(&record)
				return err
			}
			record.Success = true
			sess.Insert(&record)
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (mg *Migrator) exec(m Migration, sess *xorm.Session) error {
	mg.Logger.Info("Executing migration", "id", m.Id())

	condition := m.GetCondition()
	if condition != nil {
		sql, args := condition.Sql(mg.Dialect)

		if sql != "" {
			mg.Logger.Debug("Executing migration condition sql", "id", m.Id(), "sql", sql, "args", args)
			results, err := sess.SQL(sql, args...).Query()
			if err != nil {
				mg.Logger.Error("Executing migration condition failed", "id", m.Id(), "error", err)
				return err
			}

			if !condition.IsFulfilled(results) {
				mg.Logger.Warn("Skipping migration: Already executed, but not recorded in migration log", "id", m.Id())
				return nil
			}
		}
	}

	var err error
	if codeMigration, ok := m.(CodeMigration); ok {
		mg.Logger.Debug("Executing code migration", "id", m.Id())
		err = codeMigration.Exec(sess, mg)
	} else {
		sql := m.Sql(mg.Dialect)
		mg.Logger.Debug("Executing sql migration", "id", m.Id(), "sql", sql)
		_, err = sess.Exec(sql)
	}

	if err != nil {
		mg.Logger.Error("Executing migration failed", "id", m.Id(), "error", err)
		return err
	}

	return nil
}

type dbTransactionFunc func(sess *xorm.Session) error

func (mg *Migrator) inTransaction(callback dbTransactionFunc) error {
	var err error

	sess := mg.x.NewSession()
	defer sess.Close()

	if err = sess.Begin(); err != nil {
		return err
	}

	err = callback(sess)

	if err != nil {
		sess.Rollback()
		return err
	} else if err = sess.Commit(); err != nil {
		return err
	}

	return nil
}

`
	ServiceSqlstoreMigratorMySQLDialect = `
package migrator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Mysql struct {
	BaseDialect
}

func NewMysqlDialect(engine *xorm.Engine) *Mysql {
	d := Mysql{}
	d.BaseDialect.dialect = &d
	d.BaseDialect.engine = engine
	d.BaseDialect.driverName = MYSQL
	return &d
}

func (db *Mysql) SupportEngine() bool {
	return true
}

func (db *Mysql) Quote(name string) string {
	return " `+"`\""+ "+name+\"`" +`"
}

func (db *Mysql) AutoIncrStr() string {
	return "AUTO_INCREMENT"
}

func (db *Mysql) BooleanStr(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func (db *Mysql) SqlType(c *Column) string {
	var res string
	switch c.Type {
	case DB_Bool:
		res = DB_TinyInt
		c.Length = 1
	case DB_Serial:
		c.IsAutoIncrement = true
		c.IsPrimaryKey = true
		c.Nullable = false
		res = DB_Int
	case DB_BigSerial:
		c.IsAutoIncrement = true
		c.IsPrimaryKey = true
		c.Nullable = false
		res = DB_BigInt
	case DB_Bytea:
		res = DB_Blob
	case DB_TimeStampz:
		res = DB_Char
		c.Length = 64
	case DB_NVarchar:
		res = DB_Varchar
	default:
		res = c.Type
	}

	var hasLen1 = (c.Length > 0)
	var hasLen2 = (c.Length2 > 0)

	if res == DB_BigInt && !hasLen1 && !hasLen2 {
		c.Length = 20
		hasLen1 = true
	}

	if hasLen2 {
		res += "(" + strconv.Itoa(c.Length) + "," + strconv.Itoa(c.Length2) + ")"
	} else if hasLen1 {
		res += "(" + strconv.Itoa(c.Length) + ")"
	}

	switch c.Type {
	case DB_Char, DB_Varchar, DB_NVarchar, DB_TinyText, DB_Text, DB_MediumText, DB_LongText:
		res += " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
	}

	return res
}

func (db *Mysql) UpdateTableSql(tableName string, columns []*Column) string {
	var statements = []string{}

	statements = append(statements, "DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")

	for _, col := range columns {
		statements = append(statements, "MODIFY "+col.StringNoPk(db))
	}

	return "ALTER TABLE " + db.Quote(tableName) + " " + strings.Join(statements, ", ") + ";"
}

func (db *Mysql) IndexCheckSql(tableName, indexName string) (string, []interface{}) {
	args := []interface{}{tableName, indexName}
	sql := "SELECT 1 FROM " + db.Quote("INFORMATION_SCHEMA") + "." + db.Quote("STATISTICS") + " WHERE " + db.Quote("TABLE_SCHEMA") + " = DATABASE() AND " + db.Quote("TABLE_NAME") + "=? AND " + db.Quote("INDEX_NAME") + "=?"
	return sql, args
}

func (db *Mysql) ColumnCheckSql(tableName, columnName string) (string, []interface{}) {
	args := []interface{}{tableName, columnName}
	sql := "SELECT 1 FROM " + db.Quote("INFORMATION_SCHEMA") + "." + db.Quote("COLUMNS") + " WHERE " + db.Quote("TABLE_SCHEMA") + " = DATABASE() AND " + db.Quote("TABLE_NAME") + "=? AND " + db.Quote("COLUMN_NAME") + "=?"
	return sql, args
}

func (db *Mysql) CleanDB() error {
	tables, _ := db.engine.DBMetas()
	sess := db.engine.NewSession()
	defer sess.Close()

	for _, table := range tables {
		if _, err := sess.Exec("set foreign_key_checks = 0"); err != nil {
			return fmt.Errorf("failed to disable foreign key checks")
		}
		if _, err := sess.Exec("drop table " + table.Name + " ;"); err != nil {
			return fmt.Errorf("failed to delete table: %v, err: %v", table.Name, err)
		}
		if _, err := sess.Exec("set foreign_key_checks = 1"); err != nil {
			return fmt.Errorf("failed to disable foreign key checks")
		}
	}

	return nil
}

func (db *Mysql) isThisError(err error, errcode uint16) bool {
	if driverErr, ok := err.(*mysql.MySQLError); ok {
		if driverErr.Number == errcode {
			return true
		}
	}

	return false
}

func (db *Mysql) IsUniqueConstraintViolation(err error) bool {
	return db.isThisError(err, mysqlerr.ER_DUP_ENTRY)
}

func (db *Mysql) IsDeadlock(err error) bool {
	return db.isThisError(err, mysqlerr.ER_LOCK_DEADLOCK)
}

`
	ServiceSqlstoreMigratorSqliteDialect = `
package migrator

import (
	"fmt"

	"github.com/go-xorm/xorm"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type Sqlite3 struct {
	BaseDialect
}

func NewSqlite3Dialect(engine *xorm.Engine) *Sqlite3 {
	d := Sqlite3{}
	d.BaseDialect.dialect = &d
	d.BaseDialect.engine = engine
	d.BaseDialect.driverName = SQLITE
	return &d
}

func (db *Sqlite3) SupportEngine() bool {
	return false
}

func (db *Sqlite3) Quote(name string) string {
	return " `+"`\""+ "+name+\"`" +`"
}

func (db *Sqlite3) AutoIncrStr() string {
	return "AUTOINCREMENT"
}

func (db *Sqlite3) BooleanStr(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func (db *Sqlite3) DateTimeFunc(value string) string {
	return "datetime(" + value + ")"
}

func (db *Sqlite3) SqlType(c *Column) string {
	switch c.Type {
	case DB_Date, DB_DateTime, DB_TimeStamp, DB_Time:
		return DB_DateTime
	case DB_TimeStampz:
		return DB_Text
	case DB_Char, DB_Varchar, DB_NVarchar, DB_TinyText, DB_Text, DB_MediumText, DB_LongText:
		return DB_Text
	case DB_Bit, DB_TinyInt, DB_SmallInt, DB_MediumInt, DB_Int, DB_Integer, DB_BigInt, DB_Bool:
		return DB_Integer
	case DB_Float, DB_Double, DB_Real:
		return DB_Real
	case DB_Decimal, DB_Numeric:
		return DB_Numeric
	case DB_TinyBlob, DB_Blob, DB_MediumBlob, DB_LongBlob, DB_Bytea, DB_Binary, DB_VarBinary:
		return DB_Blob
	case DB_Serial, DB_BigSerial:
		c.IsPrimaryKey = true
		c.IsAutoIncrement = true
		c.Nullable = false
		return DB_Integer
	default:
		return c.Type
	}
}

func (db *Sqlite3) IndexCheckSql(tableName, indexName string) (string, []interface{}) {
	args := []interface{}{tableName, indexName}
	sql := "SELECT 1 FROM " + db.Quote("sqlite_master") + " WHERE " + db.Quote("type") + "='index' AND " + db.Quote("tbl_name") + "=? AND " + db.Quote("name") + "=?"
	return sql, args
}

func (db *Sqlite3) DropIndexSql(tableName string, index *Index) string {
	quote := db.Quote
	//var unique string
	idxName := index.XName(tableName)
	return fmt.Sprintf("DROP INDEX %v", quote(idxName))
}

func (db *Sqlite3) CleanDB() error {
	return nil
}

func (db *Sqlite3) isThisError(err error, errcode int) bool {
	if driverErr, ok := err.(sqlite3.Error); ok {
		if int(driverErr.ExtendedCode) == errcode {
			return true
		}
	}

	return false
}

func (db *Sqlite3) IsUniqueConstraintViolation(err error) bool {
	return db.isThisError(err, int(sqlite3.ErrConstraintUnique))
}

func (db *Sqlite3) IsDeadlock(err error) bool {
	return false // No deadlock
}

`
	ServiceSqlstoreMigratorTypes = `
package migrator

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
)

const (
	POSTGRES = "postgres"
	SQLITE   = "sqlite3"
	MYSQL    = "mysql"
	MSSQL    = "mssql"
)

type Migration interface {
	Sql(dialect Dialect) string
	Id() string
	SetId(string)
	GetCondition() MigrationCondition
}

type CodeMigration interface {
	Migration
	Exec(sess *xorm.Session, migrator *Migrator) error
}

type SQLType string

type ColumnType string

const (
	DB_TYPE_STRING ColumnType = "String"
)

type Table struct {
	Name        string
	Columns     []*Column
	PrimaryKeys []string
	Indices     []*Index
}

const (
	IndexType = iota + 1
	UniqueIndex
)

type Index struct {
	Name string
	Type int
	Cols []string
}

func (index *Index) XName(tableName string) string {
	if index.Name == "" {
		index.Name = strings.Join(index.Cols, "_")
	}

	if !strings.HasPrefix(index.Name, "UQE_") &&
		!strings.HasPrefix(index.Name, "IDX_") {
		if index.Type == UniqueIndex {
			return fmt.Sprintf("UQE_%v_%v", tableName, index.Name)
		}
		return fmt.Sprintf("IDX_%v_%v", tableName, index.Name)
	}
	return index.Name
}

var (
	DB_Bit       = "BIT"
	DB_TinyInt   = "TINYINT"
	DB_SmallInt  = "SMALLINT"
	DB_MediumInt = "MEDIUMINT"
	DB_Int       = "INT"
	DB_Integer   = "INTEGER"
	DB_BigInt    = "BIGINT"

	DB_Enum = "ENUM"
	DB_Set  = "SET"

	DB_Char       = "CHAR"
	DB_Varchar    = "VARCHAR"
	DB_NVarchar   = "NVARCHAR"
	DB_TinyText   = "TINYTEXT"
	DB_Text       = "TEXT"
	DB_MediumText = "MEDIUMTEXT"
	DB_LongText   = "LONGTEXT"
	DB_Uuid       = "UUID"

	DB_Date       = "DATE"
	DB_DateTime   = "DATETIME"
	DB_Time       = "TIME"
	DB_TimeStamp  = "TIMESTAMP"
	DB_TimeStampz = "TIMESTAMPZ"

	DB_Decimal = "DECIMAL"
	DB_Numeric = "NUMERIC"

	DB_Real   = "REAL"
	DB_Float  = "FLOAT"
	DB_Double = "DOUBLE"

	DB_Binary     = "BINARY"
	DB_VarBinary  = "VARBINARY"
	DB_TinyBlob   = "TINYBLOB"
	DB_Blob       = "BLOB"
	DB_MediumBlob = "MEDIUMBLOB"
	DB_LongBlob   = "LONGBLOB"
	DB_Bytea      = "BYTEA"

	DB_Bool = "BOOL"

	DB_Serial    = "SERIAL"
	DB_BigSerial = "BIGSERIAL"
)

`

)
