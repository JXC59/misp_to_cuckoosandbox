package misp

type RelatedAttribute struct {
	Attribute    Attribute `json:"Attribute"`
	AttributeTag []Tag     `json:"AttributeTag"`
}

type Attribute struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	Category           string `json:"category"`
	ToIDS              bool   `json:"to_ids"`
	UUID               string `json:"uuid"`
	EventID            string `json:"event_id"`
	Distribution       string `json:"distribution"`
	Timestamp          string `json:"timestamp"`
	Comment            string `json:"comment"`
	SharingGroupID     string `json:"sharing_group_id"`
	Deleted            bool   `json:"deleted"`
	DisableCorrelation bool   `json:"disable_correlation"`
	ObjectID           string `json:"object_id"`
	Value              string `json:"value"`

	Tags []Tag `json:"Tag,omitempty"`

	ObjectRelation *string `json:"object_relation"`
	FirstSeen      *string `json:"first_seen"`
	LastSeen       *string `json:"last_seen"`

	Galaxies         []Galaxy    `json:"Galaxy"`
	ShadowAttributes []Attribute `json:"ShadowAttribute"`
}

type Galaxy struct {
	ID          string `json:"id"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Icon        string `json:"icon"`
	Namespace   string `json:"namespace"`

	GalaxyClusters []GalaxyCluster `json:"GalaxyCluster"`
}

type GalaxyCluster struct {
	ID             string `json:"id"`
	UUID           string `json:"uuid"`
	CollectionUUID string `json:"collection_uuid"`
	Type           string `json:"type"`
	Value          string `json:"value"`
	TagName        string `json:"tag_name"`
	Description    string `json:"description"`
	GalaxyID       string `json:"galaxy_id"`
	Source         string `json:"source"`
	Version        string `json:"version"`
	TagID          string `json:"tag_id"`
	Local          bool   `json:"local"`
	Meta           Meta   `json:"meta"`

	Authors []string `json:"authors"`
}

type Meta struct {
	AttributionConfidence    []string `json:"attribution-confidence"`
	CfrSuspectedStateSponsor []string `json:"cfr-suspected-state-sponsor"`
	CfrSuspectedVictims      []string `json:"cfr-suspected-victims"`
	CfrTargetCategory        []string `json:"cfr-target-category"`
	CfrTypeOfIncident        []string `json:"cfr-type-of-incident"`
	Country                  []string `json:"country"`
	Refs                     []string `json:"refs"`
	Synonyms                 []string `json:"synonyms"`
}

type Tag struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Colour     string `json:"colour,omitempty"`
	Exportable bool   `json:"exportable,omitempty"`
	UserID     string `json:"user_id,omitempty"`
	HideTag    bool   `json:"hide_tag,omitempty"`
	Local      int    `json:"local,omitempty"`

	NumericalValue *string `json:"numerical_value,omitempty"`
}

type RelatedTag struct {
	Tag Tag `json:"Tag"`
}

type Event struct {
	ID                 string `json:"id"`
	OrgcID             string `json:"orgc_id"`
	OrgID              string `json:"org_id"`
	Date               string `json:"date"`
	ThreatLevelID      string `json:"threat_level_id"`
	Info               string `json:"info"`
	Published          bool   `json:"published"`
	UUID               string `json:"uuid"`
	AttributeCount     string `json:"attribute_count"`
	Analysis           string `json:"analysis"`
	Timestamp          string `json:"timestamp"`
	Distribution       string `json:"distribution"`
	ProposalEmailLock  bool   `json:"proposal_email_lock"`
	Locked             bool   `json:"locked"`
	PublishTimestamp   string `json:"publish_timestamp"`
	SharingGroupID     string `json:"sharing_group_id"`
	DisableCorrelation bool   `json:"disable_correlation"`
	ExtendsUUID        string `json:"extends_uuid"`

	Org  Organisation `json:"Org"`
	Orgc Organisation `json:"Orgc"`

	Attributes    []Attribute    `json:"Attribute"`
	RelatedEvents []RelatedEvent `json:"RelatedEvent"`
	Galaxies      []Galaxy       `json:"Galaxy"`
	Tags          []Tag          `json:"Tag"`

	ShadowAttributes []any    `json:"ShadowAttribute"`
	Objects          []Object `json:"Object"`
}

type RelatedEvent struct {
	Event Event `json:"Event"`
}

type Organisation struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	UUID  string `json:"uuid"`
	Local bool   `json:"local"`
}

type Object struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	MetaCategory    string  `json:"meta-category"`
	Description     string  `json:"description"`
	TemplateUUID    string  `json:"template_uuid"`
	TemplateVersion string  `json:"template_version"`
	EventID         string  `json:"event_id"`
	UUID            string  `json:"uuid"`
	Timestamp       string  `json:"timestamp"`
	Distribution    string  `json:"distribution"`
	SharingGroupID  string  `json:"sharing_group_id"`
	Comment         string  `json:"comment"`
	Deleted         bool    `json:"deleted"`
	FirstSeen       *string `json:"first_seen"`
	LastSeen        *string `json:"last_seen"`

	Attributes []Attribute `json:"Attribute"`
}
