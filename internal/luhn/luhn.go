package luhn

func IsValid(number int) bool {
	checksum := number % 10
	number /= 10

	var sum int
	double := false

	for number > 0 {
		double = !double

		digit := number % 10
		number /= 10

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	return (sum+checksum)%10 == 0
}
