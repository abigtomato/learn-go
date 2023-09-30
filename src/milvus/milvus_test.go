package milvus

import (
	"context"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"testing"
)

//type Request struct {
//	Model            string     `json:"model"`             // 模型
//	Messages         *[]Message `json:"messages"`          // 消息
//	MaxTokens        uint       `json:"max_tokens"`        // 生成结果时的最大单词数，不能超过模型的上下文长度
//	Temperature      float64    `json:"temperature"`       // 随机因子，控制结果的随机性，如果希望结果更有创意可以尝试 0.9，或者希望有固定结果可以尝试 0.0
//	TopP             int        `json:"top_p"`             // 随机因子2，一个可用于代替 temperature 的参数，对应机器学习中 nucleus sampling（核采样），如果设置 0.1 意味着只考虑构成前 10% 概率质量的 tokens
//	FrequencyPenalty int        `json:"frequency_penalty"` // 重复度惩罚因子，是 -2.0 ~ 2.0 之间的数字，正值会根据新 tokens 在文本中的现有频率对其进行惩罚，从而降低模型逐字重复同一行的可能性
//	PresencePenalty  int        `json:"presence_penalty"`  // 控制主题的重复度，是 -2.0 ~ 2.0 之间的数字，正值会根据到目前为止是否出现在文本中来惩罚新 tokens，从而增加模型谈论新主题的可能性
//}

func TestClient(t *testing.T) {
	mClient, err := client.NewGrpcClient(context.Background(), "192.168.17.12:19530")
	if err != nil {
		log.Fatal("NewGrpcClient Failed: ", err.Error())
		return
	}

	f, err := excelize.OpenFile("./pg.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := f.GetRows("SheetJS")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, row := range rows {
		id, _ := strconv.ParseInt(row[0], 10, 64)
		title := row[1]
		abstract := row[2]
		content := row[3]
		_, err = mClient.Insert(
			context.Background(),
			"KnowledgeVector",
			"",
			entity.NewColumnInt64("knowledge_id", []int64{id}),
			entity.NewColumnVarChar("knowledge_title", []string{title}),
			entity.NewColumnVarChar("knowledge_abstract", []string{abstract}),
			entity.NewColumnVarChar("knowledge_content", []string{content}),
			entity.NewColumnFloatVector("vector_content", 1536, [][]float32{}),
		)
		if err != nil {
			log.Fatal("failed to insert data:", err.Error())
		}
		break
	}
}
