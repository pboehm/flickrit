package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/encoder"
	"github.com/pboehm/flickrit/api"
	"net/http"
)

func main() {
	api := &api.API{
		ApiKey: "122cc483be92bd806b696e7d458596ac",
		DataCacheRenewalInterval: 15,
	}
	api.Setup()
	go api.CyclicCacheRenewal()

	m := martini.Classic()
	m.Map(api)

	m.Use(cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowCredentials: true,
	}))

	m.Use(func(c martini.Context, w http.ResponseWriter) {
		c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})

	m.Get("/photos/:name", func(params martini.Params,
		enc encoder.Encoder) (int, []byte) {
		name, ok := params["name"]

		if ok {
			photos, _ := api.GetPhotosForUser(name)
			return http.StatusOK, encoder.Must(enc.Encode(photos))
		}

		return 404, []byte("Not Found")
	})
	m.Run()
}
