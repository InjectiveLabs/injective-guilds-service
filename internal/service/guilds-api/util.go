package guildsapi

import (
	"strconv"

	svc "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
)

func modelGuildToResponse(m *model.Guild) *svc.Guild {
	return &svc.Guild{
		ID:                         m.ID.Hex(),
		Name:                       m.Name,
		Description:                m.Description,
		MasterAddress:              m.MasterAddress.String(),
		SpotBaseRequirement:        strconv.Itoa(m.SpotBaseRequirement),
		SpotQuoteRequirement:       strconv.Itoa(m.SpotQuoteRequirement),
		DerivativeQuoteRequirement: strconv.Itoa(m.DerivativeQuoteRequirement),
		StakingRequirement:         strconv.Itoa(m.StakingRequirement),
		Capacity:                   m.Capacity,
		MemberCount:                m.MemberCount,
	}
}
