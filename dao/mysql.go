package dao

import (
	"database/sql"
	"fmt"
	"time"
	"webapi_erc20/common/config"
	"webapi_erc20/common/logs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SqlSession *gorm.DB
var dbsql *sql.DB

const (
	dbDriver = "mysql"
	dbURLFmt = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
)

func InitMySql(cfg *config.SugaredConfig) (err error) {

	dbURL := fmt.Sprintf(dbURLFmt, cfg.Mysql.User, cfg.Mysql.Password, cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.Database)
	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	SqlSession = db

	dbsql, err = db.DB()
	if err != nil {
		panic(err)
	}

	//连接池
	dbsql.SetMaxIdleConns(50)
	dbsql.SetMaxOpenConns(300)
	dbsql.SetConnMaxLifetime(300 * time.Second)

	return
}

func Close() {
	dbsql.Close()
}

func GetDB() *sql.DB {
	return dbsql
}

func ConnectUrl(cfg *config.SugaredConfig) string {

	dbURL := fmt.Sprintf(dbURLFmt, cfg.Mysql.User, cfg.Mysql.Password, cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.Database)
	logs.Debugf("dbURL=%s", dbURL)
	return dbURL
}

func CreateTable() {

	//創建表格
	GetAddressInstance().CreateTable()
	GetTokenInstance().CreateTable()
	GetBlockHeightInstance().CreateTable()
	GetTransactionInstance().CreateTable()
	GetWithdrawInstance().CreateTable()
}
