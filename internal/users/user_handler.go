package users

import "github.com/gofiber/fiber/v2"

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var (
		request RegisterRequest
	)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response, err := h.service.Register(&request)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var (
		request LoginRequest
	)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.service.Login(&request)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    user.Token,
		MaxAge:   60 * 60 * 24,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logout successful",
	})
}
