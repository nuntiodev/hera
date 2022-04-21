package measurement_repository

type ActiveMeasurement struct {
}

type MeasurementRepository interface {
	RecordActiveMeasurement()
}