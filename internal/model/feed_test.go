// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package model // import "miniflux.app/v2/internal/model"

import (
	"fmt"
	"os"
	"testing"
	"time"

	"miniflux.app/v2/internal/config"
)

func TestFeedCategorySetter(t *testing.T) {
	feed := &Feed{}
	feed.WithCategoryID(int64(123))

	if feed.Category == nil {
		t.Fatal(`The category field should not be null`)
	}

	if feed.Category.ID != int64(123) {
		t.Error(`The category ID must be set`)
	}
}

func TestFeedErrorCounter(t *testing.T) {
	feed := &Feed{}
	feed.WithTranslatedErrorMessage("Some Error")

	if feed.ParsingErrorMsg != "Some Error" {
		t.Error(`The error message must be set`)
	}

	if feed.ParsingErrorCount != 1 {
		t.Error(`The error counter must be set to 1`)
	}

	feed.ResetErrorCounter()

	if feed.ParsingErrorMsg != "" {
		t.Error(`The error message must be removed`)
	}

	if feed.ParsingErrorCount != 0 {
		t.Error(`The error counter must be set to 0`)
	}
}

func TestFeedCheckedNow(t *testing.T) {
	feed := &Feed{}
	feed.FeedURL = "https://example.org/feed"
	feed.CheckedNow()

	if feed.SiteURL != feed.FeedURL {
		t.Error(`The site URL must not be empty`)
	}

	if feed.CheckedAt.IsZero() {
		t.Error(`The checked date must be set`)
	}
}

func TestFeedScheduleNextCheckDefault(t *testing.T) {
	var err error
	parser := config.NewParser()
	config.Opts, err = parser.ParseEnvironmentVariables()
	if err != nil {
		t.Fatalf(`Parsing failure: %v`, err)
	}

	feed := &Feed{}
	weeklyCount := 10
	feed.ScheduleNextCheck(weeklyCount)

	if feed.NextCheckAt.IsZero() {
		t.Error(`The next_check_at must be set`)
	}

	if feed.NextCheckAt.After(time.Now().Add(time.Minute * time.Duration(config.Opts.PollingFrequency()))) {
		t.Error(`The next_check_at should not be after the now + polling frequency`)
	}
}

func TestFeedScheduleNextCheckEntryCountBasedMaxInterval(t *testing.T) {
	maxInterval := 5
	minInterval := 1
	os.Clearenv()
	os.Setenv("POLLING_SCHEDULER", "entry_frequency")
	os.Setenv("SCHEDULER_ENTRY_FREQUENCY_MAX_INTERVAL", fmt.Sprintf("%d", maxInterval))
	os.Setenv("SCHEDULER_ENTRY_FREQUENCY_MIN_INTERVAL", fmt.Sprintf("%d", minInterval))

	var err error
	parser := config.NewParser()
	config.Opts, err = parser.ParseEnvironmentVariables()
	if err != nil {
		t.Fatalf(`Parsing failure: %v`, err)
	}
	feed := &Feed{}
	weeklyCount := maxInterval * 100
	feed.ScheduleNextCheck(weeklyCount)

	if feed.NextCheckAt.IsZero() {
		t.Error(`The next_check_at must be set`)
	}

	if feed.NextCheckAt.After(time.Now().Add(time.Minute * time.Duration(maxInterval))) {
		t.Error(`The next_check_at should not be after the now + max interval`)
	}
}

func TestFeedScheduleNextCheckEntryCountBasedMinInterval(t *testing.T) {
	maxInterval := 500
	minInterval := 100
	os.Clearenv()
	os.Setenv("POLLING_SCHEDULER", "entry_frequency")
	os.Setenv("SCHEDULER_ENTRY_FREQUENCY_MAX_INTERVAL", fmt.Sprintf("%d", maxInterval))
	os.Setenv("SCHEDULER_ENTRY_FREQUENCY_MIN_INTERVAL", fmt.Sprintf("%d", minInterval))

	var err error
	parser := config.NewParser()
	config.Opts, err = parser.ParseEnvironmentVariables()
	if err != nil {
		t.Fatalf(`Parsing failure: %v`, err)
	}
	feed := &Feed{}
	weeklyCount := minInterval / 2
	feed.ScheduleNextCheck(weeklyCount)

	if feed.NextCheckAt.IsZero() {
		t.Error(`The next_check_at must be set`)
	}

	if feed.NextCheckAt.Before(time.Now().Add(time.Minute * time.Duration(minInterval))) {
		t.Error(`The next_check_at should not be before the now + min interval`)
	}
}

func TestFeedScheduleNextCheckEntryFrequencyFactor(t *testing.T) {
	factor := 2
	os.Clearenv()
	os.Setenv("POLLING_SCHEDULER", "entry_frequency")
	os.Setenv("SCHEDULER_ENTRY_FREQUENCY_FACTOR", fmt.Sprintf("%d", factor))

	var err error
	parser := config.NewParser()
	config.Opts, err = parser.ParseEnvironmentVariables()
	if err != nil {
		t.Fatalf(`Parsing failure: %v`, err)
	}
	feed := &Feed{}
	weeklyCount := 7
	feed.ScheduleNextCheck(weeklyCount)

	if feed.NextCheckAt.IsZero() {
		t.Error(`The next_check_at must be set`)
	}

	if feed.NextCheckAt.After(time.Now().Add(time.Minute * time.Duration(config.Opts.SchedulerEntryFrequencyMaxInterval()/factor))) {
		t.Error(`The next_check_at should not be after the now + factor * count`)
	}
}
