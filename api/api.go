package api

import (
	"encoding/json"
	flickr "github.com/pboehm/go-flickr"
	"sync"
	"time"
)

type UserData struct {
	Username, NSID string
	Photos         []FlickrPhoto
	Ttl            int
}

type API struct {
	ApiKey     string
	PhotoCache map[string]*UserData
	Mutex      sync.RWMutex

	UserDataDefaultTTL   int
	WatchdogTTLDecrement int
	WatchdogInterval     int
}

func (self *API) Setup() {
	if self.PhotoCache == nil {
		self.PhotoCache = map[string]*UserData{}
	}
}

func (self *API) GetPhotosForUser(username string) ([]FlickrPhoto, error) {

	self.Mutex.RLock()
	data, ok := self.PhotoCache[username]
	self.Mutex.RUnlock()

	if !ok {
		newdata, err := self.getUserData(username)
		if err != nil {
			return nil, err
		}

		self.Mutex.Lock()
		self.PhotoCache[username] = newdata
		self.Mutex.Unlock()

		return newdata.Photos, nil

	} else {
		return data.Photos, nil
	}
}

func (self *API) getUserData(username string) (*UserData, error) {
	data := &UserData{Username: username, Ttl: self.UserDataDefaultTTL}

	err := self.setNSID(data)
	if err != nil {
		return nil, err
	}

	data.Photos = self.getPhotos(data)

	return data, nil
}

func (self *API) setNSID(data *UserData) error {

	r := &flickr.Request{
		ApiKey: self.ApiKey,
		Method: "flickr.people.findByUsername",
		Args: map[string]string{
			"username": data.Username,
			"format":   "json",
		},
	}

	resp, err := r.Execute()
	if err != nil {
		return err
	}
	resp = resp[14 : len(resp)-1]

	var res NSIDResponse

	err = json.Unmarshal([]byte(resp), &res)
	data.NSID = res.User.Nsid

	return nil
}

func (self *API) getPhotos(data *UserData) []FlickrPhoto {

	r := &flickr.Request{
		ApiKey: self.ApiKey,
		Method: "flickr.people.getPublicPhotos",
		Args: map[string]string{
			"user_id":  data.NSID,
			"extras":   "date_taken,owner_name,views,url_z,url_o,count_faves,count_comments",
			"format":   "json",
			"per_page": "500",
		},
	}

	resp, err := r.Execute()
	if err != nil {
		return nil
	}
	resp = resp[14 : len(resp)-1]

	var res PhotosResponse
	err = json.Unmarshal([]byte(resp), &res)
	if err != nil {
		return nil
	}

	for i := 0; i < len(res.Photos.Photo); i++ {
		res.Photos.Photo[i].GenerateExtraMembers()
	}

	return res.Photos.Photo
}

func (self *API) CacheCleanup() {

	for {
		self.Mutex.Lock()
		for user, data := range self.PhotoCache {

			data.Ttl -= self.WatchdogTTLDecrement
			if data.Ttl <= 0 {
				delete(self.PhotoCache, user)
			}
		}
		self.Mutex.Unlock()

		time.Sleep(time.Duration(self.WatchdogInterval) * time.Second)
	}
}
