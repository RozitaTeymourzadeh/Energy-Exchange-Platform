package p4

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"sort"
	"strings"
)


/* Show()
*
* To Print blockChain
*
 */
func (bc *BlockChain) Show() string {
	rs := ""
	var idList []int
	for id := range bc.Chain {
		idList = append(idList, int(id))
	}
	sort.Ints(idList)
	for _, id := range idList {
		var hashs []string
		for _, block := range bc.Chain[int32(id)] {
			hashs = append(hashs, block.Header.Hash+"<="+block.Header.ParentHash)
		}
		sort.Strings(hashs)
		rs += fmt.Sprintf("%v: ", id)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}
	sum := sha3.Sum256([]byte(rs))
	rs = fmt.Sprintf("This is the BlockChain: %s\n", hex.EncodeToString(sum[:])) + rs
	return rs
}

/* Canonical
*
* To print Canonical branch
*
 */
func (blockChain *BlockChain) Canonical() string {
	rs := ""
	forksBlocks:= blockChain.GetLatestBlocks()
	for i, currentBlock := range forksBlocks {
		height := blockChain.Length
		rs += "\n"
		rs += fmt.Sprintf("Chain # %d:\n ", i)
		for  height > 0{
			rs += fmt.Sprintf("height=%d, timestamp=%d, hash=%s, parentHash=%s, size=%d , value=%s\n",
				currentBlock.Header.Height, currentBlock.Header.Timestamp, currentBlock.Header.Hash,
				currentBlock.Header.ParentHash, currentBlock.Header.Size, currentBlock.Value)
			currentBlock, _ = blockChain.GetBlock(currentBlock.Header.Height-1, currentBlock.Header.ParentHash)
			height = height - 1
		}
	}
	rs += "\n"
	fmt.Println(rs)
 return rs
}

/* GetBlock()
*
* Return block from the blocklist in fork
*
 */
func (blockChain *BlockChain) GetBlock(height int32, hash string) (Block, bool) {
	isAvali := false
	block:= Block{}
	blocks := blockChain.Chain[height]
	lenBlock := len(blocks)
	if lenBlock != 0{
		for i:=0; i < lenBlock; i++{
			if blocks[i].Header.Hash == hash {
				block = blocks[i]
				isAvali = true
				return block, isAvali
			}
		}
	}
	return block, isAvali
}


/*-------------------------STRUCT---------------------------------------------------*/
/* Struct data structure for variables
/*-------------------------STRUCT---------------------------------------------------*/

/* BlockChain struct
*
* To Define blockChain variables
*
 */
type BlockChain struct {
	Chain  map[int32][]Block `json:"chain"`
	Length int32             `json:"length"`
}

/*-------------------------INITIALIZATION---------------------------------------------------*/
/* Initialize blockChain
/*-------------------------INITIALIZATION---------------------------------------------------*/

/* Initial
*
* To Initialize blockChain
*
*/
func (blockChain *BlockChain) Initial() {
	blockChain.Chain = make(map[int32][]Block)
	blockChain.Length = 0
}

/* NewBlockChain
*
* To Create empty New Block Chain
*
*/
func NewBlockChain() BlockChain{
	return BlockChain{
		Chain: make(map[int32][]Block),
		Length: 0,
	}
}

/*-------------------------MASTER---------------------------------------------------*/
/* Main function
/*-------------------------MASTER---------------------------------------------------*/

/* GetLatestBlocks()
*
* Returns the list of blocks of height "BlockChain.length"
*
 */
func (blockChain *BlockChain) GetLatestBlocks() []Block {
	height := blockChain.Length
	blocks := blockChain.Chain[height]
	return blocks
}

/* GetParentBlock()
*
* Takes a block as the parameter, and returns its parent block.
*
 */
func (blockChain *BlockChain) GetParentBlock(block Block) Block {

	parentBlock:= Block{}
	blocks := blockChain.Chain[block.Header.Height-1]
	lenBlock := len(blocks)
	if lenBlock != 0{
		for i:=0; i < lenBlock; i++{
			if blocks[i].Header.Hash == block.Header.ParentHash {
				parentBlock = blocks[i]
				return parentBlock
			}
		}
	}
	return parentBlock
}


/* Get
*
* To return blocks in chain with certain height
* @input: height int32
* @output: blockChain.Chain[height]
 */
func (blockChain *BlockChain) Get(height int32) ([]Block, bool) {
	found := false
	if blockChain.Chain[height] != nil{
		found = true
		return blockChain.Chain[height], found
	}
	return blockChain.Chain[height], found
}


/* ConvertIntToString()
*
* Function to convert Int32 to String
*
 */

func ConvertIntToString(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

/* Insert
*
* To insert block into blockchain
*
 */
func (blockChain *BlockChain) Insert(block Block) {
	if block.Header.Height > blockChain.Length {
		blockChain.Length = block.Header.Height
	}
	heightBlocks := blockChain.Chain[block.Header.Height]
	if heightBlocks == nil { // return empty block if heght is zero
		heightBlocks = []Block{}
	}
	for _, heightBlock := range heightBlocks { // find simmilar hash in blockchain
		if heightBlock.Header.Hash == block.Header.Hash {
			return
		}
	}
	// append to blockChain
	blockChain.Chain[block.Header.Height] = append(heightBlocks, block)
}

/*-------------------------JSON HELPER---------------------------------------------------*/
/* Serialize and decerialization
/*-------------------------JSON HELPER---------------------------------------------------*/

/* EncodeToJSON
*
* To encode block into Json block
* @input: jsonString string
* @output: string, error
*
 */
func (blockChain *BlockChain) EncodeToJSON() (string, error) {
	jsonBytes, err := json.Marshal(blockChain)
	return string(jsonBytes), err
}

/* DecodeFromJSON
*
* To decerialize JSON to blockChain
* @input: jsonString string
* @output: blockChain
*
 */
// Deserializes the given JSON to a block chain instance.
func (blockChain *BlockChain) DecodeFromJSON(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), blockChain)
}

/* UnmarshalJSON
* Interitted from golang library
* To decerialize blockChain as Json type
* @input: data []byte
* @output: nill
*
 */
func (blockChain *BlockChain) UnmarshalJSON(data []byte) error {
	blocks := make([]Block, 0)
	err := json.Unmarshal(data, &blocks)
	if err != nil {
		return err
	}
	blockChain.Initial()
	for _, block := range blocks {
		blockChain.Insert(block) // update blockChain by insertion
	}
	return nil
}

/* MarshalJSON
* Interitted from golang library
* To serialize  blockChain as Json type
*
 */
func (blockChain *BlockChain) MarshalJSON() ([]byte, error) {
	blocks := make([]Block, 0)
	for _, v := range blockChain.Chain {
		blocks = append(blocks, v...)
	}
	return json.Marshal(blocks)
}


func (blockChain *BlockChain) GetEventInfornation(eventId string) string {
	rs := ""
	forksBlocks:= blockChain.GetLatestBlocks()
	for i, currentBlock := range forksBlocks {
		height := blockChain.Length
		rs += "\n"
		rs += fmt.Sprintf("Chain # %d:\n ", i+1)
		for  height > 0{
			for _, valueObject := range currentBlock.Value.db {
				if strings.Contains(valueObject.String(), eventId) {
					fmt.Println("TransactionObject:", valueObject.String())
					rs += fmt.Sprintf("Value=%s\n", valueObject.String());
				}else{
					fmt.Println("eventId:", eventId," does not exist in our BlockChain!")
				}
			}
			currentBlock, _= blockChain.GetBlock(currentBlock.Header.Height-1, currentBlock.Header.ParentHash)
			height = height - 1
		}
	}
	rs += "\n"
	fmt.Println(rs)
	return rs
}