package templates

import "html/template"

// für templates: https://www.practical-go-lessons.com/chap-32-templates
const pageHtml = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Kalendereintrag</title>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC" crossorigin="anonymous">
</head>
<body>
	<button type="button" class="btn btn-primary">Primary</button>
</body>
</html>`

var TemplTest = template.Must(template.New("page").Parse(pageHtml))
