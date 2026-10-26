package processor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/xueaaaa/go-uptime/internal/check/model"
	"github.com/xueaaaa/go-uptime/internal/check/service"
	model2 "github.com/xueaaaa/go-uptime/internal/site/model"
	service2 "github.com/xueaaaa/go-uptime/internal/site/service"
)

type fakeCheckService struct {
	service.CheckService
	created     model.Check
	createCalls int
	createErr   error
	calls       *[]string
}

func (s *fakeCheckService) Create(_ context.Context, check model.Check) (uuid.UUID, error) {
	s.createCalls++
	s.created = check
	*s.calls = append(*s.calls, "create")
	return uuid.Nil, s.createErr
}

type fakeSiteService struct {
	service2.SiteService
	updated     model.Check
	updateCalls int
	calls       *[]string
}

func (s *fakeSiteService) UpdateByCheck(_ context.Context, check model.Check) error {
	s.updateCalls++
	s.updated = check
	*s.calls = append(*s.calls, "update")
	return nil
}

func TestProcessorSavesCheckAndUpdatesSite(t *testing.T) {
	check := model.Check{
		SiteID:     uuid.New(),
		Status:     model2.Unavailable,
		StatusCode: 500,
		Error:      "internal error",
		CheckedAt:  time.Now(),
	}
	results := make(chan model.Check, 1)
	results <- check
	close(results)

	var calls []string
	checkSvc := &fakeCheckService{calls: &calls}
	siteSvc := &fakeSiteService{calls: &calls}
	processor := NewProcessor(results, siteSvc, checkSvc)

	if err := processor.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if checkSvc.createCalls != 1 || checkSvc.created != check {
		t.Fatalf("expected check to be saved once, got %d saves with check %+v", checkSvc.createCalls, checkSvc.created)
	}
	if siteSvc.updateCalls != 1 || siteSvc.updated != check {
		t.Fatalf("expected site to be updated once with the check, got %d updates with check %+v", siteSvc.updateCalls, siteSvc.updated)
	}
	if len(calls) != 2 || calls[0] != "create" || calls[1] != "update" {
		t.Fatalf("expected create then update, got %v", calls)
	}
}

func TestProcessorDoesNotUpdateSiteWhenSavingCheckFails(t *testing.T) {
	check := model.Check{SiteID: uuid.New(), Status: model2.Unavailable}
	results := make(chan model.Check, 1)
	results <- check
	close(results)

	var calls []string
	saveErr := errors.New("save failed")
	checkSvc := &fakeCheckService{createErr: saveErr, calls: &calls}
	siteSvc := &fakeSiteService{calls: &calls}
	processor := NewProcessor(results, siteSvc, checkSvc)

	if err := processor.Run(context.Background()); !errors.Is(err, saveErr) {
		t.Fatalf("expected save error %v, got %v", saveErr, err)
	}
	if checkSvc.createCalls != 1 || checkSvc.created != check {
		t.Fatalf("expected one save attempt with the check, got %d attempts with check %+v", checkSvc.createCalls, checkSvc.created)
	}
	if siteSvc.updateCalls != 0 {
		t.Fatalf("expected 0 site updates after save failure, got %d", siteSvc.updateCalls)
	}
	if len(calls) != 1 || calls[0] != "create" {
		t.Fatalf("expected only create call, got %v", calls)
	}
}
