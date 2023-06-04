package usecase

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"vehicles/packages/domain/models"
)

type carRecommendation struct {
	car                 *models.Car
	recommendationValue float64
}

func generateResultOfFuzzyAlgorithm(cars *[]models.Car, priorities []string) *[]models.Car {

	carRecs := make([]carRecommendation, 0, len(*cars))

	for i := 0; i < len(*cars); i++ {
		fmt.Println((*cars)[i].Model)
		carRecs = append(carRecs, carRecommendation{car: &(*cars)[i], recommendationValue: performFuzzyAlgorithm(&(*cars)[i], priorities)})
	}

	sort.Slice(carRecs, func(i, j int) bool {
		return carRecs[i].recommendationValue > carRecs[j].recommendationValue
	})

	var sortedCars []models.Car
	for i := 0; i < len(carRecs); i++ {
		sortedCars = append(sortedCars, *carRecs[i].car)
	}

	return &sortedCars
}

func performFuzzyAlgorithm(car *models.Car, priorities []string) float64 {

	seq := []string{"экономичность", "динамика", "управляемость", "комфорт", "безопасность"}

	var combs [][]string
	for i := 1; i <= len(seq); i++ {
		combinations := generateCombinations(seq)
		for _, c := range combinations {
			if len(c) == i {
				combs = append(combs, c)
			}
		}
	}

	rulesFilesMap := make(map[string]string)
	all := 0
	j := 1
	for _, comb := range combs {
		allPermutations := permute(comb)

		all += len(allPermutations)

		for _, permutation := range allPermutations {

			rulesFilesMap[fmt.Sprintf("%v", permutation)] = fmt.Sprintf("%d_rules.txt", j)
			j++
		}

	}

	dir := "/home/daniel/Desktop/THESIS/MAIN_THESIS/packages/usecase/usecase/rules"
	filename := filepath.Join(dir, rulesFilesMap[fmt.Sprintf("%v", priorities)])
	fmt.Println((filename))
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Считываем данные из файла
	var currentSet []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			// Добавляем набор правил в срез currentSet
			continue
		} else {
			// Добавляем правило в текущий набор правил
			currentSet = append(currentSet, line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	var rules [][]string
	// Выводим считанные данные
	var result []string
	for _, r := range currentSet {
		parts := strings.Split(r, " ")

		for i := 0; i < len(parts); i += 2 {
			if i+1 < len(parts) {
				result = append(result, parts[i]+" "+parts[i+1])
			} else {
				result = append(result, parts[i])
			}
		}
		rules = append(rules, result)
		result = []string{}
	}

	var valuesRecommendation []int
	for _, rule := range rules {
		value, _ := strconv.Atoi(rule[len(rule)-1])
		valuesRecommendation = append(valuesRecommendation, value)
	}

	fuel_consumption := car.TrimLevel.MixedFuelConsumption
	timeOfAcceleration0To100kmh := car.TrimLevel.Acceleration0To100kmh

	handling_coef := calculateHandlingCoefficient(car.Engine.Power, car.Specification.FrontTrackWidth, car.Specification.RearTrackWidth, car.Drive,
		car.Suspension, car.Tires, car.ActiveSafetyAndMotionControlSystem.ABSPresent, car.ActiveSafetyAndMotionControlSystem.ESPPresent,
		car.ActiveSafetyAndMotionControlSystem.EBDPresent, car.ActiveSafetyAndMotionControlSystem.BASPresent,
		car.ActiveSafetyAndMotionControlSystem.TCSPresent, car.Brakes.FrontBrakes, car.Brakes.BackBrakes,
		car.TrimLevel.MassKg, car.Specification.Wheelbase, car.Specification.LengthMeters, car.Specification.WidthMeters,
		car.Specification.HeightMeters, car.Specification.GroundClearance, car.Specification.DragCoefficient,
	)

	comfort_coef := calculateComfortCoefficient(car.Suspension, car.Gearbox, car.CabinMicroclimate, car.InteriorDesign, car.ElectricOptions, car.MultimediaSystem, car.Light, car.Specification.PowerSteeringType, car.TrimLevel.CarAlarm, car.TrimLevel.TrunkVolumeLiters)

	safety_coef := calculateSafetyCoefficient(car.Specification.CrashTestEstimate, car.ActiveSafetyAndMotionControlSystem, car.Airbags, car.Brakes)

	valuesOfMemebershipFunction := [][]float64{}
	var params [3]float64
	for _, rule := range rules {
		values := []float64{}

		for _, set := range rule {
			parts := strings.Split(set, " ")

			switch parts[0] {
			case "безопасность":
				params = provideFunctionParameters(parts[0], parts[1])
				safety := calculateMembershipFunctionValueForLinguisticVariables(params, parts[1], safety_coef)
				values = append(values, safety)

			case "управляемость":
				params = provideFunctionParameters(parts[0], parts[1])
				handling := calculateMembershipFunctionValueForLinguisticVariables(params, parts[1], handling_coef)
				values = append(values, handling)

			case "комфорт":
				params = provideFunctionParameters(parts[0], parts[1])
				comfort := calculateMembershipFunctionValueForLinguisticVariables(params, parts[1], comfort_coef)
				values = append(values, comfort)

			case "динамика":
				params = provideFunctionParameters(parts[0], parts[1])
				dynamics := calculateMembershipFunctionValueForLinguisticVariables(params, parts[1], timeOfAcceleration0To100kmh)
				values = append(values, dynamics)

			case "экономичность":
				params = provideFunctionParameters(parts[0], parts[1])
				economy := calculateMembershipFunctionValueForLinguisticVariables(params, parts[1], fuel_consumption)
				values = append(values, economy)

			}
		}
		valuesOfMemebershipFunction = append(valuesOfMemebershipFunction, values)
	}

	var minValuesOfMemebershipFunction []float64
	for _, val := range valuesOfMemebershipFunction {
		minValuesOfMemebershipFunction = append(minValuesOfMemebershipFunction, findMin(val))
	}

	var recommendationValue float64
	if len(minValuesOfMemebershipFunction) == len(valuesRecommendation) {
		recommendationValue = defuzzyficate(minValuesOfMemebershipFunction, valuesRecommendation)
	} else {
		fmt.Println("Ошибка, количество степеней рекомендации должно быть равно количеству правил.")
	}

	return recommendationValue
}

func permute(strings []string) [][]string {
	var result [][]string

	if len(strings) == 0 {
		return [][]string{{}}
	}

	for i := 0; i < len(strings); i++ {
		rest := make([]string, len(strings)-1)
		copy(rest[0:], strings[0:i])
		copy(rest[i:], strings[i+1:])

		for _, permutation := range permute(rest) {
			currentPerm := append([]string{strings[i]}, permutation...)
			result = append(result, currentPerm)
		}
	}

	return result
}

func generateCombinations(arr []string) [][]string {
	var result [][]string
	var gen func(int, []string)
	gen = func(n int, ss []string) {
		if n == len(arr) {
			return
		}
		for i := n; i < len(arr); i++ {
			result = append(result, append([]string{}, ss...))
			result[len(result)-1] = append(result[len(result)-1], arr[i])
			gen(i+1, result[len(result)-1])
		}
	}
	gen(0, []string{})
	return result
}

func calculateHandlingCoefficient(power, frontTrackWidth, backTrackWidth float64, dr models.DriveTypes, sps models.Suspensions,
	ts models.TiresTypes, absPresent, espPresent, ebdPresent, basPresent, tcsPresent, frontBrakes, backBrakes string,
	massKg, wheelbase, lengthMeters, widthMeters, heightMeters, groundClearance, dragCoefficient float64) float64 {

	var power_ float64 = power * 735.5 // перевод из лошадиных сил в ватты
	var driveTypeCoefficient float64
	switch dr.DriveType {
	case "Передний(FF)", "Передний":
		driveTypeCoefficient = 0.9
	case "Полный (4WD)", "Полный":
		driveTypeCoefficient = 1
	case "Задний(FR)", "Задний":
		driveTypeCoefficient = 0.7
	}

	var frontStabilizerCoefficient float64 = 1.0
	var backStabilizerCoefficient float64 = 1.0

	if sps.FrontStabilizer == "Есть" {
		frontStabilizerCoefficient = 1.2
	}

	if sps.BackStabilizer == "Есть" {
		backStabilizerCoefficient = 1.2
	}

	var frontSuspensionCoefficient float64 = 1.0
	var backSuspensionCoefficient float64 = 1.0
	switch sps.FrontSuspension {
	case "Многорычажная, независимая":
		frontSuspensionCoefficient = 1.8
	case "Независимая, на двойных поперечных рычагах":
		frontSuspensionCoefficient = 1.9
	case "Пневматическая":
		frontSuspensionCoefficient = 1.7
	case "Независимая, амортизационная стойка типа МакФерсон":
		frontSuspensionCoefficient = 1.6
	case "Полузависимая, торсионная балка":
		frontSuspensionCoefficient = 1.5
	case "Листовая, пружинная":
		frontSuspensionCoefficient = 1.4
	case "Зависимая, пружинная":
		frontSuspensionCoefficient = 1.3
	}

	switch sps.BackSuspension {
	case "Многорычажная, независимая":
		backSuspensionCoefficient = 1.9
	case "Независимая, на двойных поперечных рычагах":
		backSuspensionCoefficient = 1.8
	case "Пневматическая":
		backSuspensionCoefficient = 1.7
	case "Независимая, амортизационная стойка типа МакФерсон":
		backSuspensionCoefficient = 1.6
	case "Полузависимая, торсионная балка":
		backSuspensionCoefficient = 1.4
	case "Листовая, пружинная":
		backSuspensionCoefficient = 1.3
	case "Зависимая, пружинная":
		backSuspensionCoefficient = 1.2
	}

	frontTiresWidth_ := float64(ts.FrontTiresWidth) / 1000
	backTiresWidth_ := float64(ts.BackTiresWidth) / 1000

	backTiresAspectRatio := float64(ts.BackTiresAspectRatio)
	frontTiresAspectRatio := float64(ts.FrontTiresAspectRatio)

	frontTiresProfile := frontTiresWidth_ * (frontTiresAspectRatio / 100)
	backTiresProfile := backTiresWidth_ * (backTiresAspectRatio / 100)

	frontTiresDiameter_ := float64(ts.FrontTiresRimDiameter)*0.0254 + 2*frontTiresProfile // перевод в миллиметры из дюймов inches
	backTiresDiameter_ := float64(ts.BackTiresRimDiameter)*0.0254 + 2*backTiresProfile    // перевод в миллиметры из дюймов inches

	var absCoefficient float64 = 0
	var espCoefficient float64 = 0
	var ebdCoefficient float64 = 0
	var basCoefficient float64 = 0
	var tcsCoefficient float64 = 0
	if absPresent == "Есть" {

		absCoefficient = 0.064
	}
	if espPresent == "Есть" {

		espCoefficient = 0.07
	}
	if ebdPresent == "Есть" {

		ebdCoefficient = 0.056
	}
	if basPresent == "Есть" {

		basCoefficient = 0.059
	}
	if tcsPresent == "Есть" {

		tcsCoefficient = 0.051
	}

	var frontBrakesCoefficient float64
	var backBrakesCoefficient float64

	switch frontBrakes {
	case "Дисковые вентилируемые", "Дисковые":
		frontBrakesCoefficient = 0.7
	case "Барабанные":
		frontBrakesCoefficient = 0.5
	}

	switch backBrakes {
	case "Дисковые вентилируемые", "Дисковые":
		backBrakesCoefficient = 0.6
	case "Барабанные":
		backBrakesCoefficient = 0.4
	}

	// Информация о диаметре шин может быть использована для расчета общей высоты автомобиля, что также может влиять на его управляемость. Общая высота автомобиля определяется как сумма диаметра колес и высоты профиля шин.

	// Формула для расчета общей высоты автомобиля будет выглядеть следующим образом:

	// Общая высота автомобиля = диаметр передних колес + диаметр задних колес + (высота профиля передней шины * 2) + (высота профиля задней шины * 2)

	efficientFrontTrackWidth := frontTrackWidth + 0.5*(frontTiresWidth_-backTiresWidth_)/frontTiresAspectRatio
	efficientBackTrackWidth := backTrackWidth + 0.5*(backTiresWidth_-frontTiresWidth_)/backTiresAspectRatio
	handlingCoefficient := ((power_ * (efficientFrontTrackWidth + efficientBackTrackWidth) / 2) * driveTypeCoefficient *
		(frontSuspensionCoefficient*frontStabilizerCoefficient + backSuspensionCoefficient*backStabilizerCoefficient) *
		(frontTiresWidth_*frontTiresDiameter_ + backTiresWidth_*backTiresDiameter_) * (frontBrakesCoefficient + backBrakesCoefficient)) /
		(massKg * wheelbase * (lengthMeters + widthMeters + heightMeters) * groundClearance * dragCoefficient)

	sizeCoefficient := 50.0
	handlingCoefficient = handlingCoefficient + handlingCoefficient*(absCoefficient+espCoefficient+ebdCoefficient+basCoefficient+tcsCoefficient) - sizeCoefficient

	return handlingCoefficient
}

func calculateComfortCoefficient(sps models.Suspensions, gb models.Gearboxes, mc models.CabinMicroclimateTypes, idn models.InteriorDesigns, eo models.SetOfElectricOptions, ms models.MultimediaSystems, ls models.Lights, powerSteeringType, carAlarm string, trunkVolume float64) float64 {

	var frontSuspensionCoefficient float64
	switch sps.FrontSuspension {
	case "Многорычажная, независимая":
		if sps.FrontStabilizer == "Есть" {
			frontSuspensionCoefficient = 3.8
		} else {
			frontSuspensionCoefficient = 3
		}
	case "Независимая, на двойных поперечных рычагах":
		if sps.FrontStabilizer == "Есть" {
			frontSuspensionCoefficient = 3.6
		} else {
			frontSuspensionCoefficient = 2.8
		}
	case "Пневматическая":
		if sps.FrontStabilizer == "Есть" {
			frontSuspensionCoefficient = 4
		} else {
			frontSuspensionCoefficient = 3.2
		}
	case "Независимая, амортизационная стойка типа МакФерсон":
		if sps.FrontStabilizer == "Есть" {
			frontSuspensionCoefficient = 3.4
		} else {
			frontSuspensionCoefficient = 2.6
		}
	case "Полузависимая, торсионная балка":
		if sps.FrontStabilizer == "Есть" {
			frontSuspensionCoefficient = 2.8
		} else {
			frontSuspensionCoefficient = 2
		}
	case "Листовая, пружинная":
		if sps.FrontStabilizer == "Есть" {
			frontSuspensionCoefficient = 2.4
		} else {
			frontSuspensionCoefficient = 1.6
		}
	case "Зависимая, пружинная":
		if sps.FrontStabilizer == "Есть" {
			frontSuspensionCoefficient = 1.8
		} else {
			frontSuspensionCoefficient = 1
		}
	}

	var backSuspensionCoefficient float64
	switch sps.BackSuspension {
	case "Многорычажная, независимая":
		if sps.BackStabilizer == "Есть" {
			backSuspensionCoefficient = 3.8
		} else {
			backSuspensionCoefficient = 3
		}
	case "Независимая, на двойных поперечных рычагах":
		if sps.BackStabilizer == "Есть" {
			backSuspensionCoefficient = 3.6
		} else {
			backSuspensionCoefficient = 2.8
		}
	case "Пневматическая":
		if sps.BackStabilizer == "Есть" {
			backSuspensionCoefficient = 4
		} else {
			backSuspensionCoefficient = 3.2
		}
	case "Независимая, амортизационная стойка типа МакФерсон":
		if sps.BackStabilizer == "Есть" {
			backSuspensionCoefficient = 3.4
		} else {
			backSuspensionCoefficient = 2.6
		}
	case "Полузависимая, торсионная балка":
		if sps.BackStabilizer == "Есть" {
			backSuspensionCoefficient = 2.8
		} else {
			backSuspensionCoefficient = 2
		}
	case "Листовая, пружинная":
		if sps.BackStabilizer == "Есть" {
			backSuspensionCoefficient = 2.4
		} else {
			backSuspensionCoefficient = 1.6
		}
	case "Зависимая, пружинная":
		if sps.BackStabilizer == "Есть" {
			backSuspensionCoefficient = 1.8
		} else {
			backSuspensionCoefficient = 1
		}
	}

	var powerSteeringTypeCoefficient float64
	switch powerSteeringType {
	case "Электроусилитель руля", "Гидроусилитель руля", "Электрогидроусилитель руля":
		powerSteeringTypeCoefficient = 2
	}

	var gearboxCoefficient float64
	switch gb.Gearbox {
	case "АКПП 6", "АКПП 5", "Вариатор":
		gearboxCoefficient = 4
	}

	var climateCoefficient float64
	if mc.AirConditioner == "Есть" && mc.ClimateControl == "Есть" {
		climateCoefficient = 3
	} else if mc.AirConditioner == "Есть" && mc.ClimateControl != "Есть" {
		climateCoefficient = 2
	} else if mc.AirConditioner != "Есть" && mc.ClimateControl == "Есть" {
		climateCoefficient = 2
	}

	var interiorCoefficient float64
	if idn.Upholstery == "Кожаная" {
		interiorCoefficient = 0.2962962962962963
	}

	var lightsCoefficient float64
	if ls.Headlights != "Галогенные фары" {
		lightsCoefficient += 0.8888888888888888
	}
	if ls.FrontFogLights == "Есть" {
		lightsCoefficient += 0.2962962962962963
	}
	if ls.BackFogLights == "Есть" {
		lightsCoefficient += 0.2962962962962963
	}
	if ls.LightSensor == "Есть" {
		lightsCoefficient += 0.2962962962962963
	}

	var electricOptionsCoefficient float64
	if eo.ElectricFrontSideWindowsLifts == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.ElectricBackSideWindowsLifts == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.ElectricHeatingOfFrontSeats == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.ElectricHeatingOfBackSeats == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.ElectricHeatingOfSteeringWheel == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.ElectricHeatingOfWindshield == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.ElectricHeatingOfRearWindow == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.ElectricHeatingOfMirrors == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.ElectricDriveOfDriverSeat == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.ElectricDriveOfFrontSeats == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.ElectricDriveOfSideMirrors == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.ElectricTrunkOpener == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if eo.RainSensor == "Есть" {
		electricOptionsCoefficient += 0.2962962962962963
	}

	var trunkVolumeCoefficient float64
	if trunkVolume > 500 {
		trunkVolumeCoefficient = 0.8888888888888888
	}

	var carAlarmCoefficient float64
	if carAlarm == "Есть" {
		carAlarmCoefficient = 0.5925925925925926
	}

	var multimediaCoefficient float64
	if ms.OnBoardComputer == "Есть" {
		multimediaCoefficient += 0.2962962962962963
	}
	if ms.MP3Support == "Есть" {
		multimediaCoefficient += 0.2962962962962963
	}
	if ms.HandsFreeSupport == "Есть" {
		multimediaCoefficient += 0.2962962962962963
	}

	size_coef := 0.8
	comfortCoefficient := (frontSuspensionCoefficient + backSuspensionCoefficient + powerSteeringTypeCoefficient +
		gearboxCoefficient + climateCoefficient + interiorCoefficient + lightsCoefficient + electricOptionsCoefficient +
		trunkVolumeCoefficient + carAlarmCoefficient + multimediaCoefficient) * size_coef
	return comfortCoefficient
}

func calculateSafetyCoefficient(crashTestEstimate float64, sys models.ActiveSafetyAndMotionControlSystems, as models.SetOfAirbags, bs models.BrakesTypes) float64 {

	var controlSystemCoefficient float64
	if sys.ABSPresent == "Есть" {
		controlSystemCoefficient += 3
	}
	if sys.ESPPresent == "Есть" {
		controlSystemCoefficient += 1
	}
	if sys.EBDPresent == "Есть" {
		controlSystemCoefficient += 1
	}
	if sys.BASPresent == "Есть" {
		controlSystemCoefficient += 1
	}
	if sys.TCSPresent == "Есть" {
		controlSystemCoefficient += 1
	}

	var airbagsCoefficient float64
	if as.DriverAirbag == "Есть" {
		airbagsCoefficient += 1
	}
	if as.FrontPassengerAirbag == "Есть" {
		airbagsCoefficient += 1
	}
	if as.SideAirbags == "Есть" {
		airbagsCoefficient += 1
	}
	if as.CurtainAirbags == "Есть" {
		airbagsCoefficient += 1
	}

	var frontBrakesCoefficient float64
	switch bs.FrontBrakes {
	case "Дисковые вентилируемые", "Дисковые":
		frontBrakesCoefficient = 2
	}

	var backBrakesCoefficient float64
	switch bs.BackBrakes {
	case "Дисковые вентилируемые", "Дисковые":
		backBrakesCoefficient = 2
	}

	safetyCoefficient := crashTestEstimate + controlSystemCoefficient + airbagsCoefficient + frontBrakesCoefficient + backBrakesCoefficient
	return safetyCoefficient
}

func sigmoid_function(x float64, L float64, k float64, x0 float64) float64 {
	return L / (1 + math.Exp(-k*(x-x0)))
}

func gaussian_function(x, amp, cen, wid float64) float64 {
	return amp / (math.Sqrt(2*math.Pi) * wid) * math.Exp(-math.Pow(x-cen, 2)/(2*math.Pow(wid, 2)))
}

func calculateMembershipFunctionValueForLinguisticVariables(params [3]float64, term string, value float64) float64 {

	var valueOfMembershipFunction float64
	switch term {
	case "низкий":
		valueOfMembershipFunction = sigmoid_function(value, params[0], params[1], params[2])
	case "средний":
		valueOfMembershipFunction = gaussian_function(value, params[0], params[1], params[2])
	case "высокий":
		valueOfMembershipFunction = sigmoid_function(value, params[0], params[1], params[2])
	}

	if valueOfMembershipFunction < 0 {
		valueOfMembershipFunction = 0
	} else if valueOfMembershipFunction > 1 {
		valueOfMembershipFunction = 1
	}
	return valueOfMembershipFunction
}

type Parameters struct {
	params_low     [3]float64
	params_average [3]float64
	params_high    [3]float64
}

func provideFunctionParameters(variable, term string) [3]float64 {
	data := make(map[string]Parameters)
	data["экономичность"] = Parameters{[3]float64{1.043723139993038, 0.5194913435480255, 11.165188013054621},
		[3]float64{2.2900397063026374, 9.43414665796981, 2.4138470099365112},
		[3]float64{1.949834151590793, -0.3804532441502327, 5.188639378787266}}
	//
	data["динамика"] = Parameters{[3]float64{1.0231319819933777, 0.5016231903133455, 13.44547910618538},
		[3]float64{4.59765854168931, 10.810654375352698, 3.529577571232097},
		[3]float64{1.1836613715914706, -0.3792245799359442, 7.418367289135995}}

	data["управляемость"] = Parameters{[3]float64{1.2027418825694678, -0.07501884950616336, 22.4088071782321},
		[3]float64{31.124751770295614, 44.305695848946904, 19.042293165507264},
		[3]float64{1.3245201627428753, 0.07214432798281176, 69.8908258450921}}

	data["комфорт"] = Parameters{[3]float64{1.1834768835495801, -0.2870468773928149, 5.7724209993240825},
		[3]float64{5.0422852289029185, 10.10928688292464, 4.210219040980836},
		[3]float64{1.2492396969207602, 0.27009941927484593, 14.702080359730674}}

	data["безопасность"] = Parameters{[3]float64{1.3490219429573107, -0.21934076866027877, 4.73473258614666},
		[3]float64{5.450821590257078, 10.048764235757659, 4.185552288427339},
		[3]float64{1.2799644032509998, 0.3119397892443973, 15.65152657451626}}

	switch variable {
	case "экономичность":
		if term == "низкий" {
			return data["экономичность"].params_low
		} else if term == "средний" {
			return data["экономичность"].params_average
		} else if term == "высокий" {
			return data["экономичность"].params_high
		}
	case "динамика":
		if term == "низкий" {
			return data["динамика"].params_low
		} else if term == "средний" {
			return data["динамика"].params_average
		} else if term == "высокий" {
			return data["динамика"].params_high
		}
	case "управляемость":
		if term == "низкий" {
			return data["управляемость"].params_low
		} else if term == "средний" {
			return data["управляемость"].params_average
		} else if term == "высокий" {
			return data["управляемость"].params_high
		}
	case "комфорт":
		if term == "низкий" {
			return data["комфорт"].params_low
		} else if term == "средний" {
			return data["комфорт"].params_average
		} else if term == "высокий" {
			return data["комфорт"].params_high
		}
	case "безопасность":
		if term == "низкий" {
			return data["безопасность"].params_low
		} else if term == "средний" {
			return data["безопасность"].params_average
		} else if term == "высокий" {
			return data["безопасность"].params_high
		}
	}
	return [3]float64{}
}

func findMin(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}

	min := numbers[0]
	for _, num := range numbers {
		if num < min {
			min = num
		}
	}
	return min
}

func defuzzyficate(minValuesOfMemebershipFunction []float64, valuesRecommendation []int) float64 {
	var centersOfMassOfTheAreaUnderTheGraph []float64
	var areasOfTriangles []float64
	for i := 0; i < len(minValuesOfMemebershipFunction); i++ {
		centersOfMassOfTheAreaUnderTheGraph = append(centersOfMassOfTheAreaUnderTheGraph, calculateCenterOfMassOfTheAreaUnderTheGraph(minValuesOfMemebershipFunction[i], valuesRecommendation[i]))
		areasOfTriangles = append(areasOfTriangles, calculateAreaOfTriangle(minValuesOfMemebershipFunction[i]))
	}

	var enumerator float64
	for i := range areasOfTriangles {
		enumerator += areasOfTriangles[i] * centersOfMassOfTheAreaUnderTheGraph[i]
	}

	var denominator float64
	for i := 0; i < len(areasOfTriangles); i++ {
		denominator += areasOfTriangles[i]
	}

	result := enumerator / denominator

	return result
}
func calculateCenterOfMassOfTheAreaUnderTheGraph(minValueOfMemebershipFunction float64, valueReommendation int) float64 {
	x_values := setX(float64(valueReommendation), 1)
	area := trapezoidalRule(func(x float64) float64 {
		return setY(x, x_values[0], x_values[1], x_values[2], minValueOfMemebershipFunction)
	}, x_values[0], x_values[2], 10000)

	numerator := trapezoidalRule(func(x float64) float64 {
		return float64(valueReommendation) * setY(x, x_values[0], x_values[1], x_values[2], minValueOfMemebershipFunction)
	}, x_values[0], x_values[2], 10000)
	return numerator / area
}

func calculateAreaOfTriangle(minValueOfMemebershipFunction float64) float64 {

	var baseOfTriangle float64 = 2
	area := 0.5 * baseOfTriangle * minValueOfMemebershipFunction
	return area
}

func setY(x, leftBound, center, rightBound float64, yMax ...float64) float64 {
	var yMaxValue float64 = 1
	if len(yMax) > 0 {
		yMaxValue = yMax[0]
	}

	if x <= leftBound || x >= rightBound {
		return 0
	}
	if x <= center && x >= leftBound {
		return x - leftBound - (1 - yMaxValue)
	}
	if leftBound <= x && x <= rightBound {
		return rightBound - x - (1 - yMaxValue)
	}

	return 0
}

func setX(center, deviation float64) []float64 {
	x := []float64{}
	x = append(x, center-deviation)
	x = append(x, center)
	x = append(x, center+deviation)

	return x
}

func trapezoidalRule(f func(float64) float64, a, b float64, n int) float64 {
	h := (b - a) / float64(n)
	s := 0.5 * (f(a) + f(b))
	for i := 1; i < n; i++ {
		s += f(a + float64(i)*h)
	}
	return h * s
}
