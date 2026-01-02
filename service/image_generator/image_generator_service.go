package image_generator

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Out-Of-India-Theory/oit-go-commons/logging"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/service/image_uploader"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/service/supply"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
	"image"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ImageGeneratorService struct {
	logger        *zap.Logger
	supplyService supply.Service
	imageUploader image_uploader.Service
}

func InitImageGeneratorService(ctx context.Context, supplyService supply.Service, imageUploader image_uploader.Service) *ImageGeneratorService {
	return &ImageGeneratorService{
		logger:        logging.WithContext(ctx),
		supplyService: supplyService,
		imageUploader: imageUploader,
	}
}

func (s *ImageGeneratorService) GenerateImage(ctx context.Context, supplyId int) error {
	supplyDetails, err := s.supplyService.GetSupplyDetails(ctx, supplyId)
	if err != nil {
		return err
	}

	prashnaTranslations := map[string]string{
		"en": fmt.Sprintf("Your Jyotish Consultation\nhas been assigned to\n%s", supplyDetails.Data.NameV1["en"]),
		"hi": fmt.Sprintf("आपकी ज्योतिष परामर्श सेवा\n%s को\nसौंप दी गई है।", supplyDetails.Data.NameV1["hi"]),
		"kn": fmt.Sprintf("ನಿಮ್ಮ ಜ್ಯೋತಿಷ್ಯ ಸಲಹೆ\n%s ಗೆ\nಹಂಚಲಾಗಿದೆ", supplyDetails.Data.NameV1["kn"]),
		"gu": fmt.Sprintf("તમારી જ્યોતિષ પરામર્શ સેવા\n%s ને\nસોંપવામાં આવી છે।", supplyDetails.Data.NameV1["gu"]),
		"ta": fmt.Sprintf("உங்கள் ஜோதிட ஆலோசனை\n%s க்கு\nஒதுக்கப்பட்டுள்ளது", supplyDetails.Data.NameV1["ta"]),
		"te": fmt.Sprintf("మీ జ్యోతిష సంప్రదింపు సేవ\n%s కు\nకేటాయించబడింది", supplyDetails.Data.NameV1["te"]),
		"mr": fmt.Sprintf("आपली ज्योतिष सल्ला सेवा\n%s कडे\nसोपविण्यात आली आहे।", supplyDetails.Data.NameV1["mr"]),
	}

	jyotishaTranslations := map[string]string{
		"en": fmt.Sprintf("Your Jyotisha Consultation\nhas been assigned to\n%s", supplyDetails.Data.NameV1["en"]),
		"hi": fmt.Sprintf("आपकी ज्योतिष परामर्श सेवा\n%s को\nसौंप दी गई है।", supplyDetails.Data.NameV1["hi"]),
		"kn": fmt.Sprintf("ನಿಮ್ಮ ಜ್ಯೋತಿಷ್ಯ ಸಲಹೆ\n%s ಗೆ\nಹಂಚಲಾಗಿದೆ", supplyDetails.Data.NameV1["kn"]),
		"gu": fmt.Sprintf("તમારી જ્યોતિષ પરામર્શ સેવા\n%s ને\nસોંપવામાં આવી છે।", supplyDetails.Data.NameV1["gu"]),
		"ta": fmt.Sprintf("உங்கள் ஜோதிட ஆலோசனை\n%s க்கு\nஒதுக்கப்பட்டுள்ளது", supplyDetails.Data.NameV1["ta"]),
		"te": fmt.Sprintf("మీ జ్యోతిష సంప్రదింపు సేవ\n%s కు\nకేటాయించబడింది", supplyDetails.Data.NameV1["te"]),
		"mr": fmt.Sprintf("आपली ज्योतिष सल्ला सेवा\n%s कडे\nसोपविण्यात आली आहे।", supplyDetails.Data.NameV1["mr"]),
	}

	experienceText := map[string]string{
		"en": "%d yrs experience",
		"hi": "%d वर्षों का अनुभव",
		"kn": "%d ವರ್ಷದ ಅನುಭವ",
		"gu": "%d વર્ષનો અનુભવ",
		"ta": "%d வருட அனுபவம்",
		"te": "%d సంవత్సరాల అనుభవం",
		"mr": "%d वर्षांचा अनुभव",
	}

	checkStatusText := map[string]string{
		"en": "CHECK ORDER STATUS",
		"hi": "ऑर्डर स्थिति देखें",
		"kn": "ಆರ್ಡರ್ ಸ್ಥಿತಿಯನ್ನು ಪರಿಶೀಲಿಸಿ",
		"gu": "ઓર્ડર સ્થિતિ તપાસો",
		"ta": "ஆர்டர் நிலையை சரிபார்க்கவும்",
		"te": "ఆర్డర్ స్థితిని తనిఖీ చేయండి",
		"mr": "ऑर्डर स्थिती तपासा",
	}

	fontMap := map[string]string{
		"en": "assets/fonts/Ubuntu-M.ttf",
		"hi": "assets/fonts/NotoSansDevanagari-Regular.ttf",
		"mr": "assets/fonts/NotoSansDevanagari-Regular.ttf",
		"kn": "assets/fonts/NotoSansKannada-Regular.ttf",
		"gu": "assets/fonts/NotoSansGujarati-Regular.ttf",
		"ta": "assets/fonts/NotoSansTamil-Regular.ttf",
		"te": "assets/fonts/NotoSansTelugu-Regular.ttf",
	}

	bgPath, _ := filepath.Abs("assets/images/background.png")

	personImg, err := downloadImage(supplyDetails.Data.ImageWithoutBackground)
	if err != nil {
		return fmt.Errorf("failed to download supply image_generator: %w", err)
	}

	personImgBytes, err := imageToPNGBytes(personImg)
	if err != nil {
		return err
	}

	personBase64 := base64.StdEncoding.EncodeToString(personImgBytes)

	for _, lang := range supplyDetails.Data.Languages {
		fontAbsPath, err := filepath.Abs(fontMap[lang])
		if err != nil {
			return err
		}

		// ---- PRASHNA IMAGE ----
		prashnaValues := map[string]string{
			"TITLE_TEXT":      prashnaTranslations[lang],
			"EXPERIENCE_TEXT": fmt.Sprintf(experienceText[lang], supplyDetails.Data.YearsOfExperience),
			"CTA_TEXT":        checkStatusText[lang],
			"FONT_PATH":       "file://" + fontAbsPath,
			"BG_PATH":         "file://" + bgPath,
			"SUPPLY_IMAGE":    "data:image_generator/png;base64," + personBase64,
			"LANG":            lang,
		}

		imgBytes, err := GenerateHTMLToImage(ctx, "assets/template.html", prashnaValues)
		if err != nil {
			return fmt.Errorf("prashna image_generator failed (%s): %w", lang, err)
		}

		if _, err := s.imageUploader.UploadToS3(ctx, fmt.Sprintf("%d_prashna_%s", supplyId, lang), imgBytes); err != nil {
			return err
		}

		// ---- JYOTISHA IMAGE ----
		jyotishaValues := map[string]string{
			"TITLE_TEXT":      jyotishaTranslations[lang],
			"EXPERIENCE_TEXT": fmt.Sprintf(experienceText[lang], supplyDetails.Data.YearsOfExperience),
			"CTA_TEXT":        checkStatusText[lang],
			"FONT_PATH":       "file://" + fontAbsPath,
			"BG_PATH":         "file://" + bgPath,
			"SUPPLY_IMAGE":    "data:image_generator/png;base64," + personBase64,
			"LANG":            lang,
		}

		imgBytes, err = GenerateHTMLToImage(ctx, "assets/template.html", jyotishaValues)
		if err != nil {
			return fmt.Errorf("jyotisha image_generator failed (%s): %w", lang, err)
		}

		if _, err := s.imageUploader.UploadToS3(ctx, fmt.Sprintf("%d_jyotisha_%s", supplyId, lang), imgBytes); err != nil {
			return err
		}
	}
	return nil
}

func downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get image_generator from URL: %w", err)
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image_generator: %w", err)
	}
	return img, nil
}

func imageToPNGBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GenerateHTMLToImage(ctx context.Context, htmlPath string, values map[string]string) ([]byte, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("allow-file-access-from-files", true),
		//chromedp.Flag("disable-dev-shm-usage", true),
		//chromedp.Flag("single-process", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	// Read template
	htmlData, err := os.ReadFile(htmlPath)
	if err != nil {
		return nil, err
	}

	html := string(htmlData)
	for k, v := range values {
		html = strings.ReplaceAll(html, "{{"+k+"}}", v)
	}

	// Write temp HTML
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("render_%d.html", time.Now().UnixNano()))
	if err := os.WriteFile(tmpFile, []byte(html), 0644); err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile)

	absPath, _ := filepath.Abs(tmpFile)

	var buf []byte
	err = chromedp.Run(ctx,
		chromedp.EmulateViewport(1120, 800, chromedp.EmulateScale(2)), // MUST match CSS
		chromedp.Navigate("file://"+absPath),
		chromedp.Sleep(2*time.Second),
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
