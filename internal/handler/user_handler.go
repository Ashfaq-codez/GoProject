package handler

import (
	"database/sql"
	"time"
	"user-api/db/sqlc"
	"user-api/internal/logger"
	"user-api/internal/models"
	"user-api/internal/repository"
	"user-api/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler struct {
	repo     *repository.UserRepository
	validate *validator.Validate
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo:     repo,
		validate: validator.New(),
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed: " + err.Error()})
	}

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format. Use YYYY-MM-DD"})
	}

	user, err := h.repo.Queries().CreateUser(c.Context(), db.CreateUserParams{
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		logger.Log.Error("Failed to create user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format("2006-01-02"),
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.repo.Queries().GetUser(c.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		logger.Log.Error("Failed to fetch user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	age := service.CalculateAge(user.Dob)

	return c.Status(fiber.StatusOK).JSON(models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format("2006-01-02"),
		Age:  &age,
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format. Use YYYY-MM-DD"})
	}

	user, err := h.repo.Queries().UpdateUser(c.Context(), db.UpdateUserParams{
		ID:   int32(id),
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		logger.Log.Error("Failed to update user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.Status(fiber.StatusOK).JSON(models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format("2006-01-02"),
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = h.repo.Queries().DeleteUser(c.Context(), int32(id))
	if err != nil {
		logger.Log.Error("Failed to delete user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// Bonus: Pagination Support
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	users, err := h.repo.Queries().ListUsers(c.Context(), db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		logger.Log.Error("Failed to list users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list users"})
	}

	var response []models.UserResponse
	for _, u := range users {
		age := service.CalculateAge(u.Dob)
		response = append(response, models.UserResponse{
			ID:   u.ID,
			Name: u.Name,
			Dob:  u.Dob.Format("2006-01-02"),
			Age:  &age,
		})
	}

	// Return empty array instead of null if no users exist
	if response == nil {
		response = []models.UserResponse{}
	}

	return c.Status(fiber.StatusOK).JSON(response)
}