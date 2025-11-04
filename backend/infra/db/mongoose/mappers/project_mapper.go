package mappers

import (
	"time"

	"ai_hub.com/app/core/domain/projectdomain"
)

type ProjectDoc struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Status      string    `bson:"status"`
	APIKey      string    `bson:"apiKey"`
	OwnerID     string    `bson:"ownerId"`
	AdminAccess []string  `bson:"adminAccess"`
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

func ProjectFromDoc(doc *ProjectDoc) (*projectdomain.Project, error) {
	if doc == nil {
		return nil, nil
	}
	// ID
	id, err := projectdomain.NewProjectID(doc.ID)
	if err != nil {
		return nil, err
	}
	// Name
	name, err := projectdomain.NewProjectName(doc.Name)
	if err != nil {
		return nil, err
	}
	// KeyEnc
	keyEnc, err := projectdomain.NewHashedAPIKey(doc.APIKey)
	if err != nil {
		return nil, err
	}
	// OwnerID
	owner, err := projectdomain.NewOwnerID(doc.OwnerID)
	if err != nil {
		return nil, err
	}
	// AdminAccess
	accesses := make([]projectdomain.AccessAdminID, 0, len(doc.AdminAccess))
	for _, a := range doc.AdminAccess {
		acc, err := projectdomain.NewAccessAdminID(a)
		if err != nil {
			return nil, err
		}
		accesses = append(accesses, acc)
	}

	props := projectdomain.RestoreProps{
		ID:          id,
		Name:        name,
		Status:      doc.Status,
		APIKey:      keyEnc,
		OwnerID:     owner,
		AdminAccess: accesses,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}

	p, err := projectdomain.Restore(props)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func ProjectToPersistence(pj *projectdomain.Project) (*ProjectDoc, error) {
	p := pj.ToPrimitives()

	createdAt, err := time.Parse(time.RFC3339, p.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &ProjectDoc{
		ID:          p.ID,
		Name:        p.Name,
		Status:      p.Status,
		APIKey:      p.APIKey,
		OwnerID:     p.OwnerID,
		AdminAccess: p.AdminAccess,
		CreatedAt:   createdAt.UTC(),
		UpdatedAt:   updatedAt.UTC(),
	}, nil
}
