package misp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SearchEventsRequest struct {
	Limit int `json:"limit,omitempty"`
	Page  int `json:"page,omitempty"`

	// filters
	Tag  string   `json:"tag,omitempty"`
	Tags []string `json:"tags,omitempty"` // tags (e.g. "tlp:white")

	Last      string `json:"last,omitempty"`      // Last events within the last duration (e.g. 1d, 1h, 1m)
	From      string `json:"from,omitempty"`      // Events from a specific date (e.g. 2021-01-01)
	To        string `json:"to,omitempty"`        // Events to a specific date (e.g. 2021-01-01)
	Timestamp string `json:"timestamp,omitempty"` // Events after a specific timestamp (e.g. 1612137600)
}
type SearchEventsResponse struct {
	Response []RelatedEvent `json:"response"`
}

func (c *Client) SearchEvents(ctx context.Context, req *SearchEventsRequest) ([]RelatedEvent, error) {
	if req.Tag == "" && len(req.Tags) == 1 {
		req.Tag = req.Tags[0]
		req.Tags = nil
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequestExpectOK(ctx, http.MethodPost, "/events/restSearch",
		bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	var body SearchEventsResponse
	if err := json.NewDecoder(respBody).Decode(&body); err != nil {
		return nil, err
	}
	return body.Response, nil
}

func (c *Client) GetEvent(ctx context.Context, eventID string) (*Event, error) {
	if eventID == "" || strings.Contains(eventID, "/") {
		return nil, fmt.Errorf("invalid event id: %q", eventID)
	}

	respBody, err := c.doRequestExpectOK(ctx, http.MethodGet, "/events/view/"+eventID, nil)
	if err != nil {
		return nil, err
	}

	var body RelatedEvent
	if err := json.NewDecoder(respBody).Decode(&body); err != nil {
		return nil, err
	}
	return &body.Event, nil
}
