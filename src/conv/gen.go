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
	 	set: function() {
            parts = [];
			{{range .SetLines}}
			{{.}}
			{{end}}			
            $('#main').html(parts.join("\n"));
			{{range .FixLines}}			
			$('#{{.Name}}').click(function() {
                {{.Code}}
            });
			{{end}}
        },
        redisplay: function() {
			{{range .RedisplayLines}}
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

type outData struct {
	Pages []*outPage
}
type outPage struct {
	Name           string
	SetLines       []string
	FixLines       []fix
	RedisplayLines []string
}
type fix struct {
	Name string
	Code string
}

var theOutData outData

//func makeTestData() { //not used
//	pageZzz := &outPage{Name: "zzz", SetLines: []string{"aaa", "bbb"}, RedisplayLines: []string{"xxx", "yyy"}}
//	pageBill := &outPage{Name: "bill", SetLines: []string{"zzz", "xxx"}, RedisplayLines: []string{"ccc", "vvv"}}
//	theOutData = outData{Pages: []*outPage{pageZzz, pageBill}}
//}
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
<script src='jquery.js'></script>
<script>
		var gd = {};
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
    cp.set();
    console.log('displayed ' + currentPage);
};
var displayPage = function() {
    cp.redisplay();
    console.log('redisplayed ' + currentPage);
};`)))
}
func genEnd(w io.Writer) {
	w.Write([]byte(compress(`setPage('start');
displayPage();
console.log('script loaded');</script>
</body>
</html>`)))
}
