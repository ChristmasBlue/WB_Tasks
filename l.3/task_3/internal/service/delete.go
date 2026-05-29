package service

func (s *Service) DeleteCommentById(id int) error {
	return s.storage.DeleteCommentById(id)
}
