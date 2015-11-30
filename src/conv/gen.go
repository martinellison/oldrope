// gen.go
package main

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"text/template"
	"time"
)

//func main() {
//	makeTemplate()
//	makeTestData()
//	expandTemplate(os.Stdout)
//}

var templ *template.Template

func makeTemplate() {
	var err error
	templateText := compress(`pages = {
	{{range .Pages}} {{.Name}}: {
	 	init: function() {
			{{range .InitLines}}
			{{.}}
			{{end}}
		},
		display: function(parts) {
			{{range .SetLines}}
			{{.}}
			{{end}}
        	},
		fix: function(parts) {
			{{range .FixLines}}
			{{.}}
			{{end}}
        	},
	},
	{{end}}
	};`)
	templ, err = template.New("page").Parse(templateText)
	if err != nil {
		log.Fatalf("template def error: %v", err)
	}
}

var whiteSpaceRegex *regexp.Regexp

func init() {
	whiteSpaceRegex = regexp.MustCompile(`[\s]+`)
}
func compress(inStr string) string { return whiteSpaceRegex.ReplaceAllLiteralString(inStr, " ") }

var theOutData outData

func expandTemplate(w io.Writer) {
	err := templ.Execute(w, theOutData)
	if err != nil {
		log.Fatalf("template exp error: %v", err)
	}
}
func genStart(w io.Writer) {
	w.Write([]byte(compress(
		`<!DOCTYPE html>
<html>
<head>
<meta charset='UTF-8'/>
 <style>a {
    color: blue;
    text-decoration: underline;
    cursor: pointer;
}
html, body {
    color: black;
    font-family: Georgia, serif;
}</style>
</head>
<body>
<div id='main'> </div>
`)))
}

func genHeader(w io.Writer) {
	w.Write([]byte(fmt.Sprintf("/* created by program on %s */", time.Now())))
}
func genJsStart(w io.Writer) {
	w.Write([]byte(compress(
		`var gd = {};
var ld = {};
var currentPage = 'start';
var cp;
var pages;
var displayPage = function() { var parts = [];
    cp = pages[currentPage];
    if (!cp) console.error('unknown page: ' + currentPage);
    cp.display(parts);
    setHtml('main',parts.join("\n"));
    cp.fix();
    console.log('displayed ' + currentPage);
};
var setPage = function(pageName) {
    console.log('displaying page: ' + pageName);
    currentPage = pageName;
    ld = {};
	df = {};
	displayPage();
};
var setHtml=function(id,text){var elt=document.getElementById(id); if(!elt)alert('no '+id);elt.innerHTML = text;};
var setClick=function(id,fn){var elt=document.getElementById(id); if(!elt)console.log('no '+id);else elt.onclick=fn;};
`)))
}
func genJsEnd(w io.Writer) {
	w.Write([]byte(compress(`setPage('start');
displayPage();
console.log('script loaded');
`)))
}
func genEnd(w io.Writer) {
	w.Write([]byte(compress(`</body>
</html>`)))
}
