package or

func Or(channels ...<-chan interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func(ch chan interface{}) {
		i := 0
		for {
			select {
			case <-channels[i]:
				close(ch)
				return
			default:
				i = (i + 1) % len(channels)
			}
		}
	}(ch)
	return ch
}
