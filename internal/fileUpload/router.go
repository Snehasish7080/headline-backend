package fileUpload

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
	"github.com/zone/headline/config"
)

func AddFileRoutes(app *fiber.App, storage *FileUploadStorage, env config.EnvVars) {

	// store := filestore.FileStore{
	// 	Path: "/home/zone/Projects/upload",
	// }

	s3Config := &aws.Config{
		Region:      aws.String("eu-north-1"),
		Credentials: credentials.NewStaticCredentials(env.S3_ACCESS_KEY, env.S3_SECRET_KEY, ""),
	}
	s3Locker := s3store.New(env.S3_BUCKET, s3.New(session.Must(session.NewSession()), s3Config))

	composer := tusd.NewStoreComposer()
	// store.UseIn(composer)
	s3Locker.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})

	if err != nil {
		panic(fmt.Errorf("Unable to create handler: %s", err))
	}

	go func() {
		for {
			event := <-handler.CompleteUploads
			fmt.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()

	app.Post("/files", adaptor.HTTPHandlerFunc(handler.PostFile))
	app.Head("/files/:id", adaptor.HTTPHandlerFunc(handler.HeadFile))
	app.Patch("/files/:id", adaptor.HTTPHandlerFunc(handler.PatchFile))
	app.Get("/files/:id", adaptor.HTTPHandlerFunc(handler.GetFile))

}
