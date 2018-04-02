package utils

import "time"

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

// StationList -- this is not a complete list. Use for testing.
func StationList() []StationType {
	s := []StationType{

		{90004, "30th Street Station"},
		{90404, "Airport Terminal A"},
		{90403, "Airport Terminal B"},
		{90402, "Airport Terminal C-D"},
		{90401, "Airport Terminal E-F"},
		{90208, "Allegheny"},
		{90804, "Allen Lane"},
		{90526, "Ambler"},
		{90313, "Angora"},
		{90309, "Clifton-Aldan"},
		{90533, "Colmar"},
		{90225, "Conshohocken"},
		{90706, "Cornwells Heights"},
		{90414, "Crestmont"},
		{90704, "Croydon"},
		{90409, "Elkins Park"},
		{90504, "Exton"},
		{90407, "Fern Rock TC"},
		{90312, "Fernwood"},
		{90214, "Folcroft"},
		{90005, "Suburban Station"},
		{90305, "Swarthmore"},
	}

	return s
}
