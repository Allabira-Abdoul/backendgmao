package postgres

import (
	"log"
	"fmt"

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
		&domain.Supplier{},
		&domain.ModelSupplier{},
		&domain.EquipmentModelPartRequirement{},
		&domain.PartConsumptionLog{},
		&domain.EquipmentModelMaintenanceRule{},
		&domain.EquipmentInstanceMaintenanceState{},
	)

	log.Println("Seeding GSE equipment data...")
	
	// Create a generic supplier for now
	sup1 := createSupplier(db, "Aviation GSE Supplier", "contact@gse.com")

	// ============================================
	// ELEVATEURS ET TRACTEURS AGRICOLES
	// ============================================
	catTracteurs := "ELEVATEURS ET TRACTEURS AGRICOLES"
	modFourche7 := createEquipmentModel(db, "ELEVATEUR A FOURCHE 7 T", catTracteurs, "GSE")
	createEquipmentInstance(db, "XH2F001", modFourche7.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modFourche25 := createEquipmentModel(db, "ELEVATEUR A FOURCHE 2.5 T", catTracteurs, "GSE")
	createEquipmentInstance(db, "XH2F002", modFourche25.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modFourche10 := createEquipmentModel(db, "ELEVATEUR A FOURCHE 10 T", catTracteurs, "GSE")
	createEquipmentInstance(db, "XH1F123", modFourche10.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modTracteurAgricole := createEquipmentModel(db, "TRACTEUR AGRICOLE", catTracteurs, "GSE")
	
	// Add flexible maintenance rules for TRACTEUR AGRICOLE
	h250 := 250.0
	m3 := 3
	createEquipmentModelMaintenanceRule(db, modTracteurAgricole.ID, "Minor Maintenance", &h250, &m3)
	
	h2000 := 2000.0
	m12 := 12
	createEquipmentModelMaintenanceRule(db, modTracteurAgricole.ID, "Major Maintenance", &h2000, &m12)

	createEquipmentInstance(db, "XE2A001", modTracteurAgricole.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XE2A002", modTracteurAgricole.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XE2A003", modTracteurAgricole.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modTricycle := createEquipmentModel(db, "TRICYCLE", catTracteurs, "GSE")
	createEquipmentInstance(db, "XE2G001", modTricycle.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XE2G002", modTricycle.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modNacelle := createEquipmentModel(db, "NACELLE", catTracteurs, "GSE")
	createEquipmentInstance(db, "XE2D001", modNacelle.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	// ============================================
	// TAPIS A BAGAGES
	// ============================================
	catTapis := "TAPIS A BAGAGES"
	modTapis := createEquipmentModel(db, "TAPIS A BAGAGES", catTapis, "GSE")
	createEquipmentInstance(db, "XH2D001", modTapis.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2D002", modTapis.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2D003", modTapis.ID, &sup1.ID, "OPERATIONAL", "PARK GSE") // Assuming XH2D00 -> 003
	createEquipmentInstance(db, "XH2D004", modTapis.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	// ============================================
	// AMBULIFT ET BUS
	// ============================================
	catAmbulift := "AMBULIFT ET BUS"
	modBus := createEquipmentModel(db, "BUS", catAmbulift, "GSE")
	createEquipmentInstance(db, "XH2R001", modBus.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modAmbulift := createEquipmentModel(db, "AMBULIFT", catAmbulift, "GSE")
	createEquipmentInstance(db, "XH2U001", modAmbulift.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2U002", modAmbulift.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	// ============================================
	// ENGINS NON MOTORISES
	// ============================================
	catNonMoto := "ENGINS NON MOTORISES"
	modChariotBagages := createEquipmentModel(db, "CHARIOTS A BAGAGES", catNonMoto, "GSE")
	for i := 1; i <= 20; i++ {
		createEquipmentInstance(db, fmt.Sprintf("XH2I%03d", i), modChariotBagages.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	}

	modChariotPortePalettes := createEquipmentModel(db, "CHARIOTS PORTE PALETTES", catNonMoto, "GSE")
	for i := 1; i <= 20; i++ {
		createEquipmentInstance(db, fmt.Sprintf("XH2K%03d", i), modChariotPortePalettes.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	}

	modChariotPorteConteneurs := createEquipmentModel(db, "CHARIOTS PORTE CONTENEURS", catNonMoto, "GSE")
	for i := 1; i <= 10; i++ {
		createEquipmentInstance(db, fmt.Sprintf("XH2Y%03d", i), modChariotPorteConteneurs.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	}

	modBarres := createEquipmentModel(db, "BARRES DE TRACTAGE", catNonMoto, "GSE")
	for i := 1; i <= 5; i++ {
		createEquipmentInstance(db, fmt.Sprintf("XH2B%03d", i), modBarres.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	}

	modTranspalettes := createEquipmentModel(db, "TRANSPALETTES", catNonMoto, "GSE")
	for i := 1; i <= 7; i++ {
		createEquipmentInstance(db, fmt.Sprintf("XH2Q%02d", i), modTranspalettes.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	}

	modChariotBasculement := createEquipmentModel(db, "CHARIOTS A BASCULEMENT", catNonMoto, "GSE")
	createEquipmentInstance(db, "XH2W001", modChariotBasculement.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2W002", modChariotBasculement.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modEscarbot := createEquipmentModel(db, "ESCARBOT TECHNIQUE", catNonMoto, "GSE")
	for i := 1; i <= 5; i++ {
		createEquipmentInstance(db, fmt.Sprintf("XH2N%03d", i), modEscarbot.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	}

	modChariotCarburant := createEquipmentModel(db, "CHARIOT A CARBURANT", catNonMoto, "GSE")
	createEquipmentInstance(db, "XH2J001", modChariotCarburant.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modBancsPalette := createEquipmentModel(db, "BANCS DE STOCKAGE PALETTE", catNonMoto, "GSE")
	for i := 1; i <= 10; i++ {
		createEquipmentInstance(db, fmt.Sprintf("XH2O%04d", i), modBancsPalette.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	}

	modBancsContenaire := createEquipmentModel(db, "BANCS DE STOCKAGE CONTENAIRE", catNonMoto, "GSE")
	for i := 1; i <= 5; i++ {
		createEquipmentInstance(db, fmt.Sprintf("XH2O10%d", i), modBancsContenaire.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	}

	// ============================================
	// GROUPES (G4 / G5)
	// ============================================
	createEquipmentModel(db, "GROUPE G4", "GROUPES", "GSE")
	createEquipmentModel(db, "GROUPE G5", "GROUPES", "GSE")

	// ============================================
	// TRACTEURS DE MANUTENTION
	// ============================================
	catTracteursManutention := "TRACTEURS DE MANUTENTION"
	modTracteurTD1800 := createEquipmentModel(db, "TRACTEUR TD 1800", catTracteursManutention, "GSE")
	createEquipmentInstance(db, "XH2T001", modTracteurTD1800.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modTracteurK32 := createEquipmentModel(db, "TRACTEUR K32", catTracteursManutention, "GSE")
	createEquipmentInstance(db, "XH2T003", modTracteurK32.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2T004", modTracteurK32.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2T005", modTracteurK32.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2T006", modTracteurK32.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2T010", modTracteurK32.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modTracteurK40 := createEquipmentModel(db, "TRACTEUR K40", catTracteursManutention, "GSE")
	createEquipmentInstance(db, "XH2T007", modTracteurK40.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2T008", modTracteurK40.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2T009", modTracteurK40.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	// ============================================
	// PUSH-BACKS
	// ============================================
	catPushBacks := "PUSH-BACKS"
	modPushTMX400 := createEquipmentModel(db, "PUSHBACK TMX 400", catPushBacks, "GSE")
	createEquipmentInstance(db, "XH2P002", modPushTMX400.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modPushTPX200 := createEquipmentModel(db, "PUSHBACK TPX 200", catPushBacks, "GSE")
	createEquipmentInstance(db, "XH2P003", modPushTPX200.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modPushTMX150 := createEquipmentModel(db, "PUSHBACK TMX150", catPushBacks, "GSE")
	createEquipmentInstance(db, "XH2P004", modPushTMX150.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modPushTD18000 := createEquipmentModel(db, "PUSHBACK TD 18000", catPushBacks, "GSE")
	createEquipmentInstance(db, "XH1P001", modPushTD18000.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	// ============================================
	// LOADERS
	// ============================================
	catLoaders := "LOADERS"
	modLoaderPEB7 := createEquipmentModel(db, "LOADER PEB 7", catLoaders, "GSE")
	createEquipmentInstance(db, "XH2L001", modLoaderPEB7.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2L003", modLoaderPEB7.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modLoaderPEB14 := createEquipmentModel(db, "LOADER PEB 14", catLoaders, "GSE")
	createEquipmentInstance(db, "XH2L004", modLoaderPEB14.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modLoaderLAM7 := createEquipmentModel(db, "LOADER LAM7", catLoaders, "GSE")
	createEquipmentInstance(db, "XH2L005", modLoaderLAM7.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	// ============================================
	// GROUPES
	// ============================================
	catGroupes := "GROUPES"
	modGroupeGPU := createEquipmentModel(db, "GROUPE GPU", catGroupes, "GSE")
	createEquipmentInstance(db, "XH1G001", modGroupeGPU.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2G002", modGroupeGPU.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2G003", modGroupeGPU.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modGroupeASU := createEquipmentModel(db, "GROUPE ASU", catGroupes, "GSE")
	createEquipmentInstance(db, "XH2C001", modGroupeASU.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2C002", modGroupeASU.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2C003", modGroupeASU.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	// ============================================
	// CAMIONS ET TONNES
	// ============================================
	catCamions := "CAMIONS ET TONNES"
	modTonneVide := createEquipmentModel(db, "TONNE A VIDE TOILETTE", catCamions, "GSE")
	createEquipmentInstance(db, "XH2V001", modTonneVide.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modCamionVide := createEquipmentModel(db, "CAMION VIDE TOILETTE", catCamions, "GSE")
	createEquipmentInstance(db, "XH2V002", modCamionVide.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modTonneEau := createEquipmentModel(db, "TONNE A EAU POTABLE", catCamions, "GSE")
	createEquipmentInstance(db, "XH2S001", modTonneEau.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modCamionHotelier := createEquipmentModel(db, "CAMION HOTELIER", catCamions, "GSE")
	createEquipmentInstance(db, "XH2H001", modCamionHotelier.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2H002", modCamionHotelier.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modCamionPubelle := createEquipmentModel(db, "CAMION PUBELLE", catCamions, "GSE")
	createEquipmentInstance(db, "XE2E001", modCamionPubelle.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modCamionBalayeur := createEquipmentModel(db, "CAMION BALAYEUR", catCamions, "GSE")
	createEquipmentInstance(db, "XE2E002", modCamionBalayeur.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modTranspMater := createEquipmentModel(db, "TRANSP MATER DEGOM", catCamions, "GSE")
	createEquipmentInstance(db, "XE2Y001", modTranspMater.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	// ============================================
	// PASSERELLES
	// ============================================
	catPasserelles := "PASSERELLES"
	modPasserelleAuto := createEquipmentModel(db, "PASSERELLE AUTO-TRACTEE", catPasserelles, "GSE")
	createEquipmentInstance(db, "XH2E002", modPasserelleAuto.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2E003", modPasserelleAuto.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	modPasserelleTractee := createEquipmentModel(db, "PASSERELLE TRACTEE", catPasserelles, "GSE")
	createEquipmentInstance(db, "XH2E004", modPasserelleTractee.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")
	createEquipmentInstance(db, "XH2E005", modPasserelleTractee.ID, &sup1.ID, "OPERATIONAL", "PARK GSE")

	log.Println("Seeding GSE airport equipment data completed")
}

func ptr(f float64) *float64 {
	return &f
}

func createSupplier(db *gorm.DB, name, contact string) domain.Supplier {
	var s domain.Supplier
	db.Where(domain.Supplier{Name: name}).
		Assign(domain.Supplier{ContactInfo: contact}).
		FirstOrCreate(&s)
	return s
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

func createEquipmentInstance(db *gorm.DB, code string, modelID uuid.UUID, supplierID *uuid.UUID, status, location string) domain.EquipmentInstance {
	var i domain.EquipmentInstance
	db.Where(domain.EquipmentInstance{Code: code}).
		Assign(domain.EquipmentInstance{EquipmentModelID: modelID, SupplierID: supplierID, Status: status, Location: location}).
		FirstOrCreate(&i)
	return i
}

func createPartInstance(db *gorm.DB, eqInstID, partModelID uuid.UUID, supplierID *uuid.UUID, sn string) domain.PartInstance {
	var p domain.PartInstance
	db.Where(domain.PartInstance{SerialNumber: sn}).
		Assign(domain.PartInstance{EquipmentInstanceID: &eqInstID, PartModelID: partModelID, SupplierID: supplierID, Status: "INSTALLED"}).
		FirstOrCreate(&p)
	return p
}

func createModelSupplier(db *gorm.DB, supplierID uuid.UUID, eqModelID, partModelID *uuid.UUID, ref, doc string) domain.ModelSupplier {
	var ms domain.ModelSupplier
	db.Where(domain.ModelSupplier{SupplierID: supplierID, EquipmentModelID: eqModelID, PartModelID: partModelID}).
		Assign(domain.ModelSupplier{SupplierReferenceCode: ref, TechnicalDocReference: doc}).
		FirstOrCreate(&ms)
	return ms
}

func createEquipmentModelPartRequirement(db *gorm.DB, eqModelID, partModelID uuid.UUID, qty int) domain.EquipmentModelPartRequirement {
	var req domain.EquipmentModelPartRequirement
	db.Where(domain.EquipmentModelPartRequirement{EquipmentModelID: eqModelID, PartModelID: partModelID}).
		Assign(domain.EquipmentModelPartRequirement{Quantity: qty}).
		FirstOrCreate(&req)
	return req
}

func createPartConsumptionLog(db *gorm.DB, partModelID uuid.UUID, qty int, notes string) domain.PartConsumptionLog {
	var log domain.PartConsumptionLog
	db.Where(domain.PartConsumptionLog{PartModelID: partModelID, Notes: notes}).
		Assign(domain.PartConsumptionLog{QuantityUsed: qty}).
		FirstOrCreate(&log)
	return log
}

func createEquipmentModelMaintenanceRule(db *gorm.DB, eqModelID uuid.UUID, name string, intervalHours *float64, intervalMonths *int) domain.EquipmentModelMaintenanceRule {
	var rule domain.EquipmentModelMaintenanceRule
	db.Where(domain.EquipmentModelMaintenanceRule{EquipmentModelID: eqModelID, RuleName: name}).
		Assign(domain.EquipmentModelMaintenanceRule{IntervalHours: intervalHours, IntervalMonths: intervalMonths}).
		FirstOrCreate(&rule)
	return rule
}
