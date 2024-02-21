# Project Design (TODO)

The overall architecture is based on domain-driven development and ports-adapters.

Project structure:

```sh
internal/
  common/
  infra/
  stock/
  usecase/
```

## Domains

There are X domains in this project:

- common
- stock
- user
- trade
- screener

### Stock

The app is mainly dealing with stock time series data and behaviours around them.

### Example

```bash
project_root/
├── cmd/              # Contains application entry points
|   ├── main.go        # Application entry point
|   └── server.go      # HTTP server entry point (optional)
├── configs/           # Configuration files
|   ├── app.yaml        # General application configuration
|   └── database.yaml  # Database connection details
├── pkg/               # Shared utilities and libraries
|   ├── errors/         # Error handling utilities
|   |   ├── errors.go    # Custom error types and handling functions
|   |   └── utils.go      # Helper functions for error handling
|   └── utils/          # Common functions and helper code
|       ├── logging.go   # Logging utilities
|       └── validation.go # Data validation functions
├── internal/          # Internal application components
|   ├── domain/       # Domain logic and entities
|   |   ├── entities/   # Domain models with business logic
|   |   |   ├── user.go      # User entity with fields and methods
|   |   |   └── product.go   # Product entity with fields and methods
|   |   ├── events/      # Domain events and handlers
|   |   |   ├── user_created.go # Event representing user creation
|   |   |   └── product_updated.go # Event representing product update
|   |   |   └── handlers.go    # Handlers for domain events
|   |   ├── repositories/ # Abstracted access to persistence
|   |   |   ├── interface.go   # Repository interface definition
|   |   |   └── user_repository.go # User repository implementation
|   |   |   └── product_repository.go # Product repository implementation
|   |   └── services/     # Domain-specific operations
|   |       ├── user_service.go # User-related operations
|   |       └── product_service.go # Product-related operations
|   ├── application/  # Application services and use cases
|   |   ├── usecases/        # Use case implementations
|   |   |   ├── register_user.go # Use case for user registration
|   |   |   └── update_product.go # Use case for product update
|   |   └── services/        # Application services
|   |       ├── user_app_service.go # User application service
|   |       └── product_app_service.go # Product application service
|   ├── infrastructure/ # Infrastructure implementation details
|   |   ├── adapters/     # Adapters to specific technologies
|   |   |   └── db_adapter.go   # Database adapter
|   |   └── persistence/  # Persistence logic (e.g., repository implementations)
|   |       ├── user_repo_impl.go  # User repository implementation using DB adapter
|   |       └── product_repo_impl.go # Product repository implementation using DB adapter
└── test/              # Unit and integration tests
    ├── unit/            # Unit tests
    |   ├── entities_test.go # Tests for domain entities
    |   └── services_test.go # Tests for domain services
    └── integration/      # Integration tests
        ├── repositories_test.go # Tests for repository implementations
        └── application_test.go # Tests for application services
```

```go
// user.go
package entities

import (
  "time"
)

// User represents a user in the system
type User struct {
  ID        int64  `json:"id"`
  Username  string `json:"username"`
  Email     string `json:"email"`
  FirstName string `json:"first_name"`
  LastName  string `json:"last_name"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new User with the provided information
func NewUser(username, email, firstName, lastName string) (*User, error) {
  // Add validation and business logic here

  return &User{
    Username:  username,
    Email:     email,
    FirstName: firstName,
    LastName:  lastName,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }, nil
}

// ChangeEmail changes the user's email address
func (u *User) ChangeEmail(newEmail string) error {
  // Add validation and business logic here

  u.Email = newEmail
  u.UpdatedAt = time.Now()
  return nil
}

// Other methods specific to User entity behavior
```

```go
// interface.go
package repositories

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
  // GetUserByID retrieves a user by their ID
  GetByID(int64) (*User, error)

  // CreateUser creates a new user
  CreateUser(*User) error

  // UpdateUser updates an existing user
  UpdateUser(*User) error

  // DeleteUser deletes a user
  DeleteUser(int64) error
}
```

```go
// user_repository.go
// This user_repository.go file implements the UserRepository interface for
// a specific data storage technology (replace project/internal/infrastructure/adapters
// with your actual adapter package).
package repositories

import (
  "errors"
  "fmt"
  "project/internal/domain/entities"
  "project/internal/infrastructure/adapters" // replace with your specific adapter
)

// UserRepositoryImpl implements the UserRepository interface
type UserRepositoryImpl struct {
  dbAdapter adapters.DatabaseAdapter // Adapter to the underlying database
}

// NewUserRepository creates a new UserRepositoryImpl instance
func NewUserRepository(dbAdapter adapters.DatabaseAdapter) UserRepository {
  return &UserRepositoryImpl{dbAdapter: dbAdapter}
}

// GetByID retrieves a user by their ID
func (repo *UserRepositoryImpl) GetByID(id int64) (*entities.User, error) {
  // Implement logic to fetch user data from the database using the adapter
  // ...
  if user == nil {
    return nil, errors.New("user not found")
  }
  return user, nil
}

// CreateUser creates a new user
func (repo *UserRepositoryImpl) CreateUser(user *entities.User) error {
  // Implement logic to insert user data into the database using the adapter
  // ...
  return nil
}

// UpdateUser updates an existing user
func (repo *UserRepositoryImpl) UpdateUser(user *entities.User) error {
  // Implement logic to update user data in the database using the adapter
  // ...
  return nil
}

// DeleteUser deletes a user
func (repo *UserRepositoryImpl) DeleteUser(id int64) error {
  // Implement logic to delete user data from the database using the adapter
  // ...
  return nil
}
```

```go
// user_service.go
package services

import (
  "errors"
  "project/internal/domain/entities"
  "project/internal/domain/repositories"
)

// UserService defines operations related to user management
type UserService interface {
  // GetUser retrieves a user by their ID
  GetUser(int64) (*entities.User, error)

  // RegisterUser creates a new user
  RegisterUser(*entities.User) error
}

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
  userRepository repositories.UserRepository
}

// NewUserService creates a new UserServiceImpl instance
func NewUserService(userRepository repositories.UserRepository) UserService {
  return &UserServiceImpl{userRepository: userRepository}
}

// GetUser retrieves a user by their ID
func (s *UserServiceImpl) GetUser(id int64) (*entities.User, error) {
  return s.userRepository.GetByID(id)
}

// RegisterUser creates a new user
func (s *UserServiceImpl) RegisterUser(user *entities.User) error {
  // Add validation and business logic here, potentially interacting with other services
  if err := s.userRepository.CreateUser(user); err != nil {
    return err
  }
  return nil
}
```

```go
// register_user.go
package usecases

import (
  "errors"
  "project/internal/domain/entities"
)

// RegisterUserRequest defines the data required for user registration
type RegisterUserRequest struct {
  Username string `json:"username"`
  Email    string `json:"email"`
  Password string `json:"password"`
}

// RegisterUserUseCase defines the interface for registering a new user
type RegisterUserUseCase interface {
  Execute(*RegisterUserRequest) error
}

// RegisterUserUseCaseImpl implements the RegisterUserUseCase interface
type RegisterUserUseCaseImpl struct {
  userService UserService // Dependency injection for UserService
}

// NewRegisterUserUseCase creates a new RegisterUserUseCaseImpl instance
func NewRegisterUserUseCase(userService UserService) RegisterUserUseCase {
  return &RegisterUserUseCaseImpl{userService: userService}
}

// Execute registers a new user
func (uc *RegisterUserUseCaseImpl) Execute(req *RegisterUserRequest) error {
  // Validate user data
  if err := validateUserData(req); err != nil {
    return err
  }

  // Hash password before persisting
  hashedPassword, err := hashPassword(req.Password)
  if err != nil {
    return errors.New("failed to hash password")
  }

  // Create a new user entity
  user := &entities.User{
    Username: req.Username,
    Email:    req.Email,
    Password: hashedPassword,
  }

  // Register the user using the UserService
  if err := uc.userService.RegisterUser(user); err != nil {
    return err
  }

  return nil
}

// validateUserData performs validation on the registration request
func validateUserData(req *RegisterUserRequest) error {
  // Implement your specific validation logic here
  return nil
}

// hashPassword hashes the provided password using a secure algorithm
func hashPassword(password string) (string, error) {
  // Implement your password hashing logic here
  // ...
  return hashedPassword, nil
}
```

```go
// user_app_service.go
package services

import (
  "errors"
  "project/internal/application/usecases"
)

// UserAppservice defines application-level user operations
type UserAppservice interface {
  // RegisterUser handles user registration with potential use case orchestration
  RegisterUser(*usercases.RegisterUserRequest) error
}

// UserAppserviceImpl implements the UserAppservice interface
type UserAppserviceImpl struct {
  registerUserUseCase usecases.RegisterUserUseCase
}

// NewUserAppservice creates a new UserAppserviceImpl instance
func NewUserAppservice(registerUserUseCase usecases.RegisterUserUseCase) UserAppservice {
  return &UserAppserviceImpl{registerUserUseCase: registerUserUseCase}
}

// RegisterUser handles user registration with potential use case orchestration
func (s *UserAppserviceImpl) RegisterUser(req *usercases.RegisterUserRequest) error {
  // Perform any necessary validation or pre-processing
  if err := s.registerUserUseCase.Execute(req); err != nil {
    return err
  }
  return nil
}
```

```go
// db_adapter.go
package adapters

// DatabaseAdapter defines the interface for interacting with the database
type DatabaseAdapter interface {
  // Connect establishes a connection to the database
  Connect() error

  // Disconnect closes the connection to the database
  Disconnect() error

  // ExecuteQuery executes a query on the database
  ExecuteQuery(string, map[string]interface{}) (*Result, error)

  // BeginTransaction starts a database transaction
  BeginTransaction() (*Transaction, error)

  // CommitTransaction commits a database transaction
  CommitTransaction(*Transaction) error

  // RollbackTransaction rolls back a database transaction
  RollbackTransaction(*Transaction) error
}

// Result represents the result of a database query
type Result struct {
  // Rows affected by the query
  RowsAffected int64

  // Rows returned by the query
  Rows []map[string]interface{}
}

// Transaction represents a database transaction
type Transaction struct {
  // ... transaction specific details
}
```

```go
// user_repo_impl.go
package repositories

import (
  "errors"
  "fmt"
  "project/internal/domain/entities"
  "project/internal/infrastructure/adapters" // replace with your specific adapter
)

// UserRepositoryImpl implements the UserRepository interface using the database adapter
type UserRepositoryImpl struct {
  dbAdapter adapters.DatabaseAdapter
}

// NewUserRepository creates a new UserRepositoryImpl instance
func NewUserRepository(dbAdapter adapters.DatabaseAdapter) UserRepository {
  return &UserRepositoryImpl{dbAdapter: dbAdapter}
}

// GetByID retrieves a user by their ID
func (repo *
```

```go
// product_repo_impl.go
package repositories

import (
  "errors"
  "fmt"
  "project/internal/domain/entities"
  "project/internal/infrastructure/adapters" // replace with your specific adapter
)

// ProductRepositoryImpl implements the ProductRepository interface using the database adapter
type ProductRepositoryImpl struct {
  dbAdapter adapters.DatabaseAdapter
}

// NewProductRepository creates a new ProductRepositoryImpl instance
func NewProductRepository(dbAdapter adapters.DatabaseAdapter) ProductRepository {
  return &ProductRepositoryImpl{dbAdapter: dbAdapter}
}

// GetByID retrieves a product by its ID
func (repo *ProductRepositoryImpl) GetByID(id int64) (*entities.Product, error) {
  // Implement logic to fetch product data from the database using the adapter
  // ...
  if product == nil {
    return nil, errors.New("product not found")
  }
  return product, nil
}

// CreateProduct creates a new product
func (repo *ProductRepositoryImpl) CreateProduct(product *entities.Product) error {
  // Implement logic to insert product data into the database using the adapter
  // ...
  return nil
}

// UpdateProduct updates an existing product
func (repo *ProductRepositoryImpl) UpdateProduct(product *entities.Product) error {
  // Implement logic to update product data in the database using the adapter
  // ...
  return nil
}

// DeleteProduct deletes a product
func (repo *ProductRepositoryImpl) DeleteProduct(id int64) error {
  // Implement logic to delete product data from the database using the adapter
  // ...
  return nil
}
```
