package utils

var testGetParseMapData = []byte(
	`{"Elkins Park Departures: March 30, 2018, 2:01 pm":
[{"Northbound":[{"direction":"N","path":"R4N","train_id":"436",
"origin":"Airport Terminal E-F","destination":"Warminster",
"line":"Warminster","status":"2 min","service_type":"LOCAL",
"next_station":"30th Street Station","sched_time":"2018-03-30 14:31:00.000",
"depart_time":"2018-03-30 14:31:00.000","track":"2","track_change":null,
"platform":"","platform_change":null},{"direction":"N","path":"R4N",
"train_id":"438","origin":"Airport Terminal E-F","destination":"Glenside",
"line":"Airport","status":"On Time","service_type":"LOCAL",
"next_station":null,"sched_time":"2018-03-30 15:01:00.000",
"depart_time":"2018-03-30 15:01:00.000","track":"2",
"track_change":null,"platform":"","platform_change":null},
{"direction":"N","path":"R4N","train_id":"440","origin":
"Airport Terminal E-F","destination":"Warminster","line":"Airport",
"status":"On Time","service_type":"LOCAL","next_station":null,
"sched_time":"2018-03-30 15:31:00.000",
"depart_time":"2018-03-30 15:31:00.000","track":"2",
"track_change":null,"platform":"","platform_change":null}]},
{"Southbound":[{"direction":"S","path":"R4S","train_id":"443",
"origin":"Glenside","destination":"Airport","line":"Warminster",
"status":"On Time","service_type":"LOCAL","next_station":null,
"sched_time":"2018-03-30 14:28:00.000",
"depart_time":"2018-03-30 14:28:00.000","track":"1","track_change":null,
"platform":"","platform_change":null},{"direction":"S","path":"R4S",
"train_id":"445","origin":"Warminster","destination":"Airport",
"line":"Warminster","status":"On Time","service_type":"LOCAL",
"next_station":null,"sched_time":"2018-03-30 14:57:00.000",
"depart_time":"2018-03-30 14:57:00.000","track":"1","track_change":null,
"platform":"","platform_change":null},{"direction":"S","path":"R4S",
"train_id":"447","origin":"Glenside","destination":"Airport",
"line":"Warminster","status":"On Time","service_type":"LOCAL",
"next_station":null,"sched_time":"2018-03-30 15:28:00.000",
"depart_time":"2018-03-30 15:28:00.000","track":"1","track_change":null,
"platform":"","platform_change":null}]}]}`)