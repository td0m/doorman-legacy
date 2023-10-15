# [WIP] Doorman

My take on simplified access control. Implemented in < 1000 lines of Go.

[![asciicast](https://asciinema.org/a/3y6N8aJoBnQGHmKb2kn3hkOQl.svg)](https://asciinema.org/a/3y6N8aJoBnQGHmKb2kn3hkOQl)

Essentially just a directed acyclic graph db built on top of Postgres with a few added constraints to enforce good structure and conventions.

<!-- ## Philosophy -->
<!---->
<!-- My first encounter with access control systems was -->
<!-- [Gatekeeper](https://github.com/uatuko/gatekeeper); an open source project we -->
<!-- bootstrapped at work. -->
<!---->
<!-- Honestly, I liked a lot of concepts behind it, and it felt much easier to -->
<!-- understand and limited in scope (in a good way) compared to other solutions out -->
<!-- there ([OpenFGA](https://openfga.dev), [Ory Keto](https://ory.sh/keto)). -->
<!---->
<!-- After some time using it, I had an idea for how the concepts could be more generalized -->
<!-- and the codebase simplified. So I set out to build a solution in less than 1,000 lines. -->

<!-- ## Structure and Validation -->
<!---->
<!-- ```mermaid -->
<!-- flowchart LR -->
<!--   subgraph users -->
<!--     bob -->
<!--     alice -->
<!--   end -->
<!--   subgraph collections -->
<!--     management_department -->
<!--   end -->
<!---->
<!--   subgraph roles -->
<!--     order_manager -->
<!--   end -->
<!---->
<!--   subgraph permissions -->
<!--     orders.update -->
<!--     orders.refund -->
<!--     orders.cancel -->
<!--   end -->
<!---->
<!--   subgraph orders -->
<!--     order_a -->
<!--   end -->
<!---->
<!--   alice --> management_department -->
<!--   management_department --> order_manager -->
<!--   order_manager --> orders.update -->
<!--   order_manager --> orders.refund -->
<!--   order_manager --> orders.cancel -->
<!---->
<!--   bob --> order_a -->
<!--   bob --> order_manager -->
<!-- ``` -->

## Quick Start

```
go get github.com/td0m/doorman/cmd/
```

## Usage

Doorman can be used via one of the following:
 - JSON API
 - gRPC service
 - In-Process Go library

TODO: cli via json or grpc?

TODO: rebuild cache endpoint.

## Performance

TODO: some benchmarks, initial test with 100,000,000 relations (1,000,000 users, 1,000,000 resources) and sub ms responses looks good though.

## Custom entities

TODO: explain how custom entity types are supported.
