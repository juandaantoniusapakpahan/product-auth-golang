package entity

import (
	"encoding/json"
	"net/http"
	"product-auth/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

type Owner struct {
	UserId primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
}

type Resume struct {
	ID                 primitive.ObjectID  `bson:"_id,omitempty" json:"_id"`
	UserId             Owner               `bson:"owner" json:"owner"`
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

func AddResume(name, email string, payload []byte) (error, int) {
	rsm := new(Resume)
	err := json.Unmarshal(payload, &rsm)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	err, u := CheckUser(name, email)

	objID, err := primitive.ObjectIDFromHex(u.ID.Hex())
	if err != nil {
		return err, http.StatusInternalServerError
	}

	rsm.UserId = Owner{objID}
	db, err := database.Connect()
	if err != nil {
		return err, http.StatusInternalServerError
	}

	_, err = db.Collection("resumes").InsertOne(ctx, rsm)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusCreated
}

func EditResume(name, email, rsmId string, payload []byte) (int, error) {
	userPayload := new(Resume)
	resumeId, err := primitive.ObjectIDFromHex(rsmId)

	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = json.Unmarshal(payload, &userPayload)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	db, err := database.Connect()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err, user := CheckUser(name, email)
	if err != nil {
		return http.StatusNotFound, err
	}
	userId, err := primitive.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	filter := bson.M{}
	filter["_id"] = resumeId
	filter["owner._id"] = userId

	_, err = db.Collection("resumes").UpdateOne(ctx, filter, bson.M{"$set": bson.M{"name": userPayload.Name,
		"position": userPayload.Position, "profile": userPayload.Profile, "jobhistory": userPayload.JobHistory,
		"education_history": userPayload.EducationHistory, "internship_history": userPayload.InternshipsHistory,
		"detail_info": userPayload.DetailInfo, "media_sosial": userPayload.MediaSosial,
		"skill_history": userPayload.SkillsHistory, "hobbies": userPayload.Hobbies}})

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return 0, nil
}

func GetResume(name, email string) (*Resume, int, error) {
	err, user := CheckUser(name, email)
	if err != nil {
		return &Resume{}, http.StatusNotFound, err
	}

	userId, err := primitive.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		return &Resume{}, http.StatusInternalServerError, err
	}

	db, err := database.Connect()
	if err != nil {
		return &Resume{}, http.StatusInternalServerError, err
	}

	resultResume := new(Resume)
	filter := bson.M{"owner._id": userId}
	err = db.Collection("resumes").FindOne(ctx, filter).Decode(&resultResume)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &Resume{}, http.StatusNotFound, err
		}
		panic(err)
	}
	return resultResume, http.StatusOK, nil
}
