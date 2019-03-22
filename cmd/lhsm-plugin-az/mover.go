// Copyright (c) 2018 DDN. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
    "context"
	"fmt"
	"net/url"
    "os"
	"path"
	"time"

    "github.com/pkg/errors"

	"github.com/Azure/azure-storage-blob-go/azblob"
    "github.com/whamcloud/go-lustre"
    "github.com/whamcloud/go-lustre/fs"
    "github.com/whamcloud/go-lustre/status"

	"github.com/edwardsp/lemur/dmplugin"
	"github.com/intel-hpdd/logging/debug"
	"github.com/pborman/uuid"
)

// Mover is an az data mover
type Mover struct {
	name  string
	creds *azblob.SharedKeyCredential
	cfg   *archiveConfig
}

// AzMover returns a new *Mover
func AzMover(cfg *archiveConfig, creds *azblob.SharedKeyCredential, archiveID uint32) *Mover {
	return &Mover{
		name:  fmt.Sprintf("az-%d", archiveID),
		creds: creds,
		cfg:   cfg,
	}
}

func newFileID() string {
	return uuid.New()
}

func (m *Mover) destination(id string) string {
	return path.Join(m.cfg.Prefix,
		"o",
		id)
}

/*
func (m *Mover) newUploader() *s3manager.Uploader {
	// can configure stuff here with custom setters
	var partSize = func(u *s3manager.Uploader) {
		u.PartSize = m.cfg.UploadPartSize
	}
	return s3manager.NewUploaderWithClient(m.s3Svc, partSize)

}
*/

/*
func (m *Mover) newDownloader() *s3manager.Downloader {
	return s3manager.NewDownloaderWithClient(m.s3Svc)
}
*/

// Start signals the mover to begin any asynchronous processing (e.g. stats)
func (m *Mover) Start() {
	debug.Printf("%s started", m.name)
}

func (m *Mover) fileIDtoContainerPath(fileID string) (string, string, error) {
	var container, path string

	u, err := url.ParseRequestURI(fileID)
	if err == nil {
		if u.Scheme != "az" {
			return "", "", errors.Errorf("invalid URL in file_id %s", fileID)
		}
		path = u.Path[1:]
		container = u.Host
	} else {
		path = m.destination(fileID)
		container = m.cfg.Container
	}
	debug.Printf("Parsed %s -> %s / %s", fileID, container, path)
	return container, path, nil
}

// Archive fulfills an HSM Archive request
func (m *Mover) Archive(action dmplugin.Action) error {
	debug.Printf("%s id:%d archive %s %s", m.name, action.ID(), action.PrimaryPath(), action.UUID())
    rate.Mark(1)
	start := time.Now()

	fileID := newFileID()
	fileKey := m.destination(fileID)

    fid, _ := lustre.ParseFid(action.UUID())
    rootDir, _ := fs.MountRoot(action.PrimaryPath())
    names, _ := status.FidPathnames(rootDir, fid)
    debug.Printf("FILENAME: %s",names[0])

    p := azblob.NewPipeline(m.creds, azblob.PipelineOptions{})
    cURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", m.cfg.AzStorageAccount, m.cfg.Container))
    containerURL := azblob.NewContainerURL(*cURL, p)
    ctx := context.Background()
    blobURL := containerURL.NewBlockBlobURL(fileKey)

    debug.Printf("\naction.PrimaryPath() = %s\naction.WritePath() = %s", action.PrimaryPath(), action.PrimaryPath())
    file, _ := os.Open(action.PrimaryPath())
    defer file.Close()

    _, err := azblob.UploadFileToBlockBlob(
        ctx,
        file,
        blobURL,
        azblob.UploadToBlockBlobOptions{
		    BlockSize:   m.cfg.UploadPartSize,
		    Parallelism: uint16(m.cfg.NumThreads),
        })
	if err != nil {
		return errors.Wrap(err, "upload failed")
	}

	debug.Printf("%s id:%d Archived %d bytes in %v from %s to %s/%s", m.name, action.ID(), action.Length(),
		time.Since(start),
		action.PrimaryPath(),
		cURL, fileKey)

	u := url.URL{
		Scheme: "az",
		Host:   cURL.String(),
		Path:   fileKey,
	}

	action.SetUUID(fileID)
	action.SetURL(u.String())
	action.SetActualLength(action.Length())
	return nil
}

// Restore fulfills an HSM Restore request
func (m *Mover) Restore(action dmplugin.Action) error {
	debug.Printf("%s id:%d restore %s %s", m.name, action.ID(), action.PrimaryPath(), action.UUID())
	rate.Mark(1)

	start := time.Now()
	if action.UUID() == "" {
		return errors.Errorf("Missing file_id on action %d", action.ID())
	}
	container, srcObj, err := m.fileIDtoContainerPath(action.UUID())
	if err != nil {
		return errors.Wrap(err, "fileIDtoContainerPath failed")
	}

    p := azblob.NewPipeline(m.creds, azblob.PipelineOptions{})
    cURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", m.cfg.AzStorageAccount, container))
    containerURL := azblob.NewContainerURL(*cURL, p)
    ctx := context.Background()
    blobURL := containerURL.NewBlobURL(srcObj)

    blobProp, err := blobURL.GetProperties(ctx, azblob.BlobAccessConditions{})
	if err != nil {
		return errors.Wrapf(err, "GetProperties on %s failed", srcObj)
	}
    contentLen := blobProp.ContentLength()
	debug.Printf("obj %s, size %d", srcObj, contentLen)

    debug.Printf("\naction.PrimaryPath() = %s\naction.WritePath() = %s", action.PrimaryPath(), action.WritePath())
    file, _ := os.Create(action.WritePath())
    defer file.Close()
    err = azblob.DownloadBlobToFile(
        ctx, blobURL, 0, 0, file,
        azblob.DownloadFromBlobOptions{
            BlockSize: m.cfg.UploadPartSize,
            Parallelism: uint16(m.cfg.NumThreads),
        })

	if err != nil {
		return errors.Errorf("az.Download() of %s failed: %s", srcObj, err)
	}

	debug.Printf("%s id:%d Restored %d bytes in %v from %s to %s", m.name, action.ID(), contentLen,
		time.Since(start),
		srcObj,
		action.PrimaryPath())
	action.SetActualLength(contentLen)
    return nil
}

// Remove fulfills an HSM Remove request
func (m *Mover) Remove(action dmplugin.Action) error {
	debug.Printf("%s id:%d remove %s %s", m.name, action.ID(), action.PrimaryPath(), action.UUID())
    rate.Mark(1)
	if action.UUID() == "" {
		return errors.New("Missing file_id")
	}

	container, srcObj, err := m.fileIDtoContainerPath(string(action.UUID()))

    p := azblob.NewPipeline(m.creds, azblob.PipelineOptions{})
    cURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", m.cfg.AzStorageAccount, container))
    containerURL := azblob.NewContainerURL(*cURL, p)
    ctx := context.Background()
    blobURL := containerURL.NewBlobURL(srcObj)
	_, err = blobURL.Delete(ctx,
        "",
        azblob.BlobAccessConditions{})
	return errors.Wrap(err, "delete object failed")
    return nil
}
