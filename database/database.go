package database

import (
	"gopkg.in/mgo.v2"
	"time"
)

const (
	host     = "localhost"
	database = "curvefever_api"
	username = "curvefever"
	password = "foobar"
)

var session *mgo.Session

func init() {
	var err error
	session, err = mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{host},
		Database: database,
		Username: username,
		Password: password,
		Timeout:  5 * time.Second,
	})

	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
}

//Returns an instance of database connection
//Automatically closes session after execution
func GetDatabase() *mgo.Database {
	sess := session.Copy()
	//defer sess.Close()
	return sess.DB(database)
}
