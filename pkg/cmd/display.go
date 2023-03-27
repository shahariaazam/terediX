package cmd

import (
	"fmt"
	"log"
	"net/http"
	"teredix/pkg/config"
	"teredix/pkg/storage"
	"teredix/pkg/visualize"
	"teredix/pkg/visualize/cytoscape"

	"github.com/spf13/cobra"
)

func NewDisplayCommand() *cobra.Command {
	var cfgFile string

	cmd := cobra.Command{
		Use:   "display",
		Short: "Display resource graph",
		Long:  "Display resource graph",
		RunE: func(cmd *cobra.Command, args []string) error {
			appConfig, err := config.Load(cfgFile)
			if err != nil {
				return err
			}

			err = config.Validate(appConfig)
			if err != nil {
				return err
			}

			st := storage.BuildStorage(appConfig)
			c := cytoscape.NewCytoscapa(st)
			v := visualize.NewVisualizer(c)
			html, err := v.Render()
			if err != nil {
				return err
			}

			http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
				fmt.Fprint(writer, html)
			})
			log.Println("Displaying resource graph at http://localhost:8989")
			err = http.ListenAndServe(":8989", nil)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "a valid yaml file is required")

	return &cmd
}