package entity

import (
	"errors"
	"time"
)

// Node represents a node entity in the business domain
type Node struct {
	id          int
	url         string
	domainName  string
	title       string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewNode creates a new node entity with validation
func NewNode(url, domainName, title, description string) (*Node, error) {
	if url == "" {
		return nil, errors.New("node URL cannot be empty")
	}

	if len(url) > 2048 {
		return nil, errors.New("node URL cannot exceed 2048 characters")
	}

	if domainName == "" {
		return nil, errors.New("domain name cannot be empty")
	}

	if len(title) > 255 {
		return nil, errors.New("node title cannot exceed 255 characters")
	}

	if len(description) > 1000 {
		return nil, errors.New("node description cannot exceed 1000 characters")
	}

	now := time.Now()
	return &Node{
		url:         url,
		domainName:  domainName,
		title:       title,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Getters - immutable from outside
func (n *Node) ID() int              { return n.id }
func (n *Node) URL() string          { return n.url }
func (n *Node) DomainName() string   { return n.domainName }
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
	return n.url != "" &&
		len(n.url) <= 2048 &&
		n.domainName != "" &&
		len(n.title) <= 255 &&
		len(n.description) <= 1000
}
