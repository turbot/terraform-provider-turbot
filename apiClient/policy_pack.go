package apiClient

import "fmt"

// Policy pack reads used by the policy pack / smart folder attachment resources.
//
// These deliberately avoid the generic `resource(id:)` query. Guardrails authorizes
// `resource(id:)` against the resource being read, and policy packs live at the Turbot root by
// default, so `resource(id:)` on a pack requires a grant wherever the pack sits. Attaching a pack
// only requires permissions on the attachment *target*, and the Guardrails console reads packs
// through `policyPack(id:)`, which is authorized for such an identity. Using the same queries as
// the console keeps Terraform working for callers that hold permissions on the target alone.

// ReadPolicyPackIdentity resolves a policy pack by id OR aka and returns its Turbot metadata
// (numeric id + akas). Uses the `policyPack(id:)` query so it does not require a grant on the
// pack itself. resourceType names the concept in error messages, so callers managing smart
// folders report "smart folder" rather than "policy pack".
func (client *Client) ReadPolicyPackIdentity(policyPackAka, resourceType string) (*TurbotResourceMetadata, error) {
	query := readPolicyPackIdentityQuery(policyPackAka)
	responseData := &PolicyPackIdentityResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, policyPackAka, resourceType)
	}
	return &responseData.PolicyPack.Turbot, nil
}

// ReadAttachedPolicyPacks returns the policy packs attached to the given resource, read from
// the resource side via `resource(id:).attachedSmartFolders`. The caller necessarily holds
// permissions on the attachment target, so this is authorized wherever the attach itself is.
//
// The list is DIRECT attachments only — it does not include packs attached to an ancestor.
// (Pack policies still propagate down the hierarchy; the attachment list does not.) Verified
// against a live workspace: with a pack attached only at a grandparent folder, the child and
// grandchild both report an empty list. Detaching from a child while an ancestor attachment
// remains correctly empties the child's list. `Exists` depends on this: were the list the
// effective set instead, a nested resource would report an attachment nothing created.
//
// `attachedSmartFolders` accepts no paging arguments, so it returns whatever the server default
// allows. The second return value reports whether the server truncated the list, which callers
// need because only ABSENCE is ambiguous under truncation - see PolicyPackAttached.
//
// Caveat: truncation cannot be forced, since the field takes no arguments, so a truncated
// response has never been observed. What is confirmed is that `paging.next` decodes and is empty
// for complete lists. If the server ever truncates WITHOUT setting a cursor, this will not catch it.
func (client *Client) ReadAttachedPolicyPacks(resourceAka string) ([]TurbotResourceMetadata, bool, error) {
	query := readAttachedPolicyPacksQuery(resourceAka)
	responseData := &AttachedPolicyPacksResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, false, client.handleReadError(err, resourceAka, "attached policy packs")
	}

	attached := responseData.Resource.AttachedSmartFolders
	packs := make([]TurbotResourceMetadata, 0, len(attached.Items))
	for _, item := range attached.Items {
		packs = append(packs, item.Turbot)
	}
	return packs, attached.Paging.Next != "", nil
}

// PolicyPackAttached reports whether policyPack (matched by numeric id or by any aka) is
// currently attached to resourceAka.
func (client *Client) PolicyPackAttached(resourceAka, policyPack string) (bool, error) {
	attached, truncated, err := client.ReadAttachedPolicyPacks(resourceAka)
	if err != nil {
		return false, err
	}
	// Found it: truncation is irrelevant, a later page cannot un-attach it.
	if policyPackInList(attached, policyPack) {
		return true, nil
	}
	// Absent AND truncated is the only ambiguous case. Erroring on truncation alone would make
	// every attachment on a heavily-attached resource unrefreshable, including ones sitting in the
	// page we already read - strictly worse than the churn it is meant to prevent, because churn
	// self-corrects and a refresh error blocks every plan until someone hand-edits state.
	if truncated {
		return false, fmt.Errorf("error reading attached policy packs: resource %s reports more attachments than a single page returns, so attachment state cannot be determined reliably", resourceAka)
	}
	return false, nil
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
