package p4

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math/rand"
	"time"
)

/*-------------------------STRUCT---------------------------------------------------*/
/* Struct data structure for variables
/*-------------------------STRUCT---------------------------------------------------*/

/* Block Header struct
*
* To Define block header using Block header struct
*
 */
type BlockHeader struct {
	Size       int32
	ParentHash string
	Height     int32
	Timestamp  int64
	Hash       string
	Nonce      string
}

/* Block struct
*
* Block struct
 */
type Block struct {
	Header BlockHeader
	Value  MerklePatriciaTrie
}

/*-------------------------INITIALIZATION---------------------------------------------------*/
/* Data initialization
/*-------------------------INITIALIZATION---------------------------------------------------*/

/* Initial
*
* To initialize MPT height and parentHash
* @input: value MerklePatriciaTrie, height int32, parentHash string
* @output: nill
*
 */
func (block *Block) Initial(height int32, parentHash string, value MerklePatriciaTrie, nonce string) {
	block.Header.Height = height
	block.Header.Timestamp = time.Now().Unix()
	block.Header.ParentHash = parentHash
	block.Header.Size = int32(len([]byte(fmt.Sprintf("%v", value))))
	block.Value = value
	block.Header.Nonce = nonce
	//hashConverter := sha3.New256()
	//hashStr := string(block.Header.Height) + string(block.Header.Timestamp) + block.Header.ParentHash + block.Value.root + string(block.Header.Size)
	//block.Header.Hash = hex.EncodeToString(hashConverter.Sum([]byte(hashStr)))
	block.Header.Hash = block.Hash()
}

/*-------------------------JSON HELPER---------------------------------------------------*/
/* JSON feature
/*-------------------------JSON HELPER---------------------------------------------------*/

/* UnmarshalJSON
* Intehrited from golang library
* To encodes a block instance into a JSON format string
* @input: an instanse of block
* @output: a string of JSON format
*
 */
func (block *Block) UnmarshalJSON(input []byte) error {
	SymmetricBlockJson := BlockJson{}
	err := json.Unmarshal(input, &SymmetricBlockJson)
	if err != nil {
		return err
	}
	block.Header.Height = SymmetricBlockJson.Height
	block.Header.Timestamp = SymmetricBlockJson.Timestamp
	block.Header.Hash = SymmetricBlockJson.Hash
	block.Header.ParentHash = SymmetricBlockJson.ParentHash
	block.Header.Size = SymmetricBlockJson.Size
	block.Header.Nonce = SymmetricBlockJson.Nonce
	mpt := MerklePatriciaTrie{}
	mpt.Initial()
	for k, v := range SymmetricBlockJson.MPT {
		mpt.Insert(k, v)
	}
	block.Value = mpt
	return nil
}

/* EncodeToJSON
* Inherited from golang library
* To encodes a block instance into a JSON format string
* @input: an instanse of block
* @output: a string of JSON format
*
 */
func (block *Block) EncodeToJSON() (string, error) {
	jsonBytes, err := json.Marshal(block)
	return string(jsonBytes), err
}

/* DecodeFromJSON
*
* To take a string that represents the JSON value of a block as an input, and decodes the input string back to a block instance.
* @input:  a string of JSON format
* @output: an instanse of block
*
 */
func (block *Block) DecodeFromJSON(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), block)
}

/* MarshalJSON
*
* To hash MPT with the SHA3-256 encoded value of this string and update MPT value upon the
* insertion
* @input:  block
* @output: updated block
*
 */
func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(BlockJson{
		Height:     block.Header.Height,
		Timestamp:  block.Header.Timestamp,
		Size:       block.Header.Size,
		Nonce:      block.Header.Nonce,
		Hash:       block.Header.Hash,
		ParentHash: block.Header.ParentHash,
		MPT:        block.Value.LeafList(),
	})
}

/* BlockJson
*
* BlockJson struct for Block struct
*
 */
type BlockJson struct {
	Height     int32             `json:"height"`
	Timestamp  int64             `json:"timeStamp"`
	Size       int32             `json:"size"`
	Nonce      string            `json:"nonce"`
	Hash       string            `json:"hash"`
	ParentHash string            `json:"parentHash"`
	MPT        map[string]string `json:"mpt"`
}

func (block *Block) GetHash() string {
	return block.Header.Hash
}

/* Hash
*
* To hash the block
*
 */
func (block *Block) Hash() string {
	var hashStr string
	hashStr = string(block.Header.Height) + string(block.Header.Timestamp) + string(block.Header.ParentHash) +
		string(block.Value.root) + string(block.Header.Size) + string(block.Header.Nonce)
	sum := sha3.Sum256([]byte(hashStr))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

/* StringRandom()
*
* To generate Random String
* https://www.calhoun.io/creating-random-strings-in-go/
*
 */
const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
func StringRandom(length int) string {
	return StringWithCharset(length, charset)
}
