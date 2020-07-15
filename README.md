# Prometheus Exporter for Content Security Policy violations

`csp_exporter` is a web server that listens for reports of [Content Security
Policy][csp] violations and exposes the reports as Prometheus metrics so that
you can incorporate CSP violations into your normal monitoring and alerting
process.

## How to Use

First, download and build `csp_exporter`. If you have the Go compiler tools
installed, this should be a familiar process. Either `git clone
https://git.burwell.io/csp_exporter` and run `go build`, or `go get
bnbl.io/csp_exporter`.

Next, start the server. You'll probably want to run it under some kind of
supervisor depending on your OS. By default, a web server accepting CSP reports
is started on port 80 and a Prometheus metrics server is started on port 9477.
You can override these by setting the `COLLECTOR_BIND_ADDR` and `PROM_BIND_ADDR`
environment variables, respectively.

Once you have the server running, add the appropriate `report-uri` directive to
your content security policy. For example, you might add the following header to
your HTTP responses:

```
Content-Security-Policy: default-src 'none'; report-uri https://csp-exporter.example.com/report/csp/mysite
```

Note the `/report/csp/mysite` path. The `csp_exporter` accepts reports sent to
`/report/csp/<app>`, where `<app>` can be any URL path fragment. Whatever the
`<app>` is set to will be included as the value for the `app` label in your
metrics; this allows you to use `csp_exporter` to collect violation reports for
different websites, test different policy versions, etc.

Finally, configure Prometheus to scrape metrics from `csp_exporter` by adding
something like the following to your `prometheus.yml`:

```yml
scrape_configs:
- job_name: "csp"
  static_configs:
  - targets: ["cspexporter.intra.example.com:9477"]
```

You will now start to accumulate `csp_violation_reports_total` metrics in your
Prometheus system. The labels are derived from the fields provided in the
violation reports and should allow for very granular queries. If you are not
interested in high granularity or are concerned with recording many discrete
time series, you may want to add a [`relabel_config`][relabel] to your CSP
scrape job to drop some of the labels.

## Contributing

Send patches to <ben@benburwell.com>. For instructions about how to mail
patches, see [`git-send-email(1)`][gitsendemail], [the Pro Git book][progit], or
<https://git-send-email.io>.

## License

MIT

[csp]: https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
[relabel]: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config
[gitsendemail]: https://git-scm.com/docs/git-send-email
[progit]: https://git-scm.com/book/en/v2/Distributed-Git-Contributing-to-a-Project#_project_over_email
