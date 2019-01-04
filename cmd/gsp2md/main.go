package main

import (
	"fmt"
	"github.com/argcv/gsp2md/pkg/assets"
	"github.com/argcv/gsp2md/pkg/gs"
	"github.com/argcv/gsp2md/pkg/md"
	"github.com/argcv/gsp2md/pkg/utils"
	"github.com/argcv/gsp2md/version"
	"github.com/argcv/stork/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"path"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/sheets/v4"
)

var (
	binCleanName = path.Clean(os.Args[0])
	versionMsg   = fmt.Sprintf("%v version \"%s (%s)\" %s\n", binCleanName, version.Version, version.GitHash, version.BuildDate)
	rootCmd      = &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if verbose, err := cmd.Flags().GetBool("verbose"); err == nil {
				if verbose {
					log.Verbose()
					log.Debug(versionMsg)
					log.Debug("verbose mode: ON")
				}
			}

			// init config
			conf, _ := cmd.Flags().GetString("config")
			if e := assets.LoadConfig(conf); e != nil {
				return e
			}

			// init system variables
			rand.Seed(time.Now().Unix())
			utils.SetMaxProcs()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := assets.LoadGsp2MdConfig()
			if err != nil {
				return err
			}
			clientId := cfg.GoogleSpreadsheets.Client
			clientSecret := cfg.GoogleSpreadsheets.Secret

			log.Infof("client: %v, secret: %v", clientId, clientSecret)
			// If modifying these scopes, delete your previously saved token.json.
			config := &oauth2.Config{
				ClientID:     clientId,
				ClientSecret: clientSecret,
				RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
				Scopes: []string{
					"https://www.googleapis.com/auth/spreadsheets.readonly",
				},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://accounts.google.com/o/oauth2/auth",
					TokenURL: "https://www.googleapis.com/oauth2/v3/token",
				},
			}

			client, err := gs.GetClient(config)
			if err != nil {
				log.Fatalf("Login failed: %v", err)
			}

			srv, err := sheets.New(client)
			if err != nil {
				log.Fatalf("Unable to retrieve Sheets client: %v", err)
			}

			for _, in := range cfg.Input {
				spreadsheetId, e := gs.Url2SpreadsheetId(in.Url)
				if e != nil {
					log.Errorf("Parse url %v failed: %v", in.Url, e)
					continue
				}

				for _, readRange := range in.Ranges {
					log.Infof("id: %v, range: %v", spreadsheetId, readRange)
					resp, e := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
					if e != nil {
						log.Fatalf("Unable to retrieve data from sheet: %v", e)
						continue
					}

					if len(resp.Values) == 0 {
						fmt.Println("No data found.")
					} else {
						table := md.Table{}
						for i, row := range resp.Values {
							sarow := []md.Cell{}
							for _, r := range row {
								sarow = append(sarow, md.NewPlainTextCell(fmt.Sprint(r)))
							}
							if i == 0 {
								table.Headers = sarow
							} else {
								table.Rows = append(table.Rows, md.Row{
									Columns: sarow,
								})
							}
						}
						fmt.Printf("Output:\n%v", table.String())
					}

				}
			}

			return nil
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Prints the version of %s", binCleanName),
		// do not execute any persistent actions
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(versionMsg)
		},
	}
)

func init() {
	rootCmd.AddCommand(
		versionCmd,
	)
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode")
	rootCmd.PersistentFlags().StringP("config", "c", "", "explicit assign a configuration file")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "log verbose")

	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Execute Failed!!! %v", err.Error())
		os.Exit(1)
	}
}
