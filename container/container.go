package container

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/service"
	"github.com/spf13/viper"
)

// Container holds all application dependencies
type Container struct {
	// Database
	db *sqlx.DB
	sq squirrel.StatementBuilderType
	
	// Repositories
	memberRepo repo.MemberRepository
	eventRepo  repo.EventRepository
	clientRepo repo.ClientRepository
	roleRepo   repo.RoleRepository
	dbManager  repo.DatabaseManager
	
	// Services
	memberService service.MemberServiceInterface
	logtoService  service.LogtoServiceInterface
	eventService  service.EventServiceInterface
	clientService service.ClientServiceInterface
}

// NewContainer creates and initializes a new dependency injection container
func NewContainer() *Container {
	container := &Container{}
	container.initializeRepositories()
	container.initializeServices()
	return container
}

// initializeRepositories sets up all repository instances
func (c *Container) initializeRepositories() {
	// Get database connection from global state (for now, during migration)
	c.db = repo.GetDB() // We'll need to add this method to implementations.go
	c.sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	
	// Initialize repositories
	c.memberRepo = repo.NewMemberRepository(c.db, c.sq)
	c.eventRepo = repo.NewEventRepository(c.db, c.sq)
	c.clientRepo = repo.NewClientRepository(c.db, c.sq)
	c.roleRepo = repo.NewRoleRepository(c.db, c.sq)
	c.dbManager = repo.NewDatabaseManager()
}

// initializeServices sets up all service instances with their dependencies
func (c *Container) initializeServices() {
	// Initialize Logto service
	logtoEndpoint := viper.GetString("logto.endpoint")
	c.logtoService = service.NewLogtoService(logtoEndpoint)
	
	// Initialize member service with dependencies
	c.memberService = service.NewMemberService(
		c.memberRepo,
		c.roleRepo,
		c.logtoService,
	)
	
	// TODO: Initialize other services
	// c.eventService = service.NewEventService(c.eventRepo, c.memberService)
	// c.clientService = service.NewClientService(c.clientRepo)
}

// Getter methods for services

// MemberService returns the member service instance
func (c *Container) MemberService() service.MemberServiceInterface {
	return c.memberService
}

// LogtoService returns the logto service instance
func (c *Container) LogtoService() service.LogtoServiceInterface {
	return c.logtoService
}

// EventService returns the event service instance
func (c *Container) EventService() service.EventServiceInterface {
	return c.eventService
}

// ClientService returns the client service instance
func (c *Container) ClientService() service.ClientServiceInterface {
	return c.clientService
}

// Getter methods for repositories (for testing)

// MemberRepository returns the member repository instance
func (c *Container) MemberRepository() repo.MemberRepository {
	return c.memberRepo
}

// EventRepository returns the event repository instance
func (c *Container) EventRepository() repo.EventRepository {
	return c.eventRepo
}

// ClientRepository returns the client repository instance
func (c *Container) ClientRepository() repo.ClientRepository {
	return c.clientRepo
}

// RoleRepository returns the role repository instance
func (c *Container) RoleRepository() repo.RoleRepository {
	return c.roleRepo
}

// DatabaseManager returns the database manager instance
func (c *Container) DatabaseManager() repo.DatabaseManager {
	return c.dbManager
}