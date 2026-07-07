package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parseData := strings.Split(data, ",")

	if len(parseData) != 3 {
		return 0, ``, 0, errors.New("invalid data")
	}

	steps, err := strconv.Atoi(parseData[0])

	if err != nil {
		return 0, ``, 0, err
	}

	if steps <= 0 {
		return 0, ``, 0, errors.New("invalid steps")
	}

	duration, err := time.ParseDuration(parseData[2])

	if err != nil {
		return 0, ``, 0, err
	}

	if duration <= 0 {
		return 0, ``, 0, errors.New("invalid duration")
	}

	activity := parseData[1]

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	length := float64(height) * stepLengthCoefficient

	totalDistance := float64(steps) * length / float64(mInKm)

	return totalDistance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	totalDistance := distance(steps, height)

	return totalDistance / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, trainingType, duration, err := parseTraining(data)

	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64
	switch trainingType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return ``, errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return ``, err
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainingType, duration.Hours(), dist, speed, calories), nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0.0, errors.New("invalid data")
	}
	meanUserSpeed := meanSpeed(steps, height, duration)

	return (weight * meanUserSpeed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	totalSpentCalories, err := RunningSpentCalories(steps, weight, height, duration)

	if err != nil {
		return 0.0, err
	}

	return totalSpentCalories * walkingCaloriesCoefficient, nil

}
