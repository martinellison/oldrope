// gen.go
package main

import (
	"io"
	"log"
	"regexp"
	"text/template"
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
	 	set: function(parts) {
			{{range .SetLines}}
			{{.}}
			{{end}}		
		},
		fix: function() {	
			{{range .FixLines}}			
			setClick('{{.Name}}',function() {
                	{{.Code}}
            		});
			{{end}}
        	},
        	redisplay: function() {
			{{range .RedisplayLines}}
			{{.}}
			{{end}}
		},
		refix: function() {
			{{range .Refixes}}			
			setClick('{{.Name}}',function() {
                	{{.Code}}
            		});
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

type outData struct {
	Pages []*outPage
}
type outPage struct {
	Name           string
	SetLines       []string
	FixLines       []fix
	RedisplayLines []string
	Refixes        []fix
}
type fix struct {
	Name string
	Code string
}

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

//<script>
func genJsStart(w io.Writer) {
	w.Write([]byte(compress(
		`var gd = {};
var ld = {};
var currentPage = 'start';
var cp;
var pages;
var setPage = function(pageName) {
    console.log('displaying page: ' + pageName);
    currentPage = pageName;
    cp = pages[currentPage];
    if (!cp) console.error('unknown page: ' + currentPage);
    ld = {};
    var parts = [];
    cp.set(parts);
    setHtml('main',parts.join("\n"));
    cp.fix();
    console.log('displayed ' + currentPage);
};
var displayPage = function() {
    cp.redisplay();
    console.log('redisplayed ' + currentPage);
};
var setHtml=function(id,text){document.getElementById(id).innerHTML = text;};
var setClick=function(id,fn){document.getElementById(id).onclick=fn;};
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
