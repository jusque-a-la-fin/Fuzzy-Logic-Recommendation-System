package gateway

import (
	"context"
	"encoding/json"
	"strconv"
	"vehicles/packages/domain/models"

	"github.com/lib/pq"
)

func (sr *selectionRepository) SelectCars(sl *models.Selection) (*[]models.Car, error) {

	query := `SELECT makes.make, models.model, generations.generation, steering_wheel_positions.position, power_steering_types.power_steering, body_types.body, specifications.length_meters, specifications.width_meters, specifications.height_meters, specifications.ground_clearance, specifications.drag_coefficient, specifications.front_track_width, specifications.rear_track_width, specifications.wheelbase, specifications.crash_test_estimate, specifications.year, engines.fuel_used, engines.engine_type, engines.capacity, engines.power, engines.max_torque, gearboxes.gearbox, drive_types.drive_type, suspensions.front_stabilizer, suspensions.back_stabilizer, suspensions.front_suspension, suspensions.back_suspension, tires.rear_tires_width, tires.front_tires_width, tires.front_tires_aspect_ratio, tires.rear_tires_aspect_ratio, tires.front_tires_rim_diameter, tires.rear_tires_rim_diameter, brakes.front_brakes, brakes.back_brakes, brakes.parking_brake, active_safety_and_motion_control_systems.abs_system, active_safety_and_motion_control_systems.esp_system, active_safety_and_motion_control_systems.ebd_system,
	active_safety_and_motion_control_systems.bas_system,
	active_safety_and_motion_control_systems.tcs_system, active_safety_and_motion_control_systems.front_parking_sensor, active_safety_and_motion_control_systems.back_parking_sensor, active_safety_and_motion_control_systems.rear_view_camera, active_safety_and_motion_control_systems.cruise_control, colors.color, lights.headlights, lights.light_sensor, lights.front_fog_lights, lights.back_fog_lights, interior_design.upholstery, cabin_microclimate.air_conditioner, cabin_microclimate.climate_control, electric_options.electric_front_side_windows_lifts, electric_options.electric_back_side_windows_lifts,
	electric_options.electric_heating_of_front_seats, electric_options.electric_heating_of_back_seats, electric_options.electric_heating_of_steering_wheel, electric_options.electric_heating_of_windshield, electric_options.electric_heating_of_rear_window, electric_options.electric_heating_of_mirrors, electric_options.electric_drive_of_driver_seat, electric_options.electric_drive_of_front_seats,
	electric_options.electric_drive_of_side_mirrors, electric_options.electric_trunk_opener, electric_options.rain_sensor, airbags.driver_airbag, airbags.front_passenger_airbag, airbags.side_airbags, airbags.curtain_airbags, multimedia_system.on_board_computer, multimedia_system.mp3_support, multimedia_system.hands_free_support, trim_levels.trim_level, trim_levels.acceleration_0_to_100_km_h, trim_levels.max_speed_kmh, trim_levels.city_fuel_consumption, trim_levels.highway_fuel_consumption, trim_levels.mixed_fuel_consumption, trim_levels.number_of_seats, trim_levels.trunk_volume_liters, trim_levels.weight_kg, trim_levels.car_alarm, offerings.price, offerings.mileage_km, offerings.photo_urls
		FROM makes
		INNER JOIN countries ON makes.country_id = countries.id
		INNER JOIN models ON makes.id = models.make_id
		INNER JOIN generations ON models.id = generations.model_id
		INNER JOIN specifications ON generations.id = specifications.generation_id
		INNER JOIN trim_levels ON specifications.id = trim_levels.specification_id
		INNER JOIN engines ON trim_levels.engine_id = engines.id
		INNER JOIN gearboxes ON trim_levels.gearbox_id = gearboxes.id
		INNER JOIN drive_types ON trim_levels.drive_type_id = drive_types.id
		INNER JOIN suspensions ON specifications.id = suspensions.id
		INNER JOIN tires ON specifications.id = tires.id
		INNER JOIN brakes ON specifications.id = brakes.id
		INNER JOIN active_safety_and_motion_control_systems ON trim_levels.active_safety_and_motion_control_systems_id = active_safety_and_motion_control_systems.id
		INNER JOIN colors ON trim_levels.color_id = colors.id
		INNER JOIN lights ON trim_levels.lights_id = lights.id
		INNER JOIN cabin_microclimate ON trim_levels.cabin_microclimate_id = cabin_microclimate.id
		INNER JOIN electric_options ON trim_levels.electric_options_id = electric_options.id
		INNER JOIN multimedia_system ON trim_levels.multimedia_system_id = multimedia_system.id
		INNER JOIN offerings ON trim_levels.id = offerings.trim_level_id
		LEFT JOIN steering_wheel_positions ON specifications.steering_wheel_position_id = steering_wheel_positions.id
		LEFT JOIN power_steering_types ON specifications.power_steering_type_id = power_steering_types.id
		LEFT JOIN body_types ON specifications.body_type_id = body_types.id
		LEFT JOIN interior_design ON trim_levels.interior_design_id = interior_design.id
		LEFT JOIN airbags ON trim_levels.airbags_id = airbags.id`

	//Добавляем условие WHERE по странам
	whereClause := ""
	args := make([]interface{}, 0)

	if sl.MinPrice != 0 || sl.MaxPrice != 0 {
		if whereClause == "" {
			whereClause += "WHERE "
		} else {
			whereClause += " AND "
		}

		if sl.MinPrice != 0 && sl.MaxPrice != 0 {
			whereClause += "offerings.price BETWEEN $1 AND $2"
			args = append(args, sl.MinPrice, sl.MaxPrice)
		} else if sl.MinPrice != 0 {
			whereClause += "offerings.price >= $1"
			args = append(args, sl.MinPrice)
		} else if sl.MaxPrice != 0 {
			whereClause += "offerings.price <= $1"
			args = append(args, sl.MaxPrice)
		}
	}

	for i, man := range sl.Manufacturers {
		if man == "Другие" {
			sl.Manufacturers[i] = "Чехия"
		}
	}

	if len(sl.Manufacturers) != 0 {
		if whereClause == "" {
			whereClause += "WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "countries.country IN ("
		for i, m := range sl.Manufacturers {
			args = append(args, m)
			whereClause += "$" + strconv.Itoa(len(args))
			if i < len(sl.Manufacturers)-1 {
				whereClause += ", "
			}
		}
		whereClause += ")"
	}

	if whereClause != "" {
		query += " " + whereClause
	}

	rows, err := sr.vehiclesDB.Query(query, args...)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	cars := []models.Car{}
	for rows.Next() {
		var car models.Car
		err := rows.Scan(&car.Make, &car.Model, &car.Generation, &car.Specification.SteeringWheelPosition, &car.Specification.PowerSteeringType,
			&car.Specification.Body, &car.Specification.LengthMeters, &car.Specification.WidthMeters, &car.Specification.HeightMeters,
			&car.Specification.GroundClearance, &car.Specification.DragCoefficient, &car.Specification.FrontTrackWidth,
			&car.Specification.RearTrackWidth, &car.Specification.Wheelbase, &car.Specification.CrashTestEstimate, &car.Specification.Year, &car.Engine.FuelUsed,
			&car.Engine.EngineType, &car.Engine.Capacity, &car.Engine.Power, &car.Engine.MaxTorque, &car.Gearbox.Gearbox,
			&car.Drive.DriveType, &car.Suspension.FrontStabilizer, &car.Suspension.BackStabilizer, &car.Suspension.FrontSuspension,
			&car.Suspension.BackSuspension, &car.Tires.BackTiresWidth, &car.Tires.FrontTiresWidth, &car.Tires.FrontTiresAspectRatio,
			&car.Tires.BackTiresAspectRatio, &car.Tires.FrontTiresRimDiameter, &car.Tires.BackTiresRimDiameter, &car.Brakes.FrontBrakes,
			&car.Brakes.BackBrakes, &car.Brakes.ParkingBrake, &car.ActiveSafetyAndMotionControlSystem.ABSPresent,
			&car.ActiveSafetyAndMotionControlSystem.ESPPresent, &car.ActiveSafetyAndMotionControlSystem.EBDPresent,
			&car.ActiveSafetyAndMotionControlSystem.BASPresent, &car.ActiveSafetyAndMotionControlSystem.TCSPresent,
			&car.ActiveSafetyAndMotionControlSystem.FrontParkingSensor, &car.ActiveSafetyAndMotionControlSystem.BackParkingSensor,
			&car.ActiveSafetyAndMotionControlSystem.RearViewCamera, &car.ActiveSafetyAndMotionControlSystem.CruiseControl, &car.TrimLevel.Color,
			&car.Light.Headlights, &car.Light.LightSensor, &car.Light.FrontFogLights, &car.Light.BackFogLights, &car.InteriorDesign.Upholstery,
			&car.CabinMicroclimate.AirConditioner, &car.CabinMicroclimate.ClimateControl, &car.ElectricOptions.ElectricFrontSideWindowsLifts,
			&car.ElectricOptions.ElectricBackSideWindowsLifts, &car.ElectricOptions.ElectricHeatingOfFrontSeats, &car.ElectricOptions.ElectricHeatingOfBackSeats,
			&car.ElectricOptions.ElectricHeatingOfSteeringWheel, &car.ElectricOptions.ElectricHeatingOfWindshield,
			&car.ElectricOptions.ElectricHeatingOfRearWindow, &car.ElectricOptions.ElectricHeatingOfMirrors, &car.ElectricOptions.ElectricDriveOfDriverSeat,
			&car.ElectricOptions.ElectricDriveOfFrontSeats, &car.ElectricOptions.ElectricDriveOfSideMirrors, &car.ElectricOptions.ElectricTrunkOpener,
			&car.ElectricOptions.RainSensor, &car.Airbags.DriverAirbag, &car.Airbags.FrontPassengerAirbag, &car.Airbags.SideAirbags, &car.Airbags.CurtainAirbags,
			&car.MultimediaSystem.OnBoardComputer, &car.MultimediaSystem.MP3Support, &car.MultimediaSystem.HandsFreeSupport, &car.TrimLevel.Level,
			&car.TrimLevel.Acceleration0To100kmh, &car.TrimLevel.MaxSpeedkmh, &car.TrimLevel.CityFuelConsumption, &car.TrimLevel.HighwayFuelConsumption,
			&car.TrimLevel.MixedFuelConsumption, &car.TrimLevel.NumberOfSeats, &car.TrimLevel.TrunkVolumeLiters, &car.TrimLevel.MassKg, &car.TrimLevel.CarAlarm,
			&car.Offering.Price, &car.Offering.Mileagekm,
			pq.Array(&car.Offering.PhotoURLs))

		cars = append(cars, car)

		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// выполняем запрос к бд vehicles
	return &cars, nil
}

func (sr *selectionRepository) LoadCarsData(cars []models.Car) error {

	carsJSON, err := json.Marshal(cars)
	if err != nil {
		return err
	}
	ctx := context.Background()
	if err = sr.rdb.Set(ctx, "cars", string(carsJSON), 0).Err(); err != nil {
		return err
	}

	return nil
}

func (sr *selectionRepository) GetCarsData() ([]models.Car, error) {

	ctx := context.Background()
	val, err := sr.rdb.Get(ctx, "cars").Result()
	if err != nil {
		return nil, err
	}

	// Преобразование JSON-строки в массив моделей Car
	cars := []models.Car{}
	err = json.Unmarshal([]byte(val), &cars)
	if err != nil {
		return nil, err
	}
	return cars, nil
}
