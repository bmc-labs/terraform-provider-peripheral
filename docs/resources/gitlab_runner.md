---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "peripheral_gitlab_runner Resource - peripheral"
subcategory: ""
description: |-
  GitLabRunner resource
---

# peripheral_gitlab_runner (Resource)

GitLabRunner resource



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) GitLabRunner ID as provided by GitLab
- `image` (String) Docker image for GitLabRunner
- `token` (String) Token for GitLabRunner registration
- `url` (String) URL of GitLab instance for GitLabRunner

### Optional

- `description` (String) Description of GitLabRunner
- `run_untagged` (Boolean) Allow untagged jobs
- `tag_list` (String) Comma-separated list of tags for GitLabRunner
