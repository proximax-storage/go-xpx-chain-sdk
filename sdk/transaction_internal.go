package sdk

type hashDto string

func (dto *hashDto) Hash() (Hash, error) {
	return StringToHash(string(*dto))
}

type signatureDto string

func (dto *signatureDto) Signature() (Signature, error) {
	return StringToSignature(string(*dto))
}

type transactionStatusDTOs []*transactionStatusDTO

func (t *transactionStatusDTOs) toStruct() ([]*TransactionStatus, error) {
	dtos := *t
	statuses := make([]*TransactionStatus, 0, len(dtos))

	for _, dto := range dtos {
		status, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}
