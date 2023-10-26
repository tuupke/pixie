# CuProxy, a simple ipp 1.0-2.2 compliant banner proxy

This project was started to solve a problem faced within ICPC-style programming contests: "How to determine which team
submitted which printjob?"

In these programming contests it is customary for teams to be able to print their sourcecode so they can continue
debugging while a teammate has the use of the (single) team machine. These prints should be kept private to the teams
that issued the prints to prevent cheating. In practice the only way of ensuring this is to have the organisation 
handing out the print jobs instead of teams retrieving the jobs themselves. 

The main problem with this is that the printer operator needs to know which print should be brought to which team. In 
linux-based contest environments the printing services are commonly setup using CUPS. CUPS has built-in support for
banner pages printing most - if not all - of the required information when setup properly. Though a program can override
this behavior which would result in a 'mystery-print' without a banner page.

CuProxy can be used to prevent this issue. It is a small IPP proxy that listens for jobs, constructs data, creates a
banner-page and ap- or prepends the banner to the document. The only problem with this is that it requires a printer
that supports printing PDF documents.

When CuProxy sees a printing request it will - in parallel to proxy-ing the request - start by constructing the banner
page. There are two ways of getting data on the banner-pages. The first is by telling CuProxy the data using key-value 
pairs in the printer-url. The second - and much more powerful - method is to use webhooks. It is important to note that
since 


# Installation

Simply download a binary or `go install` the proxy. The binary is statically compiled and will therefor work on all Linux
AMD64 installations.

```go install github.com/tuupke/pixie/cuproxy@latest```

# Running

### systemd service
The following systemd service can be used to ensure the proxy is running.

# Data on the banner
There are two ways of getting data printed on the banner-page: (1) using key-value pairs in the ipp-url from the client to CuProxy, and (2) using webhooks.

## KV in url
The following example registers two keys `foo`, and `bar` with the values of `fooz` and `foobar` respectively.

```ipp://ppp:6631/foo=fooz/bar=foobar```

## Webhooks
Webhooks are arguably the most powerful way of getting data for the banner-page. The implementation allows for the concurrent, or sequential, retrieval of data.


The following config example will be used to show how to configure and use the webhooks. Assume a value for `WEBHOOKS_TO_CALL` of:

```
misc;POST;https://www.domjudge.org/api/non-existing-endpoint&&user;GET;https://www.domjudge.org/demoweb/api/user?strict=false|team;GET;https://jury:jury@www.domjudge.org/demoweb/api/teams/{{team_id}}?strict=false
```

In this example, 2 sets of http-requests are made in parallel, the first - named 'misc' - executes a POST request to a non-existing endpoint while the second f

To specify a single http call use the format: `<name>;<http verb>;<url>`.

This variable first gets split by `&&`, where every split element is part will be called in parallel! Within every part


# Configuration

CUProxy is configured using the environment. For convenience all environment settings can be stored in a file called `.env` and stored in the same directory as the binary.

### Global settings

| Variable            | Type      | Default | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
|---------------------|-----------|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `PRINT_TO`          | `String`  | ""      | The IPP url of where the actual printer is located.                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `LOG_LEVEL`         | `String`  | ""      | The IPP url of where the actual printer is located.                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `LISTEN`            | `String`  | ":631"  | IP + port where to listen on. Defaults to 0.0.0.0:631 (CUPS port).                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| `PDF_LOCATION`      | `String`  | "/tmp"  | Where to store the rendered banner pages.                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `BANNER_APPEND`     | `Boolean` | false   | Whether to append the banner instead of prepending it.                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| `BANNER_ON_BACK`    | `Boolean` | false   | Whether to put the banner on the back of a page. Assumes but does not check whether a duplexer is installed!                                                                                                                                                                                                                                                                                                                                                                                         |
| `PRINT_KEYS`        | `String`  | "*"     | Comma seperated list of keys which values should be printed on the banner page. The special value "*" prints all keys in alphabetical order. If the key is unknown, it will not be listed on the banner page!                                                                                                                                                                                                                                                                                        |
| `DUMP_IPP_CONTENTS` | `String`  | ""      | The location on disk where to store proxied IPP messages. Leave empty to disable. Does nothing when `DUMP_ORIGINAL` and `DUMP_REPLACEMENTS` are both false. The dumped files have the following filenames: `<seq-id>-<dir>-<type>.bin` where `seq-id` is an incrementing integer uniquely identifying the request; `dir` the "direction", is it the request ("req"), or is it the printers response (res); and `type` depicts whether it is the original ("orig"), or the modified request ("repl"). |
| `DUMP_ORIGINAL`     | `Boolean` | false   | Whether to dump the original contents. Does nothing when `DUMP_IPP_CONTENTS` is empty.                                                                                                                                                                                                                                                                                                                                                                                                               |
| `DUMP_REPLACEMENTS` | `Boolean` | false   | Whether to dump the replaced contents. Does nothing when `DUMP_IPP_CONTENTS` is empty.                                                                                                                                                                                                                                                                                                                                                                                                               |

### Webhook settings

| Variable               | Type      | Default       | Description                                                                                                                                                                                 |
|------------------------|-----------|---------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `BASIC_AUTH_IN_DATA`   | `Boolean` | false         | The IPP url of where the actual printer is located.                                                                                                                                         |
| `BASIC_AUTH_USERNAME`  | `String`  | "ba_password" | IP + port where to listen on. Defaults to 0.0.0.0:631 (CUPS port).                                                                                                                          |
| `BASIC_AUTH_PASSWORD`  | `String`  | "ba_username" | Where to store the rendered banner pages.                                                                                                                                                   |
| `IMAGE_KEY`            | `String`  | false         | Whether to append the banner instead of prepending it.                                                                                                                                      |
| `HTTP_REQUEST_NONCE`   | `String`  | ""            | A nonce which is added as `X-Pixie-Nonce` header. Defaults to a randomized string of 32 characters.                                                                                         |
| `DOWNLOAD_DIR`         | `String`  | "/tmp"        | Whether to put the banner on the back of a page. Assumes but does not check whether a duplexer is installed!                                                                                |
| `WEBHOOK_KEY_TEMPLATE` | `String`  | ""            | Whether to put the banner on the back of a page. Assumes but does not check whether a duplexer is installed!                                                                                |
| `WEBHOOKS_TO_CALL`     | `String`  | ""            | Comma seperated list of keys which values should be printed on the banner page. The special value "*" prints all keys in alphabetical order. If the key is unknown, it will not be printed! |

### PDF settings
`gofpdf` is used for rendering the banner page, while `pdfcpu` is used to pre- or append the banner page to the actual print. The following variables can be set.

| Variable            | Type      | Default | Description                                                                                                                                                  |
|---------------------|-----------|---------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `PDF_UNIT`          | `String`  | "mm"    | The unit to use for the other settings. Supported values are: millimeter `mm`, centimeter `cm`, point `pt`, inch `in`.                                       |
| `PDF_PAGE_SIZE`     | `String`  | "A4"    | The size of the PDF to generate. Supported values are: `A3`, `A4`, `A5`, `Letter`, `Legal`, and `Tabloid`.                                                   |
| `PDF_FONT_DIR`      | `String`  | "."     | The file system location in which font resources will be found. Currently unused as different fonts are not supported. Every banner is rendered in Arial 12. |
| `PDF_LANDSCAPE`     | `boolean` | false   | Whether to render the banner in landscape.                                                                                                                   |
| `PDF_LEFT_MARGIN`   | `integer` | 10      | The number of `PDF_UNIT` units to leave blank at the left of the banner. Note, no right margin can be set, long lines will not be split in two!              |
| `PDF_TOP_MARGIN`    | `integer` | 15      | The number of `PDF_UNIT` units to leave blank at the top of the banner.                                                                                      |
| `PDF_BOTTOM_MARGIN` | `integer` | 6       | The number of `PDF_UNIT` units to leave blank at the bottom of the banner.                                                                                   |
| `PDF_LINE_HEIGHT`   | `integer` | 3       | The number of `PDF_UNIT` units to leave between two lines.                                                                                                   |

## Datasources

### URL params (key-value)

### HTTP(s) requests


