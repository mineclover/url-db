package entity

import (
	"errors"
	"time"
)

// Node represents a node entity in the business domain
type Node struct {
	id          int
	content     string // This is the URL field in database
	domainID    int
	title       string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewNode creates a new node entity with validation
func NewNode(url, title, description string, domainID int) (*Node, error) {
	if url == "" {
		return nil, errors.New("node URL cannot be empty")
	}

	if len(url) > 2048 {
		return nil, errors.New("node URL cannot exceed 2048 characters")
	}

	if domainID <= 0 {
		return nil, errors.New("domain ID must be positive")
	}

	if len(title) > 255 {
		return nil, errors.New("node title cannot exceed 255 characters")
	}

	if len(description) > 1000 {
		return nil, errors.New("node description cannot exceed 1000 characters")
	}

	now := time.Now()
	return &Node{
		content:     url, // Store URL in content field
		domainID:    domainID,
		title:       title,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Getters - immutable from outside
func (n *Node) ID() int              { return n.id }
func (n *Node) Content() string      { return n.content }
func (n *Node) URL() string          { return n.content } // Alias for content
func (n *Node) DomainID() int        { return n.domainID }
func (n *Node) Title() string        { return n.title }
func (n *Node) Description() string  { return n.description }
func (n *Node) CreatedAt() time.Time { return n.createdAt }
func (n *Node) UpdatedAt() time.Time { return n.updatedAt }

// Setters for internal use (e.g., by repository)
func (n *Node) SetID(id int) { n.id = id }

// Business logic methods
func (n *Node) UpdateTitle(title string) error {
	if len(title) > 255 {
		return errors.New("node title cannot exceed 255 characters")
	}

	n.title = title
	n.updatedAt = time.Now()
	return nil
}

func (n *Node) UpdateDescription(description string) error {
	if len(description) > 1000 {
		return errors.New("node description cannot exceed 1000 characters")
	}

	n.description = description
	n.updatedAt = time.Now()
	return nil
}

func (n *Node) UpdateContent(title, description string) error {
	if err := n.UpdateTitle(title); err != nil {
		return err
	}

	if err := n.UpdateDescription(description); err != nil {
		return err
	}

	return nil
}

// IsValid checks if the node is in a valid state
func (n *Node) IsValid() bool {
	return n.content != "" &&
		len(n.content) <= 2048 &&
		n.domainID > 0 &&
		len(n.title) <= 255 &&
		len(n.description) <= 1000
}

// SetTimestamps sets creation and update timestamps (for repository usage)
func (n *Node) SetTimestamps(createdAt, updatedAt time.Time) {
	n.createdAt = createdAt
	n.updatedAt = updatedAt
}
