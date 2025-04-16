package casbin

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

type Enforcer struct {
	enforcer *casbin.Enforcer
	mu       sync.RWMutex
}

type Config struct {
	ModelPath  string
	PolicyPath string
	ModelText  string
}

func NewEnforcer(config Config) (*Enforcer, error) {
	var m model.Model
	var err error

	if config.ModelText != "" {
		m, err = model.NewModelFromString(config.ModelText)
		if err != nil {
			return nil, fmt.Errorf("failed to create model from string: %w", err)
		}
	} else if config.ModelPath != "" {
		m, err = model.NewModelFromFile(config.ModelPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create model from file: %w", err)
		}
	} else {
		return nil, fmt.Errorf("either model path or model text must be provided")
	}

	var e *casbin.Enforcer
	if config.PolicyPath != "" {
		adapter := fileadapter.NewAdapter(config.PolicyPath)
		e, err = casbin.NewEnforcer(m, adapter)
		if err != nil {
			return nil, fmt.Errorf("failed to create enforcer: %w", err)
		}
	} else {
		e, err = casbin.NewEnforcer(m)
		if err != nil {
			return nil, fmt.Errorf("failed to create enforcer: %w", err)
		}
	}

	return &Enforcer{
		enforcer: e,
	}, nil
}

func (e *Enforcer) LoadPolicy() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.LoadPolicy()
}

func (e *Enforcer) SavePolicy() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.SavePolicy()
}

func (e *Enforcer) AddPolicy(params ...interface{}) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.AddPolicy(params...)
}

func (e *Enforcer) RemovePolicy(params ...interface{}) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.RemovePolicy(params...)
}

func (e *Enforcer) AddGroupingPolicy(params ...interface{}) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.AddGroupingPolicy(params...)
}

func (e *Enforcer) RemoveGroupingPolicy(params ...interface{}) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.RemoveGroupingPolicy(params...)
}

func (e *Enforcer) Enforce(params ...interface{}) (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.Enforce(params...)
}

func (e *Enforcer) GetAllRoles() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetAllRoles()
}

func (e *Enforcer) GetAllObjects() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetAllObjects()
}

func (e *Enforcer) GetAllSubjects() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetAllSubjects()
}

func (e *Enforcer) GetAllActions() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetAllActions()
}

func (e *Enforcer) GetPolicy() [][]string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetPolicy()
}

func (e *Enforcer) GetGroupingPolicy() [][]string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetGroupingPolicy()
}

func (e *Enforcer) GetFilteredPolicy(fieldIndex int, fieldValues ...string) [][]string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetFilteredPolicy(fieldIndex, fieldValues...)
}

func (e *Enforcer) GetFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) [][]string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetFilteredGroupingPolicy(fieldIndex, fieldValues...)
}

func (e *Enforcer) HasPolicy(params ...interface{}) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.HasPolicy(params...)
}

func (e *Enforcer) HasGroupingPolicy(params ...interface{}) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.HasGroupingPolicy(params...)
}

func (e *Enforcer) GetRolesForUser(user string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetRolesForUser(user)
}

func (e *Enforcer) GetUsersForRole(role string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetUsersForRole(role)
}

func (e *Enforcer) GetPermissionsForUser(user string) [][]string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetPermissionsForUser(user)
}

func (e *Enforcer) GetImplicitRolesForUser(user string, domain ...string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetImplicitRolesForUser(user, domain...)
}

func (e *Enforcer) GetImplicitPermissionsForUser(user string, domain ...string) [][]string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetImplicitPermissionsForUser(user, domain...)
}

func (e *Enforcer) AddRoleForUser(user string, role string, domain ...string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.AddRoleForUser(user, role, domain...)
}

func (e *Enforcer) DeleteRoleForUser(user string, role string, domain ...string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.DeleteRoleForUser(user, role, domain...)
}

func (e *Enforcer) DeleteRolesForUser(user string, domain ...string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.DeleteRolesForUser(user, domain...)
}

func (e *Enforcer) DeleteUser(user string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.DeleteUser(user)
}

func (e *Enforcer) DeleteRole(role string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.DeleteRole(role)
}

func (e *Enforcer) DeletePermission(permission ...string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.DeletePermission(permission...)
}

func (e *Enforcer) AddPermissionForUser(user string, permission ...string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.AddPermissionForUser(user, permission...)
}

func (e *Enforcer) DeletePermissionForUser(user string, permission ...string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.DeletePermissionForUser(user, permission...)
}

func (e *Enforcer) DeletePermissionsForUser(user string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.DeletePermissionsForUser(user)
}

func (e *Enforcer) GetRoleManager() casbin.RoleManager {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetRoleManager()
}

func (e *Enforcer) SetRoleManager(rm casbin.RoleManager) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.enforcer.SetRoleManager(rm)
}

func (e *Enforcer) EnableAutoSave(autoSave bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.enforcer.EnableAutoSave(autoSave)
}

func (e *Enforcer) EnableAutoBuildRoleLinks(autoBuild bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.enforcer.EnableAutoBuildRoleLinks(autoBuild)
}

func (e *Enforcer) BuildRoleLinks() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.enforcer.BuildRoleLinks()
}
