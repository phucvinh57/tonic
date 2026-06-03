package docs

// https://swagger.io/specification/#path-item-object
type PathItemObject struct {
	Summary     string                 `json:"summary,omitempty"`
	Description string                 `json:"description,omitempty"`
	Get         *OperationObject       `json:"get,omitempty"`
	Query       *OperationObject       `json:"query,omitempty"`
	Put         *OperationObject       `json:"put,omitempty"`
	Post        *OperationObject       `json:"post,omitempty"`
	Delete      *OperationObject       `json:"delete,omitempty"`
	Options     *OperationObject       `json:"options,omitempty"`
	Head        *OperationObject       `json:"head,omitempty"`
	Patch       *OperationObject       `json:"patch,omitempty"`
	Trace       *OperationObject       `json:"trace,omitempty"`
	Servers     []ServerObject         `json:"servers,omitempty"`
	Parameters  []ParameterOrReference `json:"parameters,omitempty"`
}

type PathItemOrReference struct {
	*PathItemObject  `json:",inline,omitempty"`
	*ReferenceObject `json:",inline,omitempty"`
}
