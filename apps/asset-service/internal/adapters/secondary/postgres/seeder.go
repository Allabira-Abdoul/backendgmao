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
	)

	var count int64
	db.Model(&domain.EquipmentModel{}).Count(&count)
	if count > 0 {
		return
	}

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
	createPartInstance(db, pt1.ID, tugEngine.ID)
	createPartInstance(db, pt1.ID, tugTire.ID)
	createPartInstance(db, pt1.ID, tugTire.ID)
	createPartInstance(db, pt1.ID, tugTire.ID)
	createPartInstance(db, pt1.ID, tugTire.ID)
	createPartInstance(db, pt1.ID, tugPin.ID)

	// Pushback Tractor PT2 (In Maintenance)
	pt2 := createEquipmentInstance(db, "PT-002", pushbackModel.ID, "DOWN", "Maintenance Hangar", now.AddDate(-1, -6, 0), 125000.0)
	createPartInstance(db, pt2.ID, tugEngine.ID)
	createPartInstance(db, pt2.ID, tugPin.ID)

	// GPU 1
	gpu1 := createEquipmentInstance(db, "GPU-1A", gpuModel.ID, "OPERATIONAL", "Gate 12", now.AddDate(0, -2, 0), 45000.0)
	createPartInstance(db, gpu1.ID, gpuGenerator.ID)
	createPartInstance(db, gpu1.ID, gpuCable.ID)

	// Jetway 1
	jet1 := createEquipmentInstance(db, "GATE-10-BRIDGE", jetwayModel.ID, "OPERATIONAL", "Terminal 1 Gate 10", now.AddDate(-5, 0, 0), 850000.0)
	createPartInstance(db, jet1.ID, jetwayCanopy.ID)
	createPartInstance(db, jet1.ID, jetwayConsole.ID)

	// Runway Sweeper
	sweep1 := createEquipmentInstance(db, "SWP-R1", sweeperModel.ID, "OPERATIONAL", "Airfield Garage", now.AddDate(-2, -3, 0), 210000.0)
	createPartInstance(db, sweep1.ID, sweeperBrush.ID)
	createPartInstance(db, sweep1.ID, sweeperBrush.ID)

	// Create instances for other models to avoid "declared and not used" errors
	createEquipmentInstance(db, "BLT-100", beltLoaderModel.ID, "OPERATIONAL", "Apron 2", now, 60000.0)
	createEquipmentInstance(db, "DEICE-1", deiceModel.ID, "IN_STOCK", "Winter Garage", now, 150000.0)
	createEquipmentInstance(db, "SCAN-X1", scannerModel.ID, "OPERATIONAL", "Terminal 1 Security", now, 120000.0)
	createEquipmentInstance(db, "CAR-ARR1", carouselModel.ID, "OPERATIONAL", "Arrivals Hall B", now, 250000.0)
	createEquipmentInstance(db, "ILS-RWY09", ilsModel.ID, "OPERATIONAL", "Runway 09", now, 1500000.0)

	log.Println("Seeding airport equipment data completed")
}

func createEquipmentModel(db *gorm.DB, name, category, desc string) domain.EquipmentModel {
	m := domain.EquipmentModel{
		ID:          uuid.New(),
		Name:        name,
		Category:    category,
		Description: desc,
	}
	db.Create(&m)
	return m
}

func createPartModel(db *gorm.DB, name, category string, qty int) domain.PartModel {
	p := domain.PartModel{
		ID:            uuid.New(),
		Name:          name,
		Category:      category,
		SpareQuantity: qty,
	}
	db.Create(&p)
	return p
}

func createEquipmentInstance(db *gorm.DB, code string, modelID uuid.UUID, status, location string, date time.Time, value float64) domain.EquipmentInstance {
	i := domain.EquipmentInstance{
		ID:               uuid.New(),
		Code:             code,
		EquipmentModelID: modelID,
		Status:           status,
		Location:         location,
		PurchaseDate:     date,
		PurchaseValue:    value,
	}
	db.Create(&i)
	return i
}

func createPartInstance(db *gorm.DB, eqInstID, partModelID uuid.UUID) domain.PartInstance {
	p := domain.PartInstance{
		ID:                  uuid.New(),
		EquipmentInstanceID: eqInstID,
		PartModelID:         partModelID,
		Status:              "OPERATIONAL",
	}
	db.Create(&p)
	return p
}
