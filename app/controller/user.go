package controller

import (
	"log"

	"github.com/aicacia/pickup/app/middleware"
	"github.com/aicacia/pickup/app/model"
	"github.com/aicacia/pickup/app/repository"
	"github.com/gofiber/fiber/v2"
)

// GetCurrentUser
//
//	@Summary		Get current user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.UserST
//	@Failure		400	{object}	model.ErrorST
//	@Failure		401	{object}	model.ErrorST
//	@Failure		403	{object}	model.ErrorST
//	@Failure		500	{object}	model.ErrorST
//	@Router			/user [get]
//	@Security		Authorization
func GetCurrentUser(c *fiber.Ctx) error {
	return c.JSON(model.UserFromUserRow(*middleware.GetUser(c)))
}

// PatchCompleteUser
//
//	@Summary		complete current user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			complete	body    model.CompleteUserST	true	"complete user"
//	@Success		200	{object}	model.UserST
//	@Failure		400	{object}	model.ErrorST
//	@Failure		401	{object}	model.ErrorST
//	@Failure		403	{object}	model.ErrorST
//	@Failure		500	{object}	model.ErrorST
//	@Router			/user/complete [patch]
//	@Security		Authorization
func PatchCompleteUser(c *fiber.Ctx) error {
	var complete model.CompleteUserST
	if err := c.BodyParser(&complete); err != nil {
		log.Printf("error parsing body: %v\n", err)
		return model.NewError(400).AddError("complete", "invalid", "body").Send(c)
	}
	if !complete.TermsOfServiceAcknowledge {
		return model.NewError(400).AddError("terms_of_service_acknowledge", "required", "body").Send(c)
	}
	user, err := repository.UpdateUser(middleware.GetUser(c).Id, repository.UpdateUserST{
		Username:                  &complete.Username,
		TermsOfServiceAcknowledge: &complete.TermsOfServiceAcknowledge,
		Birthdate:                 &complete.Birthdate,
		Sex:                       &complete.Sex,
	})
	if err != nil {
		log.Printf("error updating user: %v\n", err)
		return model.NewError(500).AddError("application", "internal", "database").Send(c)
	}
	if user == nil {
		return model.NewError(500).AddError("application", "internal", "database").Send(c)
	}
	return c.JSON(model.UserFromUserRow(*user))
}

// PatchUpdateUser
//
//	@Summary		update current user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			update	body    model.UpdateUserST	true	"complete user"
//	@Success		200	{object}	model.UserST
//	@Failure		400	{object}	model.ErrorST
//	@Failure		401	{object}	model.ErrorST
//	@Failure		403	{object}	model.ErrorST
//	@Failure		500	{object}	model.ErrorST
//	@Router			/user [patch]
//	@Security		Authorization
func PatchUpdateUser(c *fiber.Ctx) error {
	var update model.UpdateUserST
	if err := c.BodyParser(&update); err != nil {
		log.Printf("error parsing body: %v\n", err)
		return model.NewError(400).AddError("complete", "invalid", "body").Send(c)
	}
	user, err := repository.UpdateUser(middleware.GetUser(c).Id, repository.UpdateUserST{
		Username:  &update.Username,
		Birthdate: &update.Birthdate,
		Sex:       &update.Sex,
	})
	if err != nil {
		log.Printf("error updating user: %v\n", err)
		return model.NewError(500).AddError("application", "internal", "database").Send(c)
	}
	if user == nil {
		return model.NewError(500).AddError("application", "internal", "database").Send(c)
	}
	return c.JSON(model.UserFromUserRow(*user))
}
