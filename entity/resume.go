package entity

import (
	"gopkg.in/mgo.v2/bson"
)

type EmploymentHistory struct {
	Position       string `bson:"position"`
	CompanyName    string `bson:"company_name"`
	ComanyLocation string `bson:"company_location"`
	Discription    string `bson:"discription"`
	// StartedAt      time.Time `bson:"started_at"`
	// EndedAt        time.Time `bson:"ended_at"`
}

type Education struct {
	EducationName       string `bson:"education_name"`
	InstitutionName     string `bson:"institution_name"`
	InstitutionLocation string `bson:"institution_location"`
	Description         string `bson:"description"`
	// StartedAt           time.Time `bson:"started_at"`
	// EndedAt             time.Time `bson:"ended_at"`
}

type Internships struct {
	PositionName        string `bson:"position_name"`
	InstitutionName     string `bson:"institution_name"`
	InstitutionLocation string `bson:"institution_locations"`
	Description         string `bson:"description"`
	// StartedAt      time.Time `bson:"started_at"`
	// EndedAt        time.Time `bson:"ended_at"`
}

type Medsos struct {
	MedsosName string `bson:"medsos_name"`
	Url        string `bson:"url"`
}

type Skills struct {
	NameSkil string  `bson:"name_skil"`
	Level    float32 `bson:"level"`
}

type Details struct {
	City    string `bson:"city"`
	Country string `bson:"country"`
	Phone   string `bson:"phone"`
	Email   string `bson:"email"`
}

type Resume struct {
	ID                 bson.ObjectId `bson:"_id, omitempty"`
	UserId             bson.ObjectId `bson:"_id, omitempty"`
	Name               string        `json:"name" bson:"name"`
	Position           string        `json:"position" bson:"position"`
	Profile            string        `json:"profile" bson:"profile"`
	JobHistory         []EmploymentHistory
	EducationHistory   []Education
	InternshipsHistory []Internships
	DetailInfo         Details
	MediaSosial        []Medsos
	SkillsHistory      []Skills
	Hobbies            []string
}
