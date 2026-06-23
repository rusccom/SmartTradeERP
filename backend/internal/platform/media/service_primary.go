package media

import "context"

// PrimaryThumbsByOwners returns ownerID -> thumbnail URL for the primary image
// of each owner, in one query, so list views can show thumbnails without N+1.
func (s *Service) PrimaryThumbsByOwners(
	ctx context.Context,
	tenantID string,
	ownerType string,
	ownerIDs []string,
) (map[string]string, error) {
	if s.objects == nil {
		return nil, ErrStorageNotConfigured
	}
	if len(ownerIDs) == 0 {
		return map[string]string{}, nil
	}
	refs, err := s.repo.PrimaryByOwners(ctx, tenantID, ownerType, ownerIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(refs))
	for _, ref := range refs {
		result[ref.OwnerID] = thumbURL(s.objects.PublicURL(ref.ObjectKey))
	}
	return result, nil
}
