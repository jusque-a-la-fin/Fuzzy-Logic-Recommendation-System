package models

type Car struct {
	Make                               string `db:"make"`
	Model                              string `db:"model"`
	Generation                         string `db:"generation"`
	Suspension                         Suspensions
	Specification                      Specifications
	Engine                             Engines
	Gearbox                            Gearboxes
	Drive                              DriveTypes
	Tires                              TiresTypes
	Brakes                             BrakesTypes
	ActiveSafetyAndMotionControlSystem ActiveSafetyAndMotionControlSystems
	Light                              Lights
	InteriorDesign                     InteriorDesigns
	CabinMicroclimate                  CabinMicroclimateTypes
	ElectricOptions                    SetOfElectricOptions
	Airbags                            SetOfAirbags
	MultimediaSystem                   MultimediaSystems
	TrimLevel                          TrimLevels
	Offering                           Offerings
}

type Suspensions struct {
	FrontStabilizer string `db:"front_stabilizer"`
	BackStabilizer  string `db:"back_stabilizer"`
	FrontSuspension string `db:"front_suspension"`
	BackSuspension  string `db:"back_suspension"`
}

type Specifications struct {
	SteeringWheelPosition string  `db:"position"`
	PowerSteeringType     string  `db:"power_steering"`
	Body                  string  `db:"body"`
	LengthMeters          float64 `db:"length_meters"`
	WidthMeters           float64 `db:"width_meters"`
	HeightMeters          float64 `db:"height_meters"`
	GroundClearance       float64 `db:"ground_clearance"`
	DragCoefficient       float64 `db:"drag_coefficient"`
	FrontTrackWidth       float64 `db:"front_track_width"`
	RearTrackWidth        float64 `db:"rear_track_width"`
	Wheelbase             float64 `db:"wheelbase"`
	CrashTestEstimate     float64 `db:"crash_test_estimate"`
	Year                  int     `db:"year"`
}

type Engines struct {
	FuelUsed   string  `db:"fuel_used"`
	EngineType string  `db:"engine_type"`
	Capacity   float64 `db:"capacity"`
	Power      float64 `db:"power"`
	MaxTorque  string  `db:"max_torque"`
}

type Gearboxes struct {
	Gearbox string `db:"gearbox"`
}

type DriveTypes struct {
	DriveType string `db:"drive_type"`
}

type TiresTypes struct {
	BackTiresWidth        int `db:"rear_tires_width"`
	FrontTiresWidth       int `db:"front_tires_width"`
	FrontTiresAspectRatio int `db:"front_tires_aspect_ratio"`
	BackTiresAspectRatio  int `db:"rear_tires_aspect_ratio"`
	FrontTiresRimDiameter int `db:"front_tires_rim_diameter"`
	BackTiresRimDiameter  int `db:"rear_tires_rim_diameter"`
}

type BrakesTypes struct {
	FrontBrakes  string `db:"front_brakes"`
	BackBrakes   string `db:"back_brakes"`
	ParkingBrake string `db:"parking_brake"`
}

type ActiveSafetyAndMotionControlSystems struct {
	ABSPresent         string `db:"abs_system"`
	ESPPresent         string `db:"esp_system"`
	EBDPresent         string `db:"ebd_system"`
	BASPresent         string `db:"bas_system"`
	TCSPresent         string `db:"tcs_system"`
	FrontParkingSensor string `db:"front_parking_sensor"`
	BackParkingSensor  string `db:"back_parking_sensor"`
	RearViewCamera     string `db:"rear_view_camera"`
	CruiseControl      string `db:"cruise_control"`
}

type Lights struct {
	Headlights     string `db:"headlights"`
	LightSensor    string `db:"light_sensor"`
	FrontFogLights string `db:"front_fog_lights"`
	BackFogLights  string `db:"back_fog_lights"`
}

type InteriorDesigns struct {
	Upholstery string `db:"upholstery"`
}

type CabinMicroclimateTypes struct {
	AirConditioner string `db:"air_conditioner"`
	ClimateControl string `db:"climate_control"`
}

type SetOfElectricOptions struct {
	ElectricFrontSideWindowsLifts  string `db:"electric_front_side_windows_lifts"`
	ElectricBackSideWindowsLifts   string `db:"electric_back_side_windows_lifts"`
	ElectricHeatingOfFrontSeats    string `db:"electric_heating_of_front_seats"`
	ElectricHeatingOfBackSeats     string `db:"electric_heating_of_back_seats"`
	ElectricHeatingOfSteeringWheel string `db:"electric_heating_of_steering_wheel"`
	ElectricHeatingOfWindshield    string `db:"electric_heating_of_windshield"`
	ElectricHeatingOfRearWindow    string `db:"electric_heating_of_rear_window"`
	ElectricHeatingOfMirrors       string `db:"electric_heating_of_mirrors"`
	ElectricDriveOfDriverSeat      string `db:"electric_drive_of_driver_seat"`
	ElectricDriveOfFrontSeats      string `db:"electric_drive_of_front_seats"`
	ElectricDriveOfSideMirrors     string `db:"electric_drive_of_side_mirrors"`
	ElectricTrunkOpener            string `db:"electric_trunk_opener"`
	RainSensor                     string `db:"rain_sensor"`
}

type SetOfAirbags struct {
	DriverAirbag         string `db:"driver_airbag"`
	FrontPassengerAirbag string `db:"front_passenger_airbag"`
	SideAirbags          string `db:"side_airbags"`
	CurtainAirbags       string `db:"curtain_airbags"`
}

type MultimediaSystems struct {
	OnBoardComputer  string `db:"on_board_computer"`
	MP3Support       string `db:"mp3_support"`
	HandsFreeSupport string `db:"hands_free_support"`
}

type TrimLevels struct {
	Level                  string  `db:"trim_level"`
	Acceleration0To100kmh  float64 `db:"acceleration_0_to_100_km_h"`
	MaxSpeedkmh            float64 `db:"max_speed_kmh"`
	CityFuelConsumption    float64 `db:"city_fuel_consumption"`
	HighwayFuelConsumption float64 `db:"highway_fuel_consumption"`
	MixedFuelConsumption   float64 `db:"mixed_fuel_consumption"`
	NumberOfSeats          int     `db:"number_of_seats"`
	TrunkVolumeLiters      float64 `db:"trunk_volume_liters"`
	MassKg                 float64 `db:"weight_kg"`
	CarAlarm               string  `db:"car_alarm"`
	Color                  string  `db:"color"`
}

type Offerings struct {
	Price     float64  `db:"price"`
	Mileagekm int      `db:"mileage_km"`
	PhotoURLs []string `db:"photo_urls"`
}
