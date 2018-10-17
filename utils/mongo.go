package utils

import (
	"fmt"
	"github.com/memory/setting"
	"gopkg.in/mgo.v2"
)

type MMMongo struct {
	session *mgo.Session
}

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func (this *MMMongo) _getSession() bool {

	//log.Println("_getSession")

	settings := setting.MMSettingsSington()

	session, err := mgo.Dial(fmt.Sprintf("%s:%s", settings.Mongo.Host, settings.Mongo.Port))

	if err != nil {
		//panic(err) //直接终止程序运行
		logger := settings.MMLogger
		logger.Printf("Mongo._getSession failed: %s ", err)

		return false
	} else {
		//最大连接池默认为4096

		this.session = session
		//this.session=session.Clone()
		return true
	}

}
func (this *MMMongo) CloseSession() {
	if this.session != nil {
		this.session.Close()
	}

}
func (this *MMMongo) getSession() bool {

	state := false
	if this.session == nil {
		//重新连接

		state = this._getSession()
	} else {

		err := this.session.Ping()

		if err != nil {
			//重新连接
			state = this._getSession()

		} else {
			//正常
			state = true
		}

	}

	//time.Sleep(time.Second*60)
	//
	//this.getSession()

	return state

}

//公共方法，获取collection对象
//func (this *MMMongo) witchCollection(collection string, query func(*mgo.Collection) error) error {
//
//	session := this._getSession()
//	defer session.Close()
//
//
//	return query(session.DB(setting.MMSettingsSington().Mongo.Db).C(collection))
//}

//func (this *MMMongo) Find(collection string) []map[string]string {
//
//	var results []map[string]string
//
//	query := func(c *mgo.Collection) error {
//		return c.Find(nil).All(&results)
//	}
//	err := this.witchCollection(collection, query)
//	if err != nil {
//		//
//
//	}
//
//	//result := Person{}
//	//err = c.Find(bson.M{"name": "Ale"}).One(&result)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	return results
//}

func (this *MMMongo) Insert(collection string, item interface{}) bool {

	//query := func(c *mgo.Collection) error {
	//	return c.Insert(
	//		&item,
	//	)
	//}
	//err := this.witchCollection(collection, query)

	if this.getSession() {

		err := this.session.DB(setting.MMSettingsSington().Mongo.Db).C(collection).Insert(&item)

		if err != nil {

			logger := setting.MMSettingsSington().MMLogger
			logger.Printf("Mongo.Insert failed: %s （collection=%s）（item=%s）", err, collection, item)

			return false
		} else {
			return true
		}

	} else {
		return false

	}

}
