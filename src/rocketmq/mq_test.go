/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rocketmq

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/apache/rocketmq-clients/golang"
)

func TestSendMessage(t *testing.T) {
	producer, err := golang.NewProducer(&golang.Config{
		Endpoint:      "47.92.142.215:8081",
		ConsumerGroup: "GID_BANNIU_TEST",
	}, golang.WithTopics("BANNIU_TEST"))
	if err != nil {
		log.Fatal(err)
	}

	err = producer.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer func(producer golang.Producer) {
		err := producer.GracefulStop()
		if err != nil {
		}
	}(producer)

	for i := 0; i < 10; i++ {
		msg := &golang.Message{
			Topic: "BANNIU_TEST",
			Body:  []byte("this is a message : " + strconv.Itoa(i)),
		}
		msg.SetKeys("a", "b")
		msg.SetTag("ab")
		resp, err := producer.Send(context.TODO(), msg)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(resp); i++ {
			fmt.Printf("%#v\n", resp[i])
		}
		time.Sleep(time.Second * 1)
	}
}
