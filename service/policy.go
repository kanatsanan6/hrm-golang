package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kanatsanan6/hrm/model"
)

var policies = map[string]map[string][]string{
	"user_management": {
		"read":   {"admin"},
		"invite": {"admin"},
		"delete": {"admin"},
	},
	"leave": {
		"read":    {"admin"},
		"create":  {"admin", "member"},
		"approve": {"admin"},
	},
}

type PolicyInterface interface {
	Authorize(c *fiber.Ctx, subject string, action string) bool
	Export(c *fiber.Ctx) []map[string]string
}

type Policy struct{}

func NewPolicy() PolicyInterface {
	return &Policy{}
}

func (p *Policy) Authorize(c *fiber.Ctx, subject string, action string) bool {
	user := c.Locals("user").(model.User)

	policySubject, ok := policies[subject]
	if !ok {
		return false
	}

	policyAction, ok := policySubject[action]
	if !ok {
		return false
	}

	return contains(policyAction, user.Role)
}

type PolicyType struct {
	Subject string `json:"subject"`
	Action  string `json:"action"`
}

func (p *Policy) Export(c *fiber.Ctx) []map[string]string {
	user := c.Locals("user").(model.User)
	result := []map[string]string{}

	for subject, actions := range policies {
		for action, roles := range actions {
			for _, role := range roles {
				if role == user.Role {
					obj := map[string]string{"subject": subject, "action": action}
					result = append(result, obj)
					break
				}
			}
		}
	}

	return result
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
