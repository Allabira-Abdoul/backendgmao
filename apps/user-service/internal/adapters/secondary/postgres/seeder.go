package postgres

import (
	"log"
	"os"

	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/pkg/auth"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Seed creates the default roles and admin user on first startup.
// This is idempotent — it skips creation if the data already exists.
func Seed(db *gorm.DB) {
	seedRoles(db)
	seedAdminUser(db)
	seedDefaultManager(db)
	seedDefaultTechnician(db)
	seedDefaultOperator(db)
	seedDefaultPlanner(db)
	seedDefaultAuditor(db)
}

func seedRoles(db *gorm.DB) {
	roles := []struct {
		Name        string
		Description string
		Privileges  []string
	}{
		{
			Name:        "Administrator",
			Description: "Full system access — all privileges granted",
			Privileges:  domain.AllPrivileges(),
		},
		{
			Name:        "Manager",
			Description: "Operational management — approval, analytics, and oversight",
			Privileges: []string{
				domain.PrivilegeUserView, domain.PrivilegeUserCreate, domain.PrivilegeUserUpdate,
				domain.PrivilegeRoleView,
				domain.PrivilegeAssetView, domain.PrivilegeAssetCreate, domain.PrivilegeAssetUpdate,
				domain.PrivilegeWorkOrderView, domain.PrivilegeWorkOrderCreate, domain.PrivilegeWorkOrderUpdate,
				domain.PrivilegeWorkOrderAssign, domain.PrivilegeWorkOrderApprove, domain.PrivilegeWorkOrderClose,
				domain.PrivilegeMaintenanceView, domain.PrivilegeMaintenancePlanCreate,
				domain.PrivilegeMaintenancePlanUpdate, domain.PrivilegeMaintenanceSchedule,
				domain.PrivilegeInventoryView,
				domain.PrivilegeAnalyticsView, domain.PrivilegeAnalyticsExport,
			},
		},
		{
			Name:        "Technician",
			Description: "Field technician — maintenance and asset operations",
			Privileges: []string{
				domain.PrivilegeAssetView, domain.PrivilegeAssetUpdate,
				domain.PrivilegeWorkOrderView, domain.PrivilegeWorkOrderUpdate, domain.PrivilegeWorkOrderClose,
				domain.PrivilegeMaintenanceView,
				domain.PrivilegeInventoryView, domain.PrivilegeInventoryUpdate,
			},
		},
		{
			Name:        "Operator",
			Description: "Operator — perform maintenance tasks",
			Privileges: []string{
				domain.PrivilegeAssetView,
				domain.PrivilegeWorkOrderView, domain.PrivilegeWorkOrderUpdate, domain.PrivilegeWorkOrderClose,
				domain.PrivilegeMaintenanceView,
				domain.PrivilegeInventoryView, domain.PrivilegeInventoryUpdate,
			},
		},
		{
			Name:        "Planner",
			Description: "Planner — maintenance planning and scheduling",
			Privileges: []string{
				domain.PrivilegeAssetView,
				domain.PrivilegeWorkOrderView, domain.PrivilegeWorkOrderCreate, domain.PrivilegeWorkOrderUpdate,
				domain.PrivilegeMaintenanceView, domain.PrivilegeMaintenancePlanCreate,
				domain.PrivilegeMaintenancePlanUpdate, domain.PrivilegeMaintenanceSchedule,
				domain.PrivilegeInventoryView,
				domain.PrivilegeAnalyticsView, domain.PrivilegeAnalyticsExport,
			},
		},
		{
			Name:        "Auditor",
			Description: "Auditor — perform audit tasks",
			Privileges: []string{
				domain.PrivilegeAssetView,
				domain.PrivilegeWorkOrderView,
				domain.PrivilegeMaintenanceView,
				domain.PrivilegeInventoryView,
				domain.PrivilegeSystemAuditView,
				domain.PrivilegeAuditLogView, domain.PrivilegeAuditLogExport, domain.PrivilegeAuditLogImport,
			},
		},
	}

	for _, r := range roles {
		var existing domain.Role
		result := db.Where("name = ?", r.Name).First(&existing)
		if result.Error == nil {
			log.Printf("Seeder: Role '%s' already exists, skipping", r.Name)
			continue
		}

		role := domain.Role{
			Name:        r.Name,
			Description: r.Description,
		}

		if err := db.Create(&role).Error; err != nil {
			log.Printf("Seeder: Failed to create role '%s': %v", r.Name, err)
			continue
		}

		// Set privileges
		rolePrivileges := make([]domain.RolePrivilege, 0, len(r.Privileges))
		for _, p := range r.Privileges {
			rolePrivileges = append(rolePrivileges, domain.RolePrivilege{
				RoleID:    role.ID,
				Privilege: p,
			})
		}

		if err := db.Create(&rolePrivileges).Error; err != nil {
			log.Printf("Seeder: Failed to set privileges for role '%s': %v", r.Name, err)
			continue
		}

		log.Printf("Seeder: Created role '%s' with %d privileges", r.Name, len(r.Privileges))
	}
}

func seedAdminUser(db *gorm.DB) {
	// Check if any admin user already exists
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 0 {
		log.Println("Seeder: Users already exist, skipping admin user creation")
		return
	}

	// Get the Administrator role
	var adminRole domain.Role
	if err := db.Where("name = ?", "Administrator").First(&adminRole).Error; err != nil {
		log.Printf("Seeder: Cannot find Administrator role, skipping admin user: %v", err)
		return
	}

	// Get admin password from env or use a default
	adminPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "Admin@2026!"
		log.Println("Seeder: WARNING — Using default admin password. Set DEFAULT_ADMIN_PASSWORD env var in production!")
	}

	hashedPassword, err := auth.HashPassword(adminPassword)
	if err != nil {
		log.Printf("Seeder: Failed to hash admin password: %v", err)
		return
	}

	adminUser := domain.User{
		ID:           uuid.New(),
		FullName:     "System Administrator",
		Email:        "admin@gmao.local",
		Password:     hashedPassword,
		Status:       domain.StatusActive,
		RoleID:       adminRole.ID,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		log.Printf("Seeder: Failed to create admin user: %v", err)
		return
	}

	log.Println("Seeder: Created default admin user (admin@gmao.local)")
}

func seedDefaultManager(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 1 {
		log.Println("Seeder: Users already exist, skipping manager user creation")
		return
	}

	var managerRole domain.Role
	if err := db.Where("name = ?", "Manager").First(&managerRole).Error; err != nil {
		log.Printf("Seeder: Cannot find Manager role, skipping manager user: %v", err)
		return
	}

	managerPassword := "Manager@2026!"
	hashedPassword, err := auth.HashPassword(managerPassword)
	if err != nil {
		log.Printf("Seeder: Failed to hash manager password: %v", err)
		return
	}

	managerUser := domain.User{
		ID:           uuid.New(),
		FullName:     "System Manager",
		Email:        "manager@gmao.local",
		Password:     hashedPassword,
		Status:       domain.StatusActive,
		RoleID:       managerRole.ID,
	}

	if err := db.Create(&managerUser).Error; err != nil {
		log.Printf("Seeder: Failed to create manager user: %v", err)
		return
	}

	log.Println("Seeder: Created default manager user (manager@gmao.local)")
}

func seedDefaultTechnician(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 2 {
		log.Println("Seeder: Users already exist, skipping technician user creation")
		return
	}

	var technicianRole domain.Role
	if err := db.Where("name = ?", "Technician").First(&technicianRole).Error; err != nil {
		log.Printf("Seeder: Cannot find Technician role, skipping technician user: %v", err)
		return
	}

	technicianPassword := "Technician@2026!"
	hashedPassword, err := auth.HashPassword(technicianPassword)
	if err != nil {
		log.Printf("Seeder: Failed to hash technician password: %v", err)
		return
	}

	technicianUser := domain.User{
		ID:           uuid.New(),
		FullName:     "System Technician",
		Email:        "technician@gmao.local",
		Password:     hashedPassword,
		Status:       domain.StatusActive,
		RoleID:       technicianRole.ID,
	}

	if err := db.Create(&technicianUser).Error; err != nil {
		log.Printf("Seeder: Failed to create technician user: %v", err)
		return
	}

	log.Println("Seeder: Created default technician user (technician@gmao.local)")
}

func seedDefaultOperator(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 3 {
		log.Println("Seeder: Users already exist, skipping operator user creation")
		return
	}

	var operatorRole domain.Role
	if err := db.Where("name = ?", "Operator").First(&operatorRole).Error; err != nil {
		log.Printf("Seeder: Cannot find Operator role, skipping operator user: %v", err)
		return
	}

	operatorPassword := "Operator@2026!"
	hashedPassword, err := auth.HashPassword(operatorPassword)
	if err != nil {
		log.Printf("Seeder: Failed to hash operator password: %v", err)
		return
	}

	operatorUser := domain.User{
		ID:           uuid.New(),
		FullName:     "System Operator",
		Email:        "operator@gmao.local",
		Password:     hashedPassword,
		Status:       domain.StatusActive,
		RoleID:       operatorRole.ID,
	}

	if err := db.Create(&operatorUser).Error; err != nil {
		log.Printf("Seeder: Failed to create operator user: %v", err)
		return
	}

	log.Println("Seeder: Created default operator user (operator@gmao.local)")
}

func seedDefaultPlanner(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 4 {
		log.Println("Seeder: Users already exist, skipping planner user creation")
		return
	}

	var plannerRole domain.Role
	if err := db.Where("name = ?", "Planner").First(&plannerRole).Error; err != nil {
		log.Printf("Seeder: Cannot find Planner role, skipping planner user: %v", err)
		return
	}

	plannerPassword := "Planner@2026!"
	hashedPassword, err := auth.HashPassword(plannerPassword)
	if err != nil {
		log.Printf("Seeder: Failed to hash planner password: %v", err)
		return
	}

	plannerUser := domain.User{
		ID:           uuid.New(),
		FullName:     "System Planner",
		Email:        "planner@gmao.local",
		Password:     hashedPassword,
		Status:       domain.StatusActive,
		RoleID:       plannerRole.ID,
	}

	if err := db.Create(&plannerUser).Error; err != nil {
		log.Printf("Seeder: Failed to create planner user: %v", err)
		return
	}

	log.Println("Seeder: Created default planner user (planner@gmao.local)")
}

func seedDefaultAuditor(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 5 {
		log.Println("Seeder: Users already exist, skipping auditor user creation")
		return
	}

	var auditorRole domain.Role
	if err := db.Where("name = ?", "Auditor").First(&auditorRole).Error; err != nil {
		log.Printf("Seeder: Cannot find Auditor role, skipping auditor user: %v", err)
		return
	}

	auditorPassword := "Auditor@2026!"
	hashedPassword, err := auth.HashPassword(auditorPassword)
	if err != nil {
		log.Printf("Seeder: Failed to hash auditor password: %v", err)
		return
	}

	auditorUser := domain.User{
		ID:           uuid.New(),
		FullName:     "System Auditor",
		Email:        "auditor@gmao.local",
		Password:     hashedPassword,
		Status:       domain.StatusActive,
		RoleID:       auditorRole.ID,
	}

	if err := db.Create(&auditorUser).Error; err != nil {
		log.Printf("Seeder: Failed to create auditor user: %v", err)
		return
	}

	log.Println("Seeder: Created default auditor user (auditor@gmao.local)")
}