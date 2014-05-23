package main

import (
	"encoding/json"
	"fmt"
	flickr "github.com/pboehm/go-flickr"
)

var UserPhotoCache = map[string]*UserData{}

type UserData struct {
	Username, NSID string
	Data           string
}

type NSIDResponse struct {
	User struct {
		Nsid string `json:"nsid"`
	} `json:"user"`
}

func (self *UserData) GetNSID() error {

	r := &flickr.Request{
		ApiKey: "122cc483be92bd806b696e7d458596ac",
		Method: "flickr.people.findByUsername",
		Args: map[string]string{
			"username": self.Username,
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
	self.NSID = res.User.Nsid

	return nil
}

type PhotosResponse struct {
	Photos struct {
		Page    int64  `json:"page"`
		Pages   int64  `json:"pages"`
		Perpage int64  `json:"perpage"`
		Total   string `json:"total"`
		Photo   []struct {
			Title     string `json:"title"`
			Datetaken string `json:"datetaken"`
			Views     string `json:"views"`
			Ownername string `json:"ownername"`
			UrlO      string `json:"url_o"`
			UrlZ      string `json:"url_z"`
		} `json:"photo"`
	} `json:"photos"`
}

func (self *UserData) GetPhotos() error {
	r := &flickr.Request{
		ApiKey: "122cc483be92bd806b696e7d458596ac",
		Method: "flickr.people.getPublicPhotos",
		Args: map[string]string{
			"user_id":  self.NSID,
			"extras":   "date_taken,owner_name,views,url_z,url_o",
			"format":   "json",
			"per_page": "50",
		},
	}

	resp, err := r.Execute()
	if err != nil {
		return err
	}

	self.Data = resp[14 : len(resp)-1]

	resp = resp[14 : len(resp)-1]

	var res PhotosResponse

	err = json.Unmarshal([]byte(resp), &res)

	for _, photo := range res.Photos.Photo {
		fmt.Printf("%+v\n", photo)
	}
	return nil
}

func GetUserPhotos(username string) (*UserData, error) {
	data := &UserData{Username: username}

	err := data.GetNSID()
	if err != nil {
		return nil, err
	}

	err = data.GetPhotos()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetRecentPhotos(username string) (string, error) {
	data, ok := UserPhotoCache[username]

	if !ok {
		fmt.Printf("Getting photo data for %s\n", username)
		newdata, err := GetUserPhotos(username)
		if err != nil {
			return "", err
		}

		UserPhotoCache[username] = newdata
		return newdata.Data, nil

	} else {
		return data.Data, nil
	}

}

func main() {
	GetRecentPhotos("phboehm")
	fmt.Println("\n\n\n")
	GetRecentPhotos("mperlet")
}
