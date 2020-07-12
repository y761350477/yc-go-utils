package y_md

func Markdown1Title(s string) string {
	return `# ` + s
}

func Markdown2Title(s string) string {
	return `## ` + s
}

func Markdown3Title(s string) string {
	return `### ` + s
}

func MarkdownImg(url, text string) string {
	return `![` + text + `](` + url + `)`
}

func MarkdownLink(url, text string) string {
	return `[` + text + `](` + url + `)`
}
