package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func WitohutDoneButFast(in In, stages ...Stage) Out {
	var stage Stage

	if len(stages) < 1 {
		return in
	}
	if len(stages) > 1 {
		stage, stages = stages[0], stages[1:]
		return WitohutDoneButFast(stage(in), stages...)
	}
	return stages[0](in)
}

func Demultiplex(done In, pipe In, stage Stage) Out {
	bus := make(Bi)
	go func() {
		defer close(bus)
		for {
			select {
			case <-done:
				return
			case v, ok := <-pipe:
				if !ok {
					return
				}
				bus <- v
			}
		}
	}()
	return stage(bus)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if done == nil {
		return WitohutDoneButFast(in, stages...)
	}
	pipe := in
	for _, stage := range stages {
		pipe = Demultiplex(done, pipe, stage)
	}
	return pipe
}
