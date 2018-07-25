<h1 align="center">hugo-utils</h1>

<h5 align="center">Tools for managing <a href="https://gohugo.io">Hugo</a> websites.</h5>

<br/>

`hugo-utils` contains a set of tools that I use to manage my own blog: [https://ops.tips](https://ops.tips).

Although `hugo` by itself makes the production a breeze, when you start having a bunch of content, it might become hard to manage the consistency between how content is tagged and managed as a whole.

This set of auxiliary tools comes handy for those who want to make sure that:

- every post is properly tagged; 
- their content is always up to date (remind you if old blog posts need attention) (TODO); and
- taxonomy terms (like, a specific category) have metadata alright (TODO).

## Install

### Binaries

- **darwin** [amd64](https://github.com/cirocosta/hugo-utils/releases/download/0.0.3/hugo-utils_0.0.3_darwin_amd64.tar.gz)
- **linux** [amd64](https://github.com/cirocosta/hugo-utils/releases/download/0.0.3/hugo-utils_0.0.3_linux_amd64.tar.gz)

### Using Go

```sh
go get -u -v github.com/cirocosta/hugo-utils
```

## Commands

### List

```sh
NAME:
   hugo-utils list - lists all content under a given path.

USAGE:
   hugo-utils list [command options] [format]

DESCRIPTION:
   The 'list' command iterates over each content file (*.md)
   found under a given root directory (--directory), then prints
   to 'stdout' a description of each.

   The default formatting displays the following attributes for
   each page: title, keywords, tags, categories, slug, date.

   A custom format can also be specified following Go template
   rules. In this case, the render state contains:
   - {{ . }}: the current page in the page traversal; and
   - {{ .Pages }}: the list of all pages found.

EXAMPLES:

   Display every property of the pages under a given
   section that lives under "./content/blog" using the default
   formatting:

     hugo-utils \
       --directory=./content/blog

   Display the text of every page in a given section
   that lives under "./content/blog" and their keywords:

     hugo-utils \
       --directory=./content/blog \
       '{{ .Title }} - {{ .Keywords }}'

   Display the path to the files that don't have keywords
   specified:

     hugo-utils \
       --directory=./content/blog \
       '{{ if eq (len .Keywords) 0 }} {{ .Path }} {{ end }}'


OPTIONS:
   --directory value  path to the directory where contents exist (.md)
   --type value       content type to list entries by (pages|tags|categories) (default: "pages")
   --sort value       thing to sort by (title|date|lastmod) (default: "lastmod")
   --draft            only show drafts
```

### Update

```sh
NAME:
   hugo-utils update - updates the frontmatter of a page.

   The 'update' command takes care of updating the frontmatter
   of a given content page (e.g., /content/blog/mypost.md).

   Taking a desired update in the form of 'yaml', it parses the
   content page and applies to it the merge between the original
   frontmatter and the updated frontmatter.

   When no updated yaml is passed, the default frontmatter is
   applied (e.g., a post without 'tags' would now have 'tags: []').

EXAMPLES:

   Update the contents of page1.md with the defaults of the FrontMatter
   object from './hugo':

     cat ./page1.md
       ---
       title: 'page1'
       ---
       body

     hugo-utils update --filepath ./page1.md
     cat ./page1.md
       ---
       title: page1
       description: ""
       slug: ""
       image: ""
       date: 0001-01-01T00:00:00Z
       lastmod: 0001-01-01T00:00:00Z
       draft: false
       tags: []
       categories: []
       keywords: []
       ---
       body

   Update the contents of page1.md with the defaults of the FrontMatter
   object from './hugo' merged with a custom set of fields that we defined
   with a 'yaml' provided in the positional arguments:

     hugo-utils update --filepath ./page1.md 'tags: ["tag3"]'
     cat ./page1.md
       ---
       title: page1
       ... (other fields)
       tags:
       - tag1
       - tag2
       ---
       body


USAGE:
   hugo-utils update [command options] [yaml]

OPTIONS:
   --filepath value  path to the page file
```

tip: Use this command together with `find` and `xargs` to perform updates across a great number of files:

```sh
find . -name "*.md" | xargs -I {} -P 4 hugo-utils update --filepath={}
```

