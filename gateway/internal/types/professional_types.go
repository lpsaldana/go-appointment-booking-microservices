package types

type CreateProfessionalRequest struct {
	Name       string `json:"name"`
	Profession string `json:"profession"`
	Contact    string `json:"contact"`
}
