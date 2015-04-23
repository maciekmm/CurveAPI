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
var httpClient *http.Client

func init() {
	database.GetDatabase().C("users").Create(&mgo.CollectionInfo{})
	database.GetDatabase().C("users").EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	})

	httpClient = &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
}

// Loads player profile
func loadProfile(url string) (*models.Profile, error) {
	var err error
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "CurveAPI")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("User not found")
	}

	defer res.Body.Close()

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

func getUserProfile(id string, fresh bool, url string, getFromCache func(string) (*models.Profile, error)) (*models.Profile, error) {
	// Get user profile from database
	profile, err := getFromCache(id)
	// Check if it failed or if profile needs refresh while requesting most recent data
	if err != nil || (fresh && profile.NeedsRefresh()) {
		// Tries to load user profile from CurveFever site
		profile, err = loadProfile(url)
		if err != nil {
			return nil, err
		}
		go func() {
			err = upsertUserProfile(profile)
			if err != nil {
				log.Fatalln(err)
			}
		}()
		return profile, err
	}

	// If the requested data does not have to be fresh but it should be refreshed do it in goroutine
	if !fresh && profile.NeedsRefresh() {
		go func() {
			profile, err := loadProfile(url)
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
	return profile, nil
}

// Gets user profile by its unique name
// As currently in CurveFever api there's no way to get player data by id directly, it's slower than using GetUserProfile
// Keep in mind it's not fetching data by player name, but rather by it's unique id visible on http://curvefever.com/users/{name}
// Fresh indicates whether to wait for new data if it's possible
// CurveFever api is updated at 6AM UTC
// If in need of the most actual data, use fresh=true otherwise when caring about efficiency use fresh=false
func GetUserProfileByName(id string, fresh bool) (*models.Profile, error) {
	return getUserProfile(id, fresh, "http://curvefever.com/achtung/username/"+id+"/json", getCachedUserProfileByName)
}

// Gets user profile by its ID
// Fresh indicates whether to wait for new data if it's possible
// CurveFever api is updated at 6AM UTC
// If in need of the most actual data, use fresh=true otherwise when caring about efficiency use fresh=false
func GetUserProfile(id int, fresh bool) (*models.Profile, error) {
	return getUserProfile(strconv.Itoa(id), fresh, "http://curvefever.com/achtung/user/"+strconv.Itoa(id)+"/json", getCachedUserProfile)
}

func upsertUserProfile(profile *models.Profile) error {
	collection := database.GetDatabase().C("users")
	_, err := collection.UpsertId(profile.UID, profile)
	return err
}

func getCachedUserProfile(id string) (*models.Profile, error) {
	numericId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	var profile models.Profile
	collection := database.GetDatabase().C("users")
	err = collection.FindId(numericId).One(&profile)
	return &profile, err
}

func getCachedUserProfileByName(id string) (*models.Profile, error) {
	var profile models.Profile
	collection := database.GetDatabase().C("users")
	err := collection.Find(bson.M{"name": id}).One(&profile)
	return &profile, err
}
