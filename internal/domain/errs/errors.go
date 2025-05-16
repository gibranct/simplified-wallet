package errs

import "errors"

var (
	ErrTransactionNotAllowed          = errors.New("transaction not allowed")
	ErrTransactionInvalidAmount       = errors.New("transaction invalid amount")
	ErrNotEnoughMoney                 = errors.New("not enough money")
	ErrTransactionInvalidSender       = errors.New("transaction invalid sender")
	ErrSenderNotFound                 = errors.New("sender not found")
	ErrReceiverNotFound               = errors.New("receiver not found")
	ErrEitherDocumentMustBeProvided   = errors.New("either document must be provided")
	ErrCNPJMustBeProvidedForMerchant  = errors.New("cnpj must be provided for merchant user type")
	ErrCPFMustBeProvidedForCommonUser = errors.New("cpf must be provided for common user type")
	ErrMerchantCannotHaveCPF          = errors.New("merchant user type cannot have cpf")
	ErrCommonCannotHaveCNPJ           = errors.New("common user type cannot have cnpj")
	ErrZeroOrNegativeAmount           = errors.New("amount must be a positive number")
	ErrInsufficientBalance            = errors.New("insufficient balance")
	ErrMerchantCannotSendMoney        = errors.New("merchant user type cannot send money")
	ErrNameLength                     = errors.New("name length must be between 3 and 50 characters")
	ErrEmailAlreadyRegistered         = errors.New("email already registered")
	ErrCPFAlreadyRegistered           = errors.New("cpf already registered")
	ErrCNPJAlreadyRegistered          = errors.New("cnpj already registered")
	ErrUserTypeNotFound               = errors.New("user type not found")
)
