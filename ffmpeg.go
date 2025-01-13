package helpers

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MediaKind, medya türünü belirtir.
type MediaKind string

const (
	ImageKind MediaKind = "image"
	VideoKind MediaKind = "video"
)

// MediaOptions, medya dönüşümü için seçenekleri içerir.
type MediaOptions struct {
	Width       int
	Height      int
	Crop        string // w:h:x:y
	Quality     int    // 0-100
	ResizeScale int    // 0-100
	CutVideo    string // start:duration
	Format      string // webp, avif, jpeg, png, webm, gif
	FrameTime   string // ss (snapshot time)
	StartTime   int    // GIF veya kesit başlangıç süresi (saniye)
	EndTime     int    // GIF veya kesit bitiş süresi (saniye)
	FrameRate   int    // Frame rate (FPS)
	Bitrate     string // Bitrate (örneğin, "1M" for 1 Mbps)
	CRF         int    // Constant Rate Factor (0-51, düşük değerler daha yüksek kalite)
}

// MediaConverter, medya dönüşümü için zincirleme API sunar.
type MediaConverter struct {
	kind    MediaKind
	input   string
	options MediaOptions
}

// NewConverter, yeni bir MediaConverter örneği oluşturur.
func NewConverter(kind MediaKind, input string) *MediaConverter {
	return &MediaConverter{
		kind:    kind,
		input:   input,
		options: MediaOptions{Quality: 90}, // Varsayılan kalite
	}
}

// FromFiber, Fiber context'inden parametreleri alarak MediaConverter oluşturur.
// Uygun olmayan parametreler için uyarı mesajı döndürür ve bu parametreleri görmezden gelir.
func FromFiber(c *fiber.Ctx, kind MediaKind, input string) (*MediaConverter, []string) {
	converter := NewConverter(kind, input)
	warnings := make([]string, 0)

	// Sorgu parametrelerini koşullu olarak işle
	if width := c.QueryInt("width"); width > 0 {
		converter.Resize(width, 0)
	}
	if height := c.QueryInt("height"); height > 0 {
		converter.Resize(0, height)
	}
	if quality := c.QueryInt("quality"); quality > 0 {
		converter.Quality(quality)
	}
	if format := c.Query("format"); format != "" {
		converter.Format(format)
	}

	// Video özelinde parametreler
	if kind == VideoKind {
		if fps := c.QueryInt("fps"); fps > 0 {
			converter.SetFrameRate(fps)
		}
		if start := c.QueryInt("start"); start > 0 {
			converter.SetDuration(start, 0)
		}
		if end := c.QueryInt("end"); end > 0 {
			converter.SetDuration(0, end)
		}
		if bitrate := c.Query("bitrate"); bitrate != "" {
			converter.SetBitrate(bitrate)
		}
		if crf := c.QueryInt("crf"); crf > 0 {
			converter.SetCRF(crf)
		}
	} else {
		// Video özelinde parametreler ImageKind için uygun değilse uyarı ekle
		if c.Query("fps") != "" {
			warnings = append(warnings, "fps parameter is ignored for images")
		}
		if c.Query("start") != "" {
			warnings = append(warnings, "start parameter is ignored for images")
		}
		if c.Query("end") != "" {
			warnings = append(warnings, "end parameter is ignored for images")
		}
		if c.Query("bitrate") != "" {
			warnings = append(warnings, "bitrate parameter is ignored for images")
		}
		if c.Query("crf") != "" {
			warnings = append(warnings, "crf parameter is ignored for images")
		}
	}

	// Kırpma parametresi (hem video hem görüntü için geçerli)
	if crop := c.Query("crop"); crop != "" {
		converter.Crop(crop)
	}

	return converter, warnings
}

// Resize, medyanın boyutunu ayarlar.
func (mc *MediaConverter) Resize(width, height int) *MediaConverter {
	mc.options.Width = width
	mc.options.Height = height
	return mc
}

func (mc *MediaConverter) ResizeScale(scale int) *MediaConverter {
	mc.options.ResizeScale = scale
	return mc
}

// Crop, medyayı kırpar.
func (mc *MediaConverter) Crop(dimensions string) *MediaConverter {
	mc.options.Crop = dimensions
	return mc
}

// Quality, medyanın kalitesini ayarlar.
func (mc *MediaConverter) Quality(quality int) *MediaConverter {
	mc.options.Quality = quality
	return mc
}

// Format, çıktı formatını ayarlar.
func (mc *MediaConverter) Format(format string) *MediaConverter {
	mc.options.Format = format
	return mc
}

// SetDuration, GIF veya kesit süresini ayarlar (başlangıç ve bitiş saniyeleri).
func (mc *MediaConverter) SetDuration(start, end int) *MediaConverter {
	mc.options.StartTime = start
	mc.options.EndTime = end
	return mc
}

// SetFrameRate, videonun frame rate'ini ayarlar.
func (mc *MediaConverter) SetFrameRate(fps int) *MediaConverter {
	mc.options.FrameRate = fps
	return mc
}

// SetBitrate, videonun bitrate'ini ayarlar.
func (mc *MediaConverter) SetBitrate(bitrate string) *MediaConverter {
	mc.options.Bitrate = bitrate
	return mc
}

// SetCRF, videonun CRF (Constant Rate Factor) değerini ayarlar.
func (mc *MediaConverter) SetCRF(crf int) *MediaConverter {
	mc.options.CRF = crf
	return mc
}

// ExtractFrame, videodan belirli bir zamanda bir kare alır.
func (mc *MediaConverter) ExtractFrame(time string) *MediaConverter {
	mc.options.FrameTime = time
	return mc
}

// Convert, medya dosyasını dönüştürür ve kaydeder.
func (mc *MediaConverter) Convert(outputDir string, outPutFileName string) (string, error) {
	// Çıktı dosyasının yolunu oluştur
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.%s", outPutFileName, mc.getOutputExtension()))

	// FFmpeg komutunu oluştur ve çalıştır
	args := mc.buildFFmpegArgs(outputFile)
	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v, output: %s", err, string(output))
	}

	return outputFile, nil
}

// ToGIF, videoyu GIF'e dönüştürür.
func (mc *MediaConverter) ToGIF(outputDir, outPutFileName string) (string, error) {
	// Çıktı dosyasının yolunu oluştur
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.gif", outPutFileName))

	// FFmpeg komutunu oluştur
	args := mc.buildGIFArgs(outputFile)

	// FFmpeg komutunu çalıştır
	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v, output: %s", err, string(output))
	}

	return outputFile, nil
}

// getOutputExtension, çıktı dosyasının uzantısını belirler.
func (mc *MediaConverter) getOutputExtension() string {
	if mc.options.Format != "" {
		return mc.options.Format
	}
	if mc.kind == VideoKind {
		return "webm"
	}
	return "webp"
}

// buildFFmpegArgs, FFmpeg komut satırı argümanlarını oluşturur.
func (mc *MediaConverter) buildFFmpegArgs(outputFile string) []string {
	args := []string{"-i", mc.input}

	// Video kesme
	if mc.options.CutVideo != "" {
		args = append(args, "-ss", strings.Split(mc.options.CutVideo, ":")[0])
		args = append(args, "-t", strings.Split(mc.options.CutVideo, ":")[1])
	}

	// Videodan kare alma
	if mc.options.FrameTime != "" {
		args = append(args, "-ss", mc.options.FrameTime, "-vframes", "1")
	}
	if mc.options.ResizeScale > 0 {
		args = append(args, "-vf", fmt.Sprintf("scale=iw*%d/100:ih*%d/100", mc.options.ResizeScale, mc.options.ResizeScale))
	}
	// Boyutlandırma ve kırpma
	if mc.options.Width > 0 || mc.options.Height > 0 {
		args = append(args, "-vf", fmt.Sprintf("scale=%d:%d", mc.options.Width, mc.options.Height))
	}
	if mc.options.Crop != "" {
		args = append(args, "-vf", fmt.Sprintf("crop=%s", mc.options.Crop))
	}

	// Frame rate
	if mc.options.FrameRate > 0 {
		args = append(args, "-r", strconv.Itoa(mc.options.FrameRate))
	}

	// Bitrate
	if mc.options.Bitrate != "" {
		args = append(args, "-b:v", mc.options.Bitrate)
	}

	// CRF (Constant Rate Factor)
	if mc.options.CRF > 0 {
		args = append(args, "-crf", strconv.Itoa(mc.options.CRF))
	}

	// Kalite
	if mc.options.Quality > 0 {
		args = append(args, "-q:v", strconv.Itoa(mc.options.Quality))
	}

	// Çıktı formatı
	args = append(args, "-y", outputFile)
	return args
}

// buildGIFArgs, GIF oluşturmak için FFmpeg argümanlarını hazırlar.
func (mc *MediaConverter) buildGIFArgs(outputFile string) []string {
	args := []string{"-i", mc.input}

	// Süre ayarı (start ve end)
	if mc.options.StartTime > 0 || mc.options.EndTime > 0 {
		args = append(args, "-ss", strconv.Itoa(mc.options.StartTime))
		if mc.options.EndTime > mc.options.StartTime {
			args = append(args, "-t", strconv.Itoa(mc.options.EndTime-mc.options.StartTime))
		}
	}

	// Boyutlandırma
	if mc.options.Width > 0 || mc.options.Height > 0 {
		args = append(args, "-vf", fmt.Sprintf("scale=%d:%d", mc.options.Width, mc.options.Height))
	}

	// Frame rate
	if mc.options.FrameRate > 0 {
		args = append(args, "-r", strconv.Itoa(mc.options.FrameRate))
	}

	// Çıktı formatı ve dosya yolu
	args = append(args, "-y", outputFile)
	return args
}

/*
Usage Example


package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Fiber uygulamasını başlat
	app := fiber.New()

	// Logger middleware'ini ekle
	app.Use(logger.New())

	// Medya dönüşümü endpoint'i
	app.Get("/convert", func(c *fiber.Ctx) error {
		// Giriş dosyasını ve medya türünü al
		inputFile := c.Query("input")
		kind := media.VideoKind // Varsayılan olarak video kabul ediyoruz

		// Giriş dosyasının varlığını kontrol et
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Input file does not exist",
			})
		}

		// FromFiber ile MediaConverter oluştur
		converter := media.FromFiber(c, kind, inputFile)

		// Dönüşümü gerçekleştir
		outputFile, err := converter.Convert("./output")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Başarılı yanıt döndür
		return c.JSON(fiber.Map{
			"message":    "Conversion successful",
			"outputFile": outputFile,
		})
	})

	// Uygulamayı başlat
	log.Fatal(app.Listen(":3000"))
}

*/
