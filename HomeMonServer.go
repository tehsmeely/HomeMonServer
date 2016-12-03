
package main

import (
	auth "github.com/abbot/go-http-auth"
	"fmt"
	"net"
	"net/http"
	"bufio"
	"github.com/spf13/viper"
)

// https://systembash.com/a-simple-go-tcp-server-and-tcp-client/
var (
	monConnection 	net.Conn
	daemonAddress 	string
	daemonPort 		string
)

func MonConnect() bool {
	var err error
	monConnection, err = net.Dial("tcp", daemonAddress+":"+daemonPort)
	if err != nil{
		fmt.Println("Failed to connect to daemon", monConnection)
		return false
	}else{
		return true
	}
}

	// SecretKeys := map[string]string{
	// 	//"john" : "b98e16cbc3d01734b264adba7baa3bf9", //hello
	// 	"john"		: "f87db0562f16741c05fd9417e488a39d", //MD5(john:monserve.com:hello)
	// 	"WebHost1" 	: "021e16e58e49bf136876a2a163dd8bc5", //MD5(WebHost1:monserve.com:vcCCbSBIwm31vMpfk5Jt)
	// }



func handle(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	fmt.Fprintf(w, "Out Of Bounds. Use /api", r.Username)
}

func daemonRead(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	numRetries := 3
	var ready bool
	if monConnection == nil {
		fmt.Println("monConnection is nil, reconnecting")
		var success bool
		for i:=0; i<numRetries; i++ {
			success = MonConnect()
			if success {
				_, err := fmt.Fprintf(monConnection, "\n")
				if err == nil {
					success = true
					break
				}
			}
		}
		if !success {
			fmt.Fprintf(w, "FAIL")
			return
		}
	}

	//Write to connection, triggering returned info. Catch issues
	fmt.Println("MonConnection: ", monConnection)
	_, err := fmt.Fprintf(monConnection, "\n")
	fmt.Println("Send Err:", err)
	if err != nil {
		fmt.Println("Reconnecting to daemon (Send)")
		for i:=0; i<numRetries; i++ {
			success := MonConnect()
			if success {
				_, err := fmt.Fprintf(monConnection, "\n")
				if err == nil {
					success = true
					break
				}
			}
		}
	} else {
		ready = true
	}
	//Receive from connection and write to http response. Catch issues
	if ready {
		message, err := bufio.NewReader(monConnection).ReadString('\n')
		if err != nil {
			fmt.Println("Reconnecting to daemon (Recv)")
			for i:=0; i<numRetries; i++ {
				success := MonConnect()
				if success {
					message, err := bufio.NewReader(monConnection).ReadString('\n')
					if err == nil {
						fmt.Fprintf(w, message)
						return
					}
				}
			}
			fmt.Fprintf(w, "FAIL")
		} else {
			fmt.Fprintf(w, message)
		}
	} else {
		fmt.Fprintf(w, "FAIL")
	}
}


func main() {
	//Read Server Config. some viper help from http://karloespiritu.com/handling-configuration-files-in-go/
	viper.SetConfigType("yaml")
	viper.SetConfigName("server")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/HomeMonServer/")
	err := viper.ReadInConfig()
	var servePort string
	if err != nil {
		fmt.Println("Config file not found...", err)
		return 
	}
	servePort = viper.GetString("server.port") // Port to listen on
	daemonAddress = viper.GetString("daemon.address") // address of daemon
	daemonPort = viper.GetString("daemon.port") // port of daemon
	secrets := auth.HtdigestFileProvider("/etc/HomeMonServer/basic.htdigest")// a map of 'user:monserve.com:MD5(user:monserve.com:password)' for digest auth
	authenticator := auth.NewDigestAuthenticator("monserve.com", secrets)
	http.HandleFunc("/", authenticator.Wrap(handle))
	http.HandleFunc("/api", authenticator.Wrap(daemonRead))
	//Establish connection to daemon
	success := MonConnect()
	if success {
		fmt.Println("Listening on port ", servePort)
		http.ListenAndServe(":"+servePort, nil)
	} else {
		fmt.Println("Failed to connect to daemon")
	}
}
