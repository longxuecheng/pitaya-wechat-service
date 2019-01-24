package response

type RegionDTO struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
	Type     int    `json:"type"`
	AgencyID int    `json:"agency_id"`
}
