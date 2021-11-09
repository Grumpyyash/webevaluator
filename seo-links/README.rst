pylinkvalidator
===============

:Version: 0.3

pylinkvalidator is a standalone and pure python link validator and crawler that
traverses a web site and reports errors (e.g., 500 and 404 errors) encountered.
The crawler can also download resources such as images, scripts and
stylesheets.

pylinkvalidator's performance can be improved by installing additional libraries
that require a C compiler, but these libraries are optional.

We created pylinkvalidator so that it could be executed in environments without
access to a compiler (e.g., Microsoft Windows, some posix production
environments) or with an old version of python (e.g., Centos).

pylinkvalidator is highly modular and has many configuration options, but the
only required parameter is the starting url: pylinkvalidate.py
http://www.example.com/

pylinkvalidator can also be used programmatically by calling one of the functions
in ``pylinkvalidator.api``

.. image:: https://api.travis-ci.org/bartdag/pylinkvalidator.png


Quick Start
-----------

Install pylinkvalidator with pip or easy_install:

::

  pip install pylinkvalidator


Crawl all pages from a site and show progress:

::

  pylinkvalidate.py -P http://www.example.com/


Requirements
------------

pylinkvalidator does not require external libraries if executed with python 2.x.
It requires beautifulsoup4 if executed with python 3.x. It has been tested on
python 2.6, python 2.7, and python 3.6.

For production use, it is strongly recommended to use lxml or html5lib because
the default HTML parser provided by python is not very lenient.


Optional Requirements
---------------------

These libraries can be installed to enable certain modes in pylinkvalidator:

lxml
  beautifulsoup can use lxml to speed up the parsing of HTML pages. Because
  lxml requires C libraries, this is only an optional requirement.

html5lib
  beautifulsoup can use html5lib to process incorrect or strange markup. It is
  slower than lxml, but believed to be more lenient.

gevent
  this non-blocking io library enables pylinkvalidator to use green threads
  instead of processes or threads. gevent could potentially speed up the
  crawling speed on web sites with many small pages.

cchardet
  this library speeds up the detection of document encoding.


Usage
-----

This is a list of all available options. See the end of the README file for
usage examples.

::

  Usage: pylinkvalidate.py [options] URL ...

  Options:
    --version             Show program's version number and exit
    -h, --help            Show this help message and exit
    -V VERBOSE, --verbose=VERBOSE
                          Display debugging info
                            None:  --verbose=0 (default)
                            Quiet: --verbose=1
                            Info:  --verbose=2

    Crawler Options:
      These options modify the way the crawler traverses the site.

      -O, --test-outside  Fetch resources from other domains without crawling
                          them
      -H ACCEPTED_HOSTS, --accepted-hosts=ACCEPTED_HOSTS
                          Comma-separated list of additional hosts to crawl
                          (e.g., example.com,subdomain.another.com)
      -i IGNORED_PREFIXES, --ignore=IGNORED_PREFIXES
                          Comma-separated list of host/path prefixes to ignore
                          (e.g., www.example.com/ignore_this_and_after/)
      -b, --ignore-bad-tel-urls
                          ignore badly formed tel URLs missing the leading +
                          sign, e.g., tel:1234567890 - only necessary for Python
                          > 2.6
      -u USERNAME, --username=USERNAME
                          Username to use with basic HTTP authentication
      -p PASSWORD, --password=PASSWORD
                          Password to use with basic HTTP authentication
      -M, --multi         each argument is considered to be a different site
      -D HEADER, --header=HEADER
                          custom header of the form Header: Value (repeat for
                          multiple headers)
      --url-file-path=URL_FILE_PATH
                          get starting URLs from a line-separated file
      -t TYPES, --types=TYPES
                          Comma-separated values of tags to look for when
                          crawling a site. Default (and supported types):
                          a,img,link,script
      -T TIMEOUT, --timeout=TIMEOUT
                          Seconds to wait before considering that a page timed
                          out (default = 10)
      -C, --strict        Does not strip href and src attributes from
                          whitespaces
      -P, --progress      Prints crawler progress in the console
      -N, --run-once      Only crawl the first page (eq. to depth=0)
      -d DEPTH, --depth=DEPTH
                          Maximum crawl depth (default = 1)
      -e, --prefer-server-encoding
                          Prefer server encoding if specified. Else detect
                          encoding
      --check-presence=CONTENT_PRESENCE
                          Check presence of raw or HTML content on all pages.
                          e.g., <tag attr1="val">regex:content</tag>. Content
                          can be either regex:pattern or plain content
      --check-absence=CONTENT_ABSENCE
                          Check absence of raw or HTML content on all pages.
                          e.g., <tag attr1="val">regex:content</tag>. Content
                          can be either regex:pattern or plain content
      --check-presence-once=CONTENT_PRESENCE_ONCE
                          Check presence of raw or HTML content for one page:
                          path,content, e.g.,: /path,<tag
                          attr1="val">regex:content</tag>. Content can be either
                          regex:pattern or plain content. Path can be either
                          relative or absolute with domain.
      --check-absence-once=CONTENT_ABSENCE_ONCE
                          Check absence of raw or HTML content for one page:
                          path,content, e.g.,path,<tag
                          attr1="val">regex:content</tag>. Content can be either
                          regex:pattern or plain content. Path can be either
                          relative or absolute with domain.
      -S, --show-source   Show source of links (html) in the report.
      --allow-insecure-content
                          Allow insecure content for HTTPS sites with
                          certificate errors

    Performance Options:
      These options can impact the performance of the crawler.

      -w WORKERS, --workers=WORKERS
                          Number of workers to spawn (default = 1)
      -m MODE, --mode=MODE
                          Types of workers: thread (default), process, or green
      -R PARSER, --parser=PARSER
                          Types of HTML parse: html.parser (default) or lxml

    Output Options:
      These options change the output of the crawler.

      -f FORMAT, --format=FORMAT
                          Format of the report: plain (default)
      -o OUTPUT, --output=OUTPUT
                          Path of the file where the report will be printed.
      -W WHEN, --when=WHEN
                          When to print the report. error (only if a
                          crawling error occurs) or always (default)
      -E REPORT_TYPE, --report-type=REPORT_TYPE
                          Type of report to print: errors (default, summary and
                          erroneous links), summary, all (summary and all links)
      -c, --console       Prints report to the console in addition to other
                          output options such as file or email.

    Email Options:
      These options allows the crawler to send a report by email.

      -a ADDRESS, --address=ADDRESS
                          Comma-separated list of email addresses used to send a
                          report
      --from=FROM_ADDRESS
                          Email address to use in the from field of the email
                          (optional)
      -s SMTP, --smtp=SMTP
                          Host of the smtp server
      --port=PORT         Port of the smtp server (optional)
      --tls               Use TLS with the email server.
      --subject=SUBJECT   Subject of the email (optional)
      --smtp-username=SMTP_USERNAME
                          Username to use with the smtp server (optional)
      --smtp-password=SMTP_PASSWORD
                          Password to use with the smtp server (optional)

Usage Example
-------------

Crawl a site and show progress
  ``pylinkvalidate.py --progress http://example.com/``

Crawl a site starting from 2 URLs
  ``pylinkvalidate.py http://example.com/ http://example2.com/``

Crawl a site (example.com) and all pages belonging to another host
  ``pylinkvalidate.py -H additionalhost.com http://example.com/``

Report status of all links (even successful ones)
  ``pylinkvalidate.py --report-type=all http://example.com/``

Report status of all links and HTML show source of these links
  ``pylinkvalidate.py --report-type=all --show-source http://example.com/``

Only crawl starting URLs and access all linked resources
  ``pylinkvalidate.py --run-once http://example.com/``

Crawl two levels (one more than run-once) and access all linked resources
  ``pylinkvalidate.py --depth=1 http://example.com/``

Only access links (a href) and ignore images, stylesheets and scripts
  ``pylinkvalidate.py --types=a http://example.com/``

Crawl a site with 4 threads (default is one thread)
  ``pylinkvalidate.py --workers=4 http://example.com/``

Crawl a site with 4 processes (default is one thread)
  ``pylinkvalidate.py --mode=process --workers=4 http://example.com/``

Crawl a site and use LXML to parse HTML (faster, must be installed)
  ``pylinkvalidate.py --parser=LXML http://example.com/``

Print debugging info
  ``pylinkvalidate.py --verbose=2 http://example.com/``

Change User-Agent request header
  ``pylinkvalidate.py --header="User-Agent: Mozilla/5.0" http://example.com/``

Crawl multiple sites and report results per site
  ``pylinkvalidate.py --multi http://example.com/ http://www.example2.net/``

Check that all HTML pages have a body tag with a specific class:
  ``pylinkvalidate.py --check-content '<body class="test"></body>' http://example.com/``

Check that no HTML pages have a paragraph tag with a pattern:
  ``pylinkvalidate.py --check-absence '<p>regex:Hello\s+World</body>' http://example.com/``

Check that robots.txt have a Disallow none:
  ``pylinkvalidate.py --check-content-once '/robots.txt,regex:^Disallow:\s*$' http://example.com/``

Allow insecure content for HTTPS sites with certificate errors [SSL: CERTIFICATE_VERIFY_FAILED]
  ``pylinkvalidate.py --allow-insecure-content https://self-signed.example.com/``


API Usage
---------

To crawl a site from a single URL:

.. code-block:: python

  from pylinkvalidator.api import crawl
  crawled_site = crawl("http://www.example.com/")
  number_of_crawled_pages = len(crawled_site.pages)
  number_of_errors = len(crawled_sites.error_pages)


To crawl a site and pass some configuration options (the same supported by the
command line interface):


.. code-block:: python

  from pylinkvalidator.api import crawl_with_options
  crawled_site = crawl_with_options(["http://www.example.com/"], {"run-once":
      True, "workers": 10})
  number_of_crawled_pages = len(crawled_site.pages)
  number_of_errors = len(crawled_sites.error_pages)


FAQ and Troubleshooting
-----------------------

I cannot find pylinkvalidate.py on Windows with virtualenv
  This is a known problem with virtualenv on windows. The interpreter is
  different than the one used by the virtualenv. Prefix pylinkvalidate.py with the
  full path: ``python c:\myvirtualenv\Scripts\pylinkvalidate.py``

I see Exception KeyError ... module 'threading' when using --mode=green
  This output is generally harmless and is generated by gevent patching the
  python thread module. If someone knows how to make it go away, patches are
  more than welcome :-)


License
-------

This software is licensed under the `New BSD License`. See the `LICENSE` file
in the for the full license text. It includes the beautifulsoup library which
is licensed under the MIT license.
