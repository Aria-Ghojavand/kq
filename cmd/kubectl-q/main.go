package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/yourname/kubectl-q/internal/util"
    "github.com/yourname/kubectl-q/pkg/engine"
    "github.com/yourname/kubectl-q/pkg/fetcher"
    "github.com/yourname/kubectl-q/pkg/formatter"
    "github.com/yourname/kubectl-q/pkg/parser"
)

func main() {
    var (
        namespace string
        allNS     bool
        output    string
    )

    flag.StringVar(&namespace, "n", "", "Namespace (defaults to current context)")
    flag.BoolVar(&allNS, "all-namespaces", false, "Query across all namespaces")
    flag.StringVar(&output, "o", "table", "Output format: table|json|yaml|csv")
    flag.Parse()

    if flag.NArg() < 1 {
        fmt.Fprintln(os.Stderr, "usage: kubectl q \"SELECT ... FROM <resource> [WHERE ...] [ORDER BY ...]\"")
        os.Exit(2)
    }

    queryStr := flag.Arg(0)
    ctx := context.Background()

    cfg, restMapper, disco, clientset, err := util.BuildKubeClients()
    if err != nil {
        fmt.Fprintln(os.Stderr, "kube client error:", err)
        os.Exit(1)
    }

    if namespace == "" && !allNS {
        ns, err := util.DetectNamespace()
        if err == nil && ns != "" {
            namespace = ns
        }
    }

    q, err := parser.Parse(queryStr)
    if err != nil {
        fmt.Fprintln(os.Stderr, "parse error:", err)
        os.Exit(1)
    }

    rows, err := fetcher.Fetch(ctx, cfg, restMapper, disco, clientset, q, namespace, allNS)
    if err != nil {
        fmt.Fprintln(os.Stderr, "fetch error:", err)
        os.Exit(1)
    }

    res, err := engine.Execute(q, rows)
    if err != nil {
        fmt.Fprintln(os.Stderr, "engine error:", err)
        os.Exit(1)
    }

    switch strings.ToLower(output) {
    case "table":
        formatter.PrintTable(res)
    case "json":
        formatter.PrintJSON(res)
    case "yaml":
        formatter.PrintYAML(res)
    case "csv":
        formatter.PrintCSV(res)
    default:
        fmt.Fprintln(os.Stderr, "unknown -o format, use table|json|yaml|csv")
        os.Exit(2)
    }
}
