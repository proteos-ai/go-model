package functionsapi

import functionsmodel "go.proteos.ai/model/functions"

// DeployFunctionBundleRequest is the JSON metadata side of the multipart body
// for `PUT /functions/v1/function-bundles/:module_slug`: one wasm (the
// multipart `wasm` part) serving every hook, action and connector method of a
// module. Every entry is the same manifest a standalone deploy would send;
// connector methods ride Actions with scope `connector_method`. An entry's
// module_slug may be left empty — it is stamped from the bundle's — but when
// set it must match.
type DeployFunctionBundleRequest struct {
	ModuleSlug string                `json:"module_slug" validate:"required"`
	Hooks      []DeployHookRequest   `json:"hooks" validate:"dive"`
	Actions    []DeployActionRequest `json:"actions" validate:"dive"`
}

// DeployFunctionBundleResponse is what a bundle deploy returns: the blob every
// row now shares plus the rows as saved.
type DeployFunctionBundleResponse struct {
	ModuleSlug string                  `json:"module_slug"`
	FileId     string                  `json:"file_id"`
	Checksum   string                  `json:"checksum"`
	Hooks      []functionsmodel.Hook   `json:"hooks"`
	Actions    []functionsmodel.Action `json:"actions"`
}
