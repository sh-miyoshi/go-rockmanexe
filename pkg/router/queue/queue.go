package queue

var (
	allQueues = make(map[string][]interface{})
)

func Push(id string, info interface{}) {
	allQueues[id] = append(allQueues[id], info)
}

func Pop(id string) interface{} {
	acts, ok := allQueues[id]
	if !ok || len(acts) == 0 {
		return nil
	}
	res := acts[0]
	allQueues[id] = acts[1:]
	return res
}

func Delete(id string) {
	delete(allQueues, id)
}
