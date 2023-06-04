package controller

import (
	"strconv"
	"vehicles/packages/adapters"
	"vehicles/packages/usecase/usecase"
)

type selectionController struct {
	selectionUseCase usecase.SelectionInput
}

type Selection interface {
	ChoosePriorities(ctx adapters.Context) error
	SetPriorities(ctx adapters.Context) error
	ChoosePrice(ctx adapters.Context) error
	SetPrice(ctx adapters.Context) error
	ChooseManufacturers(ctx adapters.Context) error
	SetManufacturers(ctx adapters.Context) error
	GetSelection(ctx adapters.Context) error
	ShowCarCard(ctx adapters.Context) error
}

func NewSelectionController(sl usecase.SelectionInput) Selection {

	return &selectionController{sl}
}

func (sc *selectionController) ChoosePriorities(ctx adapters.Context) error {

	err := sc.selectionUseCase.ChoosePriorities()
	if err != nil {
		return err
	}
	return nil
}

func (sc *selectionController) SetPriorities(ctx adapters.Context) error {

	type priorities struct {
		Priorities []string `json:"priorities"`
	}

	prs := new(priorities)
	if err := ctx.BindJSON(&prs); err != nil {
		return err
	}

	err := sc.selectionUseCase.SetPriorities(&prs.Priorities)
	if err != nil {
		return err
	}

	return nil
}

func (sc *selectionController) ChoosePrice(ctx adapters.Context) error {
	err := sc.selectionUseCase.ChoosePrice()
	if err != nil {
		return err
	}
	return nil
}

func (sc *selectionController) SetPrice(ctx adapters.Context) error {
	type price struct {
		MinPrice  string `json:"minPrice"`
		MaxPrice  string `json:"maxPrice"`
		Deviation string `json:"deviation"`
	}

	pc := new(price)
	if err := ctx.BindJSON(&pc); err != nil {
		return err
	}

	err := sc.selectionUseCase.SetPrice(pc.MinPrice, pc.MaxPrice, pc.Deviation)
	if err != nil {
		return err
	}
	return nil
}

func (sc *selectionController) ChooseManufacturers(ctx adapters.Context) error {
	err := sc.selectionUseCase.ChooseManufacturers()
	if err != nil {
		return err
	}
	return nil
}

func (sc *selectionController) SetManufacturers(ctx adapters.Context) error {

	type manufacturers struct {
		Manufacturers []string `json:"manufacturers"`
	}

	mns := new(manufacturers)
	if err := ctx.BindJSON(&mns); err != nil {
		return err
	}

	err := sc.selectionUseCase.SetManufacturers(&mns.Manufacturers)
	if err != nil {
		return err
	}

	return nil
}

func (sc *selectionController) GetSelection(ctx adapters.Context) error {
	err := sc.selectionUseCase.GetSelection()
	if err != nil {
		return err
	}
	return nil
}

func (sc *selectionController) ShowCarCard(ctx adapters.Context) error {
	_id := ctx.Param("id")
	id, _ := strconv.Atoi(_id)

	err := sc.selectionUseCase.ShowCarCard(id)
	if err != nil {
		return err
	}
	return nil
}
