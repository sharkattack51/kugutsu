package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	flags "github.com/jessevdk/go-flags"
	"github.com/tailscale/win"
)

type Options struct {
	Mode string `short:"m" long:"mode" default:"udp" description:"listen mode UDP|HTTP"`
	Port int    `short:"p" long:"port" default:"5000" description:"listen port number"`
}

var (
	opts Options
)

const (
	MODE_UDP  = "UDP"
	MODE_HTTP = "HTTP"

	COMMAND_CLICK  = "CLICK"
	COMMAND_WCLICK = "WCLICK"
	COMMAND_KEY    = "KEY"
)

func main() {
	// option flags
	_, err := flags.Parse(&opts)
	if err != nil { // [help] also passes
		os.Exit(0)
	}

	switch strings.ToUpper(opts.Mode) {
	case MODE_UDP:
		opts.Mode = MODE_UDP
	case MODE_HTTP:
		opts.Mode = MODE_HTTP
	default:
		opts.Mode = MODE_UDP
	}

	fmt.Println("===================================================")
	fmt.Println("[[ kugutsu v0.1 ]]")
	switch opts.Mode {
	case MODE_UDP:
		fmt.Printf("%s listner start... :%d\n", opts.Mode, opts.Port)
		fmt.Println("")
		fmt.Println("Send UDP message examples:")
	case MODE_HTTP:
		fmt.Printf("%s listner start... http://%s:%d?msg={}\n", opts.Mode, GetHostIP(), opts.Port)
		fmt.Println("")
		fmt.Println("HTTP request query examples:")
	}
	fmt.Println("{ CLICK }			-> Mouse click at screen center")
	fmt.Println("{ CLICK:100:200 }		-> Mouse click at (100, 200)")
	fmt.Println("{ WCLICK:100:200 }		-> Mouse double click at (100, 200)")
	fmt.Println("{ KEY:a }			-> Press key[a]")
	fmt.Println("{ KEY:a,LSHIFT,RALT,CTRL }	-> Press key[a + left-shift + right-alt + ctrl]")
	fmt.Println("{ KEY:a:3 }			-> Press key[a] hold 3 sec")
	fmt.Println("{ KEY:a,SHIFT:1 }		-> Press key[a + shift] hold 1 sec")
	fmt.Println("{ KEY:abc }			-> Type key[a b c]")
	fmt.Println("===================================================")

	switch opts.Mode {
	case MODE_UDP:
		addr := &net.UDPAddr{
			IP:   net.ParseIP("localhost"),
			Port: opts.Port,
		}
		upd, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Fatalln(err)
		}

		buf := make([]byte, 64)
		for {
			n, _, err := upd.ReadFromUDP(buf)
			if err != nil {
				// pass
			} else {
				msg := string(buf[:n])
				go Automation(msg)
			}
		} // listen loop

	case MODE_HTTP:
		srv := &http.Server{Addr: ":" + strconv.Itoa(opts.Port)}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			if query.Has("msg") {
				go Automation(query.Get("msg"))
			}
		})
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}
}

func Automation(msg string) {
	log.Printf("message recieved... %s\n", msg)
	msgs := strings.Split(msg, ":")

	if len(msgs) == 0 {
		return
	}

	cmd, pos, keyStr, spKeys, keyHold, err := ParseCommand(msgs)
	if err != nil {
		log.Println(err)
		return
	}

	switch cmd {
	case COMMAND_CLICK:
		robotgo.Move(pos[0], pos[1])
		robotgo.Click("left", false)

	case COMMAND_WCLICK:
		robotgo.Move(pos[0], pos[1])
		robotgo.Click("left", true)

	case COMMAND_KEY:
		if _, ok := robotgo.Keycode[keyStr]; ok {
			if keyHold > 0 {
				st := time.Now()
				for time.Since(st) < (time.Duration(keyHold) * time.Second) {
					robotgo.KeyDown(keyStr, spKeys...)
				}
				robotgo.KeyUp(keyStr, spKeys...)
			} else {
				robotgo.KeyTap(keyStr, spKeys...)
			}
		} else {
			robotgo.TypeStr(keyStr)

			if keyHold > 0 {
				log.Println("strings key hold is invalid.")
			}
		}

	default:
	}
}

func ParseCommand(msg []string) (string, [2]int, string, []interface{}, int, error) {
	var cmd string
	var pos [2]int
	var keyStr string
	var spKeys []interface{}
	var keyHold int
	var err error

	if len(msg) == 0 {
		return cmd, pos, keyStr, spKeys, keyHold, errors.New("parse comannd error")
	}
	cmd = strings.ToUpper(msg[0])

	switch cmd {
	case COMMAND_CLICK, COMMAND_WCLICK:
		if len(msg) == 1 {
			w := int(win.GetSystemMetrics(win.SM_CXSCREEN))
			h := int(win.GetSystemMetrics(win.SM_CYSCREEN))
			pos = [2]int{w / 2, h / 2}
		} else if len(msg) == 3 {
			var x, y int
			if x, err = strconv.Atoi(msg[1]); err != nil {
				return cmd, pos, keyStr, spKeys, keyHold, errors.New("parse comannd error")
			}
			if y, err = strconv.Atoi(msg[2]); err != nil {
				return cmd, pos, keyStr, spKeys, keyHold, errors.New("parse comannd error")
			}
			pos = [2]int{x, y}
		}

	case COMMAND_KEY:
		if len(msg) >= 2 {
			strs := strings.Split(msg[1], ",")
			if len(strs) > 0 {
				keyStr = strs[0]
			}

			if len(strs) >= 2 {
				for i := 1; i < len(strs); i++ {
					switch strings.ToUpper(strs[i]) {
					case "SHIFT":
						spKeys = append(spKeys, "shift")

					case "LSHIFT":
						spKeys = append(spKeys, "lshift")

					case "RSHIFT":
						spKeys = append(spKeys, "rshift")

					case "ALT":
						spKeys = append(spKeys, "alt")

					case "LALT":
						spKeys = append(spKeys, "lalt")

					case "RALT":
						spKeys = append(spKeys, "ralt")

					case "CTRL":
						spKeys = append(spKeys, "ctrl")

					case "LCTRL":
						spKeys = append(spKeys, "lctrl")

					case "RCTRL":
						spKeys = append(spKeys, "rctrl")

					default:
					}
				}
			}
		}

		if len(msg) >= 3 {
			var t int
			if t, err = strconv.Atoi(msg[2]); err != nil {
				return cmd, pos, keyStr, spKeys, keyHold, errors.New("parse comannd error")
			}
			keyHold = t
		}
	}

	return cmd, pos, keyStr, spKeys, keyHold, nil
}

func GetHostIP() string {
	ip := "127.0.0.1"
	if runtime.GOOS == "windows" {
		host, _ := os.Hostname()
		addrs, _ := net.LookupIP(host)
		for _, a := range addrs {
			if ipv4 := a.To4(); ipv4 != nil {
				ip = ipv4.String()
				break
			}
		}
	} else {
		addrs, _ := net.InterfaceAddrs()
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP.String()
					break
				}
			}
		}
	}

	return ip
}
