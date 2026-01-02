package supply

type SupplyResponse struct {
	Data    SupplyData `json:"data"`
	Message string     `json:"message"`
	Status  int        `json:"status"`
}

type SupplyData struct {
	Name                   string            `json:"name"`
	NameV1                 map[string]string `json:"name_v1"`
	ImageURL               string            `json:"image_url"`
	Languages              []string          `json:"languages"`
	Specializations        []string          `json:"specializations"`
	IsActive               bool              `json:"is_active"`
	YearsOfExperience      int               `json:"years_of_experience"`
	ImageWithoutBackground string            `json:"image_without_background"`
}
