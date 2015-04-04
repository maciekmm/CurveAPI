package hooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/maciekmm/curveapi/database"
	"github.com/maciekmm/curveapi/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var userIDRegex *regexp.Regexp = regexp.MustCompile("([0-9]+)")

func init() {
	database.GetDatabase().C("users").Create(&mgo.CollectionInfo{})
	database.GetDatabase().C("users").EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})
}

func loadIdFromName(name string) (int, error) {
	var err error
	res, err := http.Get("http://curvefever.com/users/" + name)
	if err != nil {
		return -1, err
	}
	if res.StatusCode == 200 && res.Header.Get("X-Yadis-Location") != "" {
		str := userIDRegex.FindString(res.Header.Get("X-Yadis-Location"))
		id, err := strconv.Atoi(str)
		return id, err
	} else if res.StatusCode == 404 {
		err = errors.New("User not found")
	}
	return -1, errors.New("Unexpected error")
}

func loadUserProfile(id int) (*models.Profile, error) {
	var err error
	res, err := http.Get("http://curvefever.com/achtung/user/" + strconv.Itoa(id) + "/json")
	if err != nil {
		return nil, err
	}

	var m models.Profile
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	//Error suppression
	//When User is unranked in certain Rank
	//It will throw marshal error, it would default to 0 anyway
	json.Unmarshal(buf.Bytes(), &m)
	if m.UID == 0 {
		err = errors.New("User not found")
		return nil, err
	}
	m.LastUpdate = time.Now().UTC().Unix()
	return &m, nil
}

// Gets user profile by its unique name
// As currently in CurveFever api there's no way to get player data by id directly, it's slower than using GetUserProfile
// Keep in mind it's not fetching data by player name, but rather by it's unique id visible on http://curvefever.com/users/{name}
// Fresh indicates whether to wait for new data if it's possible
// CurveFever api is updated at 6AM UTC
// If in need of the most actual data, use fresh=true otherwise when caring about efficiency use fresh=false
func GetUserProfileByName(id string, fresh bool) (*models.Profile, error) {
	profile, err := getCachedUserProfileByName(id)
	if err != nil {
		id, err := loadIdFromName(id)
		if err != nil {
			return nil, err
		}
		return GetUserProfile(id, fresh)
	}
	if err != nil {
		return nil, err
	}
	handleProfileSeek(profile, fresh)
	return profile, nil
}

// Gets user profile by its ID
// Fresh indicates whether to wait for new data if it's possible
// CurveFever api is updated at 6AM UTC
// If in need of the most actual data, use fresh=true otherwise when caring about efficiency use fresh=false
func GetUserProfile(id int, fresh bool) (*models.Profile, error) {
	profile, err := getCachedUserProfile(id)
	if err != nil || (fresh && profile.NeedsRefresh()) {
		profile, err = loadUserProfile(id)
		if err != nil {
			return nil, err
		}
		go func() {
			err = upsertUserProfile(profile)
		}()
	}
	if err != nil {
		return nil, err
	}
	handleProfileSeek(profile, fresh)
	return profile, nil
}

func handleProfileSeek(profile *models.Profile, fresh bool) {
	if !fresh && profile.NeedsRefresh() {
		go func() {
			profile, err := loadUserProfile(profile.UID)
			if err != nil {
				log.Fatalln(err)
				return
			}
			err = upsertUserProfile(profile)
			if err != nil {
				log.Fatalln(err)
			}
		}()
	}
}

func upsertUserProfile(profile *models.Profile) error {
	collection := database.GetDatabase().C("users")
	_, err := collection.UpsertId(profile.UID, profile)
	return err
}

func getCachedUserProfile(id int) (*models.Profile, error) {
	var profile models.Profile
	collection := database.GetDatabase().C("users")
	err := collection.FindId(id).One(&profile)
	return &profile, err
}

func getCachedUserProfileByName(id string) (*models.Profile, error) {
	var profile models.Profile
	collection := database.GetDatabase().C("users")
	err := collection.Find(bson.M{"name": id}).One(&profile)
	return &profile, err
}
