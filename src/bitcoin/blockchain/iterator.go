package blockchain

import (
	"errors"
	"github.com/boltdb/bolt"
	"log"
)

// Iterator 区块链迭代器结构
type Iterator struct {
	db                 *bolt.DB // 数据库
	currentHashPointer []byte   // 当前Hash指针
}

// GetIterator 获取区块链的迭代器
func (c *BlockChain) GetIterator() *Iterator {
	return &Iterator{c.db, c.tail}
}

/*
Next
迭代器的迭代函数:
1. 从最后一个区块开始迭代；
2. 每次迭代都会向前改变currentHashPointer的指向，直到创世块。
*/
func (i *Iterator) Next() (*Block, error) {
	var block Block
	var err error

	_ = i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("nil chain")
		}

		// 取出当前区块，并将currentHashPointer指向前置区块，下一次调用Next()函数就会得到前一个区块开始，依此类推
		DeSerialize(bucket.Get(i.currentHashPointer), &block)
		i.currentHashPointer = block.PrevHash

		// 设置结束条件: 当前区块的前置Hash为空的字节流，也就是迭代到了创世块
		if len(block.PrevHash) == 0 {
			err = errors.New("遍历结束")
		}

		return nil
	})

	return &block, err
}
