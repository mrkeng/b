package models

import (
	"github.com/russross/blackfriday"
	"github.com/slogsdon/b/util"
	"gopkg.in/yaml.v1"
	"html/template"
	"io/ioutil"
	"strings"
	"time"
)

func SavePost(root string, form map[string][]string) error {
	filename := form["filename"][0]
	raw := form["raw"][0]
	hm, _ := ParsePostContent([]byte(raw), "md")
	categories := strings.Join(hm.Categories, "/") + "/"

	err := util.MakeDir(root + "/" + categories)
	if err != nil {
		return err
	}

	fullpath := root + "/" + categories + filename
	return util.WriteFile(fullpath, raw)
}

// GetAllPosts returns all posts from the storage system by name.
func GetAllPosts(root string) []Post {
	var posts []Post

	for _, f := range util.ReadDir(root) {
		posts = append(posts, preparePost(f))
	}

	return posts
}

// GetPost returns a single post from the storage system by name.
func GetPost(name, root string) Post {
	var post Post

	for _, f := range util.ReadDir(root) {
		if f.Filename == name {
			post = preparePost(f)
			break
		}
	}

	return post
}

// ParsePostContent parses the HeadMatter and HTML from a raw post.
func ParsePostContent(contents []byte, t string) (HeadMatter, template.HTML) {
	m, c := parseHeadMatter(contents)

	switch t {
	case "md", "mdown", "markdown":
		c = markdown(c)
	}

	return m, template.HTML(string(c))
}

// ParsePostSlugAndType parses a post's slug and type from
// its filename. The file extension is used for the post type.
// The slug is grabbed from the basename sans a prefixed date
// used for organization.
// It returns the post's slug and type.
func ParsePostSlugAndType(filename string) (string, string) {
	filenameNoDate := strings.Join(strings.Split(filename, "-")[3:], "-")
	split := strings.Split(filenameNoDate, ".")
	slug := strings.ToLower(strings.Join(split[:len(split)-1], "."))
	t := strings.ToLower(split[len(split)-1])
	return slug, t
}

func preparePost(f util.FileReading) Post {
	// Read file contents
	contents, _ := ioutil.ReadFile(f.Filename)

	// Grab slug and type from filename
	slug, t := ParsePostSlugAndType(f.Info.Name())

	// Parse our content/head matter from our file
	// Return our prepared Post
	head, formattedContents := ParsePostContent(contents, t)
	time, _ := time.Parse("2006-01-02 15:04:05", head.Date)
	return Post{
		Title:       head.Title,
		Slug:        slug,
		Content:     formattedContents,
		HeadMatter:  head,
		Filename:    f.Info.Name(),
		Directory:   strings.Replace(f.Filename, "/"+f.Info.Name(), "", 1),
		Type:        t,
		Raw:         string(contents),
		UpdatedAt:   f.Info.ModTime(),
		PublishedAt: time,
	}
}

// Represents the possible data contained within the
// head matter section of a post, fenced with leading
// and following --- lines.
type HeadMatter struct {
	Title      string   `json:"title"`
	Date       string   `json:"date"`
	Categories []string `json:"categories"`
}

func parseHeadMatter(contents []byte) (HeadMatter, []byte) {
	m := HeadMatter{}
	c := string(contents)

	if strings.Count(c, "---") >= 2 {
		split := strings.Split(c, "---")
		_ = yaml.Unmarshal([]byte(split[1]), &m)
		c = strings.Join(split[2:], "---")
	}

	return m, []byte(c)
}

func markdown(str []byte) []byte {
	// this did use blackfriday.MarkdownCommon, but it was stripping out <script>

	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_FOOTNOTES

	return blackfriday.Markdown(str, renderer, extensions)
}
