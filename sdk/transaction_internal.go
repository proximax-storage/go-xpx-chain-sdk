package sdk

type hashDto string

func (dto *hashDto) Hash() (*Hash, error) {
	s := string(*dto)

	if len(s) == 0 {
		return nil, nil
	}

	return StringToHash(s)
}

type signatureDto string

func (dto *signatureDto) Signature() (*Signature, error) {
	s := string(*dto)

	if len(s) == 0 {
		return nil, nil
	}

	return StringToSignature(s)
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
