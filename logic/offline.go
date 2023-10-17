package logic

//离线处理器

type offlineProcessor struct {
}

var OfflineProcessor = newOfflineProcessor()

func newOfflineProcessor() *offlineProcessor {
	return &offlineProcessor{}
}

func (o *offlineProcessor) Save(msg *Message) {

}
