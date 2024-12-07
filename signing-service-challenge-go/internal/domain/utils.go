package domain

func ConvertStringToAlgorithmType(input string) AlgorithmType {
	switch input {
	case "ECC":
		return AlgorithmTypeECC
	case "RSA":
		return AlgorithmTypeRSA
	default:
		return AlgorithmTypeUnknown
	}
}
