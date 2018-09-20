

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
 * Install stdlog library using make configure
 * make build
./main -log=<loglevel>
 * Specifying log level is optional, default log level is 'info'
  
# Assumptions:
 * Initially only one node is present in the ring
 * Two consecutive join/leave ring command, should have an offset of at-least stabilization period, defined in the configuration
 * Mode of leave-ring is randomized, with 50-50 chances to be 'orderly' or 'immediate'
 * The key to be inserted has to be an integer between 0 and ring size, last value excluded

### Function Description

The Function below provide the description of their working in the program
```
func read_from_channels(Node_Pointer *map[int]Node, input_channel <-chan []byte, output_channel chan<- []byte, key int)
Description: This function provides interface between our main program and the commnand function which are present in form of json messages. This function also helps in starting and continuing our go routines
 ```
 ```
func node_commands(n_1 Node, input_channel <-chan []byte, output_channel chan<- []byte)
Description: The function node_commands has the switch case commands to perform different function for the nodes present in the Deployment
```
```
func (n_1 Node) join_chord_deployment(sponsor_node int, key int)
Description: The function helps in joining a new node to deployment with the help of the sponsor nodes
```
```
func (n_1 Node) build_finger_table(target_node int, response_node int)
Description: The function helps in building the finger table for the newly joined Nodes
```
```
func (n_1 Node) find_successor(target_node int, response_node int)
Description: The function helps in finding the successor of the given target node using the node n_1 and send the data back to the response node who has requested for it.
```
```
func (n_1 Node) closest_preceding_node(target_node int)
Description: The function helps in finding the closest predecessor in the finger table of the node n_1
```
```
func (n_1 Node) find_predecessor(target_node int)
Description: The function find predecessor helps in finding the predecessor of the target node using the node n_1
```
```
func (n_1 Node) notify(target_node int)
Description: The function notify helps in notifyon the target node about it own predecessor
```
```
func (n_1 Node) stabilize()
Description: The function stabilize helps in stabilizing the state of each whenever a new node enters the deployment or leaves the current deployment
```
```
func (n_1 Node) store_data(target_node int, Communication_holder Data)
Description: The function helps in storing the data in form of key value pair using the node n_1 in the target node
```
```
func (n_1 Node) retrieve_data(target_node int, Communication_holder Data)
Description: The function helps in retrieving the data from the target node using node n_1
```
```
func (n_1 Node) leave_deployment(target_node int, Communication_holder Data)
Description: The function helps in removing a node from the current chord deployment
```
```
func (n_1 Node) remove_data(target_node int, Communication_holder Data)
Description: Function removes the data from the target node using the node n_1
```

```
func (n_1 Node) fix_finger_table()
Description: The function helps in fixing the entries of the finger tables of the nodes whenever a new node enters the system Deployment
```
```
func run_stabilize()
Description: The function run through go routine enabled in the main and run the stabilize, fix_finger_table and notify functions for each active node in the system periodically
```

```
func sending_request(ex Exchange)
Description: The helps in forming a json message that is send to the command function to execute the required messages from the nodes
```
```
func sending_reply(message_to int, channel_key int)
Description: Function sending reply helps in sending the outcome of the action functions to the response Nodes, using the output channels
```
```
func send_node_data)(n_1 Node, channel_key int)
Description: The function helps in sending back the node data across the channels
```

### Testing the Program
Kindly enter the commands in the format below for testing and running the program,please do note the we have few commented code for testing in the program itself which might help you in testing the functions.

```
ex2 := Exchange{"find-ring-successor", 1, 1, 1, "EMPTY", Communication_holder, 0}
sending_request(ex2)
time.Sleep(1 * time.Second)
err := json.Unmarshal(<-output_channel[1], &node_data)
if err != nil {
	panic(err)
}
x := *node_data
fmt.Println("The successor of the node is:", x)
 
 
ex4 := Exchange{"join-ring", 8, 8, 1, "EMPTY", Communication_holder, 0}
sending_request(ex4)
time.Sleep(1 * time.Second)
ex5 := Exchange{"join-ring", 14, 14, 8, "EMPTY", Communication_holder, 0}
sending_request(ex5)
time.Sleep(1 * time.Second)
 
 
Communication_holder.Key_1 = 1
Communication_holder.Values = "Rajat"
ex8 := Exchange{"put-data", 21, 1, 21, "EMPTY", Communication_holder, 0}
sending_request(ex8)
time.Sleep(1 * time.Second)
Communication_holder.Key_1 = 1
ex9 := Exchange{"get-data", 21, 1, 21, "EMPTY", Communication_holder, 0}
sending_request(ex9)
time.Sleep(1 * time.Second)

 ```

# Configuration Details:
 * JSON format must to be maintained
 * ring.size specifies the total number possible in the ring.
 
 * If ring.size is 5 implies total size is 2^5 = 32.
 * startup.node.id specifies the IP address of the node, that is present at the ring startup [Assumption 1]
 * stabilize.period.millis specifies the period in milli seconds, at which coordinator will send the stabilize command to all the nodes

 * Recommended : Not to reduce the the stabilization period to be less than 7 seconds
 * liveChanges specifies the live commands/changes that coordinator will manage
 * id specifies the node IP address to which the command would be sent
 * timeInMillis specifies the time offset in milli seconds of the command to be sent after
 * action specifies join-ring/leave-ring command to be send to a node
 * query specifies the hash query put/get/remove to be sent to a node
 * data option comes with query to specify the key/value pair, needed to complete the query
 * Either action or query should be used for sending the command
 * query and data must be defined together
 * For get/remove query speficy only the key
 * For put query specify in the format key=value
 
 ## Authors
 
* **Sowmya Rampatruni**
* **Rajat Patel**

