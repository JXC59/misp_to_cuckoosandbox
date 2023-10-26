package misp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kdungs/zip"
	"github.com/samber/lo"
)

// DownloadAttribute downloads an attribute from MISP.
// Caller take the responsibility to close the returned reader.
func (c *Client) DownloadAttribute(ctx context.Context, attributeID string) (io.ReadCloser, error) {
	if attributeID == "" || strings.Contains(attributeID, "/") {
		return nil, fmt.Errorf("invalid attribute id: %q", attributeID)
	}

	respBody, err := c.doRequestExpectOK(ctx, http.MethodGet, "/attributes/download/"+attributeID, nil)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// DownloadInfectedSample downloads an infected sample from MISP.
// Caller take the responsibility to close the returned reader.
func (c *Client) DownloadInfectedSample(ctx context.Context, attributeID string) (io.ReadCloser, error) {
	if attributeID == "" || strings.Contains(attributeID, "/") {
		return nil, fmt.Errorf("invalid attribute id: %q", attributeID)
	}
	resp, err := c.doRequest(ctx, http.MethodGet, "/attributes/download/"+attributeID, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to copy response body: %w", err)
	}

	return ExtractInfectedSample(bytes.NewReader(buf.Bytes()), resp.ContentLength)
}

func ExtractInfectedSample(zipReader io.ReaderAt, size int64) (io.ReadCloser, error) {
	const filenameIndicator = ".filename.txt"
	const samplePassword = "infected"

	zipFile, err := zip.NewReader(zipReader, size)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %w", err)
	}

	targetName := ""

	{
		for _, file := range zipFile.File {
			if strings.HasSuffix(file.Name, filenameIndicator) {
				targetName = strings.TrimSuffix(file.Name, filenameIndicator)
				break
			}
		}
		if targetName == "" {
			return nil, fmt.Errorf("failed to find filename indicator")
		}
	}

	targetFile, ok := lo.Find(zipFile.File, func(file *zip.File) bool {
		return file.Name == targetName
	})
	if !ok {
		return nil, fmt.Errorf("failed to find target file")
	}

	if targetFile.IsEncrypted() {
		targetFile.SetPassword(samplePassword)
	}

	return targetFile.Open()
}

type addAttributeRequest struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Data  []byte `json:"data"`
}

func (c *Client) AddAttachmentAttribute(ctx context.Context, eventID string, filename string, r io.Reader) (*Attribute, error) {
	if eventID == "" || strings.Contains(eventID, "/") {
		return nil, fmt.Errorf("invalid attribute id: %q", eventID)
	}

	rawBuf, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	req := &addAttributeRequest{
		Type:  "attachment",
		Value: filename,
		Data:  rawBuf, // automatically base64 encoded
	}

	reqBuf := new(bytes.Buffer)
	if err := json.NewEncoder(reqBuf).Encode(req); err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	resp, err := c.doRequestExpectOK(ctx, http.MethodPost, "/attributes/add/"+eventID, reqBuf)
	if err != nil {
		return nil, err
	}

	var attr RelatedAttribute
	if err := json.NewDecoder(resp).Decode(&attr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &attr.Attribute, nil
}
