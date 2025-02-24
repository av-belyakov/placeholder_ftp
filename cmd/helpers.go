package main

import (
	"fmt"
	"os"

	"github.com/av-belyakov/placeholder_ftp/internal/appname"
	"github.com/av-belyakov/placeholder_ftp/internal/appversion"
	"github.com/av-belyakov/placeholder_ftp/internal/confighandler"
)

func getInformationMessage(name string, confLocalFtp, confMainFtp confighandler.ConfigFtp) string {
	appStatus := fmt.Sprintf("%vproduction%v", Ansi_Bright_Blue, Ansi_Reset)
	nameRegObj := fmt.Sprintf("%v%s%v", Ansi_Bright_Blue, name, Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_PHFTP_MAIN")
	if ok && envValue == "development" {
		appStatus = fmt.Sprintf("%v%s%v", Ansi_Bright_Red, envValue, Ansi_Reset)
	}

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched", appname.GetAppName(), appversion.GetAppVersion())

	fmt.Printf("\n%v%v%s.%v\n", Bold_Font, Ansi_Bright_Green, msg, Ansi_Reset)
	fmt.Printf("%v%vApplication status is '%s'.%v\n", Underlining, Ansi_Bright_Green, appStatus, Ansi_Reset)
	fmt.Printf("%vName regional object:'%s'.%v\n", Ansi_Bright_Green, nameRegObj, Ansi_Reset)
	fmt.Printf("%vLocal FTP server settings:%v\n", Ansi_Bright_Green, Ansi_Reset)
	fmt.Printf("%v  ip: %v%s%v\n", Ansi_Bright_Green, Ansi_Bright_Blue, confLocalFtp.Host, Ansi_Reset)
	fmt.Printf("%v  net port: %v%d%v\n", Ansi_Bright_Green, Ansi_Bright_Magenta, confLocalFtp.Port, Ansi_Reset)
	fmt.Printf("%vMain FTP server settings:%v\n", Ansi_Bright_Green, Ansi_Reset)
	fmt.Printf("%v  ip: %v%s%v\n", Ansi_Bright_Green, Ansi_Bright_Blue, confMainFtp.Host, Ansi_Reset)
	fmt.Printf("%v  net port: %v%d%v\n\n", Ansi_Bright_Green, Ansi_Bright_Magenta, confMainFtp.Port, Ansi_Reset)

	return msg
}
