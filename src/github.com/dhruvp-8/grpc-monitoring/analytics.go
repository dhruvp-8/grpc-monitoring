package main

import (
	"gopkg.in/mgo.v2"
)

const (
	collectionName = "request_analytics"
)

type requestAnalytics struct {
	URL         string
	Method      string
	RequestTime int64
	Day         string
	Hour        int
}

type mongo struct {
	sess *mgo.Session
}

func (m mongo) Close() error {
	m.sess.Close()
	return nil
}

func (m mongo) Write(r requestAnalytics) error {
	return m.sess.DB("metrics_tb").C(collectionName).Insert(r)
}

func (m mongo) Count() (int, error) {
	return m.sess.DB("metrics_tb").C(collectionName).Count()
}
