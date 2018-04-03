package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	//"encoding/base64"
	"strings"
	//"io/ioutil"
	"strconv"
	//"bytes"
)

type immobilien struct {
	Angebot_titel     string  `json:"angebot_titel"`
	Id                int     `json:"id"`
	Angebot_id        int     `json:"angebot_id"`
	Kaufpreis         float32 `json:"kaufpreis"`
	Immobilien_adress string  `json:"immobilien_adress"`
	Qm                int     `json:"qm"`
	Zimmeranzahl      int     `json:"zimmeranzahl"`
	Ortname           string  `json:"ortname"`
	Plz               string  `json:"plz"`
	Beschreibung      string  `json:"beschreibung"`
	Media             string  `json:"media"`
}

type angebot struct {
	offer_title     string
	offer_type      string
	offer_startdate string
	offer_enddate   string
	offer_rpricing  int
	offer_bpricing  int
	offer_utilcosts int
	offer_bail      int
	offer_provision float64
	agent_nick      string
	agent_firstname string
	agent_lastname  string
	agent_agency    string
	immo_desc       string
	immo_features   string
	immo_address    string
	city_name       string
	city_plz        string
	immo_type       string
	immo_floors     int
	immo_rooms      string
	immo_condition  string
	immo_qm         int
	immo_heating    string
	immo_baujahr    int
	immo_media      string
}

type offer struct {
	created string `json:"created"`
}

func filter(res http.ResponseWriter, req *http.Request) {

	dsn := "root:@tcp(localhost:3306)/fa17g17"
	db, err := sql.Open("mysql", dsn)
	checkErr(err)

	params := req.URL.Query()

	var data []immobilien

	// query
	var sql = "SELECT a.angebot_titel, i.id, a.id as angebot_id, a.kaufpreis, i.immobilien_adress, i.qm, i.zimmeranzahl, o.ortname, o.plz, i.beschreibung, i.media"
	sql += " FROM fa17g17.immobilien i inner join fa17g17.ort o on immobilien_ort = o.id inner join fa17g17.angebot a on i.id = a.immobilien_id inner join fa17g17.user u on u.id = i.verkaufer_id"
	sql += " WHERE a.angebots_art like ('%" + params.Get("angebot_art") + "%') and o.ortname like ('%" + params.Get("ort") + "%') and i.immobilien_art like ('%" + params.Get("objektart") + "%')"
	sql += " and i.qm <= " + params.Get("qm") + " and a.mietenpreis <= " + params.Get("preis") + " and i.zimmeranzahl >= " + params.Get("zimmeranzahl") + " and i.active=1"
	sql += " and i.qm between " + params.Get("minQm") + " and " + params.Get("maxQm") + " and a.mietenpreis between " + params.Get("minPreis") + " and "
	sql += params.Get("maxPreis") + " and i.baujahr >= " + params.Get("bauJahr") + " and i.heizungs_art like('%" + params.Get("heizungs_art") + "%') and u.agency like('%"
	sql += params.Get("agencyID") + "%')"

	//fmt.Println (sql)

	rows, err := db.Query(sql)
	checkErr(err)

	var record immobilien
	for rows.Next() {
		var Angebot_titel string
		var Id int
		var Angebot_id int
		var Kaufpreis float32
		var Immobilien_adress string
		var Qm int
		var Zimmeranzahl int
		var Ortname string
		var Plz string
		var Beschreibung string
		var Media string

		err = rows.Scan(&Angebot_titel, &Id, &Angebot_id, &Kaufpreis, &Immobilien_adress, &Qm, &Zimmeranzahl, &Ortname, &Plz, &Beschreibung, &Media)
		checkErr(err)

		record.Angebot_titel = Angebot_titel
		record.Id = Id
		record.Angebot_id = Angebot_id
		record.Kaufpreis = Kaufpreis
		record.Immobilien_adress = Immobilien_adress
		record.Qm = Qm
		record.Zimmeranzahl = Zimmeranzahl
		record.Ortname = Ortname
		record.Plz = Plz
		record.Beschreibung = Beschreibung
		record.Media = Media

		data = append(data, record)
	}

	db.Close()

	//fmt.Println(data)
	var Response, _ = json.Marshal(data)
	//fmt.Println(string(Response))

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:17017")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	res.Write(Response)
}

func search(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Search Function begin!")
	dsn := "root:@tcp(localhost:3306)/fa17g17"
	db, err := sql.Open("mysql", dsn)
	checkErr(err)

	params := req.URL.Query()

	var data []immobilien

	// query
	var sql = "SELECT a.angebot_titel, i.id, a.id as angebot_id, a.kaufpreis, i.immobilien_adress, i.qm, i.zimmeranzahl, o.ortname, o.plz, i.beschreibung, i.media"
	sql += " FROM fa17g17.immobilien i inner join fa17g17.ort o on immobilien_ort = o.id inner join fa17g17.angebot a on i.id = a.immobilien_id"
	sql += " WHERE a.angebots_art like ('%" + params.Get("angebot_art") + "%') and o.ortname like ('%" + params.Get("ort") + "%') and i.immobilien_art like ('%" + params.Get("objektart") + "%')"
	sql += " and i.qm <= " + params.Get("qm") + " and a.mietenpreis <= " + params.Get("preis") + " and i.zimmeranzahl >= " + params.Get("zimmeranzahl") + " and i.active=1"

	fmt.Println(sql)
	rows, err := db.Query(sql)
	checkErr(err)

	var record immobilien
	for rows.Next() {
		var Angebot_titel string
		var Id int
		var Angebot_id int
		var Kaufpreis float32
		var Immobilien_adress string
		var Qm int
		var Zimmeranzahl int
		var Ortname string
		var Plz string
		var Beschreibung string
		var Media string

		err = rows.Scan(&Angebot_titel, &Id, &Angebot_id, &Kaufpreis, &Immobilien_adress, &Qm, &Zimmeranzahl, &Ortname, &Plz, &Beschreibung, &Media)
		checkErr(err)

		record.Angebot_titel = Angebot_titel
		record.Id = Id
		record.Angebot_id = Angebot_id
		record.Kaufpreis = Kaufpreis
		record.Immobilien_adress = Immobilien_adress
		record.Qm = Qm
		record.Zimmeranzahl = Zimmeranzahl
		record.Ortname = Ortname
		record.Plz = Plz
		record.Beschreibung = Beschreibung
		record.Media = Media

		data = append(data, record)
	}

	db.Close()
	fmt.Println("Search Function end!")
	var Response, _ = json.Marshal(data)
	fmt.Println(string(Response))

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:17017")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	res.Write(Response)
}

func estateCreate(res http.ResponseWriter, req *http.Request) {

	fmt.Println(req)

	params := req.URL.Query()

	//fmt.Println(params)
	// data
	var ownerID = params.Get("ownerID")
	var agentID = params.Get("agentID")
	//2 - Wenn User normaler Käufer ist und Anfrage stellt (active 0)
	//3 - Wenn User Makler ist und Immobilie einstellen will (active 1)
	var request = params.Get("request")
	var title = params.Get("title")
	var description = params.Get("desc")
	var condition = params.Get("condition")
	var estateType = params.Get("estateType")
	var heatingType = params.Get("heatingType")
	var baujahr = params.Get("baujahr")
	//String Array
	var features = params.Get("features")
	var address = params.Get("address")
	var postalcode = params.Get("postal")
	var city = params.Get("city")
	var floors = params.Get("floors")
	var rooms = params.Get("rooms")
	//Quadratmeter
	var size = params.Get("size")
	var offerType = params.Get("offerType")
	var offerPurpose = params.Get("offerPurpose")
	var price = params.Get("price")
	//Kaution
	var bail = params.Get("bail")
	var provision = params.Get("provision")
	//Nebenkosten
	var utilities = params.Get("utilities")
	var startdate = params.Get("startdate")
	var enddate = params.Get("enddate")
	//Dateien
	//var hasFiles 		= params.Get("HasFiles")
	//var isLocal			= params.Get("IsLocal")
	var attachments = params.Get("Attachments")
	//var fileStrings []string

	//fmt.Println ("ket qua 1 o day!")
	//fmt.Println (attachments)
	//fmt.Println ("ket qua 2 o day!")

	//for i:=0; i<=len(attachments); i++{
	//data, err := ioutil.ReadFile(attachments[i])
	//	fmt.Println (attachments[i])
	//	fmt.Println ("ket qua 3 o day!")
	//}

	if req.Method == "GET" {
		//query
		dsn := "root:@tcp(localhost:3306)/fa17g17"
		db, err := sql.Open("mysql", dsn)
		checkErr(err)

		//reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(attachments))

		//file, err := os.Create("http://localhost:17017/fa17g17/img/go.jpg")
		//checkErr(err)

		//_, err = io.Copy(file, reader)
		//checkErr(err)

		var checkCityQuery = "Select id from ort where plz ='" + postalcode + "' and ortname='" + city + "'"
		fmt.Println(checkCityQuery)

		var insert int = 1
		var cityID int = 0

		rows, err := db.Query(checkCityQuery)
		checkErr(err)

		for rows.Next() {
			err = rows.Scan(&cityID)
			fmt.Println(cityID)
			checkErr(err)
			if cityID != 0 {
				insert = 0
			}
		}

		if insert == 1 {
			var insertCityQuery = "INSERT INTO ort (ortname, beschreibung, plz) VALUES ('" + city + "','---','" + postalcode + "')"
			fmt.Println(insertCityQuery)
			insertCity, err := db.Query(insertCityQuery)
			checkErr(err)
			insertCity.Close()
		}

		var active = 0
		if request == "3" {
			active = 1
		} else {
			active = 0
		}

		var insertImmoQuery = "INSERT immobilien SET immobilien_art=?, immobilien_adress=?, immobilien_ort=?, zimmeranzahl=?, "
		insertImmoQuery += "qm=?, etage_anzahl=?, bauzustand=?, features=?, nutzungszweck=?, beschreibung=?, verkaufer_id=?, "
		insertImmoQuery += "besitzer_id=?, active=?, heizungs_art=?, baujahr=?, media=?"

		//fmt.Println(insertImmoQuery)

		stmtImmo, err := db.Prepare(insertImmoQuery)
		checkErr(err)

		resImmo, err := stmtImmo.Exec(estateType, address, cityID, rooms, size, floors, condition, features, offerPurpose, description, agentID, ownerID, active, heatingType, baujahr, attachments)
		checkErr(err)

		immoID, err := resImmo.LastInsertId()
		checkErr(err)

		//fmt.Println(immoID)

		var insertOfferQuery = "INSERT angebot SET verkaufer_id=?, immobilien_id=?, angebots_art=?, kaution=?, nebenkosten=?, mietenpreis=?,"
		insertOfferQuery += " kaufpreis=?, provision=?, startdate=?, enddate=?, angebot_titel=?, active=?"

		//fmt.Println(insertOfferQuery)

		stmtOffer, err := db.Prepare(insertOfferQuery)
		checkErr(err)

		resOffer, err := stmtOffer.Exec(agentID, immoID, offerType, bail, utilities, price, price, provision, startdate, enddate, title, active)
		checkErr(err)

		offerID, err := resOffer.LastInsertId()
		checkErr(err)

		fmt.Println(offerID)
		db.Close()
	}

	var result offer
	result.created = "true"
	fmt.Println("result is : ", result)

	var Response, _ = json.Marshal(result)
	fmt.Println(Response)

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	res.Write(Response)
}

func showData(res http.ResponseWriter, req *http.Request) {

	dsn := "root:@tcp(localhost:3306)/fa17g17"
	db, err := sql.Open("mysql", dsn)
	checkErr(err)

	params := req.URL.Query()

	var immoID = params.Get("immoID")
	fmt.Println(params)

	var selectImmoQuery = "select angebot_titel, a.angebots_art, a.startdate, a.enddate, a.mietenpreis, a.kaufpreis, a.nebenkosten, a.kaution, a.provision,"
	selectImmoQuery += " u.nickname, u.firstname, u.lastname, u.agency, i.beschreibung, i.features, i.immobilien_adress, o.ortname, o.plz,"
	selectImmoQuery += " i.immobilien_art, i.etage_anzahl, i.zimmeranzahl, i.bauzustand, i.qm, i.heizungs_art, i.baujahr, i.media"
	selectImmoQuery += " from immobilien i inner join angebot a on i.id = a.immobilien_id inner join ort o on o.id = i.immobilien_ort"
	selectImmoQuery += " inner join user u on u.id = i.verkaufer_id where i.id = " + immoID

	fmt.Println(selectImmoQuery)

	rows, err := db.Query(selectImmoQuery)
	checkErr(err)

	var A angebot

	for rows.Next() {
		err = rows.Scan(&A.offer_title, &A.offer_type, &A.offer_startdate, &A.offer_enddate, &A.offer_rpricing, &A.offer_bpricing, &A.offer_utilcosts, &A.offer_bail, &A.offer_provision, &A.agent_nick, &A.agent_firstname, &A.agent_lastname, &A.agent_agency, &A.immo_desc, &A.immo_features, &A.immo_address, &A.city_name, &A.city_plz, &A.immo_type, &A.immo_floors, &A.immo_rooms, &A.immo_condition, &A.immo_qm, &A.immo_heating, &A.immo_baujahr, &A.immo_media)
		checkErr(err)
	}

	var ratingImmoQuery = "SELECT count(k.id) as numberofcomment, sum(wert_id) as rating"
	ratingImmoQuery += " FROM fa17g17.kommentar k inner join fa17g17.angebot a on k.angebot_id = a.id"
	ratingImmoQuery += " group by a.immobilien_id"
	ratingImmoQuery += " having a.immobilien_id =" + immoID

	row1s, err := db.Query(ratingImmoQuery)
	checkErr(err)

	var rating int

	for row1s.Next() {
		var numRating int
		var value int
		err = row1s.Scan(&numRating, &value)
		checkErr(err)
		rating = value / numRating
	}

	db.Close()

	if rating == 0 {
	}

	var Price int
	if A.offer_type == "Rent" {
		Price = A.offer_rpricing
	} else {
		Price = A.offer_bpricing
	}
	//fmt.Println(A)

	var immo_image string = ""
	var immo_gmap string = "https://www.google.com/maps/embed/v1/place?key=AIzaSyAv_dsBOBDTvbpYBZOr_Lbh8BZPIJN8dhc&q="

	immo_gmap += A.immo_address + "+" + A.city_plz + "+" + A.city_name

	if A.immo_media != "" {
		s := strings.Split(A.immo_media, ",")
		for i := 0; i < len(s); i++ {
			if i == 0 {
				immo_image += "<div class='item active'> <img src='http://localhost:17017/fa17g17/img/" + s[i] + "' class='img-responsive' style 'width=100%' alt='Image'></div>"
			} else {
				immo_image += "<div class='item'> <img src='http://localhost:17017/fa17g17/img/" + s[i] + "' class='img-responsive' style 'width=100%' alt='Image'></div>"
			}
		}
	}

	var page string
	page = "<!DOCTYPE html>"
	page += "<html lang='en'>"
	page += "<head>"
	page += "<title> Angebot Profile </title>"
	page += "<meta charset='utf-8'>"
	page += "<meta name='viewport' content='width=device-width, initial-scale=1'>"
	page += "<link rel='stylesheet' href='https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css'>"
	page += "<link rel='stylesheet' href='https://netdna.bootstrapcdn.com/font-awesome/3.2.1/css/font-awesome.css'>"
	page += "<script src='https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js'></script>"
	page += "<script src='https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js'></script>"
	page += "<script src='https://cdn.jsdelivr.net/npm/js-cookie@2/src/js.cookie.min.js'></script>"
	page += "<style>"
	page += "#comment { resize: none;}"
	page += ".row.content {height: 550px}"
	page += ".sidenav { background-color: #f1f1f1; height: 100%;}"
	page += "@media screen and (max-width: 767px) { .row.content {height: auto;}}"
	page += ".carousel-control.right, .carousel-control.left { background-image: none; color: #f4511e;}"
	page += ".well { background-color: #333; border-color: #222;}"
	page += ".put-back { z-index:-1 !important; positon: relative !important;}"
	page += " </style>"
	page += "</head>"
	page += "<body style='"
	page += "background-image: url(http://localhost:17017/fa17g17/img/background.jpeg);"
	page += "background-image-size: auto;"
	page += "color: white;'>"
	page += "<div>"
	page += "<nav class='navbar navbar-inverse' role='navigation' style='padding-left:80px;'>"
	page += "<ul class='nav navbar-nav'>"
	page += "<li><a href='http://localhost:17017/fa17g17/search'>Back<span class='sr-only'>(current)</span></a></li>"
	page += "</ul>"
	page += "</nav>"
	page += "</div>"
	page += "<div class='container-fluid'>"
	page += "<div class='row content' style='color: white;'>"
	page += "<div class='col-sm-3'>"
	page += "<br>"
	page += "<div id='myCarousel_image' style= 'width:100%;' class='carousel slide text-center' data-ride='carousel'>"
	page += "<ol class='carousel-indicators'>"
	page += "<li data-target='#myCarousel_image' data-slide-to='0' class='active'></li>"
	page += "<li data-target='#myCarousel_image' data-slide-to='1'></li>"
	page += "<li data-target='#myCarousel_image' data-slide-to='2'></li>"
	page += "</ol>"
	page += "<div class='carousel-inner' role='listbox' id='picBox'>"
	page += immo_image
	page += "</div>"
	page += "<a class='left carousel-control' href='#myCarousel_image' role='button' data-slide='prev'>"
	page += "<span class='glyphicon glyphicon-chevron-left' aria-hidden='true'></span>"
	page += "<span class='sr-only'>Previous</span>"
	page += "</a>"
	page += "<a class='right carousel-control' href='#myCarousel_image' role='button' data-slide='next'>"
	page += "<span class='glyphicon glyphicon-chevron-right' aria-hidden='true'></span>"
	page += "<span class='sr-only'>Next</span>"
	page += "</a>"
	page += "</div>"
	page += "</div>"
	page += "<div class='col-sm-9'>"
	page += "<div class='well'>"
	page += "<h4 id='title'>" + A.offer_title + "</h4>"
	page += "<p id='offerType'>Rent/Buy: " + A.offer_type + "</p>"
	page += "<p id='endDate'>Duration of the Offer: " + A.offer_startdate + "-" + A.offer_enddate + "</p>"
	page += "</div>"
	page += "<div class='row'>"
	page += "<div class='col-sm-3'>"
	page += "<div class='well'>"
	page += "<h4>Pricing</h4>"
	page += "<p id='price'>Rent/Buy: " + strconv.Itoa(Price) + " € </p>"
	page += "<p id='utilCost'>Utilities: " + strconv.Itoa(A.offer_utilcosts) + " € </p>"
	page += "</div>"
	page += "</div>"
	page += "<div class='col-sm-3'>"
	page += "<div class='well'>"
	page += "<h4>Additional Costs</h4>"
	page += "<p id='bail'>Bail: " + strconv.Itoa(A.offer_bail) + " € </p>"
	page += "<p id='provision'>Provision: " + strconv.FormatFloat(A.offer_provision, 'g', 3, 64) + " % </p>"
	page += "</div>"
	page += "</div>"
	page += "<div class='col-sm-3'>"
	page += "<div class='well'>"
	page += "<h4>Contact Information</h4>"
	page += "<p id='seller'>Agent: " + A.agent_firstname + " - " + A.agent_lastname + "</p>"
	page += "<p id='agency'>Agency: " + A.agent_agency + "</p>"
	page += "</div>"
	page += "</div>"
	page += "<div class='col-sm-3'>"
	page += "<div class='well'>"
	page += "<h4>Ratings</h4>"
	page += "<p id='rating_value'>Value: " + strconv.Itoa(rating) + "/5</p>"
	page += "</div>"
	page += "</div>"
	page += "</div>"
	page += "<div class='row'>"
	page += "<div class='col-sm-12'>"
	page += "<div class='well'>"
	page += "<h4>Description</h4>"
	page += "<p id='description'>" + A.immo_desc + "</p>"
	page += "</div>"
	page += "</div>"
	page += "</div>"
	page += "<div class='row'>"
	page += "<div class='col-sm-12'>"
	page += "<div class='well'>"
	page += "<h4>Features</h4>"
	page += "<div id='autoFeatureElements'>" + A.immo_features + "</div>"
	page += "</div>"
	page += "</div>"
	page += "</div>"
	page += "</div>"
	page += "</div>"
	page += "</div>"
	page += "<div class='container-fluid'>"
	page += "<div class='row content'>"
	page += "<div id='mapDiv' class='col-sm-3 sidenav'>"
	page += "<iframe id='maps' style='margin: 0; border:0' src='" + immo_gmap + "' width='100%' height='100%' frameborder='0' allowfullscreen></iframe>"
	page += "</div>"
	page += "<div class='col-sm-9'>"
	page += "<div class='well text-center'>"
	page += "<h4>Estate Informations</h4>"
	page += "</div>"
	page += "<div class='row'>"
	page += "<div class='col-sm-6'>"
	page += "<div class='well'>"
	page += "<h4>Adress Information</h4>"
	page += "<p id='address'>Address: " + A.immo_address + "</p>"
	page += "<p id='city'>City: " + A.city_name + "</p>"
	page += "<p id='postal'>Zip Code: " + A.city_plz + "</p>"
	page += "<p id='country'>Country: German</p>"
	page += "</div>"
	page += "</div>"
	page += "<div class='col-sm-6'>"
	page += "<div class='well'>"
	page += "<h4>Estate Details</h4>"
	page += "<p id='immo_type'>Type: " + A.immo_type + "</p>"
	page += "<p id='floors'>Floors: " + strconv.Itoa(A.immo_floors) + "</p>"
	page += "<p id='rooms'>Rooms: " + A.immo_rooms + "</p>"
	page += "<p id='condition'>Condition: " + A.immo_condition + "</p>"
	page += "<p id='qm'>Square Meters: " + strconv.Itoa(A.immo_qm) + "</p>"
	page += "<p id='heating'>Heating Type: " + A.immo_heating + "</p>"
	page += "<p id='estateyear'>Year of Construction: " + strconv.Itoa(A.immo_baujahr) + "</p>"
	page += "</div>"
	page += "</div>"
	page += "</div>"
	page += "</div>"
	page += "</div>"
	page += "</div>"
	page += "<footer class='footer' style='align-items: center; padding-top:20px'>"
	page += "<div class='container'>"
	page += "SFSU Software Engineering Project, Fall 2017. For Demonstration Only"
	page += "</div>"
	page += "</footer>"
	page += "</body>"
	page += "</html>"

	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	content := page

	io.WriteString(
		res,
		content,
	)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func handlerICon(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Icon Handle here!")
}

func main() {
	fmt.Println("Server Starting here!")
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/go/favicon.ico", handlerICon)
	r.HandleFunc("/go/search", search)
	r.HandleFunc("/go/show", showData)
	r.HandleFunc("/go/filter", filter)
	r.HandleFunc("/go/estateCreate", estateCreate)
	http.ListenAndServe("127.0.0.1:3000", nil)
}
