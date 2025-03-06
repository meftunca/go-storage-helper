package helpers_tests

import (
	helpers "go-ffmpeg-helper"
	"os"
	"testing"
)

var inputFile = "~/go-workspace/src/github.com/meftunca/go-ffmpeg-helper/TestData/huge.jpg"
var outputFilePath = "~/go-workspace/src/github.com/meftunca/go-ffmpeg-helper/TestData/output"

func TestResizeImage(t *testing.T) {
	// Giriş dosyası

	// Yeni bir MediaConverter örneği oluştur
	converter := helpers.NewConverter("jpg", inputFile)

	// Resmi yeniden boyutlandır
	converter.Resize(800, 600)

	// Dönüşümü gerçekleştir
	outputFile, err := converter.Convert(outputFilePath, "resized-w800-h600")
	if err != nil {
		t.Fatalf("Resize failed: %v", err)
	}

	// Çıktı dosyasının varlığını kontrol et
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file does not exist: %s", outputFile)
	}

	// Temizleme: Çıktı dosyasını sil
	// os.Remove(outputFile)
}

func TestChangeImageQuality(t *testing.T) {
	// Giriş dosyası

	// Yeni bir MediaConverter örneği oluştur
	converter := helpers.NewConverter("jpg", inputFile)

	// Kaliteyi ayarla
	converter.Quality(50)

	// Dönüşümü gerçekleştir
	outputFile, err := converter.Convert(outputFilePath, "q50")
	if err != nil {
		t.Fatalf("Quality change failed: %v", err)
	}

	// Çıktı dosyasının varlığını kontrol et
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file does not exist: %s", outputFile)
	}

	// Temizleme: Çıktı dosyasını sil
	// os.Remove(outputFile)
}

func TestConvertImageFormat(t *testing.T) {
	// Giriş dosyası

	// Yeni bir MediaConverter örneği oluştur
	converter := helpers.NewConverter("jpg", inputFile)

	// Formatı değiştir
	converter.Format("png")

	// Dönüşümü gerçekleştir
	outputFile, err := converter.Convert(outputFilePath, "output-png")
	if err != nil {
		t.Fatalf("Format conversion failed: %v", err)
	}

	// Çıktı dosyasının varlığını kontrol et
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file does not exist: %s", outputFile)
	}

	// Temizleme: Çıktı dosyasını sil
	// os.Remove(outputFile)
}

func TestCombinedOperations(t *testing.T) {
	// Giriş dosyası

	// Yeni bir MediaConverter örneği oluştur
	converter := helpers.NewConverter("jpg", inputFile)

	// Resmi yeniden boyutlandır, kaliteyi ayarla ve formatı değiştir
	converter.Resize(800, 600).Quality(75).Format("webp")

	// Dönüşümü gerçekleştir
	outputFile, err := converter.Convert(outputFilePath, "w800-h600-q75")
	if err != nil {
		t.Fatalf("Combined operations failed: %v", err)
	}

	// Çıktı dosyasının varlığını kontrol et
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file does not exist: %s", outputFile)
	}

	// Temizleme: Çıktı dosyasını sil
	// os.Remove(outputFile)
}
