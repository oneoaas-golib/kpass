package api

import (
	"github.com/seccom/kpass/src/auth"
	"github.com/seccom/kpass/src/bll"
	"github.com/seccom/kpass/src/model"
	"github.com/seccom/kpass/src/schema"
	"github.com/seccom/kpass/src/util"
	"github.com/teambition/gear"
)

// Entry is API oject for entries
//
// @Name Entry
// @Description Entry API
// @Accepts json
// @Produces json
type Entry struct {
	models   *model.All
	entryBll *bll.Entry
}

// Init ...
func (a *Entry) Init(blls *bll.All) *Entry {
	a.models = blls.Models
	a.entryBll = blls.Entry
	return a
}

type tplEntryCreate struct {
	Name     string `json:"name" swaggo:"true,entry name,Github"`
	Category string `json:"category" swaggo:"true,entry category,Logins"`
}

func (t *tplEntryCreate) Validate() error {
	if len(t.Name) == 0 {
		return &gear.Error{Code: 400, Msg: "entry name required"}
	}
	return nil
}

// Create ...
//
// @Title Create
// @Summary Create a entry in a team
// @Description all team members can create entry
// @Param Authorization header string true "access_token"
// @Param teamID path string true "team ID"
// @Param body body tplEntryCreate true "entry body"
// @Success 200 schema.EntrySum
// @Failure 400 string
// @Failure 401 string
// @Router POST /api/teams/{teamID}/entries
func (a *Entry) Create(ctx *gear.Context) (err error) {
	TeamID, err := util.ParseOID(ctx.Param("teamID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}

	body := new(tplEntryCreate)
	if err = ctx.ParseBody(body); err != nil {
		return ctx.Error(err)
	}
	userID, err := auth.UserIDFromCtx(ctx)
	if err != nil {
		return ctx.Error(err)
	}

	entry, err := a.entryBll.Create(userID, &schema.Entry{
		TeamID:   TeamID,
		Name:     body.Name,
		Category: body.Category,
		Secrets:  []string{},
	})
	if err != nil {
		return ctx.Error(err)
	}
	return ctx.JSON(200, entry)
}

type tplEntryUpdate map[string]interface{}

// Validate ...
func (t *tplEntryUpdate) Validate() error {
	empty := true
	for key, val := range *t {
		empty = false

		switch key {
		case "name":
			v, ok := val.(string)
			if !ok || v == "" {
				return &gear.Error{Code: 400, Msg: "invalid entry name"}
			}
		case "category":
			_, ok := val.(string)
			if !ok {
				return &gear.Error{Code: 400, Msg: "invalid entry category"}
			}
		case "priority":
			v, ok := val.(float64)
			if !ok || v < 0 || v > 127 {
				return &gear.Error{Code: 400, Msg: "invalid entry priority"}
			}
		default:
			return &gear.Error{Code: 400, Msg: "invalid entry property"}
		}
	}

	if empty {
		return &gear.Error{Code: 400, Msg: "no content"}
	}
	return nil
}

// Update ...
//
// @Title Update
// @Summary Update the entry
// @Description all team members can update the entry
// @Param Authorization header string true "access_token"
// @Param entryID path string true "entry ID"
// @Param body body tplEntryUpdate true "entry body"
// @Success 200 schema.EntrySum
// @Success 204
// @Failure 400 string
// @Failure 401 string
// @Router PUT /api/entries/{entryID}
func (a *Entry) Update(ctx *gear.Context) (err error) {
	EntryID, err := util.ParseOID(ctx.Param("entryID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}

	userID, _ := auth.UserIDFromCtx(ctx)
	body := new(tplEntryUpdate)
	if err = ctx.ParseBody(body); err != nil {
		return ctx.Error(err)
	}

	entrySum, err := a.entryBll.Update(userID, EntryID, *body)
	if err != nil {
		return ctx.Error(err)
	}
	if entrySum == nil {
		return ctx.End(204)
	}
	return ctx.JSON(200, entrySum)
}

// Delete ...
//
// @Title Delete
// @Summary Delete the entry
// @Description all team members can delete the entry
// @Param Authorization header string true "access_token"
// @Param entryID path string true "entry ID"
// @Success 204
// @Failure 400 string
// @Failure 401 string
// @Router DELETE /api/entries/{entryID}
func (a *Entry) Delete(ctx *gear.Context) (err error) {
	EntryID, err := util.ParseOID(ctx.Param("entryID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}

	userID, _ := auth.UserIDFromCtx(ctx)
	if _, err = a.entryBll.Delete(userID, EntryID, true); err != nil {
		return ctx.Error(err)
	}
	return ctx.End(204)
}

// Undelete ...
//
// @Title Undelete
// @Summary Undelete the entry
// @Description all team members can undelete the entry
// @Param Authorization header string true "access_token"
// @Param entryID path string true "entry ID"
// @Success 204
// @Failure 400 string
// @Failure 401 string
// @Router POST /api/entries/{entryID}:undelete
func (a *Entry) Undelete(ctx *gear.Context) (err error) {
	EntryID, err := util.ParseOID(ctx.Param("entryID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}

	userID, _ := auth.UserIDFromCtx(ctx)
	entrySum, err := a.entryBll.Delete(userID, EntryID, false)
	if err != nil {
		return ctx.Error(err)
	}
	return ctx.JSON(200, entrySum)
}

// Find return the entry
//
// @Title Find
// @Summary Get the entry
// @Description Get the entry with all information, include secrets, files and shares.
// @Description all team members can get the entry
// @Param Authorization header string true "access_token"
// @Param entryID path string true "entry ID"
// @Success 200 schema.EntryResult
// @Failure 400 string
// @Failure 401 string
// @Router GET /api/entries/{entryID}
func (a *Entry) Find(ctx *gear.Context) error {
	EntryID, err := util.ParseOID(ctx.Param("entryID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}
	key, err := auth.KeyFromCtx(ctx)
	if err != nil {
		return ctx.Error(err)
	}
	userID, _ := auth.UserIDFromCtx(ctx)
	res, err := a.entryBll.Find(userID, key, EntryID)
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.JSON(200, res)
}

// FindByTeam return entries for current user
//
// @Title FindByTeam
// @Summary Get the team's entries list
// @Description Get the team's entries list with summary information.
// @Description all team members can get it
// @Param Authorization header string true "access_token"
// @Param teamID path string true "team ID"
// @Success 200 []schema.EntrySum
// @Failure 400 string
// @Failure 401 string
// @Router GET /api/teams/{teamID}/entries
func (a *Entry) FindByTeam(ctx *gear.Context) (err error) {
	TeamID, err := util.ParseOID(ctx.Param("teamID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}

	userID, _ := auth.UserIDFromCtx(ctx)
	res, err := a.entryBll.FindByTeam(userID, TeamID)
	if err != nil {
		return ctx.Error(err)
	}
	return ctx.JSON(200, res)
}

// DeleteFile ...
//
// @Title DeleteFile
// @Summary Delete the file
// @Description all team members can delete the file
// @Param Authorization header string true "access_token"
// @Param entryID path string true "entry ID"
// @Param fileID path string true "file ID"
// @Success 204
// @Failure 400 string
// @Failure 401 string
// @Router DELETE /api/entries/{entryID}/files/{fileID}
func (a *Entry) DeleteFile(ctx *gear.Context) error {
	EntryID, err := util.ParseOID(ctx.Param("entryID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}
	FileID, err := util.ParseOID(ctx.Param("fileID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}
	userID, err := auth.UserIDFromCtx(ctx)
	if err != nil {
		return ctx.Error(err)
	}
	if err = a.entryBll.DeleteFile(userID, EntryID, FileID); err != nil {
		return ctx.Error(err)
	}
	return ctx.End(204)
}
