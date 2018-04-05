package utils

import "time"

// URLTravelAlerts - this has all alerts
var URLTravelAlerts = "http://www3.septa.org/hackathon/Alerts/"

// TrainRRSchedules - will use for analytics
type TrainRRSchedules struct {
	TrainID     string
	DocDate     string
	Timestamp   time.Time
	RRSchedules []RRSchedules
}

// RRSchedules - node of TrainRRSchedules
type RRSchedules struct {
	Station string `json:"station"`
	SchedTM string `json:"sched_tm"`
	EstTM   string `json:"est_tm"`
	ActTM   string `json:"act_tm"`
}

// LiveViewMessage - see https://play.golang.org/p/XxsmA8a7YPj
type LiveViewMessage struct {
	Lat, Lon, TrainNo,
	Service, Dest, NextStop, Line,
	Consist string
	Late int
	Source, Track,
	TrackChange string
}

type trainRun struct {
	Direction   string `json:"direction"`
	Path        string `json:"path"`
	TrainID     string `json:"train_id"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Line        string `json:"line"`
	Status      string `json:"status"`
	ServiceType string `json:"service_type"`
	NextStation string `json:"next_station"`
	SchedTime   string `json:"sched_time"`
	DepartTime  string `json:"depart_time"`
	Track       string `json:"track"`
	TrackChange string `json:"track_change"`
}

type northRun struct {
	Northbound []trainRun `json:"Northbound"`
}

type southRun struct {
	Southbound []trainRun `json:"Southbound"`
}

// StationType -- need comment
type StationType struct {
	StationID int
	Station   string
}
