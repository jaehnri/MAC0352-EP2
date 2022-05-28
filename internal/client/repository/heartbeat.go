package repository

type HeartbeatRepository struct{}

func NewHeartbeatRepository() *HeartbeatRepository {
	return &HeartbeatRepository{}
}

func (s *HeartbeatRepository) Send(username string) {
	// TODO: implement send method
}
