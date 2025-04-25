package monitor

func (s *MonitorService) CollectAll() map[string]interface{} {
	resdata := make(map[string]interface{})
	for _, c := range s.collector {
		ms := (*c).Collect()
		for _, m := range ms {
			resdata[(*m).Name()] = map[string]interface{}{
				"Name":      (*m).Name(),
				"value":     (*m).Value(),
				"help":      (*m).Help(),
				"collector": (*c).Name(),
			}
		}
	}
	for _, m := range s.metrics {
		resdata[(*m).Name()] = map[string]interface{}{
			"Name":  (*m).Name(),
			"value": (*m).Value(),
			"help":  (*m).Help(),
		}
	}
	return resdata
}

func (s *MonitorService) GetAll() map[string]interface{} {
	resdata := make(map[string]interface{})
	for _, c := range s.collector {
		for _, m := range (*c).Metrics() {
			resdata[(*m).Name()] = map[string]interface{}{
				"Name":      (*m).Name(),
				"value":     (*m).Value(),
				"help":      (*m).Help(),
				"collector": (*c).Name(),
			}
		}
	}
	for _, m := range s.metrics {
		resdata[(*m).Name()] = map[string]interface{}{
			"Name":  (*m).Name(),
			"value": (*m).Value(),
			"help":  (*m).Help(),
		}
	}
	return resdata
}
