![fact logo](img/logo.png)


`fact` - FaaS Application Component Tracer - is a tool to collect unified monitoring, logging, and tracing information from serverless applications.

<!-- ![GIF demo](img/demo.gif) -->

Currently, supports: AWS Lambda, OpenWhisk, IBM Cloud Functions, Google Cloud Functions, Azure Functions 

**Usage**
---

### Standalone
```
Fact - FaaS Application Component Tracer - is a tool to collect unified monitoring, logging, and tracing information from serverless applications.

Usage:
  fact [command]

Available Commands:
  help        Help about any command
  tcp         starts tcp collection mode

Flags:
      --config string   config file (default is $HOME/.fact.yaml)
  -c, --continues       Writes to the Output continuously
  -f, --file string     Output File Path
  -h, --help            help for fact
  -o, --output string   Output File Format (required) - csv is the default

```

### Library

The fact library offers multiple way to extend collectors and output formats.
All libraries use the golang io.Reader/io.Writer interfaces for input and output.

If you want to use an existing collector and programmatically react to new traces you can also implement a TraceObserver:
```
    
type PrintObserver struct {
}

func (o PrintObserver) Observe(trace *fact.Trace) {
	fmt.Printf("got trace: %+v",trace);
}

func (o PrintObserver) Close() {}
...
collector := fact.NewTCPCollector(port,threads,connections)
collector.AddObserver(&PrintObserver{})
```

For a detailed overview of provided [collectors]() or [output formats]() have a look at the [docs](/docs).

**Installation Options**
---

1. Download the `fact` binary from Releases tab.
2. Select an appropriate [client library](docs/Clients.md) and follow the instructions.


**How to Contribute**
---

<!-- TODO -->

**Acknowledgements**
---
Fact is under the MIT license, for more check [License](./License).
This project was in part created in the [SMILE Project](https://ise-smile.github.io/) funded by the German Federal Ministry of Education and Research. 