package container

import (
	"testing"

	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/service"
)

func TestNewContainer(t *testing.T) {
	// This test requires database initialization
	// For now, we'll skip if database is not available
	
	// Initialize database for testing
	repo.InitDB()
	defer repo.CloseDB()
	
	container := NewContainer()
	
	// Test that all services are properly initialized
	if container.MemberService() == nil {
		t.Error("MemberService should not be nil")
	}
	
	if container.LogtoService() == nil {
		t.Error("LogtoService should not be nil")
	}
	
	// Test that all repositories are properly initialized
	if container.MemberRepository() == nil {
		t.Error("MemberRepository should not be nil")
	}
	
	if container.RoleRepository() == nil {
		t.Error("RoleRepository should not be nil")
	}
}

func TestContainerServiceTypes(t *testing.T) {
	// Initialize database for testing
	repo.InitDB()
	defer repo.CloseDB()
	
	container := NewContainer()
	
	// Test that services implement the correct interfaces
	_, ok := container.MemberService().(service.MemberServiceInterface)
	if !ok {
		t.Error("MemberService should implement MemberServiceInterface")
	}
	
	_, ok = container.LogtoService().(service.LogtoServiceInterface)
	if !ok {
		t.Error("LogtoService should implement LogtoServiceInterface")
	}
}

func TestContainerRepositoryTypes(t *testing.T) {
	// Initialize database for testing
	repo.InitDB()
	defer repo.CloseDB()
	
	container := NewContainer()
	
	// Test that repositories implement the correct interfaces
	_, ok := container.MemberRepository().(repo.MemberRepository)
	if !ok {
		t.Error("MemberRepository should implement MemberRepository interface")
	}
	
	_, ok = container.RoleRepository().(repo.RoleRepository)
	if !ok {
		t.Error("RoleRepository should implement RoleRepository interface")
	}
}