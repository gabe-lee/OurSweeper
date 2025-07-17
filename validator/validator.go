package validator

type Validator interface {
	IsValid() bool
}

type HMACSaver interface {
	LoadHMAC() []byte
	SaveHMAC(hmac []byte)
}
