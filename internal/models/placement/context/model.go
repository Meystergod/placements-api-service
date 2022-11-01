package context

import (
	partnerContext "github.com/Meystergod/placements-api-service/internal/models/partner/context"
)

type Context struct {
	IP        string `json:"ip" validate:"required,ip4_addr"`
	UserAgent string `json:"user_agent" validate:"required"`
}

func (c *Context) ToPartnerContext() partnerContext.Context {
	return partnerContext.Context{
		IP:        c.IP,
		UserAgent: c.UserAgent,
	}
}
