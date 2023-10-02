package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, s := range stages {
		bi := make(Bi)

		go func(bi Bi, out Out) {
			defer close(bi)

			for {
				select {
				case v, ok := <-out:
					if !ok {
						return
					}
					bi <- v
				case <-done:
					return
				}
			}
		}(bi, out)

		out = s(bi)
	}

	return out
}
