package converter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type N26CsvParser struct {
	logger   *log.Entry
	filepath string
}

func NewN26CsvParser(logger *log.Entry, path string) *N26CsvParser {
	return &N26CsvParser{
		logger:   logger,
		filepath: path,
	}
}

func (p *N26CsvParser) Parse() error {
	in, err := os.Open(p.filepath)

	if err != nil {
		p.logger.WithError(err).Debug("Datei konnte nicht gelesen werden")
		return err
	}

	r := csv.NewReader(in)

	records, err := r.ReadAll()

	if err != nil {
		p.logger.WithError(err).Debug("Inhalt der CSV-Datei konnte nicht geladen werden")
		return err
	}

	// 0 Datum						2019-08-31
	// 1 Empfänger					AMZNPrime DE
	// 2 Kontonummer				DE01234567891234567980 oder leer
	// 3 Transaktionstyp			MasterCard Zahlung / Gutschrift
	// 4 Verwendungszweck			Verwendungszweck oder leer
	// 5 Kategorie					Medien & Elektronik
	// 6 Betrag						-7.99 bei Belastung, positiv bei Gutschrift
	// 7 Betrag (Fremdwährung)      -7.39 bei Belastung, positiv bei Gutschrift
	// 8 Fremdwährung				EUR/USD (nur bei Belastungen gesetzt, bei Gutschriften leer, auch gesetzt wenn keine Fremdwährung)
	// 9 Wechselkurs				-1.00 oder leer, kann auch gesetzt sein wenn keine Fremdwährung

	buf := new(bytes.Buffer)

	// Ausgabe sollte sein ISO-8859-1 sein
	// Überschriften sind "Buchungstag", "Betrag", "Verwendungszweck", "Name"
	//
	// Feldbelegung wie folgt:
	// 0 Buchungstag			31.08.2019
	// 1 Betrag					[Betrag] in DE locale Formatierung
	// 2 Verwendungszweck		"-" wenn leer, ansonsten [Verwendungszweck]
	// 2 Name					[Empfänger]

	csv := csv.NewWriter(buf)
	csv.UseCRLF = true

	// Header schreiben
	err = csv.Write([]string{
		"Buchungstag",
		"Betrag",
		"Verwendungszweck",
		"Name",
	})

	if err != nil {
		p.logger.WithError(err).Debug("Fehler beim schreiben der CSV-Kopfzeile")
		return err
	}

	for i, r := range records {
		// Skip the first line, its a header
		if i == 0 {
			continue
		}

		p.logger.WithField("line", i).Debug("Verarbeite Transaktion")

		// Datum umwandeln
		rDate, err := time.Parse("2006-01-02", r[0])

		if err != nil {
			p.logger.WithError(err).WithField("date", r[0]).Debug("Datum weicht von erwartetem Wert ab")
			return err
		}

		// Betrag angleichen
		rValue := strings.ReplaceAll(r[6], ".", ",")

		// Verwendungszweck bestimmen
		rSubject := r[4]
		rSubject = strings.TrimSpace(rSubject)

		if rSubject == "" {
			rSubject = "-"
		}

		err = csv.Write([]string{
			rDate.Format("02.01.2006"),
			rValue,
			rSubject,
			r[1],
		})

		if err != nil {
			p.logger.WithError(err).Debug("Transaktion konnte nicht in CSV-Zeile umgewandelt werden")
			return err
		}
	}

	csv.Flush()

	out := new(bytes.Buffer)

	// Write output file
	iso, err := charset.NewWriter("iso-8859-1", out)

	if err != nil {
		p.logger.WithError(err).Debug("Fehler beim erstellen eines ISO-8859-1-Buffers zur Umwandlung")
		return err
	}

	n, err := iso.Write(buf.Bytes())

	if err != nil {
		p.logger.WithError(err).Debug("Fehler umwandeln der CSV-Daten in die ISO-8859-1 Kodierung")
		return err
	}

	// Determine output path
	outpath := strings.TrimSuffix(p.filepath, filepath.Ext(p.filepath))
	outpath = fmt.Sprintf("%s.sevdesk.csv", outpath)

	err = ioutil.WriteFile(outpath, out.Bytes(), 0644)

	if err != nil {
		p.logger.WithError(err).Debug("Fehler beim schreiben der Ausgabedaten")
	}

	p.logger.WithFields(log.Fields{
		"output": outpath,
		"bytes":  n,
	}).Debug("Ausgabe wurde geschrieben")

	return nil
}

func (p *N26CsvParser) ParseLine(line [][]string) {

}
