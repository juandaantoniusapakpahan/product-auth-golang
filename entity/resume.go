package entity

import (
	"gopkg.in/mgo.v2/bson"
)

type EmploymentHistory struct {
	Position       string `json:"position" bson:"position"`
	CompanyName    string `json:"company_name" bson:"company_name"`
	ComanyLocation string `json:"company_location" bson:"company_location"`
	Discription    string `json:"discription" bson:"discription"`
	// StartedAt      time.Time `bson:"started_at"`
	// EndedAt        time.Time `bson:"ended_at"`
}

type Education struct {
	EducationName       string `json:"education_name" bson:"education_name"`
	InstitutionName     string `json:"institution_name" bson:"institution_name"`
	InstitutionLocation string `json:"institution_location" bson:"institution_location"`
	Description         string `json:"description" bson:"description"`
	// StartedAt           time.Time `bson:"started_at"`
	// EndedAt             time.Time `bson:"ended_at"`
}

type Internships struct {
	PositionName        string `json:"position_name" bson:"position_name"`
	InstitutionName     string `json:"institution_name" bson:"institution_name"`
	InstitutionLocation string `json:"institution_locations" bson:"institution_locations"`
	Description         string `json:"description" bson:"description"`
	// StartedAt      time.Time `bson:"started_at"`
	// EndedAt        time.Time `bson:"ended_at"`
}

type Medsos struct {
	MedsosName string `json:"medsos_name" bson:"medsos_name"`
	Url        string `json:"url" bson:"url"`
}

type Skills struct {
	NameSkil string  `json:"name_skil" bson:"name_skil"`
	Level    float32 `json:"level" bson:"level"`
}

type Details struct {
	City    string `json:"city" bson:"city"`
	Country string `json:"country" bson:"country"`
	Phone   string `json:"phone" bson:"phone"`
	Email   string `json:"email" bson:"email"`
}

type Resume struct {
	ID                 bson.ObjectId       `bson:"_id,omitempty" json:"_id"`
	UserId             bson.ObjectId       `bson:"user_id,omitempty" json:"user_id"`
	Name               string              `json:"name" bson:"name"`
	Position           string              `json:"position" bson:"position"`
	Profile            string              `json:"profile" bson:"profile"`
	JobHistory         []EmploymentHistory `json:"jobhistory" bson:"jobhistory"`
	EducationHistory   []Education         `json:"education_history" bson:"education_history"`
	InternshipsHistory []Internships       `json:"internship_history" bson:"internship_history"`
	DetailInfo         Details             `json:"detail_info" bson:"detail_info"`
	MediaSosial        []Medsos            `json:"media_sosial" bson:"media_sosial"`
	SkillsHistory      []Skills            `json:"skill_history" bson:"skill_history"`
	Hobbies            []string            `json:"hobbies" bson:"hobbies"`
}
