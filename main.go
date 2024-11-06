package main

import (
	"ascii-art-web/ascii"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Error : page not found", http.StatusNotFound)
		return
	}
	t.Execute(w, data)
}

func restrict(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		restrictedPaths := []string{"/static", "/images"}
		for _, path := range restrictedPaths {
			if r.URL.Path == path || r.URL.Path == path+"/" {
				renderTemplate(w, "templates/403.html", nil)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderTemplate(w, "templates/400.html", nil)
		return
	}
	if r.URL.Path != "/" {
		renderTemplate(w, "templates/404.html", nil)
		return
	}
	renderTemplate(w, "templates/index.html", nil)
}

func Ascii(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		renderTemplate(w, "templates/400.html", nil)
		return
	}

	input := r.FormValue("text")
	style := r.FormValue("Style")
	if style != "standard" && style != "thinkertoy" && style != "shadow" {
		renderTemplate(w, "templates/400_invalidEntry.html", nil)
		return
	}

	PrintArt, unprintable := ascii.PrintArt(input, style)
	data := struct {
		Output             string
		UnprintableWarning bool
	}{
		Output:             PrintArt,
		UnprintableWarning: unprintable,
	}
	renderTemplate(w, "templates/index.html", data)
}
func downloadText(w http.ResponseWriter, r *http.Request) {
	output := r.FormValue("output")
	if output == "" {
		renderTemplate(w, "templates/500_NoContent.html", nil)
		return

	}
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_art.txt")
	w.Header().Set("Content-Type", "text/plain")
	_, err := io.WriteString(w, output)
	if err != nil {
		renderTemplate(w, "templates/500.html", nil)
		return

	}
}

func downloadHTML(w http.ResponseWriter, r *http.Request) {
	output := r.FormValue("output")
	if output == "" {
		renderTemplate(w, "templates/500_NoContent.html", nil)
		return
	}

	htmlContent := fmt.Sprintf("<pre>%s</pre>", output)
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_art.html")
	w.Header().Set("Content-Type", "text/html")
	_, err := io.WriteString(w, htmlContent)
	if err != nil {
		renderTemplate(w, "templates/500.html", nil)
		return

	}
}
func About(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/About.html", nil)
}
func readME(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/readme.html", nil)
}

func main() {
	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	imageFs := http.FileServer(http.Dir(filepath.Join("templates", "images")))
	http.Handle("/images/", http.StripPrefix("/images/", imageFs))

	http.HandleFunc("/", Home)
	http.HandleFunc("/ascii", Ascii)
	http.HandleFunc("/download/txt", downloadText)
	http.HandleFunc("/download/html", downloadHTML)
	http.HandleFunc("/about", About)
	http.HandleFunc("/readme", readME)

	fmt.Println("local host running : http://localhost:8080")
	http.ListenAndServe(":8080", restrict(http.DefaultServeMux))
}
