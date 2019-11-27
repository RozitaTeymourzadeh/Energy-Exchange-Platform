package driver

import (
	//"MerklePatriciaTree/p5/Blockchain_Application_P5/p4"
	"fmt"
	"sync"
)

/* SyncBlockChain Struct
*
* Sync Block Chain Struct Format
*
 */
type SyncBlockChain struct {
	bc  BlockChain
	mux sync.Mutex
}

/* NewBlockChain()
*
* Return new blockChain
*
 */
func NewSyncBlockChain() SyncBlockChain {
	return SyncBlockChain{bc: NewBlockChain()}
}

/* GetParentBlock()
* Synchronize version of block
* Takes a block as the parameter, and returns its parent block.
*
 */
func (sbc *SyncBlockChain) GetParentBlock(block Block) Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetParentBlock(block)
}

/* GetLatestBlocks()
* Synchronize version of block
* Returns the list of blocks of height "BlockChain.length"
*
 */
func (sbc *SyncBlockChain) GetLatestBlocks() []Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetLatestBlocks()
}

/* Get()
*
* Return new blockList in fork height
*
 */
func (sbc *SyncBlockChain) Get(height int32) ([]Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Get(height)
}

/* GetBlock()
*
* Return block from the blocklist in fork
*
 */
func (sbc *SyncBlockChain) GetBlock(height int32, hash string) (Block, bool) {
	return sbc.bc.GetBlock(height, hash)
}

/* Insert()
*
* Insert block in blockChain in synchronize way
*
 */
func (sbc *SyncBlockChain) Insert(block Block) {
	sbc.mux.Lock()
	sbc.bc.Insert(block)
	sbc.mux.Unlock()
}

/* CheckParentHash()
*
* It would check if the block with the given "parentHash" exists in the blockChain.
* If we have the parent block, we can insert the next block; if we don't have the parent block,
* we have to download the parent block before inserting the next block.
*
 */
func (sbc *SyncBlockChain) CheckParentHash(insertBlock Block) bool {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	if insertBlock.Header.Height > 1 {
		blocks, found := sbc.bc.Get(insertBlock.Header.Height - 1)
		if found == true {
			for _, pb := range blocks {
				if pb.Header.Hash == insertBlock.Header.ParentHash {
					return true
				}
			}
		}
	}
	return false
}

/* UpdateEntireBlockChain()
*
* It UpdateEntireBlockChain convert from BlockChainJson into BlockChain Format
*
 */
func (sbc *SyncBlockChain) UpdateEntireBlockChain(blockChainJson string) {
	blockChain := sbc.bc.DecodeFromJSON(blockChainJson)
	fmt.Println(blockChain)
}

/* BlockChainToJson()
*
* It UpdateEntireBlockChain convert from BlockChain into BlockChainJson Format
*
 */
func (sbc *SyncBlockChain) BlockChainToJson() (string, error) {
	return sbc.bc.EncodeToJSON()
}

/* GenBlock()
*
* To Generate Random Block in BlockChain
*
 */

func (sbc *SyncBlockChain) GenBlock(mpt MerklePatriciaTrie, nonce string) Block {

	var parentHash string
	var blockList []Block
	var found bool
	currHeight := sbc.bc.Length

	for currHeight >= 1 {
		blockList, found = sbc.Get(currHeight)
		if found == true {
			parentHash = blockList[0].Header.Hash
			break
		}
		currHeight--
	}
	if currHeight == 0 {
		parentHash = "gensis"
	}

	var newBlock Block
	newBlock.Initial(currHeight+1, parentHash, mpt, nonce)

	return newBlock
}

/* Show()
*
* To Show current BlockChain
*
 */
func (sbc *SyncBlockChain) Show() string {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Show()
}

/* Canonical()
*
* To Show current BlockChain after POW in synchronize way
*
 */
func (sbc *SyncBlockChain) Canonical() string {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Canonical()
}

/* GetEventInfornation()
*
* To get event information
*
 */
func (sbc *SyncBlockChain) GetEventInfornation(eventId string) string {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetEventInfornation(eventId)
}
