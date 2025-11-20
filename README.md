# Blink

Blink is a programming language documentation and reference scraping and mirroring tool.

## Usage

Crawl and process documentation using the `scrape` subcommand. It should store the processed HTML files under the `out/` directory.

``` sh
./blink scrape com.cppreference/c
./blink scrape com.cppreference/cpp
```

The contents of the `out/` directory can work statically. You can also use the `serve` subcommand to run a local file server against the `out/` directory and serve the processed documentation with some styling applied.

``` sh
./blink serve
```

## Why?

[Toph](https://toph.co), the competitive programming platform by [Furqan Software](https://furqansoftware.com), needed a way to provide contest participants access to programming language documentation and references.

In on-site contests, Internet access is often limited to toph.co and related services only. The contest organizers have to ensure that the required programming language manuals and resources are included in the participant computers for offline access.

Using Blink, Toph crawls and processes programming language documentation into a static site. Toph then provides access to it from within the contest arena.

We considered a few existing solutions, including the well-built open-source [devdocs.io](https://devdocs.io). However, it uses Ruby, requires a backend, and a few other details are not configurable and do not fit our needs.

We had to build Blink to check all the boxes.

Of course, devdocs.io provides much more documentation than what Blink does currently. But we needed just a few anyway.
