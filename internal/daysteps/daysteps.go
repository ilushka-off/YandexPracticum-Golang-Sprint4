package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	parseData := strings.Split(data, ",")

	if len(parseData) != 2 {
		return 0, 0, errors.New("invalid data")
	}

	stepsString := parseData[0]
	steps, err := strconv.Atoi(stepsString)

	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, errors.New("invalid data")
	}

	duration, err := time.ParseDuration(parseData[1])

	if err != nil {
		return 0, 0, errors.New("invalid duration")
	}

	if duration <= 0 {
		return 0, 0, errors.New("invalid duration")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)

	if err != nil {
		log.Println(err)
		return ``
	}

	distanceM := float64(steps) * stepLength
	distanceKm := distanceM / float64(mInKm)

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil {
		log.Println(err)
		return ``
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)
}
