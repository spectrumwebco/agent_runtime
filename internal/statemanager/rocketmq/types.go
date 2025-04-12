package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type MessageExt = primitive.MessageExt

const ConsumeFromLastOffset = consumer.ConsumeFromLastOffset
