package api

type NSIDResponse struct {
	User struct {
		Nsid string `json:"nsid"`
	} `json:"user"`
}

type PhotosResponse struct {
	Photos struct {
		Page    int64  `json:"page"`
		Pages   int64  `json:"pages"`
		Perpage int64  `json:"perpage"`
		Total   string `json:"total"`
		Photo   []struct {
			Title        string `json:"title"`
			Datetaken    string `json:"datetaken"`
			Views        string `json:"views"`
			Ownername    string `json:"ownername"`
			UrlO         string `json:"url_o"`
			UrlZ         string `json:"url_z"`
			FavCount     string `json:"count_faves"`
			CommentCount string `json:"count_comments"`
		} `json:"photo"`
	} `json:"photos"`
}
