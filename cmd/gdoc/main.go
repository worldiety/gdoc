package main

import (
	"flag"
	"github.com/worldiety/gdoc/internal/app"
	"log"
	"os"
	"os/exec"
)

func main() {
	var cfg app.Config
	cfg.Reset()
	cfg.Flags(flag.CommandLine)
	flag.Parse()

	buf, err := app.Apply(cfg)
	if err != nil {
		panic(err)
	}
	switch cfg.OutputFormat {
	case app.Adoc, app.Html, app.Pdf:
		file, err := os.Create("doc.adoc")
		if err != nil {
			log.Printf("file could not be created\nerror: %e", err)
		}

		_, err = file.Write(buf)
		if err != nil {
			log.Printf("could not write to file '%v'\nerror: %e", file, err)
		}

		if cfg.OutputFormat == app.Html {
			err = RenderToHtml(file.Name())
			if err != nil {
				log.Printf("file could not be created\nerror: %e", err)
			}
		}
		if cfg.OutputFormat == app.Pdf {
			err = RenderToPdf(file.Name())
			if err != nil {
				log.Printf("file could not be created\nerror: %e", err)
			}
		}
	default:
		log.Print("no output file")
	}
}

// RenderToHtml loads the file in the given path and uses the asciidoc cli tool to render and save a html file
func RenderToHtml(adocFilename string) error {
	htmlFileName := "htmlOutput.html"
	// use asciidoctor to create a html file from the adoc file
	cmd := exec.Command("asciidoctor", "-b", "html5", "-o", htmlFileName, adocFilename)
	setupCMD(cmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// RenderToPdf takes a filename of an adoc file and uses asciidoc-pdf to render and save a pdf file
func RenderToPdf(adocFileName string) error {
	// Use the asciidoctor-pdf library to generate a PDF from the adoc file
	// get commands from command line and export errors to it
	cmd := exec.Command("asciidoctor-pdf", adocFileName)
	setupCMD(cmd)

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func setupCMD(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
}
