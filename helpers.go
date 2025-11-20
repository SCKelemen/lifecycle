package lifecycle

// Helper functions to create common structures

// NewActor creates a new Actor
func NewActor(userID string, actorType ActorType) *Actor {
	return &Actor{
		UserID:    userID,
		ActorType: actorType,
	}
}

// NewHumanActor creates a new human actor
func NewHumanActor(userID string) *Actor {
	return NewActor(userID, ActorTypeHuman)
}

// NewSystemActor creates a new system actor
func NewSystemActor(systemID string) *Actor {
	return NewActor(systemID, ActorTypeSystem)
}

// NewSyntheticActor creates a new synthetic actor
func NewSyntheticActor(syntheticID string) *Actor {
	return NewActor(syntheticID, ActorTypeSynthetic)
}

// NewResource creates a new Resource
func NewResource(resourceType, id string) *Resource {
	return &Resource{
		Type: resourceType,
		ID:   id,
	}
}

