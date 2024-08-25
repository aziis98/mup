package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:embed public/*
var publicFS embed.FS

const MB = 1 << 20

var (
	maxUploadSize = pflag.Int64P("max-upload-size", "s", 100, "Maximum upload size in MB")

	port = pflag.IntP("port", "p", 5000, "Port to run the server on")
	host = pflag.StringP("host", "h", "0.0.0.0", "Host to run the server on")
)

type Map map[string]any

func filenameToSlug(filename string) string {
	ext := filepath.Ext(filename)
	filename = filename[:len(filename)-len(ext)]

	pattern := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return strings.ToLower(pattern.ReplaceAllString(filename, "-")) + ext
}

func renderTemplate(w http.ResponseWriter, data any) {
	tmpl := template.Must(template.ParseFS(publicFS, "public/index.html"))

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error rendering the template")
		log.Println(err)
	}
}

func getUploads(uploadFolder string) ([]string, error) {
	files, err := os.ReadDir(uploadFolder)
	if err != nil {
		return nil, err
	}

	uploads := []string{}
	for _, file := range files {
		uploads = append(uploads, file.Name())
	}

	slices.Reverse(uploads)

	return uploads, nil
}

func main() {
	pflag.CommandLine.Init(os.Args[0], pflag.ContinueOnError)
	pflag.Usage = func() {
		w := os.Stderr
		fmt.Fprintf(w, "Usage:\n")
		fmt.Fprintf(w, "  %s [OPTIONS] [UPLOAD_FOLDER]\n", os.Args[0])
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "A micro file uploader, the default upload folder is 'Uploads'\n")
		fmt.Fprintf(w, "\n")

		fmt.Fprintf(w, "Options:\n")
		pflag.PrintDefaults()
		fmt.Fprintf(w, "\n")
	}

	err := pflag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		if err == pflag.ErrHelp {
			os.Exit(0)
		}

		log.Fatalf("error parsing flags: %v", err)
	}

	uploadFolder := pflag.Arg(0)
	if uploadFolder == "" {
		log.Println("No upload folder specified, using 'Uploads'")
		uploadFolder = "Uploads"
	}

	if _, err := os.Stat(uploadFolder); os.IsNotExist(err) {
		os.Mkdir(uploadFolder, 0755)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		uploads, err := getUploads(uploadFolder)
		if err != nil {
			log.Println("Error getting uploads")
			log.Println(err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderTemplate(w, Map{
			"Uploads": uploads,
			"Host":    *host,
			"Port":    *port,
		})
	})

	r.Get("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, publicFS, "public/style.css")
	})

	r.Get("/uploads", func(w http.ResponseWriter, r *http.Request) {
		uploads, err := getUploads(uploadFolder)
		if err != nil {
			log.Println("Error getting uploads")
			log.Println(err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, uploads)
	})

	r.Get("/uploads/{filename}", func(w http.ResponseWriter, r *http.Request) {
		filename := chi.URLParam(r, "filename")
		http.ServeFile(w, r, path.Join(uploadFolder, filename))
	})

	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(*maxUploadSize * MB); err != nil {
			log.Println("Error parsing the form")
			log.Println(err)

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		files := r.MultipartForm.File["files"]
		for _, file := range files {
			// Open the uploaded file
			src, err := file.Open()
			if err != nil {
				log.Println("Error opening the file")
				log.Println(err)

				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer src.Close()

			filename := fmt.Sprintf("%s_%s",
				time.Now().Format("2006-01-02_15-04-05"),
				filenameToSlug(file.Filename),
			)

			// Create a destination file in the uploads folder
			dst, err := os.Create(filepath.Join(uploadFolder, filename))
			if err != nil {
				log.Println("Error creating destination file")
				log.Println(err)

				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer dst.Close()

			// Copy the uploaded file to the destination file
			if _, err := io.Copy(dst, src); err != nil {
				log.Println("Error saving the file")
				log.Println(err)

				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			log.Printf("File Uploaded: %s", filename)
			fmt.Fprintf(w, "%s\n", filename)
		}
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", *host, *port),
		Handler: r,
	}

	log.Printf("Starting server on %v...", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
