package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"authentication-service/internal/domain"
	"authentication-service/internal/repository"
)

// PostgresRepository implements the UserRepository interface
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgresRepository instance
func NewPostgresRepository(db *sqlx.DB) repository.UserRepository {
	return &PostgresRepository{
		db: db,
	}
}

// CreateUser creates a new user in the database
func (r *PostgresRepository) CreateUser(ctx context.Context, user *domain.User) error {
	// fmt.Printf("DEBUG CreateUser: Creating user with email: %s\n", user.Email)

	// Generate UUID for the user
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// KHÔNG đặt user.Active = true, để sử dụng giá trị do service truyền vào
	// để giúp flow đăng ký và OTP hoạt động đúng

	// Default role is "user" if not specified
	if user.Role == "" {
		user.Role = "user"
	}

	// fmt.Printf("DEBUG CreateUser: Generated ID: %s\n", user.ID)
	// fmt.Printf("DEBUG CreateUser: Password hash to store (FULL): %s\n", user.Password)

	query := `INSERT INTO users 
		(id, email, password_hash, first_name, last_name, username, active, role, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Password, // This should already be hashed before it gets here
		user.FirstName,
		user.LastName,
		user.Username,
		user.Active,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			// 23505 is the PostgreSQL error code for unique_violation
			if pqErr.Code == "23505" {
				// fmt.Printf("DEBUG CreateUser: Unique constraint violation: %v\n", err)
				return errors.New("user with this email already exists")
			}
		}
		// fmt.Printf("DEBUG CreateUser: Database error: %v\n", err)
		return err
	}

	// fmt.Printf("DEBUG CreateUser: User created successfully\n")
	return nil
}

// GetUserByID retrieves a user by ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, email, password_hash, first_name, last_name, username, active, role, 
		refresh_token, created_at, updated_at 
		FROM users 
		WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, refresh_token, username, first_name, last_name, 
              role, status, active, created_at, updated_at
              FROM users WHERE email = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password, // Mapped to password_hash column in the database
		&user.RefreshToken,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Status,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUser updates a user's information
func (r *PostgresRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	fmt.Printf("DEBUG UpdateUser: Updating user with ID: %s\n", user.ID)

	// Check if password_hash needs to be updated
	if user.Password != "" {
		fmt.Printf("DEBUG UpdateUser: Updating password for user ID: %s\n", user.ID)
		query := `UPDATE users 
			SET first_name = $1, last_name = $2, username = $3, active = $4, role = $5, updated_at = $6, password_hash = $7
			WHERE id = $8`

		result, err := r.db.ExecContext(
			ctx,
			query,
			user.FirstName,
			user.LastName,
			user.Username,
			user.Active,
			user.Role,
			user.UpdatedAt,
			user.Password, // Password field is mapped to password_hash in the DB
			user.ID,
		)

		if err != nil {
			fmt.Printf("DEBUG UpdateUser: Error updating user: %v\n", err)
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fmt.Printf("DEBUG UpdateUser: Error getting rows affected: %v\n", err)
			return err
		}

		if rowsAffected == 0 {
			fmt.Printf("DEBUG UpdateUser: No rows affected, user may not exist: %s\n", user.ID)
			return errors.New("user not found or not updated")
		}

		fmt.Printf("DEBUG UpdateUser: User updated successfully with new password, rows affected: %d\n", rowsAffected)
	} else {
		// If no password update, use the original query
		query := `UPDATE users 
			SET first_name = $1, last_name = $2, username = $3, active = $4, role = $5, updated_at = $6
			WHERE id = $7`

		result, err := r.db.ExecContext(
			ctx,
			query,
			user.FirstName,
			user.LastName,
			user.Username,
			user.Active,
			user.Role,
			user.UpdatedAt,
			user.ID,
		)

		if err != nil {
			fmt.Printf("DEBUG UpdateUser: Error updating user without password change: %v\n", err)
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fmt.Printf("DEBUG UpdateUser: Error getting rows affected: %v\n", err)
			return err
		}

		if rowsAffected == 0 {
			fmt.Printf("DEBUG UpdateUser: No rows affected, user may not exist: %s\n", user.ID)
			return errors.New("user not found or not updated")
		}

		fmt.Printf("DEBUG UpdateUser: User updated successfully without password change, rows affected: %d\n", rowsAffected)
	}

	return nil
}

// UpdateRefreshToken updates a user's refresh token
func (r *PostgresRepository) UpdateRefreshToken(ctx context.Context, userID, refreshToken string) error {
	query := `UPDATE users SET refresh_token = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, refreshToken, userID)
	if err != nil {
		return fmt.Errorf("failed to update refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with ID: %s", userID)
	}

	return nil
}

// UserExists checks if a user exists by email
func (r *PostgresRepository) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := r.db.GetContext(ctx, &exists, query, email)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetUserByUsername retrieves a user by username
func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	fmt.Printf("DEBUG GetUserByUsername: Looking for user with username: %s\n", username)

	var user domain.User

	query := `SELECT id, email, password_hash, first_name, last_name, username, active, role, 
		refresh_token, created_at, updated_at 
		FROM users 
		WHERE username = $1`

	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("DEBUG GetUserByUsername: User not found with username: %s\n", username)
			return nil, domain.ErrUserNotFound
		}
		fmt.Printf("DEBUG GetUserByUsername: Database error: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG GetUserByUsername: Found user with ID: %s\n", user.ID)
	return &user, nil
}

// LogoutFromAllDevices invalidates all refresh tokens for a user
func (r *PostgresRepository) LogoutFromAllDevices(ctx context.Context, userID string) error {
	// Set refresh_token to NULL (not empty string) to completely invalidate all sessions
	query := `UPDATE users SET refresh_token = NULL WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to logout from all devices: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with ID: %s", userID)
	}

	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
