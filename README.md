<h1 align="center">hugo-utils</h1>

<h5 align="center">Tools for managing <a href="https://gohugo.io">Hugo</a> websites.</h5>

<br/>

`hugo-utils` contains a set of tools that I use to manage my own blog: [https://ops.tips](https://ops.tips).

Although `hugo` by itself makes the production a breeze, when you start having a bunch of content, it might become hard to manage the consistency between how content is tagged and managed as a whole.

This set of auxiliary tools comes handy for those who want to make sure that:

- their content is always up to date (remind you if old blog posts need attention);
- every post is properly tagged (TODO); and 
- taxonomy terms (like, a specific category) have metadata alright (TODO).

## Install

### Binaries

- **darwin** [amd64](https://github.com/cirocosta/hugo-utils/releases/download/v0.0.1/hugo-utils-darwin-amd64)
- **linux** [amd64](https://github.com/cirocosta/hugo-utils/releases/download/v0.0.1/hugo-utils-linux-amd64)

### Using Go

```sh
go get -u -v github.com/cirocosta/hugo-utils
```

## Commands

### List

```sh
NAME:
   hugo-utils list - lists all content under a given path.

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


USAGE:
   hugo-utils list [command options] [format]

OPTIONS:
   --directory value  path to the directory where contents exist (.md)
   
```

