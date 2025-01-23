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
		//nolint:revive
		defer func() {
			for range in {
			}
			close(outChain)
		}()

		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}

				outChain <- value
			}
		}
	}()

	return stage(outChain)
}
