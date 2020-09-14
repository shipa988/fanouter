package sender

type QuerySenderFabric interface {
	NewQuerySender() QuerySender
}
