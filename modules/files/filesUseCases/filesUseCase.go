package filesUseCases

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/files"
)

type IFilesUseCases interface {
	UploadToGCP(req []*files.FileReq) ([]*files.FileRes, error)
	// uploadWorkers(ctx context.Context, client *storage.Client, jobsCh <-chan *files.FileReq, resultsCh chan<- *files.FileRes, errsCh chan<- error)
	DeleteFileOnGCP(req []*files.DeleteFileReq) error
}

type filesUseCase struct {
	cfg config.IConfig
}

func FilesUseCase(cfg config.IConfig) IFilesUseCases {
	return &filesUseCase{cfg: cfg}
}

type filesPub struct {
	bucket      string
	destination string
	file        *files.FileRes
}

// makePublic gives all users read access to an object.
func (f *filesPub) makePublic(ctx context.Context, client *storage.Client) error {

	acl := client.Bucket(f.bucket).Object(f.destination).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %w", err)
	}
	fmt.Printf("Blob %v is now publicly accessible.\n", f.destination)
	return nil
}

func (u *filesUseCase) uploadWorkers(ctx context.Context, client *storage.Client, jobsCh <-chan *files.FileReq, resultsCh chan<- *files.FileRes, errsCh chan<- error) {

	for job := range jobsCh {

		container, err := job.File.Open()
		if err != nil {
			errsCh <- err
			return
		}

		b, err := ioutil.ReadAll(container)

		if err != nil {
			errsCh <- err
			return
		}

		buff := bytes.NewBuffer(b)

		wc := client.Bucket(u.cfg.App().Gcpbucket()).Object(job.Destination).NewWriter(ctx)

		if _, err = io.Copy(wc, buff); err != nil {
			errsCh <- fmt.Errorf("io.Copy: %w", err)
			return
		}
		if err := wc.Close(); err != nil {
			errsCh <- fmt.Errorf("Writer.Close: %w", err)
			return
		}

		fmt.Printf("%v uploaded to %v.\n", job.FileName, job.Extension)

		newFile := &filesPub{
			file: &files.FileRes{
				FileName: job.FileName,
				Url:      fmt.Sprintf("https://storage.googleapis.com/%s/%s", u.cfg.App().Gcpbucket(), job.Destination),
			},
			bucket:      u.cfg.App().Gcpbucket(),
			destination: job.Destination,
		}

		if err := newFile.makePublic(ctx, client); err != nil {
			errsCh <- err
			return
		}

		errsCh <- nil
		resultsCh <- newFile.file

	}

}

func (u *filesUseCase) UploadToGCP(req []*files.FileReq) ([]*files.FileRes, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	jobsCh := make(chan *files.FileReq, len(req))
	resultsCh := make(chan *files.FileRes, len(req))
	errsCh := make(chan error, len(req))

	res := make([]*files.FileRes, 0)

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5

	for i := 0; i < numWorkers; i++ {
		// worker
		go u.uploadWorkers(ctx, client, jobsCh, resultsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh

		if err != nil {
			return nil, err
		}

		result := <-resultsCh
		res = append(res, result)

	}

	return res, nil
}

func (u *filesUseCase) deleteFileWorkers(ctx context.Context, client *storage.Client, jobs <-chan *files.DeleteFileReq, errs chan<- error) {

	for job := range jobs {
		o := client.Bucket(u.cfg.App().Gcpbucket()).Object(job.Destination)

		// Optional: set a generation-match precondition to avoid potential race
		// conditions and data corruptions. The request to delete the file is aborted
		// if the object's generation number does not match your precondition.
		attrs, err := o.Attrs(ctx)
		if err != nil {
			errs <- fmt.Errorf("object.Attrs: %w", err)
			return
		}
		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

		if err := o.Delete(ctx); err != nil {
			errs <- fmt.Errorf("Object(%q).Delete: %w", job.Destination, err)
			return
		}
		fmt.Printf("Blob %v deleted", job.Destination)

		errs <- nil
	}
}

func (u *filesUseCase) DeleteFileOnGCP(req []*files.DeleteFileReq) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	jobsCh := make(chan *files.DeleteFileReq, len(req))
	errsCh := make(chan error, len(req))

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5

	for i := 0; i < numWorkers; i++ {
		// worker
		// go u.uploadWorkers(ctx, client, jobsCh, resultsCh, errsCh)
		go u.deleteFileWorkers(ctx, client, jobsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh

		if err != nil {
			return err
		}
	}

	return nil

}
