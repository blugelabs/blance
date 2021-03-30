blance
======

NOTE: this is a fork of Couchbase's [blance](https://github.com/couchbase/blance) library prior to the BSL license change.

blance implements a straightforward partition assignment algorithm,
using a greedy, heuristic, functional approach.

blance provides features like multiple, user-configurable partition
states (primary, replica, read-only, etc), multi-level containment
hierarchy (shelf/rack/row/zone/datacenter awareness) with configurable
inclusion/exclusion policies, heterogeneous partition weights,
heterogeneous node weights, partition stickiness control, and multi-primary
support.

LICENSE: Apache 2.0

### Usage

See the PlanNextMap() function as a starting point.

### For developers

To get local coverage reports with heatmaps...

    go test -coverprofile=coverage.out -covermode=count && go tool cover -html=coverage.out
