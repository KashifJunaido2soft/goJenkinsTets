package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//CollectionUsers ...
const CollectionUsers = "Users"

func writeToUsersCollection(data StructKeys, session mgo.Session) error {
	s := session.Clone()
	c := s.DB(AppConfig.DbName).C(CollectionUsers)
	err := c.Insert(&data)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func getUserByKeys(data StructKeys, session mgo.Session) (User, error) {
	s := session.Clone()
	defer s.Close()
	s.SetSafe(&mgo.Safe{})
	c := s.DB(AppConfig.DbName).C(CollectionUsers)
	var res User
	//err := c.Find(bson.M{"yubikey": data.Yubikey}).One(&res)
	err := c.Find(bson.M{"yubikey": data.Yubikey}).Select(bson.M{"mobilekey": data.MobileKey}).One(&res)
	return res, err
}

func findYubikey(key string, session mgo.Session) (User, error) {
	s := session.Clone()
	defer s.Close()
	s.SetSafe(&mgo.Safe{})
	c := s.DB(AppConfig.DbName).C(CollectionUsers)
	var res User
	err := c.Find(bson.M{"yubikey": key}).One(&res)
	return res, err
}

func findMobilkey(key string, session mgo.Session) (User, error) {
	s := session.Clone()
	defer s.Close()
	s.SetSafe(&mgo.Safe{})
	c := s.DB(AppConfig.DbName).C(CollectionUsers)
	var res User
	err := c.Find(bson.M{"mobilekey": key}).One(&res)
	return res, err
}

func getAllUsers(session mgo.Session) ([]User, error) {
	s := session.Clone()
	defer s.Close()
	s.SetSafe(&mgo.Safe{})
	c := s.DB(AppConfig.DbName).C(CollectionUsers)
	var res []User
	err := c.Find(bson.M{}).All(&res)
	return res, err
}
