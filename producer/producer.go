package producer

//type Producer struct {
//	batch models.Batch
//}
//
//func newProducer(numOfItem int) (producer Producer) {
//	return Producer{
//		batch: make([]models.Item, numOfItem),
//	}
//}
//
//func Produce(numsOfItem []int) {
//	for _, numOfItem := range numsOfItem {
//		AddBuffer(newProducer(numOfItem).batch)
//		time.Sleep(time.Duration(rand.Int31n(5)) * time.Second)
//	}
//}
