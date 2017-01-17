{\rtf1\ansi\ansicpg1252\deff0\deflang1033{\fonttbl{\f0\fnil\fcharset0 Courier New;}}
{\*\generator Msftedit 5.41.21.2510;}\viewkind4\uc1\pard\lang3082\f0\fs22 /*\par
Licensed to the Apache Software Foundation (ASF) under one\par
or more contributor license agreements.  See the NOTICE file\par
distributed with this work for additional information\par
regarding copyright ownership.  The ASF licenses this file\par
to you under the Apache License, Version 2.0 (the\par
"License"); you may not use this file except in compliance\par
with the License.  You may obtain a copy of the License at\par
\par
  http://www.apache.org/licenses/LICENSE-2.0\par
\par
Unless required by applicable law or agreed to in writing,\par
software distributed under the License is distributed on an\par
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY\par
KIND, either express or implied.  See the License for the\par
specific language governing permissions and limitations\par
under the License.\par
*/\par
\par
package main\par
\par
import (\par
\tab "errors"\par
\tab "fmt"\par
\tab "strconv"\par
\tab "encoding/json"\par
\tab "time"\par
\tab "strings"\par
\par
\tab "github.com/openblockchain/obc-peer/openchain/chaincode/shim"\par
)\par
\par
// SimpleChaincode example simple Chaincode implementation\par
type SimpleChaincode struct \{\par
\}\par
\par
var marbleIndexStr = "_marbleindex"\tab\tab\tab\tab //name for the key/value that will store a list of all known marbles\par
var openTradesStr = "_opentrades"\tab\tab\tab\tab //name for the key/value that will store all open trades\par
\par
type Marble struct\{\par
\tab Name string `json:"name"`\tab\tab\tab\tab\tab //the fieldtags are needed to keep case from bouncing around\par
\tab Color string `json:"color"`\par
\tab Size int `json:"size"`\par
\tab User string `json:"user"`\par
\}\par
\par
type Description struct\{\par
\tab Color string `json:"color"`\par
\tab Size int `json:"size"`\par
\}\par
\par
type AnOpenTrade struct\{\par
\tab User string `json:"user"`\tab\tab\tab\tab\tab //user who created the open trade order\par
\tab Timestamp int64 `json:"timestamp"`\tab\tab\tab //utc timestamp of creation\par
\tab Want Description  `json:"want"`\tab\tab\tab\tab //description of desired marble\par
\tab Willing []Description `json:"willing"`\tab\tab //array of marbles willing to trade away\par
\}\par
\par
type AllTrades struct\{\par
\tab OpenTrades []AnOpenTrade `json:"open_trades"`\par
\}\par
\par
// ============================================================================================================================\par
// Main\par
// ============================================================================================================================\par
func main() \{\par
\tab err := shim.Start(new(SimpleChaincode))\par
\tab if err != nil \{\par
\tab\tab fmt.Printf("Error starting Simple chaincode: %s", err)\par
\tab\}\par
\}\par
\par
// ============================================================================================================================\par
// Init - reset all the things\par
// ============================================================================================================================\par
func (t *SimpleChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) \{\par
\tab var Aval int\par
\tab var err error\par
\par
\tab if len(args) != 1 \{\par
\tab\tab return nil, errors.New("Incorrect number of arguments. Expecting 1")\par
\tab\}\par
\par
\tab // Initialize the chaincode\par
\tab Aval, err = strconv.Atoi(args[0])\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("Expecting integer value for asset holding")\par
\tab\}\par
\par
\tab // Write the state to the ledger\par
\tab err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))\tab\tab\tab\tab //making a test var "abc", I find it handy to read/write to it right away to test the network\par
\tab if err != nil \{\par
\tab\tab return nil, err\par
\tab\}\par
\tab\par
\tab var empty []string\par
\tab jsonAsBytes, _ := json.Marshal(empty)\tab\tab\tab\tab\tab\tab\tab\tab //marshal an emtpy array of strings to clear the index\par
\tab err = stub.PutState(marbleIndexStr, jsonAsBytes)\par
\tab if err != nil \{\par
\tab\tab return nil, err\par
\tab\}\par
\tab\par
\tab var trades AllTrades\par
\tab jsonAsBytes, _ = json.Marshal(trades)\tab\tab\tab\tab\tab\tab\tab\tab //clear the open trade struct\par
\tab err = stub.PutState(openTradesStr, jsonAsBytes)\par
\tab if err != nil \{\par
\tab\tab return nil, err\par
\tab\}\par
\tab\par
\tab return nil, nil\par
\}\par
\par
// ============================================================================================================================\par
// Run - Our entry point for Invokcations\par
// ============================================================================================================================\par
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) \{\par
\tab fmt.Println("run is running " + function)\par
\par
\tab // Handle different functions\par
\tab if function == "init" \{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //initialize the chaincode state, used as reset\par
\tab\tab return t.init(stub, args)\par
\tab\} else if function == "delete" \{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //deletes an entity from its state\par
\tab\tab res, err := t.Delete(stub, args)\par
\tab\tab cleanTrades(stub)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //lets make sure all open trades are still valid\par
\tab\tab return res, err\par
\tab\} else if function == "write" \{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //writes a value to the chaincode state\par
\tab\tab return t.Write(stub, args)\par
\tab\} else if function == "init_marble" \{\tab\tab\tab\tab\tab\tab\tab\tab\tab //create a new marble\par
\tab\tab return t.init_marble(stub, args)\par
\tab\} else if function == "set_user" \{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //change owner of a marble\par
\tab\tab res, err := t.set_user(stub, args)\par
\tab\tab cleanTrades(stub)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //lets make sure all open trades are still valid\par
\tab\tab return res, err\par
\tab\} else if function == "open_trade" \{\tab\tab\tab\tab\tab\tab\tab\tab\tab //create a new trade order\par
\tab\tab return t.open_trade(stub, args)\par
\tab\} else if function == "perform_trade" \{\tab\tab\tab\tab\tab\tab\tab\tab\tab //forfill an open trade order\par
\tab\tab res, err := t.perform_trade(stub, args)\par
\tab\tab cleanTrades(stub)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //lets clean just in case\par
\tab\tab return res, err\par
\tab\} else if function == "remove_trade" \{\tab\tab\tab\tab\tab\tab\tab\tab\tab //cancel an open trade order\par
\tab\tab return t.remove_trade(stub, args)\par
\tab\}\par
\tab fmt.Println("run did not find func: " + function)\tab\tab\tab\tab\tab\tab //error\par
\par
\tab return nil, errors.New("Received unknown function invocation")\par
\}\par
\par
// ============================================================================================================================\par
// Query - Our entry point for Queries\par
// ============================================================================================================================\par
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) \{\par
\tab fmt.Println("query is running " + function)\par
\par
\tab // Handle different functions\par
\tab if function == "read" \{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //read a variable\par
\tab\tab return t.read(stub, args)\par
\tab\}\par
\tab fmt.Println("query did not find func: " + function)\tab\tab\tab\tab\tab\tab //error\par
\par
\tab return nil, errors.New("Received unknown function query")\par
\}\par
\par
// ============================================================================================================================\par
// Read - read a variable from chaincode state\par
// ============================================================================================================================\par
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) \{\par
\tab var name, jsonResp string\par
\tab var err error\par
\par
\tab if len(args) != 1 \{\par
\tab\tab return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")\par
\tab\}\par
\par
\tab name = args[0]\par
\tab valAsbytes, err := stub.GetState(name)\tab\tab\tab\tab\tab\tab\tab\tab\tab //get the var from chaincode state\par
\tab if err != nil \{\par
\tab\tab jsonResp = "\{\\"Error\\":\\"Failed to get state for " + name + "\\"\}"\par
\tab\tab return nil, errors.New(jsonResp)\par
\tab\}\par
\par
\tab return valAsbytes, nil\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //send it onward\par
\}\par
\par
// ============================================================================================================================\par
// Delete - remove a key/value pair from state\par
// ============================================================================================================================\par
func (t *SimpleChaincode) Delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) \{\par
\tab if len(args) != 1 \{\par
\tab\tab return nil, errors.New("Incorrect number of arguments. Expecting 1")\par
\tab\}\par
\tab\par
\tab name := args[0]\par
\tab err := stub.DelState(name)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //remove the key from chaincode state\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("Failed to delete state")\par
\tab\}\par
\par
\tab //get the marble index\par
\tab marblesAsBytes, err := stub.GetState(marbleIndexStr)\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("Failed to get marble index")\par
\tab\}\par
\tab var marbleIndex []string\par
\tab json.Unmarshal(marblesAsBytes, &marbleIndex)\tab\tab\tab\tab\tab\tab\tab\tab //un stringify it aka JSON.parse()\par
\tab\par
\tab //remove marble from index\par
\tab for i,val := range marbleIndex\{\par
\tab\tab fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)\par
\tab\tab if val == name\{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //find the correct marble\par
\tab\tab\tab fmt.Println("found marble")\par
\tab\tab\tab marbleIndex = append(marbleIndex[:i], marbleIndex[i+1:]...)\tab\tab\tab //remove it\par
\tab\tab\tab for x:= range marbleIndex\{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //debug prints...\par
\tab\tab\tab\tab fmt.Println(string(x) + " - " + marbleIndex[x])\par
\tab\tab\tab\}\par
\tab\tab\tab break\par
\tab\tab\}\par
\tab\}\par
\tab jsonAsBytes, _ := json.Marshal(marbleIndex)\tab\tab\tab\tab\tab\tab\tab\tab\tab //save new index\par
\tab err = stub.PutState(marbleIndexStr, jsonAsBytes)\par
\tab return nil, nil\par
\}\par
\par
// ============================================================================================================================\par
// Write - write variable into chaincode state\par
// ============================================================================================================================\par
func (t *SimpleChaincode) Write(stub *shim.ChaincodeStub, args []string) ([]byte, error) \{\par
\tab var name, value string // Entities\par
\tab var err error\par
\tab fmt.Println("running write()")\par
\par
\tab if len(args) != 2 \{\par
\tab\tab return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")\par
\tab\}\par
\par
\tab name = args[0]\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //rename for funsies\par
\tab value = args[1]\par
\tab err = stub.PutState(name, []byte(value))\tab\tab\tab\tab\tab\tab\tab\tab //write the variable into the chaincode state\par
\tab if err != nil \{\par
\tab\tab return nil, err\par
\tab\}\par
\tab return nil, nil\par
\}\par
\par
// ============================================================================================================================\par
// Init Marble - create a new marble, store into chaincode state\par
// ============================================================================================================================\par
func (t *SimpleChaincode) init_marble(stub *shim.ChaincodeStub, args []string) ([]byte, error) \{\par
\tab var err error\par
\par
\tab //   0       1       2     3\par
\tab // "asdf", "blue", "35", "bob"\par
\tab if len(args) != 4 \{\par
\tab\tab return nil, errors.New("Incorrect number of arguments. Expecting 4")\par
\tab\}\par
\par
\tab fmt.Println("- start init marble")\par
\tab if len(args[0]) <= 0 \{\par
\tab\tab return nil, errors.New("1st argument must be a non-empty string")\par
\tab\}\par
\tab if len(args[1]) <= 0 \{\par
\tab\tab return nil, errors.New("2nd argument must be a non-empty string")\par
\tab\}\par
\tab if len(args[2]) <= 0 \{\par
\tab\tab return nil, errors.New("3rd argument must be a non-empty string")\par
\tab\}\par
\tab if len(args[3]) <= 0 \{\par
\tab\tab return nil, errors.New("4th argument must be a non-empty string")\par
\tab\}\par
\tab\par
\tab size, err := strconv.Atoi(args[2])\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("3rd argument must be a numeric string")\par
\tab\}\par
\tab\par
\tab color := strings.ToLower(args[1])\par
\tab user := strings.ToLower(args[3])\par
\par
\tab str := `\{"name": "` + args[0] + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "user": "` + user + `"\}`\par
\tab err = stub.PutState(args[0], []byte(str))\tab\tab\tab\tab\tab\tab\tab\tab //store marble with id as key\par
\tab if err != nil \{\par
\tab\tab return nil, err\par
\tab\}\par
\tab\tab\par
\tab //get the marble index\par
\tab marblesAsBytes, err := stub.GetState(marbleIndexStr)\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("Failed to get marble index")\par
\tab\}\par
\tab var marbleIndex []string\par
\tab json.Unmarshal(marblesAsBytes, &marbleIndex)\tab\tab\tab\tab\tab\tab\tab //un stringify it aka JSON.parse()\par
\tab\par
\tab //append\par
\tab marbleIndex = append(marbleIndex, args[0])\tab\tab\tab\tab\tab\tab\tab\tab //add marble name to index list\par
\tab fmt.Println("! marble index: ", marbleIndex)\par
\tab jsonAsBytes, _ := json.Marshal(marbleIndex)\par
\tab err = stub.PutState(marbleIndexStr, jsonAsBytes)\tab\tab\tab\tab\tab\tab //store name of marble\par
\par
\tab fmt.Println("- end init marble")\par
\tab return nil, nil\par
\}\par
\par
// ============================================================================================================================\par
// Set User Permission on Marble\par
// ============================================================================================================================\par
func (t *SimpleChaincode) set_user(stub *shim.ChaincodeStub, args []string) ([]byte, error) \{\par
\tab var err error\par
\tab\par
\tab //   0       1\par
\tab // "name", "bob"\par
\tab if len(args) < 2 \{\par
\tab\tab return nil, errors.New("Incorrect number of arguments. Expecting 2")\par
\tab\}\par
\tab\par
\tab fmt.Println("- start set user")\par
\tab fmt.Println(args[0] + " - " + args[1])\par
\tab marbleAsBytes, err := stub.GetState(args[0])\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("Failed to get thing")\par
\tab\}\par
\tab res := Marble\{\}\par
\tab json.Unmarshal(marbleAsBytes, &res)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //un stringify it aka JSON.parse()\par
\tab res.User = args[1]\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //change the user\par
\tab\par
\tab jsonAsBytes, _ := json.Marshal(res)\par
\tab err = stub.PutState(args[0], jsonAsBytes)\tab\tab\tab\tab\tab\tab\tab\tab //rewrite the marble with id as key\par
\tab if err != nil \{\par
\tab\tab return nil, err\par
\tab\}\par
\tab\par
\tab fmt.Println("- end set user")\par
\tab return nil, nil\par
\}\par
\par
// ============================================================================================================================\par
// Open Trade - create an open trade for a marble you want with marbles you have \par
// ============================================================================================================================\par
func (t *SimpleChaincode) open_trade(stub *shim.ChaincodeStub, args []string) ([]byte, error) \{\par
\tab var err error\par
\tab var will_size int\par
\tab var trade_away Description\par
\tab\par
\tab //\tab 0        1      2     3      4      5       6\par
\tab //["bob", "blue", "16", "red", "16"] *"blue", "35*\par
\tab if len(args) < 5 \{\par
\tab\tab return nil, errors.New("Incorrect number of arguments. Expecting like 5?")\par
\tab\}\par
\tab if len(args)%2 == 0\{\par
\tab\tab return nil, errors.New("Incorrect number of arguments. Expecting an odd number")\par
\tab\}\par
\par
\tab size1, err := strconv.Atoi(args[2])\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("3rd argument must be a numeric string")\par
\tab\}\par
\par
\tab open := AnOpenTrade\{\}\par
\tab open.User = args[0]\par
\tab open.Timestamp = makeTimestamp()\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //use timestamp as an ID\par
\tab open.Want.Color = args[1]\par
\tab open.Want.Size =  size1\par
\tab fmt.Println("- start open trade")\par
\tab jsonAsBytes, _ := json.Marshal(open)\par
\tab err = stub.PutState("_debug1", jsonAsBytes)\par
\par
\tab for i:=3; i < len(args); i++ \{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //create and append each willing trade\par
\tab\tab will_size, err = strconv.Atoi(args[i + 1])\par
\tab\tab if err != nil \{\par
\tab\tab\tab msg := "is not a numeric string " + args[i + 1]\par
\tab\tab\tab fmt.Println(msg)\par
\tab\tab\tab return nil, errors.New(msg)\par
\tab\tab\}\par
\tab\tab\par
\tab\tab trade_away = Description\{\}\par
\tab\tab trade_away.Color = args[i]\par
\tab\tab trade_away.Size =  will_size\par
\tab\tab fmt.Println("! created trade_away: " + args[i])\par
\tab\tab jsonAsBytes, _ = json.Marshal(trade_away)\par
\tab\tab err = stub.PutState("_debug2", jsonAsBytes)\par
\tab\tab\par
\tab\tab open.Willing = append(open.Willing, trade_away)\par
\tab\tab fmt.Println("! appended willing to open")\par
\tab\tab i++;\par
\tab\}\par
\tab\par
\tab //get the open trade struct\par
\tab tradesAsBytes, err := stub.GetState(openTradesStr)\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("Failed to get opentrades")\par
\tab\}\par
\tab var trades AllTrades\par
\tab json.Unmarshal(tradesAsBytes, &trades)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //un stringify it aka JSON.parse()\par
\tab\par
\tab trades.OpenTrades = append(trades.OpenTrades, open);\tab\tab\tab\tab\tab\tab //append to open trades\par
\tab fmt.Println("! appended open to trades")\par
\tab jsonAsBytes, _ = json.Marshal(trades)\par
\tab err = stub.PutState(openTradesStr, jsonAsBytes)\tab\tab\tab\tab\tab\tab\tab\tab //rewrite open orders\par
\tab if err != nil \{\par
\tab\tab return nil, err\par
\tab\}\par
\tab fmt.Println("- end open trade")\par
\tab return nil, nil\par
\}\par
\par
// ============================================================================================================================\par
// Perform Trade - close an open trade and move ownership\par
// ============================================================================================================================\par
func (t *SimpleChaincode) perform_trade(stub *shim.ChaincodeStub, args []string) ([]byte, error) \{\par
\tab var err error\par
\tab\par
\tab //\tab 0\tab\tab 1\tab\tab\tab\tab\tab 2\tab\tab\tab\tab\tab 3\tab\tab\tab\tab 4\tab\tab\tab\tab\tab 5\par
\tab //[data.id, data.closer.user, data.closer.name, data.opener.user, data.opener.color, data.opener.size]\par
\tab if len(args) < 6 \{\par
\tab\tab return nil, errors.New("Incorrect number of arguments. Expecting 6")\par
\tab\}\par
\tab\par
\tab fmt.Println("- start close trade")\par
\tab timestamp, err := strconv.ParseInt(args[0], 10, 64)\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("1st argument must be a numeric string")\par
\tab\}\par
\tab\par
\tab size, err := strconv.Atoi(args[5])\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("6th argument must be a numeric string")\par
\tab\}\par
\tab\par
\tab //get the open trade struct\par
\tab tradesAsBytes, err := stub.GetState(openTradesStr)\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("Failed to get opentrades")\par
\tab\}\par
\tab var trades AllTrades\par
\tab json.Unmarshal(tradesAsBytes, &trades)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //un stringify it aka JSON.parse()\par
\tab\par
\tab for i := range trades.OpenTrades\{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //look for the trade\par
\tab\tab fmt.Println("looking at " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10) + " for " + strconv.FormatInt(timestamp, 10))\par
\tab\tab if trades.OpenTrades[i].Timestamp == timestamp\{\par
\tab\tab\tab fmt.Println("found the trade");\par
\tab\tab\tab\par
\tab\tab\tab\par
\tab\tab\tab marbleAsBytes, err := stub.GetState(args[2])\par
\tab\tab\tab if err != nil \{\par
\tab\tab\tab\tab return nil, errors.New("Failed to get thing")\par
\tab\tab\tab\}\par
\tab\tab\tab closersMarble := Marble\{\}\par
\tab\tab\tab json.Unmarshal(marbleAsBytes, &closersMarble)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //un stringify it aka JSON.parse()\par
\tab\tab\tab\par
\tab\tab\tab //verify if marble meets trade requirements\par
\tab\tab\tab if closersMarble.Color != trades.OpenTrades[i].Want.Color || closersMarble.Size != trades.OpenTrades[i].Want.Size \{\par
\tab\tab\tab\tab msg := "marble in input does not meet trade requriements"\par
\tab\tab\tab\tab fmt.Println(msg)\par
\tab\tab\tab\tab return nil, errors.New(msg)\par
\tab\tab\tab\}\par
\tab\tab\tab\par
\tab\tab\tab marble, e := findMarble4Trade(stub, trades.OpenTrades[i].User, args[4], size)\tab\tab\tab //find a marble that is suitable from opener\par
\tab\tab\tab if(e == nil)\{\par
\tab\tab\tab\tab fmt.Println("! no errors, proceeding")\par
\par
\tab\tab\tab\tab t.set_user(stub, []string\{args[2], trades.OpenTrades[i].User\})\tab\tab\tab\tab\tab\tab //change owner of selected marble, closer -> opener\par
\tab\tab\tab\tab t.set_user(stub, []string\{marble.Name, args[1]\})\tab\tab\tab\tab\tab\tab\tab\tab\tab //change owner of selected marble, opener -> closer\par
\tab\tab\tab\par
\tab\tab\tab\tab trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)\tab\tab //remove trade\par
\tab\tab\tab\tab jsonAsBytes, _ := json.Marshal(trades)\par
\tab\tab\tab\tab err = stub.PutState(openTradesStr, jsonAsBytes)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //rewrite open orders\par
\tab\tab\tab\tab if err != nil \{\par
\tab\tab\tab\tab\tab return nil, err\par
\tab\tab\tab\tab\}\par
\tab\tab\tab\}\par
\tab\tab\}\par
\tab\}\par
\tab fmt.Println("- end close trade")\par
\tab return nil, nil\par
\}\par
\par
// ============================================================================================================================\par
// findMarble4Trade - look for a matching marble that this user owns and return it\par
// ============================================================================================================================\par
func findMarble4Trade(stub *shim.ChaincodeStub, user string, color string, size int )(m Marble, err error)\{\par
\tab var fail Marble;\par
\tab fmt.Println("- start find marble 4 trade")\par
\tab fmt.Println("looking for " + user + ", " + color + ", " + strconv.Itoa(size));\par
\par
\tab //get the marble index\par
\tab marblesAsBytes, err := stub.GetState(marbleIndexStr)\par
\tab if err != nil \{\par
\tab\tab return fail, errors.New("Failed to get marble index")\par
\tab\}\par
\tab var marbleIndex []string\par
\tab json.Unmarshal(marblesAsBytes, &marbleIndex)\tab\tab\tab\tab\tab\tab\tab\tab //un stringify it aka JSON.parse()\par
\tab\par
\tab for i:= range marbleIndex\{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //iter through all the marbles\par
\tab\tab //fmt.Println("looking @ marble name: " + marbleIndex[i]);\par
\par
\tab\tab marbleAsBytes, err := stub.GetState(marbleIndex[i])\tab\tab\tab\tab\tab\tab //grab this marble\par
\tab\tab if err != nil \{\par
\tab\tab\tab return fail, errors.New("Failed to get marble")\par
\tab\tab\}\par
\tab\tab res := Marble\{\}\par
\tab\tab json.Unmarshal(marbleAsBytes, &res)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //un stringify it aka JSON.parse()\par
\tab\tab //fmt.Println("looking @ " + res.User + ", " + res.Color + ", " + strconv.Itoa(res.Size));\par
\tab\tab\par
\tab\tab //check for user && color && size\par
\tab\tab if strings.ToLower(res.User) == strings.ToLower(user) && strings.ToLower(res.Color) == strings.ToLower(color) && res.Size == size\{\par
\tab\tab\tab fmt.Println("found a marble: " + res.Name)\par
\tab\tab\tab fmt.Println("! end find marble 4 trade")\par
\tab\tab\tab return res, nil\par
\tab\tab\}\par
\tab\}\par
\tab\par
\tab fmt.Println("- end find marble 4 trade - error")\par
\tab return fail, errors.New("Did not find marble to use in this trade")\par
\}\par
\par
// ============================================================================================================================\par
// Make Timestamp - create a timestamp in ms\par
// ============================================================================================================================\par
func makeTimestamp() int64 \{\par
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))\par
\}\par
\par
// ============================================================================================================================\par
// Remove Open Trade - close an open trade\par
// ============================================================================================================================\par
func (t *SimpleChaincode) remove_trade(stub *shim.ChaincodeStub, args []string) ([]byte, error) \{\par
\tab var err error\par
\tab\par
\tab //\tab 0\par
\tab //[data.id]\par
\tab if len(args) < 1 \{\par
\tab\tab return nil, errors.New("Incorrect number of arguments. Expecting 1")\par
\tab\}\par
\tab\par
\tab fmt.Println("- start remove trade")\par
\tab timestamp, err := strconv.ParseInt(args[0], 10, 64)\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("1st argument must be a numeric string")\par
\tab\}\par
\tab\par
\tab //get the open trade struct\par
\tab tradesAsBytes, err := stub.GetState(openTradesStr)\par
\tab if err != nil \{\par
\tab\tab return nil, errors.New("Failed to get opentrades")\par
\tab\}\par
\tab var trades AllTrades\par
\tab json.Unmarshal(tradesAsBytes, &trades)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //un stringify it aka JSON.parse()\par
\tab\par
\tab for i := range trades.OpenTrades\{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //look for the trade\par
\tab\tab //fmt.Println("looking at " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10) + " for " + strconv.FormatInt(timestamp, 10))\par
\tab\tab if trades.OpenTrades[i].Timestamp == timestamp\{\par
\tab\tab\tab fmt.Println("found the trade");\par
\tab\tab\tab trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)\tab\tab\tab\tab //remove this trade\par
\tab\tab\tab jsonAsBytes, _ := json.Marshal(trades)\par
\tab\tab\tab err = stub.PutState(openTradesStr, jsonAsBytes)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //rewrite open orders\par
\tab\tab\tab if err != nil \{\par
\tab\tab\tab\tab return nil, err\par
\tab\tab\tab\}\par
\tab\tab\tab break\par
\tab\tab\}\par
\tab\}\par
\tab\par
\tab fmt.Println("- end remove trade")\par
\tab return nil, nil\par
\}\par
\par
// ============================================================================================================================\par
// Clean Up Open Trades - make sure open trades are still possible, remove choices that are no longer possible, remove trades that have no valid choices\par
// ============================================================================================================================\par
func cleanTrades(stub *shim.ChaincodeStub)(err error)\{\par
\tab var didWork = false\par
\tab fmt.Println("- start clean trades")\par
\tab\par
\tab //get the open trade struct\par
\tab tradesAsBytes, err := stub.GetState(openTradesStr)\par
\tab if err != nil \{\par
\tab\tab return errors.New("Failed to get opentrades")\par
\tab\}\par
\tab var trades AllTrades\par
\tab json.Unmarshal(tradesAsBytes, &trades)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //un stringify it aka JSON.parse()\par
\tab\par
\tab fmt.Println("# trades " + strconv.Itoa(len(trades.OpenTrades)))\par
\tab for i:=0; i<len(trades.OpenTrades); \{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //iter over all the known open trades\par
\tab\tab fmt.Println(strconv.Itoa(i) + ": looking at trade " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10))\par
\tab\tab\par
\tab\tab fmt.Println("# options " + strconv.Itoa(len(trades.OpenTrades[i].Willing)))\par
\tab\tab for x:=0; x<len(trades.OpenTrades[i].Willing); \{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //find a marble that is suitable\par
\tab\tab\tab fmt.Println("! on next option " + strconv.Itoa(i) + ":" + strconv.Itoa(x))\par
\tab\tab\tab _, e := findMarble4Trade(stub, trades.OpenTrades[i].User, trades.OpenTrades[i].Willing[x].Color, trades.OpenTrades[i].Willing[x].Size)\par
\tab\tab\tab if(e != nil)\{\par
\tab\tab\tab\tab fmt.Println("! errors with this option, removing option")\par
\tab\tab\tab\tab didWork = true\par
\tab\tab\tab\tab trades.OpenTrades[i].Willing = append(trades.OpenTrades[i].Willing[:x], trades.OpenTrades[i].Willing[x+1:]...)\tab //remove this option\par
\tab\tab\tab\tab x--;\par
\tab\tab\tab\}else\{\par
\tab\tab\tab\tab fmt.Println("! this option is fine")\par
\tab\tab\tab\}\par
\tab\tab\tab\par
\tab\tab\tab x++\par
\tab\tab\tab fmt.Println("! x:" + strconv.Itoa(x))\par
\tab\tab\tab if x >= len(trades.OpenTrades[i].Willing) \{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //things might have shifted, recalcuate\par
\tab\tab\tab\tab break\par
\tab\tab\tab\}\par
\tab\tab\}\par
\tab\tab\par
\tab\tab if len(trades.OpenTrades[i].Willing) == 0 \{\par
\tab\tab\tab fmt.Println("! no more options for this trade, removing trade")\par
\tab\tab\tab didWork = true\par
\tab\tab\tab trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)\tab\tab\tab\tab\tab //remove this trade\par
\tab\tab\tab i--;\par
\tab\tab\}\par
\tab\tab\par
\tab\tab i++\par
\tab\tab fmt.Println("! i:" + strconv.Itoa(i))\par
\tab\tab if i >= len(trades.OpenTrades) \{\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //things might have shifted, recalcuate\par
\tab\tab\tab break\par
\tab\tab\}\par
\tab\}\par
\par
\tab if(didWork)\{\par
\tab\tab fmt.Println("! saving open trade changes")\par
\tab\tab jsonAsBytes, _ := json.Marshal(trades)\par
\tab\tab err = stub.PutState(openTradesStr, jsonAsBytes)\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab\tab //rewrite open orders\par
\tab\tab if err != nil \{\par
\tab\tab\tab return err\par
\tab\tab\}\par
\tab\}else\{\par
\tab\tab fmt.Println("! all open trades are fine")\par
\tab\}\par
\par
\tab fmt.Println("- end clean trades")\par
\tab return nil\par
\}\par
}
 