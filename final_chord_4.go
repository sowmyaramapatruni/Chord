// Issue : How do we free the channel or clear it after its use
// How to send the request to same channel simultaneously

package main

import (
	"fmt"
	"math"
	"time"
	//"sync"
	"encoding/json"
)

//Global variables
const total_nodes = 6

var input_channel map[int]chan []byte
var output_channel map[int]chan []byte
var Nodes map[int]Node
var Nodes_orig_Pointer *map[int]Node

//var nodes_successor *Nodes_orig_Pointer.ex.Successor_2
var Communication_holder Data

type Exchange struct {
	Do            string
	Sponsor_node  int
	Target_node   int
	Response_node int
	Mode          string
	Data          Data
	Successor_2   int
}

type Data struct {
	key   int
	value string
}

type Node struct {
	Key          int
	successor    int
	Predecessor  int
	Status       bool
	fix          bool
	finger_table [total_nodes - 1]int
	dht          map[string]string
	ex           Exchange
}

type finger_table_function_reply struct {
	message      string
	Finger_table [total_nodes - 1]int
}

//var wg sync.WaitGroup

// func main(){
//   wg.Add(1)
//   go Coordinator()
//   wg.Wait()
// }
func zeroptr(iptr *int) {
	*iptr = 0
}

func main() {
//	var node_data *int
	input_channel = make(map[int]chan []byte, 1000)
	output_channel = make(map[int]chan []byte, 1000)
	Nodes = make(map[int]Node)
	j := 0
	Nodes[1] = Node{
		Key:          1,
		successor:    8,
		Predecessor:  30,
		Status:       true,
		finger_table: [total_nodes - 1]int{8, 8, 8, 14, 21},
		dht:          make(map[string]string),
	}

	Nodes[8] = Node{
		Key:          8,
		successor:    14,
		Predecessor:  1,
		Status:       true,
		finger_table: [total_nodes - 1]int{14, 14, 14, 21, 30},
		dht:          make(map[string]string),
	}

	Nodes[14] = Node{
		Key:          14,
		successor:    21,
		Predecessor:  8,
		Status:       true,
		finger_table: [total_nodes - 1]int{21, 21, 21, 30, 30},
		dht:          make(map[string]string),
	}

	Nodes[21] = Node{
		Key:          21,
		successor:    30,
		Predecessor:  14,
		Status:       true,
		finger_table: [total_nodes - 1]int{30, 30, 30, 30, 8},
		dht:          make(map[string]string),
	}
	Nodes[30] = Node{
		Key:          30,
		successor:    1,
		Predecessor:  21,
		Status:       true,
		finger_table: [total_nodes - 1]int{1, 1, 8, 8, 14},
		dht:          make(map[string]string),
	}
	for j < int(math.Pow(2, (total_nodes-1))) {
		if Nodes[j].Key == 0 {
			Nodes[j] = Node{
				Key:    j,
				Status: false,
				dht:    make(map[string]string),
			}
		}
		input_channel[j] = make(chan []byte, 1000)
		output_channel[j] = make(chan []byte, 1000)
		Nodes_orig_Pointer = &Nodes

		go read_from_channel(Nodes_orig_Pointer, input_channel[j], output_channel[j], j)
		j++
	}
	//go read_from_channel(Nodes_orig_Pointer, input_channel[8], output_channel[8], j)
	// ex2 := Exchange{"closest_preceding_node", 30, 21, 14, "EMPTY", Communication_holder, 0}
	// sending_request(ex2)
	// time.Sleep(1 * time.Second)
	// err := json.Unmarshal(<-output_channel[14], &node_data)
	// if err != nil {
	 //	panic(err)
	// }
	// x := *node_data
	// fmt.Println(x)
        // a:= Nodes[30].closest_preceding_node(21)
        // fmt.Println(a)
	// ex3 := Exchange{"find-ring-successor", 8, 6, 1, "EMPTY", Communication_holder, 0}
	// sending_request(ex3)
	// time.Sleep(3 * time.Second)
	// err := json.Unmarshal(<-output_channel[1], &node_data)
	// if err != nil {
	// 	panic(err)
	 //}
	// x := *node_data
	// fmt.Println(x)
	ex4 := Exchange{"join-ring", 7, 7, 1, "EMPTY", Communication_holder, 0}
	sending_request(ex4)
	time.Sleep(1 * time.Second)
        Nodes[7].stabilize()
        Nodes[8].stabilize()
        Nodes[1].stabilize()
}

func read_from_channel(Node_pointer *map[int]Node, input_channel <-chan []byte, output_channel chan<- []byte, key int) {
	n_1 := *Node_pointer
	node_commands(n_1[key], input_channel, output_channel)
	time.Sleep(1 * time.Second)
	fmt.Println("A message from the node", key, "has been sent")
}

func node_commands(n_1 Node, input_channel <-chan []byte, output_channel chan<- []byte) {
	err := json.Unmarshal(<-input_channel, &n_1.ex)
	if err != nil {
		panic(err)
	}
	switch execute := n_1.ex.Do; execute {
	case "find-ring-successor":
		fmt.Println("Node:", n_1.ex.Sponsor_node, "for node:", n_1.ex.Target_node, "and send data to node:", n_1.ex.Response_node)
		successor_1 := Node.find_successor(n_1, n_1.ex.Target_node, n_1.ex.Response_node)
		n_1.sending_reply(successor_1, n_1.ex.Response_node)
		if successor_1 > 0 {
			fmt.Println("Node", n_1.Key, "found the successor for node:", n_1.ex.Target_node)
		}
	case "get-ring-fingers":
		fmt.Println("Node:", n_1.Key, "would be now build its finger table")
		Node.build_finger_table(n_1, n_1.ex.Target_node, n_1.ex.Response_node)

	case "join-ring":
		if n_1.Status == false {
			fmt.Println(n_1.Status)
			fmt.Println(n_1.ex.Sponsor_node)
			new_node := Node.join_chord_deployment(n_1, n_1.ex.Response_node, n_1.Key)
			fmt.Println(new_node)
			n_1.sending_reply(new_node.Key, n_1.Key)
			time.Sleep(5 * time.Second)
		} else {
			fmt.Println("The sponsor node is not present in the deployment")
		}
	}
}

func (n_1 Node) find_predecessor(target_node int) int{

    var node_1 *int
    Nodes[n_1.Key].find_successor(target_node,target_node)
    err := json.Unmarshal(<-output_channel[target_node], &node_1)
    if err != nil {
       panic(err)
    }
    successor_node := *node_1

     target_node_temp := successor_node
     flag :=1
     for  { 
         Nodes[target_node_temp].find_successor(target_node_temp,target_node_temp)
         err := json.Unmarshal(<-output_channel[target_node_temp], &node_1)
         if err != nil {
             panic(err)
         }
         successor := *node_1
         fmt.Println("find_predecessor ",successor)
         if successor == target_node {
            flag =0
            fmt.Println("find_predecessor predecessor is ",target_node_temp)
            return target_node_temp
         }else {
               target_node_temp = successor
         }
         if flag ==0 {
           break
         }
    }
    return 0
}

func (n_1 Node) stabilize(){

     successor:= n_1.successor
     predecessor_of_successor:= Nodes[successor].find_predecessor(successor)
    
     if predecessor_of_successor == n_1.Key {
        fmt.Println("stabilize: No action needed")
        return
     }else if predecessor_of_successor > n_1.Key  {
             fmt.Println("stabilize: New node",predecessor_of_successor,"has joined the chord ",n_1.Key," changing its successor to ", predecessor_of_successor)
              n_1.successor = predecessor_of_successor
     }else {
         fmt.Println("stabilize: New node",n_1.Key,"has joined the chord ",n_1.Key," changing its predecessor to ", predecessor_of_successor)
         n_1.Predecessor = predecessor_of_successor
         fmt.Println("stabilize:",n_1.Key,"notifying",successor,"to change its predecessor");
         n_1.notify(successor)         
     }
     
}

func (n_1 Node) notify(target_node int){

     //Nodes[target_node].Predecessor = n_1.Key
     n_2 := Nodes[target_node]
     n_2.Predecessor = n_1.Key
     Nodes[target_node] = n_2

}
func (n_1 Node) join_chord_deployment(sponsor_node int, key int) Node {
	var reply_from_finger_table finger_table_function_reply
	fmt.Println("Node:", n_1.Key, "is joining the chord deployment using node:", sponsor_node)
	fmt.Println("Begining to find the sucessor for node", n_1.Key)
	fmt.Println("Setting the status to true")
	n_1.Status = true
	fmt.Println("Creating a finger table for this node..!!")
	ex_1 := Exchange{"get-ring-fingers", sponsor_node, n_1.Key, n_1.Key, "EMPTY", Communication_holder, n_1.ex.Successor_2}
	sending_request(ex_1)
	err := json.Unmarshal(<-output_channel[n_1.Key], &reply_from_finger_table)
	if err != nil {
		fmt.Println("The finger table creation obstructed")
		panic(err)
	}
	n_1.finger_table = reply_from_finger_table.Finger_table
        n_1.successor = n_1.finger_table[0]
        Nodes[n_1.Key] = n_1
        fmt.Println("Debug n_successor",n_1.successor)
	fmt.Println("Node", n_1.Key, ": has now joined the chord deployment")
	return n_1
}

func (n_1 Node) build_finger_table(target_node int, response_node int) {
	var node_1 *int
	var new_finger_table [total_nodes - 1]int
	for i := 0; i < (total_nodes - 1); {
		//fmt.Println("Finding the Sucessor for the node:", target_node+int(math.Pow(2, float64(i))))
		if (target_node + int(math.Pow(2, float64(i)))) == n_1.successor {
			new_finger_table[i] = n_1.successor
			i++
		} else {
			n_1.find_successor(target_node+int(math.Pow(2, float64(i))), n_1.Key)
                        fmt.Println("debug",n_1.Key)
			err := json.Unmarshal(<-output_channel[n_1.Key], &node_1)
			if err != nil {
				panic(err)
			}
			a2 := *node_1
			new_finger_table[i] = a2
			i++
		}
	}
	fmt.Println("The finger table for the requested node is created: ", new_finger_table)
	fmt.Println("Now sending the created finger_table to the node that requested !!!!")
	send_finger_table := finger_table_function_reply{"finger-table-reply", new_finger_table}
	d, err := json.Marshal(send_finger_table)
	if err != nil {
		panic(err)
	}
	fmt.Println("Node: ", n_1.Key, " sending finger table now to node: ", response_node)
	output_channel[response_node] <- d
	//return new_finger_table
}

func (n_1 Node) find_successor(target_node int, response_node int) int {
	if target_node >= int(math.Pow(2, (total_nodes-1))) {
		target_node = (target_node - int(math.Pow(2, (total_nodes-1))))
	}
	if n_1.Key < n_1.successor && target_node < n_1.successor && target_node >= n_1.Key || n_1.Key >= n_1.successor && n_1.Key <= target_node {
		successor_2 := n_1.successor
		fmt.Println("node:", n_1.successor, "is successor of", target_node)
		fmt.Println("sending this information to node:", response_node)
		n_1.sending_reply(successor_2, response_node)
		//return (successor_2)
	} else {
		fmt.Println("successor not found now finding closest Predecessor")
		successor_1 := n_1.closest_preceding_node(target_node)
		fmt.Println("Now the node", successor_1, "will help in finding the successor for node:", target_node)
		(Nodes[successor_1].find_successor(target_node, response_node))
	}
	return 0
}

func (n_1 Node) closest_preceding_node(target_node int) int {
	min := target_node
	for k := 1; k < len(n_1.finger_table); k++ {
		if min <= n_1.finger_table[k-1] && min <= n_1.finger_table[k] {
			fmt.Println("Node:", n_1.finger_table[k-1], "returned as the closest_preceder of node:", target_node)
			closest_preceder := n_1.finger_table[k-1]
			return closest_preceder
		} else {
			fmt.Println("Node ", n_1.finger_table[k], "returned as the closest_preceder of node:", target_node)
			return (n_1.finger_table[k])
		}
	}
	return 0
}

func sending_request(ex Exchange) {
	fmt.Println("Request message send to", ex.Sponsor_node)
	send, err := json.Marshal(ex)
	if err != nil {
		panic(err)
	}
	input_channel[ex.Sponsor_node] <- send
}

func (n_1 Node) sending_reply(message_to int, channel_key int) {
	fmt.Println("Sending the reply", message_to, "to channel:", channel_key)
	// ex1 := Exchange{"reply-message-successor", 0, 0, 0, "EMPTY", Communication_holder, message_to}
	reply, err := json.Marshal(message_to)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
	fmt.Println("This is the channel key", channel_key)
	output_channel[channel_key] <- reply
	//time.Sleep(10 * time.Second)
}
