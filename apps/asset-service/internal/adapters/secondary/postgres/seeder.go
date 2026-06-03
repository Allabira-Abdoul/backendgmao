package postgres

import (
	"log"
	"time"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Seed populates initial asset data with airport equipment.
func Seed(db *gorm.DB) {
	// AutoMigrate tables
	db.AutoMigrate(
		&domain.EquipmentModel{},
		&domain.PartModel{},
		&domain.EquipmentInstance{},
		&domain.PartInstance{},
		&domain.MetricThreshold{},
		&domain.Supplier{},
		&domain.ModelSupplier{},
		&domain.PartConsumptionLog{},
		&domain.Measurement{},
	)

	log.Println("Seeding airport equipment data...")

	// 1. Create Equipment Models
	// GSE
	pushbackModel := createEquipmentModel(db, "Pushback Tractor (Tug)", "GSE", "Heavy vehicle for pushing aircraft")
	beltLoaderModel := createEquipmentModel(db, "Belt Loader", "GSE", "Mobile conveyor belt for luggage")
	gpuModel := createEquipmentModel(db, "Ground Power Unit (GPU)", "GSE", "Mobile generator for aircraft power")
	deiceModel := createEquipmentModel(db, "De-icing Vehicle", "GSE", "Truck equipped with boom and heated fluid")

	// Terminal
	jetwayModel := createEquipmentModel(db, "Passenger Boarding Bridge", "TERMINAL", "Enclosed movable tunnel for boarding")
	scannerModel := createEquipmentModel(db, "Security Scanner", "TERMINAL", "X-ray machine for carry-on items")
	carouselModel := createEquipmentModel(db, "Baggage Carousel", "TERMINAL", "Looping conveyor for arrivals")

	// Airfield
	ilsModel := createEquipmentModel(db, "Instrument Landing System (ILS)", "AIRFIELD", "Ground-based radio antennas for guidance")
	sweeperModel := createEquipmentModel(db, "Runway Sweeper", "AIRFIELD", "Vehicle to remove FOD from runways")

	// 2. Create Part Models
	// Pushback Parts
	tugEngine := createPartModel(db, "Diesel Engine 500HP", "ENGINE", 2)
	tugTire := createPartModel(db, "Heavy Duty Tire", "WHEEL", 20)
	tugPin := createPartModel(db, "Towing Pin", "MECHANICAL", 15)

	// GPU Parts
	gpuGenerator := createPartModel(db, "Alternator Generator", "ELECTRICAL", 3)
	gpuCable := createPartModel(db, "Aviation Power Cable", "ELECTRICAL", 10)

	// Jetway Parts
	jetwayCanopy := createPartModel(db, "Weather Canopy", "STRUCTURAL", 5)
	jetwayConsole := createPartModel(db, "Control Console", "ELECTRONIC", 2)

	// Sweeper Parts
	sweeperBrush := createPartModel(db, "Rotary Wire Brush", "MECHANICAL", 30)

	// 3. Create Equipment Instances and assign Part Instances
	now := time.Now()

	// Pushback Tractor PT1
	pt1 := createEquipmentInstance(db, "PT-001", pushbackModel.ID, "OPERATIONAL", "Gate 10 Apron", now.AddDate(-3, 0, 0), 120000.0)
	createPartInstance(db, pt1.ID, tugEngine.ID, "SN-PT1-ENG")
	createPartInstance(db, pt1.ID, tugTire.ID, "SN-PT1-TIR1")
	createPartInstance(db, pt1.ID, tugTire.ID, "SN-PT1-TIR2")
	createPartInstance(db, pt1.ID, tugTire.ID, "SN-PT1-TIR3")
	createPartInstance(db, pt1.ID, tugTire.ID, "SN-PT1-TIR4")
	createPartInstance(db, pt1.ID, tugPin.ID, "SN-PT1-PIN")

	// Pushback Tractor PT2 (In Maintenance)
	pt2 := createEquipmentInstance(db, "PT-002", pushbackModel.ID, "DOWN", "Maintenance Hangar", now.AddDate(-1, -6, 0), 125000.0)
	createPartInstance(db, pt2.ID, tugEngine.ID, "SN-PT2-ENG")
	createPartInstance(db, pt2.ID, tugPin.ID, "SN-PT2-PIN")

	// GPU 1
	gpu1 := createEquipmentInstance(db, "GPU-1A", gpuModel.ID, "OPERATIONAL", "Gate 12", now.AddDate(0, -2, 0), 45000.0)
	createPartInstance(db, gpu1.ID, gpuGenerator.ID, "SN-GPU1-GEN")
	createPartInstance(db, gpu1.ID, gpuCable.ID, "SN-GPU1-CAB")

	// Jetway 1
	jet1 := createEquipmentInstance(db, "GATE-10-BRIDGE", jetwayModel.ID, "OPERATIONAL", "Terminal 1 Gate 10", now.AddDate(-5, 0, 0), 850000.0)
	createPartInstance(db, jet1.ID, jetwayCanopy.ID, "SN-JET1-CAN")
	createPartInstance(db, jet1.ID, jetwayConsole.ID, "SN-JET1-CON")

	// Runway Sweeper
	sweep1 := createEquipmentInstance(db, "SWP-R1", sweeperModel.ID, "OPERATIONAL", "Airfield Garage", now.AddDate(-2, -3, 0), 210000.0)
	createPartInstance(db, sweep1.ID, sweeperBrush.ID, "SN-SWP1-BRU1")
	createPartInstance(db, sweep1.ID, sweeperBrush.ID, "SN-SWP1-BRU2")

	// Create instances for other models to avoid "declared and not used" errors
	createEquipmentInstance(db, "BLT-100", beltLoaderModel.ID, "OPERATIONAL", "Apron 2", now, 60000.0)
	createEquipmentInstance(db, "DEICE-1", deiceModel.ID, "IN_STOCK", "Winter Garage", now, 150000.0)
	createEquipmentInstance(db, "SCAN-X1", scannerModel.ID, "OPERATIONAL", "Terminal 1 Security", now, 120000.0)
	createEquipmentInstance(db, "CAR-ARR1", carouselModel.ID, "OPERATIONAL", "Arrivals Hall B", now, 250000.0)
	ils1 := createEquipmentInstance(db, "ILS-RWY09", ilsModel.ID, "OPERATIONAL", "Runway 09", now, 1500000.0)

	// 4. Create Suppliers & Model Suppliers
	sup1 := createSupplier(db, "Global Aviation Parts", "contact@gap.com")
	createModelSupplier(db, sup1.ID, &pushbackModel.ID, nil, "REF-GAP-PT", "DOC-123")
	createModelSupplier(db, sup1.ID, nil, &tugEngine.ID, "REF-GAP-ENG", "DOC-ENG-123")

	// 5. Create Metric Thresholds
	createMetricThreshold(db, &pushbackModel.ID, nil, nil, nil, "Engine Temperature", nil, ptr(110.0), "Celsius")
	createMetricThreshold(db, nil, nil, &gpu1.ID, nil, "Output Voltage", ptr(110.0), ptr(125.0), "Volts")

	// 6. Create Measurements
	createMeasurement(db, &pt1.ID, nil, "Engine Temperature", 85.5, "Celsius", now)
	createMeasurement(db, &gpu1.ID, nil, "Output Voltage", 118.2, "Volts", now)

	// 7. Create Part Consumption Logs
	createPartConsumptionLog(db, tugTire.ID, 2, "Replaced 2 tires on PT-001")

	_ = ils1

	log.Println("Seeding airport equipment data completed")
}

func ptr(f float64) *float64 {
	return &f
}

func createEquipmentModel(db *gorm.DB, name, category, desc string) domain.EquipmentModel {
	var m domain.EquipmentModel
	db.Where(domain.EquipmentModel{Name: name}).
		Assign(domain.EquipmentModel{Category: category, Description: desc}).
		FirstOrCreate(&m)
	return m
}

func createPartModel(db *gorm.DB, name, category string, qty int) domain.PartModel {
	var p domain.PartModel
	db.Where(domain.PartModel{Name: name}).
		Assign(domain.PartModel{Category: category, SpareQuantity: qty, IsSerialized: true}).
		FirstOrCreate(&p)
	return p
}

func createEquipmentInstance(db *gorm.DB, code string, modelID uuid.UUID, status, location string, date time.Time, value float64) domain.EquipmentInstance {
	var i domain.EquipmentInstance
	db.Where(domain.EquipmentInstance{Code: code}).
		Assign(domain.EquipmentInstance{EquipmentModelID: modelID, Status: status, Location: location, PurchaseDate: date, PurchaseValue: value}).
		FirstOrCreate(&i)
	return i
}

func createPartInstance(db *gorm.DB, eqInstID, partModelID uuid.UUID, sn string) domain.PartInstance {
	var p domain.PartInstance
	db.Where(domain.PartInstance{SerialNumber: sn}).
		Assign(domain.PartInstance{EquipmentInstanceID: &eqInstID, PartModelID: partModelID, Status: "OPERATIONAL", CurrentLocation: "Installed"}).
		FirstOrCreate(&p)
	return p
}

func createSupplier(db *gorm.DB, name, contactInfo string) domain.Supplier {
	var s domain.Supplier
	db.Where(domain.Supplier{Name: name}).
		Assign(domain.Supplier{ContactInfo: contactInfo}).
		FirstOrCreate(&s)
	return s
}

func createModelSupplier(db *gorm.DB, supplierID uuid.UUID, eqModelID, partModelID *uuid.UUID, refCode, docRef string) domain.ModelSupplier {
	var ms domain.ModelSupplier
	db.Where(domain.ModelSupplier{SupplierID: supplierID, EquipmentModelID: eqModelID, PartModelID: partModelID}).
		Assign(domain.ModelSupplier{SupplierReferenceCode: refCode, TechnicalDocReference: docRef}).
		FirstOrCreate(&ms)
	return ms
}

func createMetricThreshold(db *gorm.DB, eqModelID, partModelID, eqInstID, partInstID *uuid.UUID, metric string, min, max *float64, unit string) domain.MetricThreshold {
	var t domain.MetricThreshold
	db.Where(domain.MetricThreshold{
		EquipmentModelID:    eqModelID,
		PartModelID:         partModelID,
		EquipmentInstanceID: eqInstID,
		PartInstanceID:      partInstID,
		MetricName:          metric,
	}).Assign(domain.MetricThreshold{MinValue: min, MaxValue: max, Unit: unit}).
		FirstOrCreate(&t)
	return t
}

func createMeasurement(db *gorm.DB, eqInstID, partInstID *uuid.UUID, metric string, value float64, unit string, recordedAt time.Time) domain.Measurement {
	var m domain.Measurement
	recordedAt = recordedAt.Truncate(time.Second)
	db.Where(domain.Measurement{
		EquipmentInstanceID: eqInstID,
		PartInstanceID:      partInstID,
		MetricName:          metric,
		RecordedAt:          recordedAt,
	}).Assign(domain.Measurement{Value: value, Unit: unit}).
		FirstOrCreate(&m)
	return m
}

func createPartConsumptionLog(db *gorm.DB, partModelID uuid.UUID, qty int, notes string) domain.PartConsumptionLog {
	var l domain.PartConsumptionLog
	db.Where(domain.PartConsumptionLog{PartModelID: partModelID, Notes: notes}).
		Attrs(domain.PartConsumptionLog{ConsumedBy: uuid.New()}).
		Assign(domain.PartConsumptionLog{QuantityUsed: qty}).
		FirstOrCreate(&l)
	return l
}
