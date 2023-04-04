package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"log"
	"math/big"
	"strings"
)

// 挖矿奖励的BTC
const reward = 12.5

// Transaction 交易结构
type Transaction struct {
	TXID      []byte     // 交易ID
	TXInputs  []TXInput  // 交易输入
	TXOutputs []TXOutput // 交易输出
}

/*
TXInput
交易输入结构:
1. 由交易创建人提供；
2. 花费一笔BTC给交易对方；
3. 由矿工校验交易的合理性。
*/
type TXInput struct {
	TXid      []byte // 引用的交易ID
	Index     int64  // 引用output的索引值(本次input的额度，由之前交易的output决定)
	Signature []byte // 数字签名
	PubKey    []byte // 付款方的公钥(为了便于网络传输，不存储原始公钥，只存储圆锥曲线的X和Y值)
}

/*
TXOutput
交易输出结构:
1. 由之前给自己转钱方所创建的input转换而来；
2. 一笔交易合理后对方的input就变成自己的output。
*/
type TXOutput struct {
	Value      float64 // 转账金额
	PubKeyHash []byte  // 收款方公钥的Hash
}

// UnLock 通过收款方地址反推出公钥Hash(因为output结构中存储的是pubKeyHash字段)
func (o *TXOutput) UnLock(address string) {
	addressByte := base58.Decode(address)

	pubKeyHash := addressByte[1 : len(addressByte)-4]

	o.PubKeyHash = pubKeyHash
}

// NewTXOutput 生成交易的Output
func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{Value: value}
	output.UnLock(address)

	return &output
}

// SetTxID 设置交易ID(也就是整个交易结构的Hash)
func (t *Transaction) SetTxID() {
	var buffer bytes.Buffer

	// 1. 对整个交易结构序列化
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(t)
	if err != nil {
		log.Panic(err)
	}

	// 2. 对序列化后的字节流执行Hash算法，得到的hash值就是交易ID
	hash := sha256.Sum256(buffer.Bytes())
	t.TXID = hash[:]
}

// IsCoinbase 判断当前交易是否是挖矿交易
func (t *Transaction) IsCoinbase() bool {
	// 挖矿交易的特点就是只有一个input
	if len(t.TXInputs) == 1 && len(t.TXInputs[0].TXid) == 0 && t.TXInputs[0].Index == -1 {
		return true
	}
	return false
}

/*
NewCoinbaseTX
创建一笔挖矿交易(参与记账，由系统奖励BTC):
1.只有一个输入和输出；
2.无需存储交易ID；
3.无需引用任何Output。
*/
func NewCoinbaseTX(address string, data string) *Transaction {
	input := TXInput{
		TXid:      []byte{},
		Index:     -1,
		Signature: nil,
		PubKey:    []byte(data), // PubKey字段由矿工自由填写，一般是矿池的名字
	}
	output := NewTXOutput(reward, address)

	// 创建交易
	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{*output},
	}
	tx.SetTxID()

	return &tx
}

/*
NewTransaction
创建一笔普通交易:
1. 付款方花费的金额由之前交易产生的可用的output决定；
2. 交易中的每一个input都必须进行数字签名(用于验证付款方的身份)。
*/
func NewTransaction(from, to string, amount float64, chain *BlockChain) *Transaction {
	// 1. 获取付款方的钱包
	wallets := NewWallets()
	wallet := wallets.WalletsMap[from]
	if wallet == nil {
		fmt.Println("没有找到该地址的钱包，交易创建失败!")
		return nil
	}

	// 2. 获取钱包中的密钥对
	privateKey := wallet.Private
	pubKey := wallet.PubKey

	// 3. 查找此次交易需要使用到的UTXO
	pubKeyHash := HashPubKey(pubKey)
	utxos, resValue := chain.FindNeedUTXOs(pubKeyHash, amount)
	if resValue < amount {
		fmt.Println("余额不足，交易失败")
		return nil
	}

	// 4. 将需要使用的UTXO转换成此次交易的inputs
	var inputs []TXInput
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			inputs = append(inputs, TXInput{
				TXid:      []byte(id), // 交易ID
				Index:     i,          // 交易ID对应交易的output的索引
				Signature: nil,
				PubKey:    pubKey,
			})
		}
	}

	// 5. 封装此次交易的output
	var outputs []TXOutput
	outputs = append(outputs, *NewTXOutput(amount, to))

	// 6. 找零操作，若本次使用的金额大于交易需要的金额，则额外生成output转账剩余的金额给自己
	if resValue > amount {
		outputs = append(outputs, *NewTXOutput(resValue-amount, from))
	}

	// 7. 封装交易结构
	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  inputs,
		TXOutputs: outputs,
	}
	tx.SetTxID()

	// 8. 查找需要进行数字签名的交易并进行数字签名
	prevTXs := chain.FindPrevTransaction(tx.TXInputs) // 从整个区块链中查找包含inputs所引用的output的交易
	tx.Sign(privateKey, prevTXs)

	return &tx
}

/*
Sign
数字签名实现:
1. 创建一笔普通交易的时候进行签名，保证付款方的身份；
2. 交易中的每一个input都要进行签名(input由之前交易的output提供)；
3. 签名需要的数据是: 当前交易的inputs所引用的output的PubKeyHash + 当前交易的Outputs。
*/
func (t *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if t.IsCoinbase() {
		return
	}

	// 1. 创建一个当前交易的副本，由交易副本完成签名的操作，最后将签名结果赋予真正的交易体
	txCopy := t.TrimmedCopy()

	// 2. 遍历txCopy的inputs，取出所有input引用的output的公钥哈希
	for i, input := range txCopy.TXInputs {
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("交易无效!")
		}

		/*
			3. 封装txCopy：
				3.1 将被引用的output的PubKeyHash存入input的PubKey字段中；
				3.2 input的PubKey字段相当于一个临时存储；
				3.3 这样操作后，txCopy中就存在被引用的output的PubKeyHash和当前交易的Outputs了。
		*/
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash

		// 4. 对封装好的txCopy做Hash处理，最终TXID就是要签名的数据
		txCopy.SetTxID()
		// 还原PubKey字段，相当于清空临时存储，好让下一个Hash正常进行
		txCopy.TXInputs[i].PubKey = nil

		/*
			5. 使用椭圆曲线算法进行数字签名：
				参数：1. 加密算法的随机数种子；2. 付款方的私钥；3. 最终签名的数据。
		*/
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, txCopy.TXID)
		if err != nil {
			log.Panic(err)
		}

		// 6. 将本次循环签名完成的input赋值到正真正的交易结构中
		t.TXInputs[i].Signature = append(r.Bytes(), s.Bytes()...)
	}
}

/*
Verify
校验数字签名:
1. 在区块生成之前校验，验证付款方的身份；
2. 交易中的所有input都需要进行校验；
*/
func (t *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if t.IsCoinbase() {
		return true
	}

	// 1. 产生交易副本
	txCopy := t.TrimmedCopy()
	for i, input := range txCopy.TXInputs {
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("交易无效!")
		}

		// 2. 交易副本承担一个临时存储和参与Hash算法的角色，存储被引用output的PubKeyHash + 本次交易的outputs
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		txCopy.SetTxID()

		// 3. 拆分数字签名
		signature := input.Signature
		r, s := &big.Int{}, &big.Int{}
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])

		// 4. 拆分公钥
		pubKey := input.PubKey
		X, Y := &big.Int{}, &big.Int{}
		X.SetBytes(pubKey[:len(pubKey)/2])
		Y.SetBytes(pubKey[len(pubKey)/2:])

		// 5. 校验数字签名
		pubKeyOrigin := &ecdsa.PublicKey{Curve: elliptic.P256(), X: X, Y: Y} // 原始公钥
		if !ecdsa.Verify(pubKeyOrigin, txCopy.TXID, r, s) {
			return false
		}
	}

	return true
}

// TrimmedCopy 创建当前交易的副本
func (t *Transaction) TrimmedCopy() Transaction {
	var newInputs []TXInput

	for _, input := range t.TXInputs {
		newInputs = append(newInputs, TXInput{input.TXid, input.Index, nil, nil})
	}

	return Transaction{t.TXID, newInputs, t.TXOutputs}
}

// 格式化打印交易
func (t *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x: ", t.TXID))

	for i, input := range t.TXInputs {
		lines = append(lines, fmt.Sprintf("	Input %d: ", i))
		lines = append(lines, fmt.Sprintf("		TXID:		%X", input.TXid))
		lines = append(lines, fmt.Sprintf("		Out:		%d", input.Index))
		lines = append(lines, fmt.Sprintf("		Signature:	%x", input.Signature))
		lines = append(lines, fmt.Sprintf("		PubKey:		%x", input.PubKey))
	}

	for i, output := range t.TXOutputs {
		lines = append(lines, fmt.Sprintf("	Output %d: ", i))
		lines = append(lines, fmt.Sprintf("		Value:  %f", output.Value))
		lines = append(lines, fmt.Sprintf("		Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}
