package main

import (
	"html/template"
	"net/http"
	"strings"
)

func uploadPageHandler(w http.ResponseWriter, r *http.Request) {
	const tpl = `
<html>
<title>Go upload</title>
<body>
<form action="{{.}}/upload" method="post" enctype="multipart/form-data">
<label for="files">Files:</label>
<input type="file" name="files" id="files" multiple> <br>
Optional FilePath:
<input type="text" name="file_path" >
<input type="submit" name="submit" value="Submit">
</form>
</body>
</html>
`
	t, _ := template.New("page").Parse(tpl)

	t.Execute(w, strings.TrimSuffix((r.RequestURI), "/upload"))
}
