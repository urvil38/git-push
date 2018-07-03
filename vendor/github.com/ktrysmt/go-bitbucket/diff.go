package bitbucket

type Diff struct {
	c *Client
}

func (d *Diff) GetDiff(do *DiffOptions) (interface{}, error) {
	urlStr := d.c.requestUrl("/repositories/%s/%s/diff/%s", do.Owner, do.Repo_slug, do.Spec)
	return d.c.execute("GET", urlStr, "")
}

func (d *Diff) GetPatch(do *DiffOptions) (interface{}, error) {
	urlStr := d.c.requestUrl("/repositories/%s/%s/patch/%s", do.Owner, do.Repo_slug, do.Spec)
	return d.c.execute("GET", urlStr, "")
}
