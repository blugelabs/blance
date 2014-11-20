//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the
//  License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an "AS
//  IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language
//  governing permissions and limitations under the License.

// The blance package provides a partition rebalancing algorithm,
// using a greedy, single-pass, heuristic, functional approach.
package blance

// A PartitionMap represents all the partitions for some logical
// resource, where the partitions are assigned to different nodes and
// with different states.  For example, partition "A...H" is assigned
// to node "x" as a "master" and to node "y" as a "replica".  And,
// partition "I...Z" is assigned node "y" as a "master" and to nodes
// "z" & "x" as "replica".
type PartitionMap map[string]*Partition // Keyed by Partition.Name.

// A Partition represents a distinct, non-overlapping subset (or a
// shard) of some logical resource.
type Partition struct {
	// The Name of a Partition must be unique within a PartitionMap.
	Name string

	// NodesByState is keyed is stateName, and value is an array of
	// node names.  For example, {"master": ["a"], "replica": ["b",
	// "c"]}.
	NodesByState map[string][]string
}

// A PartitionModel lets applications define different states for each
// partition per node, such as "master", "slave", "dead", etc.
// It is keyed by stateName, like "master", "slave", "dead", etc.
type PartitionModel map[string]*PartitionModelState

// A PartitionModelState lets applications define metadata per
// partition model state.  For example, "master" state should have
// different priority and constraints than a "slave" state.
type PartitionModelState struct {
	// Priority of zero is the highest.  e.g., "master" Priority
	// should be < than "slave" Priority, so we can define that
	// as "master" Priority of 0 and "slave" priority of 1.
	Priority int

	// A Constraint defines how many nodes the rebalancing algorithm
	// strives to assign a partition.  For example, for any given
	// partition, perhaps the application wants a 1 node to have
	// "master" state and wants 2 nodes to have "slave" state.
	Constraints int
}

// RebalancePartitions is the main entry point.
func RebalancePartitions(
	prevMap PartitionMap,
	nodesToRemove []string,
	nodesToAdd []string,
	model PartitionModel,
	// Keyed by same key as the key to partitionModel.States, e.g.,
	// "master", "slave", "dead", etc.
	modelStateConstraints map[string]int,
	partitionWeights map[string]int,
) PartitionMap {
	return rebalancePartitions(prevMap, nodesToRemove, nodesToAdd,
		model, modelStateConstraints, partitionWeights)
}