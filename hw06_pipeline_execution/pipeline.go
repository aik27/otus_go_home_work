package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	out := in

	for _, stage := range stages {
		if stage == nil {
			continue
		}
		out = exec(stage, out, done)
	}

	return out
}

func exec(stage Stage, in In, done In) Out {
	outChain := make(Bi)

	go func() {
		defer close(outChain)

		for {
			select {
			case <-done:
				// fmt.Printf("Stage %d received done signal\n", stageIndex)
				return
			case value, ok := <-in:
				if !ok {
					// fmt.Printf("Stage %d completed\n", stageIndex)
					return
				}
				select {
				case <-done:
					// fmt.Printf("Stage %d received done signal while sending\n", stageIndex)
					return
				case outChain <- value:
					// fmt.Printf("Stage %d sent element\n", stageIndex)
				}
			}
		}
	}()

	return stage(outChain)
}
