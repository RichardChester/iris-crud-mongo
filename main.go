package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//candidate type
type Candidate struct {
	ID         int           `json:"id" binding:"required"`
	Name       *PartFullName `json:"fullname"`
	Grade      string        `json:"grade" binding:"required"`
	Sector     string        `json:"sector"`
	Business   string        `json:"business"`
	Function   string        `json:"function"`
	Location   string        `json:"location"`
	Sponsor    *SponFullName `json:"fullname"`
	Duration   string        `json:"duration"`
	Background *background   `json:"background`
	WantedXP   []wantxp      `json:"wantedxp"`
	BroughtXP  []bringxp     `json:"braoughtxp"`
	HowTrans   string        `json:"howtrans"`
	CrossFun   []cross       `json:"crossfun"`
	State      string        `json:"state"`
}
type PartFullName struct {
	Firstname  string `json:"firstname"`
	Secondname string `json:"secondname"`
}
type SponFullName struct {
	Firstname  string `json:"firstname"`
	Secondname string `json:"secondname"`
}
type background struct {
	Summary    string   `json:"summary"`
	Experties  []string `json:"experties"`
	Experience []exp    `json:"experience"`
	Projects   []string `json:"projects"`
	Values     []val    `json:"values"`
}
type exp struct {
	Place    string `json:"place"`
	DateFrom string `json:"datefrom"`
	DateTo   string `json:"dateto"`
}
type val struct {
	Title string `json:"title"`
	Fluf  string `json:"fluf"`
}
type wantxp struct {
	Title string `json:"title"`
	Fluf  string `json:"fluf"`
}
type bringxp struct {
	Title string `json:"title"`
	Fluf  string `json:"fluf"`
}
type cross struct {
	Title string `json:"title"`
	Fluf  string `json:"fluf"`
}
type Placement struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Host      string `json:"host"`
	PSponsor  string `json:"psponsor"`
	SSponsor  string `json:"ssponsor"`
	Sector    string `json:"sector"`
	Barea     string `json:"barea"`
	JobFunc   string `json:"jobfunc"`
	Grade     string `json:"grade"`
	Location  string `json:"location"`
	Length    string `json:"length"`
	OppOver   string `json:"oppover"`
	SecPri    string `json:"secpri"`
	CritXP    string `json:"critxp"`
	KnowShare string `json:"knowshare"`
	KnowTran  string `json:"knowtran"`
	CrosFun   string `json:"crosfun"`
	State     string `json:"crosfun"`
}

func main() {
	app := iris.New()

	app.Logger().SetLevel("debug")
	// Recover from panics and log the panic message to the application's logger ("Warn" level).
	app.Use(recover.New())
	// logs HTTP requests to the application's logger ("Info" level)
	app.Use(logger.New())

	// Connection variables two hosts provided for local and cloud
	const (
		//Host       = "localhost:27017"
		Host          = "mongodb://<dbuser>:<dbpassword>@ds032887.mlab.com:"
		Username      = ""
		Password      = ""
		Database      = "ignite"
		CollectionCan = "Candidate"
		CollectionJob = "Placement"
	)

	// Mongo connection
	session, err := mgo.Dial(Host)
	if nil != err {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db := session.DB(Database)
	collectionC := db.C(CollectionCan)
	collectionJ := db.C(CollectionJob)

	// Create todo using POST Request body
	app.Post("/addCandi", func(ctx iris.Context) {
		// Create a new ToDo
		var candi Candidate
		// Pass the pointer of todo so it is updated with the result
		// which is the POST data
		err := ctx.ReadJSON(&candi)
		// If there is an error or no Title in the POST Body
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"message": "There was a candidate submit error"})
			return
		}
		collectionC.Insert(candi)
	})

	//create new placement
	app.Post("/addPlace", func(ctx iris.Context) {
		var Job Placement
		err := ctx.ReadJSON(&Job)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"message": "There was a placement submit error"})
			return
		}
		collectionJ.Insert(Job)
	})

	// Applicant by etc
	app.Get("/App/{WhatFieldToSearch:string}", func(ctx iris.Context) {
		title := ctx.Params().Get("WhatFieldToSearch")
		var results []Candidate
		err := collectionC.Find(bson.M{"WhatFieldToSearch": title}).All(&results)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"message": "An error occured", "error": err})
			return
		}
		ctx.JSON(iris.Map{"results": results})
	})

	// Placement by etc
	app.Get("/App/{WhatFieldToSearch:string}", func(ctx iris.Context) {
		title := ctx.Params().Get("WhatFieldToSearch")
		var results []Candidate
		err := collectionC.Find(bson.M{"WhatFieldToSearch": title}).All(&results)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"message": "An error occured", "error": err})
			return
		}
		ctx.JSON(iris.Map{"results": results})
	})

	// Get multi (cirrently set to get all state 2 apps and candis)
	app.Get("/admin/", func(ctx iris.Context) {
		var resultsCan []Candidate
		var resultsJob []Placement
		err := collectionC.Find(bson.M{"state": "2"}).All(&resultsCan)
		err1 := collectionJ.Find(bson.M{"state": "2"}).All(&resultsJob)
		if err != nil && err1 != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"message": "An error occured", "error": err})
			return
		}
		ctx.JSON(iris.Map{"results1": resultsCan})
		ctx.JSON(iris.Map{"results2": resultsJob})
	})

	// Remove Candidate
	app.Delete("/Candidate/{ID:int}", func(ctx iris.Context) {
		ID := ctx.Params().Get("ID")
		collectionC.Remove(bson.M{"ID": ID})
	})
	// Remove by Placement
	app.Delete("/Placement/{ID:int}", func(ctx iris.Context) {
		ID, _ := ctx.Params().GetBool("ID")
		collectionJ.RemoveAll(bson.M{"ID": ID})
	})

	// Run app on port 8080
	// ignore server closed errors ([ERRO] 2018/04/09 12:25 http: Server closed)
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
