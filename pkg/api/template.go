package api

import (
	"github.com/gin-gonic/gin"
	"html/template"
)

const htmlTemplate = `
<!DOCTYPE html PUBLIC>
<html>
<head>
<meta charset="EUC-KR">
<title>초간단 테이블</title>
</head>
<body>
    <table border="1">
	<th>테이블</th>
	<th>만들기</th>
	<tr><!-- 첫번째 줄 시작 -->
	    <td>첫번째 칸</td>
	    <td>두번째 칸</td>
	</tr><!-- 첫번째 줄 끝 -->
	<tr><!-- 두번째 줄 시작 -->
	    <td>첫번째 칸</td>
	    <td>두번째 칸</td>
	</tr><!-- 두번째 줄 끝 -->
    </table>
</body>
</html>
`

func ResponseHtml() gin.HandlerFunc {
	return func(c *gin.Context) {
		responseHtml(c, nil)
	}
}

func responseHtml(c *gin.Context, data *Response)  {
	t := template.New("response")
	t, err := t.Parse(htmlTemplate)
	if err != nil {}

	t.Execute(c.Writer, htmlTemplate)
	return
}