package prodcons

type ServiceRunnerI interface {
	Runner()
}

type ProdCons struct {
	records []any
}

func NewProdCons(records []any) *ProdCons {
	return &ProdCons{
		records: records,
	}
}

func (pd *ProdCons) Runner() {
	var sharedCh chan any = make(chan any)
	var doneCh chan bool = make(chan bool)
	var producer Producer = Producer{
		records: pd.records,
	}
	var consumer Consumer = Consumer{
		sharedCh,
	}

	go producer.produce(sharedCh)
	go consumer.consume(doneCh)
	<-doneCh
}
