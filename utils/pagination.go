package utils

func PageSize(data []string, page int, size int) []string {

	result := make([]string, 0)
	total := len(data)
	if total == 0 || page <= 0 {
		return result
	}

	start := (page - 1) * size
	if start >= total {
		return result
	}

	end := start + size
	if end >= total {
		end = total
	}

	result = data[start:end]

	return result
}

func PageSizeV2(total int, page int, size int) (start, end int) {
	if total == 0 || page <= 0 {
		return
	}

	start = (page - 1) * size
	if start >= total {
		start = 0
		return
	}

	end = start + size
	if end >= total {
		end = total
	}

	return
}
