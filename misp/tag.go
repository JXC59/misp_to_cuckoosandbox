package misp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

func (c *Client) AddTag(ctx context.Context, tag *Tag) (*Tag, error) {
	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(tag); err != nil {
		return nil, err
	}

	resp, err := c.doRequestExpectOK(ctx, http.MethodPost, "/tags/add", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	ret := new(RelatedTag)
	if err := json.NewDecoder(resp).Decode(ret); err != nil {
		return nil, err
	}

	return &ret.Tag, nil
}

type tagAttributeRequest struct {
	AttributeID string `json:"attribute"` // Attribute ID
	Tag         string `json:"tag"`       // Tag ID or Name
}

func (c *Client) TagAttribute(ctx context.Context, attrID string, tag string) error {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(&tagAttributeRequest{
		AttributeID: attrID,
		Tag:         tag,
	})
	if err != nil {
		return err
	}

	_, err = c.doRequestExpectOK(ctx, http.MethodPost, path.Join("/attributes/addTag"), buf)
	return err
}

type tagEventRequest struct {
	EventID string `json:"event"` // Event ID
	Tag     string `json:"tag"`   // Tag ID or Name
}

type tagEventResponse struct {
	Saved bool   `json:"saved"`
	Error string `json:"errors"`
}

func (c *Client) TagEvent(ctx context.Context, eventID string, tag string) error {
	buf, err := json.Marshal(&tagEventRequest{
		EventID: eventID,
		Tag:     tag,
	})
	if err != nil {
		return err
	}

	resp, err := c.doRequestExpectOK(ctx, http.MethodPost, path.Join("/events/addTag"), bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer resp.Close()

	var respObj tagEventResponse
	if err := json.NewDecoder(resp).Decode(&respObj); err != nil {
		return err
	}

	if respObj.Error != "" {
		return fmt.Errorf("error tagging event: %s", respObj.Error)
	}
	return nil
}

func (c *Client) UntagEvent(ctx context.Context, eventID string, tag string) error {
	buf, err := json.Marshal(&tagEventRequest{
		EventID: eventID,
		Tag:     tag,
	})
	if err != nil {
		return err
	}

	resp, err := c.doRequestExpectOK(ctx, http.MethodPost, path.Join("/events/removeTag"), bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer resp.Close()

	var respObj tagEventResponse
	if err := json.NewDecoder(resp).Decode(&respObj); err != nil {
		return err
	}

	if respObj.Error != "" {
		return fmt.Errorf("error tagging event: %s", respObj.Error)
	}
	return nil
}
