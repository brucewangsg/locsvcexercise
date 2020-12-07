package coresvc

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/brucewangsg/locsvcexercise/authsvc"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func (r *routeContext) handleLocationBooking(c *fiber.Ctx) error {
	locationID, _ := strconv.Atoi(c.Params("id"))
	currentUser := c.Locals("CurrentUser").(*authsvc.CurrentUser)

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	location := &Location{}
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", locationID).Find(&location).Error

	if err != nil {
		return errors.New("failed to book, try again later")
	}

	if location.AvailableSlot == 0 {
		tx.Rollback()
		return errors.New("no more available slot")
	}

	booking := &Booking{UserID: currentUser.ID}
	err = tx.Save(booking).Error
	if err != nil {
		tx.Rollback()
		return errors.New("only allowed to book once per location per user")
	}

	location.AvailableSlot = location.AvailableSlot - 1
	tx.Clauses(clause.Locking{Strength: "UPDATE"}).Save(location)
	c.SendStatus(http.StatusNoContent)

	return tx.Commit().Error
}
