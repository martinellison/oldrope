// Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.

// gen.go (code generator)
package main

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"text/template"
	"time"
)

type generator struct {
	/* templ is the code output template */
	templ *template.Template
}

/* makeTemplate initialises the template */
func (theGenerator *generator) makeTemplate() {
	var err error
	templateText := compress(`pages = {
	{{range .Pages}} {{.Name}}: {
	 	init: function() {
			{{.InitLines}}
		},
		display: function(parts) {
			{{.SetLines}}
        	},
		fix: function(parts) {
			{{.FixLines}}
        	},
	},
	{{end}}
	};`)
	theGenerator.templ, err = template.New("page").Parse(templateText)
	if err != nil {
		log.Fatalf("template def error: %v", err)
	}
}

/* whiteSpaceRegex is a regular expression for detecting white space */
var whiteSpaceRegex *regexp.Regexp

func init() {
	whiteSpaceRegex = regexp.MustCompile(`[\s]+`)
}

/* compress compresses a string by converting sequences of whitespace to a single space */
func compress(inStr string) string { return whiteSpaceRegex.ReplaceAllLiteralString(inStr, " ") }

// /* */ var theOutData outData

/* expandTemplate expands the output data into the template */
func (theGenerator *generator) expandTemplate(w io.Writer, theOutData outData) {
	err := theGenerator.templ.Execute(w, theOutData)
	if err != nil {
		log.Fatalf("template exp error: %v", err)
	}
}

/* genStart generates the fixed part of the output file */
func (theGenerator *generator) genStart(w io.Writer) {
	w.Write([]byte(compress(
		`<!DOCTYPE html>
<html>
<head>
<meta charset='UTF-8'/>
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
<style>
a {
    color: blue;
    text-decoration: underline;
    cursor: pointer;
}
html, body {color: black; font-family: Georgia, serif;}
</style>
</head>
<body>
<div id='main'> </div>
`)))
}

/* genHeader generates the fixed part of the output file */
func (theGenerator *generator) genHeader(w io.Writer) {
	w.Write([]byte(fmt.Sprintf("/* created by program on %s */", time.Now())))
}

/* genJsStart generates the fixed part of the output file */
func (theGenerator *generator) genJsStart(w io.Writer, startPageName string) {
	w.Write([]byte(compress(
		`var gd = {};
  var ld = {};
  var currentPage = '` + startPageName + `';
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

/* genJsEnd generates the fixed part of the output fil*/
func (theGenerator *generator) genJsEnd(w io.Writer) {
	w.Write([]byte(compress(`setPage('start');
displayPage();
console.log('script loaded');
`)))
}

/* genEnd generates the fixed part of the output file */
func (theGenerator *generator) genEnd(w io.Writer) {
	w.Write([]byte(compress(`</body>
</html>`)))
}

// This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.
