package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-ini/ini"
)

var config tomlConfig
var tomlFile string
var environment string

func main() {

	//TEST THE ARGUMENTS AND CONFIGURE ACCORDINGLY
	if len(os.Args) >= 3 {
		log.Println("First argument: ", os.Args[1])
		log.Println("Second argument: ", os.Args[2])
		if strings.Compare(os.Args[1], "prod") == -1 {
			if strings.Compare(os.Args[1], "uprod") == -1 {
				if strings.Compare(os.Args[1], "test") == -1 {
					if strings.Compare(os.Args[1], "utest") == -1 {
						panic("first argument is not Supported")
					}
				}
			}
		} else {
			environment = os.Args[1]
		}
		tomlFile = "config/" + os.Args[2]
		log.Println("Using Toml Config file:  ", tomlFile)
	} else if len(os.Args) >= 2 {
		log.Println("First argument: ", os.Args[1])
		log.Println(strings.Compare(os.Args[1], "prod"), os.Args[1] == "prod")
		//
		log.Println(os.Args[1] == "prod" || os.Args[1] == "uprod" || os.Args[1] == "test" || os.Args[1] == "utest")
		log.Println(strings.Compare(os.Args[1], "prod") == -1 || strings.Compare(os.Args[1], "uprod") == -1 || strings.Compare(os.Args[1], "test") == -1 || strings.Compare(os.Args[1], "utest") == -1)
		//
		if strings.Compare(os.Args[1], "prod") == -1 {
			if strings.Compare(os.Args[1], "uprod") == -1 {
				if strings.Compare(os.Args[1], "test") == -1 {
					if strings.Compare(os.Args[1], "utest") == -1 {
						panic("first argument is not Supported")
					}
				}
			}
		} else {
			environment = os.Args[1]
			tomlFile = "config/config.toml"
		}
	} else {
		tomlFile = "config/config.toml"
		environment = "test"
	}

	//INITIALIZE THE THE PROCESSOR WITH THE DESIRED CONFIGURATION
	if _, err := toml.DecodeFile(tomlFile, &config); err != nil {
		log.Println(err)
		return
	}

	//CREATE A NEW LOG FILE AND RENAME THE OLD ONE BY ADDING A TIME STAMP
	logTime := fmt.Sprintf("%v", time.Now().Format("2006-01-02-15-04-05"))
	//PROCESSOR LOGS
	err := os.Rename(config.Logs.Location+config.Logs.Name, config.Logs.Location+config.Logs.Name+"-"+logTime)
	if err != nil {
		log.Println("Cannot rename the log file", err)
	}

	//WINSCP LOGS
	err = os.Rename(config.Winscp.LogLoc+config.Winscp.Log, config.Winscp.LogLoc+config.Winscp.Log+"-"+logTime)
	if err != nil {
		log.Println("Cannot rename the log file", err)
	}

	_, err = os.Create(config.Logs.Location + config.Logs.Name)
	if err != nil {
		log.Println("Cannot create the log file", err)
	}
	//SETTING THE LOGS PERMISSIONS AND THE OUTPUT
	log.Println("Getting the log files and the permissions")
	perm := os.FileMode(config.Logs.Permission)
	logFile, err := os.OpenFile(config.Logs.Location+config.Logs.Name, config.Logs.Flag, perm)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	//STARTING THE PROCESS
	log.Println("............................................Starting the download Process............................................")
	log.Println("Used Command: ", os.Args)
	log.Println("Number of arguments: ", len(os.Args))
	//LOG FILE TITLE
	log.Printf("Config File Title: %s\n", config.Title)
	//CHECK THE TIME IF IT IS SAFE TO GET A TOKEN OR NOT
	t := time.Now()
	s := t.Second()
	log.Printf("Reading current seconds: %d", s)
	if int(s) >= 0 && int(s) < 6 {
		log.Println("Seconds are between zero and five it is safe to get the Token: ")
		//IT IS SAFE TO GET A TOKEN
		command()
	} else {
		//IT IS NOT SAFE TO GET A TOKEN WAITING UNTIL THE MINUTE ENDS THEN PROCEED
		log.Println("Seconds are not between zero and five it is not safe to get the Token: ", s)
		wait := (60 - s) * 1000
		log.Printf("Waiting for: %ds : then we can get the token safely\n", wait/1000)

		time.Sleep(time.Duration(wait) * time.Millisecond)
		t := time.Now()
		s := t.Second()
		if int(s) >= 0 && int(s) < 6 {
			log.Println("Now it is safe to get the Token: ", s)
			//IT IS SAFE TO GET A TOKEN
			command()

		} else {
			log.Println("Failed something went wrong")
		}

	}

}

//FUNCTION THAT HANDLES THE PROCESS
func command() {
	//GET THE TOKEN
	log.Println("calling the following command to get the token", config.Hotkey.Location+config.Hotkey.Command)
	c := exec.Command("cmd", "/C", config.Hotkey.Location+config.Hotkey.Command, "")
	if err := c.Run(); err != nil {
		log.Println("Error something went wrong with getting the token: ", err)
		panic(err)
	}
	//READING THE TOKEN FROM THE RESULTS OF THE TOKEN COMMAND
	log.Println("Reading the token from the output of the previous command")
	cfg, err := ini.Load(config.Hotkey.OutputLocation)
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	log.Println("Token:", cfg.Section("token").Key("value").String())

	//CREATING THE TRANSFER FILE FOR WINSCP TO CONNECT AND GET THE FILES AND ADD A UNIQUE PREFIX FOR EACH TYPE
	transferFile, err := os.Create(config.Winscp.Location + config.Winscp.TransFile)
	if err != nil {
		log.Println("Cannot create the winscp transfer file", err)
	}
	//defer transferFile.Close()

	host := config.Servers[environment].URL
	port := config.Servers[environment].Port
	token := cfg.Section("token").Key("value").String()

	line := fmt.Sprintf("open sftp://%s:%s@%s:%s -rawsettings ProxyMethod=3 ProxyHost=%s", config.User.Name, token, host, port, config.Winscp.Proxy)
	log.Println("Adding to the transfer file", line)
	fmt.Fprintln(transferFile, line)
	for fileName, file := range config.Files {
		log.Printf("Check if the following file should be downloaded: %s (%s, %s,%s, %s, %v)\n", fileName, file.ServerName, file.ServerLocation, file.OutputName, file.OutputLocation, file.Get)
		if file.Get {
			line = fmt.Sprintf("Adding File to the transfer: %s (%s, %s,%s, %s, %v, %s)", fileName, file.ServerName, file.ServerLocation, file.OutputName, file.OutputLocation, file.Get, file.Prefix)
			log.Println("Adding to the transfer file: ", line)
			line = fmt.Sprintf("cd %s", file.ServerLocation)
			fmt.Fprintln(transferFile, line)
			log.Println("Adding to the transfer file: ", line)
			line = fmt.Sprintf("get -nopermissions -nopreservetime %s %s%s", file.ServerName, config.Winscp.OutputLocation, file.Prefix)
			fmt.Fprintln(transferFile, line)
			log.Println("Adding to the transfer file: ", line)
		}
	}
	fmt.Fprintln(transferFile, "exit")
	log.Println("Adding to the transfer file: ", line)
	transferFile.Close()
	//CALLING THE WINSCP COMMAND
	log.Println("calling the following command to get the files", config.Winscp.Location+config.Winscp.Command+" "+config.Winscp.TransFile)
	c = exec.Command("cmd", "/C", config.Winscp.Location+config.Winscp.Command+" "+config.Winscp.TransFile+" >"+config.Winscp.LogLoc+config.Winscp.Log, "")

	if err := c.Run(); err != nil {
		log.Println("Error something went wrong with getting the files: ", err)
		panic(err)
	}
	//FINISHED THE WINSCP COMMAND

	log.Println("Starting the rename process for each downloaded file")
	//READING EACH FILE TYPE RENAME IT AND PLACE IT IN THE OUTPUT DIRECTORY FOR THE SCRIPTS TO READ IT
	for fileName, file := range config.Files {
		log.Printf("Check if the following file should be renamed: %s (%s, %s,%s, %s, %v)\n", fileName, file.ServerName, file.ServerLocation, file.OutputName, file.OutputLocation, file.Get)
		if file.Get {
			files, err := filepath.Glob(config.Winscp.OutputLocation + file.Prefix)
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range files {
				t := time.Now()
				s := fmt.Sprintf("%02d", t.Second())
				mi := fmt.Sprintf("%02d", t.Minute())
				h := fmt.Sprintf("%02d", t.Hour())
				y := fmt.Sprintf("%04d", t.Year())
				mo := fmt.Sprintf("%02d", t.Month())
				d := fmt.Sprintf("%02d", t.Day())
				time.Sleep(2000 * time.Millisecond)
				outputname := fmt.Sprintf("%s.%s-%s-%s-%s-%s-%s.001", file.OutputName, y, mo, d, h, mi, s)

				log.Printf("naming %s to %s\n", f, file.OutputLocation+outputname)
				err := os.Rename(f, file.OutputLocation+outputname)
				if err != nil {
					panic(err)
				}
			}

		}
	}
	log.Println("............................................Finished the download Process............................................")
}
