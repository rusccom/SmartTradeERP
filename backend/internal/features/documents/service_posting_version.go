package documents

import (
	"context"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/features/ledger"
)

type postingVersionInput struct {
	tenantID   string
	documentID string
	supersedes string
}

type postingData struct {
	doc   Document
	items []postingItem
}

func (s *Service) postDocumentTx(ctx context.Context, tx pgx.Tx, tenantID, documentID string) error {
	_, err := s.postDocumentVersionTx(ctx, tx, tenantID, documentID, "")
	return err
}

func (s *Service) postDocumentVersionTx(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	documentID string,
	supersedes string,
) ([]ledger.VariantSequence, error) {
	input := postingVersionInput{tenantID: tenantID, documentID: documentID, supersedes: supersedes}
	return s.postVersion(ctx, tx, input)
}

func (s *Service) postVersion(
	ctx context.Context,
	tx pgx.Tx,
	input postingVersionInput,
) ([]ledger.VariantSequence, error) {
	data, err := s.loadPostingData(ctx, tx, input)
	if err != nil {
		return nil, err
	}
	return s.writePostingData(ctx, tx, input, data)
}

func (s *Service) loadPostingData(
	ctx context.Context,
	tx pgx.Tx,
	input postingVersionInput,
) (postingData, error) {
	doc, err := s.repo.PostingDocument(ctx, tx, input.tenantID, input.documentID)
	if err != nil {
		return postingData{}, err
	}
	items, err := s.repo.PostingItems(ctx, tx, input.tenantID, input.documentID)
	return postingData{doc: doc, items: items}, err
}

func (s *Service) writePostingData(
	ctx context.Context,
	tx pgx.Tx,
	input postingVersionInput,
	data postingData,
) ([]ledger.VariantSequence, error) {
	run, err := s.newPostingRun(ctx, tx, input.tenantID, data.doc, input.supersedes)
	if err != nil {
		return nil, err
	}
	affected, err := s.postItems(run, data.items)
	if err != nil {
		return nil, err
	}
	return affected, s.finishPostingData(ctx, tx, input, affected)
}

func (s *Service) finishPostingData(
	ctx context.Context,
	tx pgx.Tx,
	input postingVersionInput,
	affected []ledger.VariantSequence,
) error {
	if err := s.ledger.RebuildAffected(ctx, tx, input.tenantID, affected); err != nil {
		return err
	}
	return s.repo.SetStatus(ctx, tx, input.tenantID, input.documentID, "posted")
}
