

# Chord - a distributed hash table
# Aim:
Write a GoLang multithreaded application that implements the CHORD protocol/distributed hash table. Assume a CHORD ring of order 2^N for some large constant N, e.g. N=32.

Chord nodes are to be implemented as goroutines in the GoLang. Each node has a randomly assigned GoLang channel variable for receiving messages from other nodes. Each channel has a unique (string) identifier; a node uses its channel's id as its own identifier. Further, each node maintains a bucket (list of) (key, value) pairs of a hash table that is distributed among the nodes that are members of the Chord ring. There is no limit on a node's bucket size besides the available memory to GoLang processes, while the keys and values are assumed to be strings.

Your main GoLang routine (aka coordinator) should spawn some Chord nodes, and then, instruct them to join/leave the Chord ring, as well as get/put/remove key-value pairs from/to the distributed hash table. You may issue such instructions at random or read/load them from a file.

For more details refer [here](https://www.csee.umbc.edu/~kalpakis/Courses/621-sp18/project/prj3.phphttps://www.google.com)

# Prerequisites:
GoLang
GO library stdlog
Linux Environment

How to execute:
⋅⋅* Install stdlog library using make configure
⋅⋅* make build
./main -log=<loglevel>
⋅⋅* Specifying log level is optional, default log level is 'info'
  
# Assumptions:
⋅⋅* Initially only one node is present in the ring
⋅⋅* Two consecutive join/leave ring command, should have an offset of at-least stabilization period, defined in the configuration
⋅⋅* Mode of leave-ring is randomized, with 50-50 chances to be 'orderly' or 'immediate'
⋅⋅* The key to be inserted has to be an integer between 0 and ring size, last value excluded

# Implemention summary:
1. Config: Defines all the events as JSON that needs to be passed as an INPUT (file name is fixed as 'config')
2. Main.go: Reads the config and starts the coordinator
3. Coordinator.go: Starts the required nodes and manages the ring
4. Nodes.go: Defines all the operations of a node
5. Fingertable.go: Defines the prototype of the finger table
6. Jsonhandler.go: Defines the prototype (structs) of all the events/message exchanged by the nodes
7. Dictionary.go: Defines all the operable methods on the global dictionary
8. DHT.go: Defines the prototype for storing the key-value pairs in a node
9. Hashing.go: Implementation of Consistent hashing using SHA-1 on the IP address/Node id
10. Enums.go: Statically defines all the events/message

# Configuration Details:
⋅⋅* JSON format must to be maintained
⋅⋅* ring.size specifies the total number possible in the ring
⋅⋅* If ring.size is 5 implies total size is 2^5 = 32
⋅⋅* startup.node.id specifies the IP address of the node, that is present at the ring startup [Assumption 1]
 ⋅⋅* stabilize.period.millis specifies the period in milli seconds, at which coordinator will send the stabilize command to all the nodes

# Recommended : Not to reduce the the stabilization period to be less than 7 seconds
⋅⋅* liveChanges specifies the live commands/changes that coordinator will manage
⋅⋅* id specifies the node IP address to which the command would be sent
⋅⋅* timeInMillis specifies the time offset in milli seconds of the command to be sent after
⋅⋅* action specifies join-ring/leave-ring command to be send to a node
⋅⋅* query specifies the hash query put/get/remove to be sent to a node
⋅⋅* data option comes with query to specify the key/value pair, needed to complete the query
⋅⋅* Either action or query should be used for sending the command
⋅⋅* query and data must be defined together
⋅⋅* For get/remove query speficy only the key
⋅⋅* For put query specify in the format key=value
