package guildsapi

import (
	svc "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
)

func modelGuildToResponse(m *model.Guild) *svc.Guild {
	return &svc.Guild{
		ID:                         m.ID.Hex(),
		Name:                       m.Name,
		Description:                m.Description,
		MasterAddress:              m.MasterAddress.String(),
		SpotBaseRequirement:        m.SpotBaseRequirement.String(),
		SpotQuoteRequirement:       m.SpotQuoteRequirement.String(),
		DerivativeQuoteRequirement: m.DerivativeQuoteRequirement.String(),
		StakingRequirement:         m.StakingRequirement.String(),
		Capacity:                   m.Capacity,
		MemberCount:                m.MemberCount,
	}
}
