# n26-to-sevdesk

Befehl um ein oder mehrere N26 CSV-Dateien in ein kompatibles Format für den manuellen Import in https://sevdesk.de/ umzuwandeln.

[![license](https://img.shields.io/github/license/adrianrudnik/n26-to-sevdesk.svg)](https://lab.klonmaschine.de/adrian.rudnik/n26-to-sevdesk/blob/master/LICENSE)
[![go report card](https://goreportcard.com/badge/github.com/adrianrudnik/n26-to-sevdesk)](https://goreportcard.com/report/github.com/adrianrudnik/n26-to-sevdesk)

## Installation

Es stehen vorkompilierte, ausführbare Dateien in den [Releases]((https://github.com/adrianrudnik/uritool/n26-to-sevdesk) zur Verfügung oder es steht einem frei diese selbst vom Quellcode zu kompilieren.

## Verwendung

Die ausführbare Datei entweder irgendwo ablegen wo sie ausgeführt werden kann. Danach alle umzuwandelnden CSV-Dateien in einen Ordner legen und den Befehl starten:

```sh
n26-to-sevdesk --verbose *.csv
```

Danach sollten alle gefundenen Dateien umgewandelt werden in gleichnamige Dateien mit der Endung `{name}.sevdesk.csv` welche für den Import nach Sevdesk vorgesehen sind.