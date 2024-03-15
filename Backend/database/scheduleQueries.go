package database

import (
	"busproject/model"
	"database/sql"
	"errors"
	"log"
	"strconv"
)

func InsertSchedule(db *sql.DB, schedule model.Schedule) error {
	sqlStatement := `INSERT INTO transport.schedule (id, bus_id, route_id, departure_time) VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(sqlStatement, schedule.Id, schedule.BusId, schedule.RouteId, schedule.DepartureTime)
	if err != nil {
		return err
	}

	// log.Println("Schedule inserted successfully")
	return nil
}

func GetAllSchedule(db *sql.DB) ([]model.Schedule, error) {
	sqlStatement := `SELECT * FROM transport.schedule;`

	res, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var schedules []model.Schedule
	for res.Next() {
		var shcedule model.Schedule
		err := res.Scan(&shcedule.Id, &shcedule.BusId, &shcedule.RouteId, &shcedule.DepartureTime)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		schedules = append(schedules, shcedule)
	}

	return schedules, nil
}

func DeleteSchedule(db *sql.DB, id string) error {
	sqlStatement := `DELETE FROM transport.schedule WHERE id=$1`

	newId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = db.Exec(sqlStatement, newId)
	if err != nil {
		return err
	}

	// log.Println("Schedule Deleted successfully")
	return nil
}

func GetUpcomingBus(db *sql.DB, source, destination int) ([]model.UpcomingBus, error) {
	// fmt.Println(source, destination)
	// if any of source or destination is 0 means they are not provided by client
	if source == 0 {
		return nil, errors.New("source need to be specified")
	}
	var sqlStatement string
	var result *sql.Rows
	var err error
	if destination == 0 {
		sqlStatement = `SELECT f.bus_id,route_id,route_name,source,destination,departure_time,b.lat,b.long,b.last_station_order FROM (SELECT bus_id,route_id,route_name,source,destination,departure_time,station_order FROM bustransportsystem WHERE station_id = $1 and status = 1) AS f LEFT JOIN transport.busstatus as b ON f.bus_id = b.bus_id WHERE (b.status = 0 AND departure_time >= current_time) OR (b.status = 1 AND b.last_station_order <= f.station_order) ORDER BY departure_time ASC;`
		result, err = db.Query(sqlStatement, source)
	} else {
		sqlStatement = `SELECT o.bus_id,route_id,route_name,source,destination,departure_time,b.lat,b.long,b.last_station_order FROM (SELECT s.bus_id,s.route_id,s.route_name,s.source,s.destination,s.departure_time FROM (SELECT * FROM bustransportsystem WHERE station_id = $1 and status = 1) as f INNER JOIN (SELECT * FROM bustransportsystem WHERE station_id = $2 and status = 1) as s ON f.route_id = s.route_id WHERE f.bus_id = s.bus_id AND f.station_order < s.station_order) as o LEFT JOIN transport.busstatus as b ON o.bus_id = b.bus_id WHERE b.status = 1 OR (b.status = 0 AND departure_time >= CURRENT_TIME) ORDER BY departure_time ASC;`
		result, err = db.Query(sqlStatement, source, destination)
	}
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("sorry no data available currently")
	}
	defer result.Close()

	var busOutput []model.UpcomingBus
	var dummy model.UpcomingBus
	for result.Next() {
		result.Scan(&dummy.Bus_id,&dummy.Route_id,&dummy.Name, &dummy.Source, &dummy.Destination, &dummy.DepartureTime, &dummy.Lat, &dummy.Long, &dummy.LastStationOrder)
		busOutput = append(busOutput, dummy)
	}
	if len(busOutput) == 0 {
		return nil, errors.New("sorry no bus available")
	}
	return busOutput, nil
}