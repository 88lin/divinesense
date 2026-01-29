package store

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/hrygo/divinesense/internal/base"
	storepb "github.com/hrygo/divinesense/proto/gen/store"
)

type Attachment struct {
	Payload       *storepb.AttachmentPayload
	MemoUID       *string
	MemoID        *int32
	Type          string
	FilePath      string
	RowStatus     string
	Filename      string
	UID           string
	OCRText       string
	ExtractedText string
	ThumbnailPath string
	Reference     string
	Blob          []byte
	CreatedTs     int64
	Size          int64
	UpdatedTs     int64
	StorageType   storepb.AttachmentStorageType
	ID            int32
	CreatorID     int32
}

type FindAttachment struct {
	MemoID         *int32
	ID             *int32
	UID            *string
	CreatorID      *int32
	Filename       *string
	FilenameSearch *string
	StorageType    *storepb.AttachmentStorageType
	Limit          *int
	Offset         *int
	MemoIDList     []int32
	Filters        []string
	GetBlob        bool
	HasRelatedMemo bool
}

type UpdateAttachment struct {
	UID           *string
	UpdatedTs     *int64
	Filename      *string
	MemoID        *int32
	Reference     *string
	Payload       *storepb.AttachmentPayload
	RowStatus     *string
	ExtractedText *string
	OCRText       *string
	ThumbnailPath *string
	ID            int32
}

type DeleteAttachment struct {
	MemoID *int32
	ID     int32
}

func (s *Store) CreateAttachment(ctx context.Context, create *Attachment) (*Attachment, error) {
	if !base.UIDMatcher.MatchString(create.UID) {
		return nil, errors.New("invalid uid")
	}
	return s.driver.CreateAttachment(ctx, create)
}

func (s *Store) ListAttachments(ctx context.Context, find *FindAttachment) ([]*Attachment, error) {
	// Defensive check for nil driver (e.g., in tests)
	if s.driver == nil {
		return nil, nil
	}

	// Set default limits to prevent loading too many attachments at once
	if find.Limit == nil && find.GetBlob {
		// When fetching blobs, we should be especially careful with limits
		defaultLimit := 10
		find.Limit = &defaultLimit
	} else if find.Limit == nil {
		// Even without blobs, let's default to a reasonable limit
		defaultLimit := 100
		find.Limit = &defaultLimit
	}

	return s.driver.ListAttachments(ctx, find)
}

func (s *Store) GetAttachment(ctx context.Context, find *FindAttachment) (*Attachment, error) {
	attachments, err := s.ListAttachments(ctx, find)
	if err != nil {
		return nil, err
	}

	if len(attachments) == 0 {
		return nil, nil
	}

	return attachments[0], nil
}

func (s *Store) UpdateAttachment(ctx context.Context, update *UpdateAttachment) error {
	if update.UID != nil && !base.UIDMatcher.MatchString(*update.UID) {
		return errors.New("invalid uid")
	}
	return s.driver.UpdateAttachment(ctx, update)
}

func (s *Store) DeleteAttachment(ctx context.Context, delete *DeleteAttachment) error {
	attachment, err := s.GetAttachment(ctx, &FindAttachment{ID: &delete.ID})
	if err != nil {
		return errors.Wrap(err, "failed to get attachment")
	}
	if attachment == nil {
		return errors.New("attachment not found")
	}

	if attachment.StorageType == storepb.AttachmentStorageType_LOCAL {
		p := filepath.FromSlash(attachment.Reference)
		if !filepath.IsAbs(p) {
			p = filepath.Join(s.profile.Data, p)
		}
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			// Log error but don't prevent database deletion
			slog.Error("failed to delete attachment file",
				"error", err,
				"path", p,
				"attachment_id", delete.ID,
			)
		}
	}

	return s.driver.DeleteAttachment(ctx, delete)
}
