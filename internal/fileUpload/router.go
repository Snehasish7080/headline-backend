package fileUpload

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
)

func AddFileRoutes(app *fiber.App, storage *FileUploadStorage) {

	store := filestore.FileStore{
		Path: "/home/zone/Projects/upload",
	}
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

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
