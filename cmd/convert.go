package cmd

import (
	"errors"
	"fmt"
	"github.com/adrianrudnik/n26-mt940-converter/converter"
	"github.com/mattn/go-zglob"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Converts the given file glob into a MT940 files",
	Args: func(cmd *cobra.Command, args []string) error {

		err := cobra.ExactArgs(1)(cmd, args)

		if err != nil {
			return err
		}

		files, err := zglob.Glob(args[0])

		if err != nil {
			return err
		}

		if len(files) == 0 {
			return errors.New("keine dateien gefunden die zum filter passen")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		files, err := zglob.Glob(args[0])

		if err != nil {
			return
		}

		fmt.Printf("%d Dateien zur Verarbeitung gefunden\n", len(files))

		var completed []bool

		for _, file := range files {
			// skip files with the pattern ".sevdesk.csv" out
			if strings.Contains(file, ".sevdesk.csv") {
				continue
			}

			flog := log.WithField("input", file)
			flog.Info("Starte Verarbeitung")

			p := converter.NewN26CsvParser(flog, file)

			err := p.Parse()

			if err != nil {
				log.WithError(err).Error("Datei konnte nicht verarbeitet werden")
				continue
			}

			flog.Info("Verarbeitung abgeschlossen")

			completed = append(completed, true)
		}

		if len(completed) > 0 {
			log.WithField("count", len(completed)).Info("Alle Dateien verarbeitet")
		}
	},
}
