package api

import (
	"fmt"
)

type NSIDResponse struct {
	User struct {
		Nsid string `json:"nsid"`
	} `json:"user"`
}

type FlickrPhoto struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Datetaken    string `json:"datetaken"`
	Views        string `json:"views"`
	Owner        string `json:"owner"`
	Ownername    string `json:"ownername"`
	UrlO         string `json:"url_o"`
	HeightO      string `json:"height_o"`
	WidthO       string `json:"width_o"`
	UrlZ         string `json:"url_z"`
	HeightZ      string `json:"height_z"`
	WidthZ       string `json:"width_z"`
	FavCount     string `json:"count_faves"`
	CommentCount string `json:"count_comments"`
	PhotoUrl     string `json:"photo_url"`
}

func (self *FlickrPhoto) GenerateExtraMembers() {
	self.PhotoUrl = fmt.Sprintf("http://www.flickr.com/photos/%s/%s",
		self.Owner, self.Id)
}

type PhotosResponse struct {
	Photos struct {
		Page    int64         `json:"page"`
		Pages   int64         `json:"pages"`
		Perpage int64         `json:"perpage"`
		Total   string        `json:"total"`
		Photo   []FlickrPhoto `json:"photo"`
	} `json:"photos"`
}
