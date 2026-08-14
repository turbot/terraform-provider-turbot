package apiClient

// Policy pack reads used by the policy pack / smart folder attachment resources.
//
// These deliberately avoid the generic `resource(id:)` query. Guardrails authorizes
// `resource(id:)` against the resource being read, and policy packs live at the Turbot
// root, so `resource(id:)` on a pack requires a root-level grant. Attaching a pack only
// requires permissions on the attachment *target*, and the Guardrails console reads packs
// through `policyPack(id:)`, which is authorized for such an identity. Using the same
// queries as the console keeps Terraform working for callers that hold permissions on the
// target resource alone.

// ReadPolicyPackIdentity resolves a policy pack by id OR aka and returns its Turbot
// metadata (numeric id + akas). Uses the `policyPack(id:)` query so it does not require a
// grant on the pack itself.
func (client *Client) ReadPolicyPackIdentity(policyPackAka string) (*TurbotResourceMetadata, error) {
	query := readPolicyPackIdentityQuery(policyPackAka)
	responseData := &PolicyPackIdentityResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, policyPackAka, "policy pack")
	}
	return &responseData.PolicyPack.Turbot, nil
}

// ReadAttachedPolicyPacks returns the policy packs attached to the given resource, read
// from the resource side via `resource(id:).attachedSmartFolders`. The caller necessarily
// holds permissions on the attachment target, so this is authorized wherever the attach
// itself is. It also scopes the answer to this one resource rather than enumerating every
// resource the pack is attached to across the whole hierarchy.
func (client *Client) ReadAttachedPolicyPacks(resourceAka string) ([]TurbotResourceMetadata, error) {
	query := readAttachedPolicyPacksQuery(resourceAka)
	responseData := &AttachedPolicyPacksResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, resourceAka, "attached policy packs")
	}

	packs := make([]TurbotResourceMetadata, 0, len(responseData.Resource.AttachedSmartFolders.Items))
	for _, item := range responseData.Resource.AttachedSmartFolders.Items {
		packs = append(packs, item.Turbot)
	}
	return packs, nil
}

// PolicyPackAttached reports whether policyPack (matched by numeric id or by any aka) is
// currently attached to resourceAka.
func (client *Client) PolicyPackAttached(resourceAka, policyPack string) (bool, error) {
	attached, err := client.ReadAttachedPolicyPacks(resourceAka)
	if err != nil {
		return false, err
	}
	return policyPackInList(attached, policyPack), nil
}

// policyPackInList reports whether policyPack identifies any pack in list, matching on the
// numeric id or on any of the pack's akas so either form works.
func policyPackInList(list []TurbotResourceMetadata, policyPack string) bool {
	for _, pack := range list {
		if pack.Id == policyPack {
			return true
		}
		for _, aka := range pack.Akas {
			if aka == policyPack {
				return true
			}
		}
	}
	return false
}
