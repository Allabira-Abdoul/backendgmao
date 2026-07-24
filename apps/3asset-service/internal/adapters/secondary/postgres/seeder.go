package postgres

import (
	"fmt"
	"log"

	"backend-gmao/apps/asset-service/internal/core/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type airport struct {
	ID        int
	IATA      string
	Name      string
	Location  string
	Multiplier float64
}

var airports = []airport{
	{1, "DLA", "Douala International Airport", "Douala, Littoral", 1.2},
	{2, "NSI", "Yaoundé-Nsimalen International Airport", "Yaoundé, Centre", 1.0},
	{3, "GOU", "Garoua International Airport", "Garoua, North", 0.5},
	{4, "MVR", "Maroua-Salak International Airport", "Maroua, Far North", 0.4},
	{5, "NGE", "Ngaoundéré Airport", "Ngaoundéré, Adamawa", 0.3},
	{6, "BPC", "Bamenda Airport", "Bamenda, North-West", 0.1},
	{7, "BTA", "Bertoua Airport", "Bertoua, East", 0.1},
}

type equipmentTemplate struct {
	Prefix      string
	Category    string
	Name        string
	IsMotorized bool
}

var templates = []equipmentTemplate{
	{"XH", "F", "ELEVATEUR A FOURCHE 7 T", true},
	{"XH", "F", "ELEVATEUR A FOURCHE 2.5 T", true},
	{"XH", "F", "ELEVATEUR A FOURCHE 10 T", true},
	{"XE", "A", "TRACTEUR AGRICOLE", true},
	{"XE", "G", "TRICYCLE", true},
	{"XE", "D", "NACELLE", true},
	{"XH", "D", "TAPIS A BAGAGES", true},
	{"XH", "R", "BUS", true},
	{"XH", "U", "AMBULIFT", true},
	{"XH", "I", "CHARIOTS A BAGAGES", false},
	{"XH", "K", "CHARIOTS PORTE PALETTES", false},
	{"XH", "Y", "CHARIOTS PORTE CONTENEURS", false},
	{"XH", "B", "BARRES DE TRACTAGE", false},
	{"XH", "Q", "TRANSPALETTES", false},
	{"XH", "W", "CHARIOTS A BASCULEMENT", false},
	{"XH", "N", "ESCARBOT TECHNIQUE", false},
	{"XH", "J", "CHARIOT A CARBURANT", false},
	{"XH", "O", "BANCS DE STOCKAGE PALETTE", false},
	{"XH", "O", "BANCS DE STOCKAGE CONTENAIRE", false},
	{"XH", "T", "TRACTEUR TD 1800", true},
	{"XH", "T", "TRACTEUR K32", true},
	{"XH", "T", "TRACTEUR K40", true},
	{"XH", "P", "PUSHBACK TMX 400", true},
	{"XH", "P", "PUSHBACK TPX 200", true},
	{"XH", "P", "PUSHBACK TMX150", true},
	{"XH", "P", "PUSHBACK TD 18000", true},
	{"XH", "L", "LOADER PEB7", true},
	{"XH", "L", "LOADER PEB 14", true},
	{"XH", "L", "LOADER LAM7", true},
	{"XH", "G", "GROUPE GPU", true},
	{"XH", "C", "GROUPE ASU", true},
	{"XH", "V", "TONNE A VIDE TOILETTE", true},
	{"XH", "V", "CAMION VIDE TOILETTE", true},
	{"XH", "S", "TONNE A EAU POTABLE", true},
	{"XH", "H", "CAMION HOTELIER", true},
	{"XE", "E", "CAMION PUBELLE", true},
	{"XE", "E", "CAMION BALAYEUR", true},
	{"XE", "Y", "TRANSP MATER DEGOM", true},
	{"XH", "E", "PASSERELLE AUTO-TRACTEE", true},
	{"XH", "E", "PASSERELLE TRACTEE", false},
}

// SeedData seeds the database with initial data for the 5-level hierarchy.
func SeedData(db *gorm.DB) {
	var count int64
	db.Model(&domain.Site{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded. Skipping.")
		return
	}

	log.Println("Seeding Asset database started...")

	// Common Inventory Items
	items := []domain.InventoryItem{
		{ID: uuid.New(), ItemType: "SPARE_PART", PartNumber: "FIL-H-001", Name: "Filtre Hydraulique Standard", Category: "Hydraulique", StockQuantity: 50, ReorderPoint: 10, UnitOfMeasure: "unit"},
		{ID: uuid.New(), ItemType: "SPARE_PART", PartNumber: "BATT-12V-100", Name: "Batterie 12V 100Ah", Category: "Electrique", StockQuantity: 20, ReorderPoint: 5, UnitOfMeasure: "unit"},
		{ID: uuid.New(), ItemType: "SPARE_PART", PartNumber: "ROUE-P-400", Name: "Roue Pleine 400mm", Category: "Train Roulant", StockQuantity: 30, ReorderPoint: 8, UnitOfMeasure: "unit"},
		{ID: uuid.New(), ItemType: "CONSUMABLE", PartNumber: "HUILE-M-15W40", Name: "Huile Moteur 15W40", Category: "Lubrifiant", StockQuantity: 200, ReorderPoint: 50, UnitOfMeasure: "litre"},
	}
	db.Create(&items)

	for _, apt := range airports {
		site := domain.Site{
			ID:          uuid.New(),
			Name:        fmt.Sprintf("Aéroport %s", apt.IATA),
			Location:    apt.Location,
			Description: apt.Name,
		}
		db.Create(&site)

		sysMot := domain.System{
			ID:          uuid.New(),
			SiteID:      site.ID,
			Name:        "GSE Motorisés",
			Description: "Équipements d'assistance en escale motorisés",
		}
		sysNonMot := domain.System{
			ID:          uuid.New(),
			SiteID:      site.ID,
			Name:        "GSE Non-Motorisés",
			Description: "Équipements tractés et chariots",
		}
		db.Create(&sysMot)
		db.Create(&sysNonMot)

		for tplIdx, tpl := range templates {
			// Calculate how many to create based on multiplier
			// Base count: at least 1 for NSI/DLA if multiplier >= 1, less for others
			qty := int(float64(2) * apt.Multiplier)
			if qty < 1 && tplIdx%int(1/apt.Multiplier+1) == 0 {
				qty = 1 // give some equipment to smaller airports
			}

			for i := 1; i <= qty; i++ {
				code := fmt.Sprintf("%s%d%s%03d", tpl.Prefix, apt.ID, tpl.Category, i)
				
				var sysID uuid.UUID
				if tpl.IsMotorized {
					sysID = sysMot.ID
				} else {
					sysID = sysNonMot.ID
				}

				asset := domain.Asset{
					ID:            uuid.New(),
					SystemID:      sysID,
					Name:          fmt.Sprintf("%s - %s", tpl.Name, code),
					Code:          code,
					Manufacturer:  "Generic GSE Manufacturer",
					Status:        "OPERATIONAL",
					RulPercentage: 100.0,
				}
				db.Create(&asset)

				// Subsystems & Components
				if tpl.IsMotorized {
					subMoteur := domain.Subsystem{ID: uuid.New(), AssetID: asset.ID, Name: "Moteur", Criticality: "HIGH"}
					subHydra := domain.Subsystem{ID: uuid.New(), AssetID: asset.ID, Name: "Hydraulique", Criticality: "MEDIUM"}
					db.Create(&subMoteur)
					db.Create(&subHydra)

					compFiltre := domain.Component{
						ID:              uuid.New(),
						SubsystemID:     subHydra.ID,
						InventoryItemID: items[0].ID,
						Name:            "Filtre Principal",
						SerialNumber:    fmt.Sprintf("SN-%s-F01", code),
					}
					db.Create(&compFiltre)
				} else {
					subRoues := domain.Subsystem{ID: uuid.New(), AssetID: asset.ID, Name: "Train Roulant", Criticality: "HIGH"}
					db.Create(&subRoues)

					compRoue := domain.Component{
						ID:              uuid.New(),
						SubsystemID:     subRoues.ID,
						InventoryItemID: items[2].ID,
						Name:            "Roue Avant Gauche",
						SerialNumber:    fmt.Sprintf("SN-%s-R01", code),
					}
					db.Create(&compRoue)
				}
			}
		}
	}

	log.Println("Seeding Asset database completed successfully.")
}
