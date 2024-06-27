---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "jamfpro_static_computer_group Resource - terraform-provider-jamfpro"
subcategory: ""
description: |-
  
---

# jamfpro_static_computer_group (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The unique name of the Jamf Pro static computer group.

### Optional

- `assignments` (Block List, Max: 1) Assignment block containing the list of computer IDs. (see [below for nested schema](#nestedblock--assignments))
- `site_id` (Number) Jamf Pro Site-related settings of the policy.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The unique identifier of the Jamf Pro static computer group.
- `is_smart` (Boolean) Computed value indicating whether the computer group is smart or static.

<a id="nestedblock--assignments"></a>
### Nested Schema for `assignments`

Required:

- `computer_ids` (List of Number) The list of computer IDs that are members of the static computer group.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)