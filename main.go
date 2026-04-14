package main

import (
	"gofr.dev/pkg/gofr"

	infraconv "pdfapi/internal/infrastructure/converter"
	infragen "pdfapi/internal/infrastructure/generator"
	infraop "pdfapi/internal/infrastructure/operation"

	"pdfapi/internal/handler"
)

func main() {
	app := gofr.New()

	// Infrastructure layer
	convSvc := infraconv.New()
	genSvc := infragen.New()
	opSvc := infraop.New()

	// Handler layer
	convH := handler.NewConverterHandler(convSvc)
	genH := handler.NewGeneratorHandler(genSvc)
	opH := handler.NewOperationHandler(opSvc)

	// Converter routes
	app.POST("/api/pdf/to-images", convH.PDFToImages)
	app.POST("/api/images/to-pdf", convH.ImagesToPDF)

	// Generator routes
	app.POST("/api/html/to-pdf", genH.HTMLToPDF)
	app.POST("/api/markdown/to-pdf", genH.MarkdownToPDF)

	// Operation routes
	app.POST("/api/pdf/merge", opH.MergePDF)
	app.POST("/api/pdf/split", opH.SplitPDF)
	app.POST("/api/pdf/compress", opH.CompressPDF)
	app.POST("/api/pdf/watermark", opH.WatermarkPDF)
	app.POST("/api/pdf/extract-text", opH.ExtractText)
	app.POST("/api/pdf/decrypt", opH.DecryptPDF)

	app.Run()
}
