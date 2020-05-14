package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	return m.sess.DB("metrics_db").C(collectionName).Insert(r)
}

func (m mongo) Count() (int, error) {
	return m.sess.DB("metrics_db").C(collectionName).Count()
}

type statsPerRoute struct {
	ID struct {
		Method string `bson:"method" json:"method"`
		URL    string `bson:"url" json:"url"`
	} `bson:"_id" json:"id"`
	NumberOfRequests int `bson:"numberOfRequests" json:"number_of_requests"`
}

func (m mongo) AverageResponseTime() (float64, error) {

	type res struct {
		AverageResponseTime float64 `bson:"averageResponseTime" json:"average_response_time"`
	}

	var ret = []res{}

	var baseMatch = bson.M{
		"$group": bson.M{
			"_id":                 nil,
			"averageResponseTime": bson.M{"$avg": "$requesttime"},
		},
	}

	err := m.sess.DB("metrics_db").C(collectionName).
		Pipe([]bson.M{baseMatch}).All(&ret)

	if len(ret) > 0 {
		return ret[0].AverageResponseTime, err
	}

	return 0, nil
}

func newMongo(addr string) (mongo, error) {
	sess, err := mgo.Dial(addr)
	if err != nil {
		return mongo{}, err
	}

	return mongo{
		sess: sess,
	}, nil
}
