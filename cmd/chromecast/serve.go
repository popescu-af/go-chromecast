package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the current folder",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := "..."
		folder := "."
		if len(args) > 0 {
			arg := args[0]
			f, err := os.Stat(arg)
			if err != nil {
				return err
			}
			if f.IsDir() {
				folder = arg
			} else {
				folder, file = filepath.Split(arg)
			}
		}

		ip, err := getOutboundIP()
		if err != nil {
			return err
		}

		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}
		port := listener.Addr().(*net.TCPAddr).Port

		addr := fmt.Sprintf("http://%s:%d/%s", ip, port, file)

		fmt.Println("You may now run:")
		fmt.Printf("%s load \"%s\"\n\n", os.Args[0], addr)
		fmt.Println("Type Ctrl+C to stop")

		var h http.Handler
		if file == "..." {
			h = http.FileServer(http.Dir(folder))
		} else {
			h = allowOneURI(file, http.FileServer(http.Dir(folder)))
		}

		return http.Serve(listener, h)
	},
}

func allowOneURI(uri string, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+uri {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	}
}

// Get preferred outbound ip of this machine
// from https://stackoverflow.com/a/37382208
func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP, nil
}
