package domain

// Privilege represents a system-defined capability in the GMAO.
// Privileges are exhaustive and follow the Oracle DB RBAC model.
// They can only be added through system updates (code changes), never by users.

// --- User Management ---
const (
	PrivilegeUserView       = "USER_VIEW"
	PrivilegeUserViewSelf   = "USER_VIEW_SELF"
	PrivilegeUserCreate     = "USER_CREATE"
	PrivilegeUserUpdate     = "USER_UPDATE"
	PrivilegeUserUpdateSelf = "USER_UPDATE_SELF"
	PrivilegeUserDelete     = "USER_DELETE"
	PrivilegeUserAssignRole = "USER_ASSIGN_ROLE"
)

// --- Role Management ---
const (
	PrivilegeRoleView   = "ROLE_VIEW"
	PrivilegeRoleCreate = "ROLE_CREATE"
	PrivilegeRoleUpdate = "ROLE_UPDATE"
	PrivilegeRoleDelete = "ROLE_DELETE"
)

// --- Team Management ---
const (
	PrivilegeTeamView   = "TEAM_VIEW"
	PrivilegeTeamCreate = "TEAM_CREATE"
	PrivilegeTeamUpdate = "TEAM_UPDATE"
	PrivilegeTeamDelete = "TEAM_DELETE"
)

// --- Asset Management ---
const (
	PrivilegeAssetView     = "ASSET_VIEW"
	PrivilegeAssetCreate   = "ASSET_CREATE"
	PrivilegeAssetUpdate   = "ASSET_UPDATE"
	PrivilegeAssetDelete   = "ASSET_DELETE"
	PrivilegeAssetTransfer = "ASSET_TRANSFER"
	PrivilegeAssetHistory  = "ASSET_HISTORY"
	PrivilegeAssetDocumentUpload   = "ASSET_DOCUMENT_UPLOAD"
	PrivilegeAssetDocumentView   = "ASSET_DOCUMENT_VIEW"
	PrivilegeAssetDocumentDelete   = "ASSET_DOCUMENT_DELETE"
	PrivilegeAssetDocumentUpdate   = "ASSET_DOCUMENT_UPDATE"
	PrivilegeAssetDocumentDownload   = "ASSET_DOCUMENT_DOWNLOAD"
)

// --- Work Order Management ---
const (
	PrivilegeWorkOrderView    = "WORKORDER_VIEW"
	PrivilegeWorkOrderCreate  = "WORKORDER_CREATE"
	PrivilegeWorkOrderUpdate  = "WORKORDER_UPDATE"
	PrivilegeWorkOrderDelete  = "WORKORDER_DELETE"
	PrivilegeWorkOrderAssign  = "WORKORDER_ASSIGN"
	PrivilegeWorkOrderApprove = "WORKORDER_APPROVE"
	PrivilegeWorkOrderClose   = "WORKORDER_CLOSE"
)

// --- Maintenance Management ---
const (
	PrivilegeMaintenanceView       = "MAINTENANCE_VIEW"
	PrivilegeMaintenancePlanCreate = "MAINTENANCE_PLAN_CREATE"
	PrivilegeMaintenancePlanUpdate = "MAINTENANCE_PLAN_UPDATE"
	PrivilegeMaintenancePlanDelete = "MAINTENANCE_PLAN_DELETE"
	PrivilegeMaintenanceSchedule   = "MAINTENANCE_SCHEDULE"
)

// --- Inventory Management ---
const (
	PrivilegeInventoryView   = "INVENTORY_VIEW"
	PrivilegeInventoryCreate = "INVENTORY_CREATE"
	PrivilegeInventoryUpdate = "INVENTORY_UPDATE"
	PrivilegeInventoryDelete = "INVENTORY_DELETE"
	PrivilegeInventoryAdjust = "INVENTORY_ADJUST"
)

// --- Analytics ---
const (
	PrivilegeAnalyticsView   = "ANALYTICS_VIEW"
	PrivilegeAnalyticsWrite  = "ANALYTICS_WRITE"
	PrivilegeAnalyticsExport = "ANALYTICS_EXPORT"
	PrivilegeAnalyticsImport = "ANALYTICS_IMPORT"
)

// --- System Administration ---
const (
	PrivilegeSystemAdmin     = "SYSTEM_ADMIN"
	PrivilegeSystemConfig    = "SYSTEM_CONFIG"
	PrivilegeSystemAuditView = "SYSTEM_AUDIT_VIEW"
)

// --- Audit Logs ---
const (
	PrivilegeAuditLogView   = "AUDIT_LOG_VIEW"
	PrivilegeAuditLogExport = "AUDIT_LOG_EXPORT"
	PrivilegeAuditLogImport = "AUDIT_LOG_IMPORT"
)

// AllPrivileges returns the exhaustive list of all system-defined privileges.
func AllPrivileges() []string {
	return []string{
		// User Management
		PrivilegeUserView, PrivilegeUserCreate, PrivilegeUserUpdate,
		PrivilegeUserDelete, PrivilegeUserAssignRole,
		// Role Management
		PrivilegeRoleView, PrivilegeRoleCreate, PrivilegeRoleUpdate, PrivilegeRoleDelete,
		// Asset Management
		PrivilegeAssetView, PrivilegeAssetCreate, PrivilegeAssetUpdate,
		PrivilegeAssetDelete, PrivilegeAssetTransfer,
		// Work Order Management
		PrivilegeWorkOrderView, PrivilegeWorkOrderCreate, PrivilegeWorkOrderUpdate,
		PrivilegeWorkOrderDelete, PrivilegeWorkOrderAssign,
		PrivilegeWorkOrderApprove, PrivilegeWorkOrderClose,
		// Maintenance Management
		PrivilegeMaintenanceView, PrivilegeMaintenancePlanCreate,
		PrivilegeMaintenancePlanUpdate, PrivilegeMaintenancePlanDelete,
		PrivilegeMaintenanceSchedule,
		// Inventory Management
		PrivilegeInventoryView, PrivilegeInventoryCreate,
		PrivilegeInventoryUpdate, PrivilegeInventoryDelete, PrivilegeInventoryAdjust,
		// Analytics
		PrivilegeAnalyticsView, PrivilegeAnalyticsWrite, PrivilegeAnalyticsExport,
		// System
		/// Can perform any action on the system
		PrivilegeSystemAdmin, 
		/// Can configure the system
		PrivilegeSystemConfig, 
		/// Can view the system audit logs
		PrivilegeSystemAuditView,
	}
}

// PrivilegesByDomain returns privileges grouped by their functional domain.
func PrivilegesByDomain() map[string][]string {
	return map[string][]string{
		"User Management": {
			PrivilegeUserView, PrivilegeUserCreate, PrivilegeUserUpdate,
			PrivilegeUserDelete, PrivilegeUserAssignRole,
		},
		"Role Management": {
			PrivilegeRoleView, PrivilegeRoleCreate, PrivilegeRoleUpdate, PrivilegeRoleDelete,
		},
		"Asset Management": {
			PrivilegeAssetView, PrivilegeAssetCreate, PrivilegeAssetUpdate,
			PrivilegeAssetDelete, PrivilegeAssetTransfer,
		},
		"Work Order Management": {
			PrivilegeWorkOrderView, PrivilegeWorkOrderCreate, PrivilegeWorkOrderUpdate,
			PrivilegeWorkOrderDelete, PrivilegeWorkOrderAssign,
			PrivilegeWorkOrderApprove, PrivilegeWorkOrderClose,
		},
		"Maintenance Management": {
			PrivilegeMaintenanceView, PrivilegeMaintenancePlanCreate,
			PrivilegeMaintenancePlanUpdate, PrivilegeMaintenancePlanDelete,
			PrivilegeMaintenanceSchedule,
		},
		"Inventory Management": {
			PrivilegeInventoryView, PrivilegeInventoryCreate,
			PrivilegeInventoryUpdate, PrivilegeInventoryDelete, PrivilegeInventoryAdjust,
		},
		"Analytics": {
			PrivilegeAnalyticsView, PrivilegeAnalyticsWrite, PrivilegeAnalyticsExport,
		},
		"System Administration": {
			PrivilegeSystemAdmin, PrivilegeSystemConfig, PrivilegeSystemAuditView,
		},
	}
}

// IsValidPrivilege checks if a privilege string is a valid system-defined privilege.
func IsValidPrivilege(privilege string) bool {
	for _, p := range AllPrivileges() {
		if p == privilege {
			return true
		}
	}
	return false
}

// ValidatePrivileges checks if all provided privileges are valid.
// Returns the list of invalid privileges found.
func ValidatePrivileges(privileges []string) []string {
	invalid := make([]string, 0)
	for _, p := range privileges {
		if !IsValidPrivilege(p) {
			invalid = append(invalid, p)
		}
	}
	return invalid
}
