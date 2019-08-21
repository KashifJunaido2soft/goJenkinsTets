package main

import (
	"github.com/jamesruan/sodium"
	"gopkg.in/mgo.v2/bson"
)

// StructKeys ..
type StructKeys struct {
	Yubikey   string `json:"yubikey"`
	MobileKey string `json:"mobileKey"`
}

// User ..
type User struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Yubikey   string        `bson:"yubikey" json:"yubikey"`
	MobileKey string        `bson:"mobileKey" json:"mobileKey"`
}

// BoxKeysModel ..
type BoxKeysModel struct {
	PrivateKey sodium.BoxSecretKey `bson:"_PrivateKey"`
	PublicKey  sodium.BoxPublicKey `bson:"_PublicKey"`
}

// Configuration object
type Configuration struct {
	DbConnection         string
	DbName               string
	IP                   string
	Port                 string
	SecurityPolicyServer string
}
