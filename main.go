package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/gojektech/heimdall/httpclient"
)

//ATHAN

//Athan grabs main data header from Json
type Athan struct {
	Data YtData `json:"data"`
}

//YtData grabs secondary header under "Athan"
type YtData struct {
	Timings YtTime `json:"timings"`
}

//YtTime Gets 3rd Header from Json file, which is where the athan times are located
type YtTime struct {
	F string `json:"Fajr"`
	D string `json:"Dhuhr"`
	A string `json:"Asr"`
	M string `json:"Maghrib"`
	I string `json:"Isha"`
}

//Config Get Config settings from config.json file
type Config struct {
	Location struct {
		City     string `json:"City"`
		Country  string `json:"Country"`
		State    string `json:"State"`
		TimeZone string `json:"TimeZone"`
		MP3      string `json:"MP3"`
	}
}

//Split API
const (
	MainAPI    string = "http://api.aladhan.com/v1/timingsByCity?city="
	CountryAPI string = "&country="
	StateAPI   string = "&state="
	MethodAPI  string = "&method="
)

var config Config

//END ATHAN
var w *astilectron.Window

func SaveConfigs(city, state, TimeZone, mp3 string) {
	// {
	// 	"location":{
	// 	   "City":"Anaheim",
	// 	   "State":"CA",
	// 	   "Country":"US",
	// 	   "TimeZone":"America/Los_Angeles"
	// 	}
	//  }

	Format := "{ \"location\":{ \"City\":" + "\"" + city + "\"," + "\"State\":" + "\"" + state + "\"," + "\"Country\":\"US\", \"TimeZone\":" + "\"" + TimeZone + "\"," + "\"MP3\": " + "\"" + mp3 + "\" } }"
	//Format := "{\n\t\"location\":{\n\t\t\"City\":" + "\"" + city + "\"" + ",\n\t\t\"State\":" + "\"" + state + "\"" + ",\n\t\t\"Country\":\"US\",\n\t\t\"TimeZone\":" + "\"" + TimeZone + "\"\n\t}\n}"
	fs, err := os.Create("config.json")
	if err != nil {
		fmt.Println(err)
	}
	fs.Write([]byte(Format))
	defer fs.Close()
}

func main() {

	//Athan
	var err error
	//Connect to Json file for settings and paramaters
	config, err = LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error importing config.json file", err)
	}

	Y := ACal()
	go func() {
		for range time.Tick(time.Second * 35) {
			//Grab Updated Config Files
			config, err := LoadConfig("config.json")
			if err != nil {
				log.Fatal("Error importing config.json file", err)
			}

			//Get Local time test
			t := time.Now()
			location, err := time.LoadLocation(config.Location.TimeZone)
			if err != nil {
				log.Fatal("Unable to get Local Location", err)
			}
			CurrentTime := fmt.Sprint(t.In(location).Format("15:04"))

			//Duhur
			pd, _ := time.Parse("15:04", Y.Data.Timings.D)
			pd = pd.Add(time.Minute * time.Duration(-30))
			pds := fmt.Sprintf(pd.Format("15:04"))

			//Asr
			pa, _ := time.Parse("15:04", Y.Data.Timings.A)
			pa = pa.Add(time.Minute * time.Duration(-30))
			pas := fmt.Sprintf(pa.Format("15:04"))

			//Magrib
			pm, _ := time.Parse("15:04", Y.Data.Timings.M)
			pm = pm.Add(time.Minute * time.Duration(-30))
			pam := fmt.Sprintf(pm.Format("15:04"))

			//Isha
			pi, _ := time.Parse("15:04", Y.Data.Timings.I)
			pi = pi.Add(time.Minute * time.Duration(-30))
			pai := fmt.Sprintf(pi.Format("15:04"))

			//Checks if its time for Fajir
			if Y.Data.Timings.F == CurrentTime {
				//fmt.Println("Time for Fajir")
				w.SendMessage("Nowfajir:"+Y.Data.Timings.F, func(m *astilectron.EventMessage) {})
				w.Show()
			}

			if pds == CurrentTime {
				w.SendMessage("Upduhur:"+Y.Data.Timings.F, func(m *astilectron.EventMessage) {})
				w.Show()
			}

			//Checks if its time for Duhur
			if Y.Data.Timings.D == CurrentTime {
				//fmt.Println("Time for Duhur")
				w.SendMessage("Nowduhur:", func(m *astilectron.EventMessage) {})
				w.Show()
			}
			if pas == CurrentTime {
				w.SendMessage("Upasr:"+Y.Data.Timings.F, func(m *astilectron.EventMessage) {})
				w.Show()
			}

			//Checks if its time for Asr
			if Y.Data.Timings.A == CurrentTime {
				//fmt.Println("Time for Asr")
				w.SendMessage("Nowasr:", func(m *astilectron.EventMessage) {})
				w.Show()
			}

			if pam == CurrentTime {
				w.SendMessage("Upmagrib:"+Y.Data.Timings.F, func(m *astilectron.EventMessage) {})
				w.Show()
			}

			//Checks if its time for Magrib
			if Y.Data.Timings.M == CurrentTime {
				w.SendMessage("Nowmagrib:", func(m *astilectron.EventMessage) {})
				w.Show()
			}

			if pai == CurrentTime {
				w.SendMessage("Upisha:"+Y.Data.Timings.F, func(m *astilectron.EventMessage) {})
				w.Show()
			}

			//Checks if time for Isha
			if Y.Data.Timings.I == CurrentTime {
				//fmt.Println("Time for Isha")
				w.SendMessage("Nowisha:", func(m *astilectron.EventMessage) {})
				w.Show()
				Y = ACal() //Recall Json Data}

			}
		} // End Loop
	}()
	//End Athan

	// Set logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create astilectron
	a, err := astilectron.New(l, astilectron.Options{
		AppName:            "Athan",
		BaseDirectoryPath:  "Athan",
		AppIconDefaultPath: "D:\\Experiments\\Toast\\icon.png",
	})
	if err != nil {
		l.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		l.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	//Create the main menu
	if w, err = a.NewWindow("athan/index.html", &astilectron.WindowOptions{
		Center:    astikit.BoolPtr(true),
		Height:    astikit.IntPtr(480), //480
		Width:     astikit.IntPtr(350),
		MaxHeight: astikit.IntPtr(480),
		MaxWidth:  astikit.IntPtr(350),
		MinHeight: astikit.IntPtr(480),
		MinWidth:  astikit.IntPtr(350),
	}); err != nil {
		l.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	//Create main menu
	if err = w.Create(); err != nil {
		l.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	//Fajir
	f := strings.Split(Y.Data.Timings.F, ":")
	fc, _ := strconv.Atoi(f[0])
	fv := fc - 12
	fc2 := strconv.Itoa(fv)
	var fl string
	if fc >= 13 {
		fl = fc2
	} else {
		fl = Y.Data.Timings.F
	}

	//Duhur
	d := strings.Split(Y.Data.Timings.D, ":")
	dc, _ := strconv.Atoi(d[0])
	dv := dc - 12
	dc2 := strconv.Itoa(dv)
	var dl string
	if dc >= 13 {
		dl = dc2
	} else {
		dl = Y.Data.Timings.D
	}

	//Asr
	a1 := strings.Split(Y.Data.Timings.A, ":")
	ac, _ := strconv.Atoi(a1[0])
	av := ac - 12
	ac2 := strconv.Itoa(av)

	//Magrib
	m1 := strings.Split(Y.Data.Timings.M, ":")
	mc, _ := strconv.Atoi(m1[0])
	mv := mc - 12
	mc2 := strconv.Itoa(mv)

	//Magrib
	i1 := strings.Split(Y.Data.Timings.I, ":")
	ic, _ := strconv.Atoi(i1[0])
	iv := ic - 12
	ic2 := strconv.Itoa(iv)

	//Send the prayer times to main page
	w.SendMessage("fajir:"+fl, func(m *astilectron.EventMessage) {})
	w.SendMessage("duhur:"+dl, func(m *astilectron.EventMessage) {})
	w.SendMessage("asr:"+ac2+":"+a1[0], func(m *astilectron.EventMessage) {})
	w.SendMessage("magrib:"+mc2+":"+m1[0], func(m *astilectron.EventMessage) {})
	w.SendMessage("isha:"+ic2+":"+i1[0], func(m *astilectron.EventMessage) {})

	//Refresh the prayer times every 5 minutes
	go func() {
		for range time.Tick(time.Minute * 5) {
			w.SendMessage("fajir:"+fl, func(m *astilectron.EventMessage) {})
			w.SendMessage("duhur:"+dl, func(m *astilectron.EventMessage) {})
			w.SendMessage("asr:"+ac2+":"+a1[0], func(m *astilectron.EventMessage) {})
			w.SendMessage("magrib:"+mc2+":"+m1[0], func(m *astilectron.EventMessage) {})
			w.SendMessage("isha:"+ic2+":"+i1[0], func(m *astilectron.EventMessage) {})
		}
	}()

	//Create a tray icon the the athan app
	var t = a.NewTray(&astilectron.TrayOptions{
		Image:   astikit.StrPtr("D:\\Experiments\\Toast\\icon.png"),
		Tooltip: astikit.StrPtr("Windows Athan"),
	})

	//Create tray
	t.Create()
	var i int = 1
	//Click on the the tray
	t.On(astilectron.EventNameTrayEventClicked, func(e astilectron.Event) (deleteListener bool) {

		if i%2 == 0 {
			w.Hide()
			i = i + 1
		} else {
			w.Show()
			i = i + 1
		}

		fmt.Println("TRAY HAS BEEN CLICKED")
		return
	})

	//Dont allow resize
	w.On(astilectron.EventNameWindowEventMinimize, func(e astilectron.Event) (deleteListener bool) {
		w.Hide()
		return
	})

	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		m.Unmarshal(&s)

		var sr *astilectron.Window

		// Process message
		if s == "edit" {
			config, err = LoadConfig("config.json")
			if err != nil {
				log.Fatal("Error importing config.json file", err)
			}

			//Edit the athan location
			if sr, err = a.NewWindow("athan/edit.html", &astilectron.WindowOptions{
				Center:    astikit.BoolPtr(true),
				Height:    astikit.IntPtr(500),
				Width:     astikit.IntPtr(400),
				MaxHeight: astikit.IntPtr(500),
				MaxWidth:  astikit.IntPtr(400),
				MinHeight: astikit.IntPtr(500),
				MinWidth:  astikit.IntPtr(400),
			}); err != nil {
				l.Fatal(fmt.Errorf("main: new window failed: %w", err))
			}

			// Create windows
			if err = sr.Create(); err != nil {
				l.Fatal(fmt.Errorf("main: creating window failed: %w", err))
			}

			//Send the prefilled data
			sr.SendMessage("SetDataState:"+config.Location.State, func(m *astilectron.EventMessage) {})
			sr.SendMessage("SetDataCity:"+config.Location.City, func(m *astilectron.EventMessage) {})
			sr.SendMessage("SetDataTimeZone:"+config.Location.TimeZone, func(m *astilectron.EventMessage) {})
			sr.SendMessage("SetDataMP3:"+config.Location.MP3, func(m *astilectron.EventMessage) {})

		}

		sr.OnMessage(func(m *astilectron.EventMessage) interface{} {
			// Unmarshal
			var s string
			m.Unmarshal(&s)

			// Process message
			if strings.Contains(s, "back") {
				clean := strings.Replace(s, "back:", "", -1)
				split := strings.Split(clean, ":")
				SaveConfigs(split[0], split[1], split[2], split[3])

				config, err = LoadConfig("config.json")
				if err != nil {
					log.Fatal("Error importing config.json file", err)
				}

				sr.SendMessage("Restart:", func(m *astilectron.EventMessage) {})
				w.SendMessage("Restart:", func(m *astilectron.EventMessage) {})
			}
			return nil
		})

		return nil
	})
	// Blocking pattern
	a.Wait()

}

//ACal API Function
func ACal() Athan {
	//var Meth = strconv.Itoa(config.Calculation.Method)
	var Meth = "2"
	var AthanAPI = MainAPI + config.Location.City + CountryAPI + config.Location.Country + StateAPI + config.Location.State + MethodAPI + Meth
	FormatAPI := fmt.Sprintf(AthanAPI)

	// Create a new HTTP client with a default timeout
	timeout := 3000 * time.Millisecond
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	// Use the clients GET method to create and execute the request
	resp, err := client.Get(FormatAPI, nil)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)

	var Y Athan
	err = json.Unmarshal(body, &Y)
	if err != nil {
		log.Fatal(err)
	}

	return Y
}

//LoadConfig file
func LoadConfig(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error importing config.json file", err)
	}
	defer configFile.Close()
	if err != nil {
		return config, err
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}
