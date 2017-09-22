package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/buffla/models"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Link)
// DB Table: Plural (links)
// Resource: Plural (Links)
// Path: Plural (/links)
// View Template Folder: Plural (/templates/links/)

// LinksResource is the resource for the link model
type LinksResource struct {
	buffalo.Resource
}

func (v LinksResource) scope(c buffalo.Context) *pop.Query {
	tx := c.Value("tx").(*pop.Connection)
	cuid := c.Session().Get("current_user_id")
	return tx.Where("user_id = ?", cuid)
}

// List gets all Links. This function is mapped to the path
// GET /links
func (v LinksResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	links := models.Links{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := v.scope(c).PaginateFromParams(c.Params())
	// You can order your list here. Just change
	err := q.All(&links)
	// to:
	// err := q.Order("created_at desc").All(links)
	if err != nil {
		return errors.WithStack(err)
	}

	// Make Links available inside the html template
	c.Set("links", links)

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("links/index.html"))
}

// Show gets the data for one Link. This function is mapped to
// the path GET /links/{link_id}
func (v LinksResource) Show(c buffalo.Context) error {
	// Allocate an empty Link
	link := &models.Link{}
	// To find the Link the parameter link_id is used.
	err := v.scope(c).Find(link, c.Param("link_id"))
	if err != nil {
		return c.Error(404, err)
	}
	// Make link available inside the html template
	c.Set("link", link)

	tx := c.Value("tx").(*pop.Connection)

	clicks := &models.ClickActivities{}
	if err := tx.RawQuery("select * from click_activity(?)", link.ID).All(clicks); err != nil {
		return errors.WithStack(err)
	}

	c.Set("clicks", clicks)
	return c.Render(200, r.HTML("links/show.html"))
}

// New renders the form for creating a new Link.
// This function is mapped to the path GET /links/new
func (v LinksResource) New(c buffalo.Context) error {
	// Make link available inside the html template
	c.Set("link", &models.Link{})
	return c.Render(200, r.HTML("links/new.html"))
}

// Create adds a Link to the DB. This function is mapped to the
// path POST /links
func (v LinksResource) Create(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)

	user := c.Value("current_user").(*models.User)

	// Allocate an empty Link
	link := &models.Link{}
	// Bind link to the html form elements
	err := c.Bind(link)
	if err != nil {
		return errors.WithStack(err)
	}

	link.UserID = user.ID

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(link)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make link available inside the html template
		c.Set("link", link)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("links/new.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Link was created successfully")
	// and redirect to the links index page
	return c.Redirect(302, "/links/%s", link.ID)
}

// Edit renders a edit form for a link. This function is
// mapped to the path GET /links/{link_id}/edit
func (v LinksResource) Edit(c buffalo.Context) error {
	// Allocate an empty Link
	link := &models.Link{}
	err := v.scope(c).Find(link, c.Param("link_id"))
	if err != nil {
		return c.Error(404, err)
	}
	// Make link available inside the html template
	c.Set("link", link)
	return c.Render(200, r.HTML("links/edit.html"))
}

// Update changes a link in the DB. This function is mapped to
// the path PUT /links/{link_id}
func (v LinksResource) Update(c buffalo.Context) error {
	// Allocate an empty Link
	link := &models.Link{}
	err := v.scope(c).Find(link, c.Param("link_id"))
	if err != nil {
		return c.Error(404, err)
	}
	// Bind Link to the html form elements
	err = c.Bind(link)
	if err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndUpdate(link)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make link available inside the html template
		c.Set("link", link)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("links/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Link was updated successfully")
	// and redirect to the links index page
	return c.Redirect(302, "/links/%s", link.ID)
}

// Destroy deletes a link from the DB. This function is mapped
// to the path DELETE /links/{link_id}
func (v LinksResource) Destroy(c buffalo.Context) error {
	// Allocate an empty Link
	link := &models.Link{}
	// To find the Link the parameter link_id is used.
	err := v.scope(c).Find(link, c.Param("link_id"))
	if err != nil {
		return c.Error(404, err)
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	err = tx.Destroy(link)
	if err != nil {
		return errors.WithStack(err)
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Link was destroyed successfully")
	// Redirect to the links index page
	return c.Redirect(302, "/links")
}

func Redirector(c buffalo.Context) error {
	link := &models.Link{}
	tx := c.Value("tx").(*pop.Connection)

	if err := tx.Where("code = ?", c.Param("code")).First(link); err != nil {
		return c.Error(404, err)
	}

	click := &models.Click{LinkID: link.ID}
	if err := tx.Create(click); err != nil {
		c.Logger().Error(err)
	}

	return c.Redirect(302, link.Link)
}
