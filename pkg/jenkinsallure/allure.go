package jenkinsallure

import (
	"context"
	"time"

	"math"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const (
	SUMMARY_AREA_SELECTOR = `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(1)`
	SUMMARY_INFO_SELECTOR = `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(1) > div.widget__body > div > div > div.widget__column.summary-widget__chart > div > svg`

	BEHAVIORS_AREA_SELECTOR = `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(4)`
	BEHAVIORS_INFO_SELECTOR = `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(4) > div.widget__body > div > h2 > span`
)

func CaptureAllureResult(url, jobName string) ([]byte, []byte, error) {
	_ctx, _cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer _cancel()

	ctx, cancel := chromedp.NewContext(_ctx)
	defer cancel()

	// capture screenshot of an element
	var summaryResult []byte
	var behaviorsResult []byte
	if err := chromedp.Run(ctx, elementScreenshot(url, SUMMARY_AREA_SELECTOR, SUMMARY_INFO_SELECTOR, &summaryResult)); err != nil {
		return nil, nil, err
	}

	if err := chromedp.Run(ctx, elementScreenshot(url, BEHAVIORS_AREA_SELECTOR, BEHAVIORS_INFO_SELECTOR, &behaviorsResult)); err != nil {
		return nil, nil, err
	}

	return summaryResult, behaviorsResult, nil
}

func ObtainAllureResult(url, jobName string) (string, string, error) {
	_ctx, _cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer _cancel()

	ctx, cancel := chromedp.NewContext(_ctx)
	defer cancel()

	// capture screenshot of an element
	var summaryResult string
	var behaviorsResult string
	if err := chromedp.Run(ctx, elementHTML(url, SUMMARY_AREA_SELECTOR, SUMMARY_INFO_SELECTOR, &summaryResult)); err != nil {
		return "", "", err
	}

	if err := chromedp.Run(ctx, elementHTML(url, BEHAVIORS_AREA_SELECTOR, BEHAVIORS_INFO_SELECTOR, &behaviorsResult)); err != nil {
		return "", "", err
	}

	return summaryResult, behaviorsResult, nil
}

func elementHTML(url, sel, sel2 string, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ScrollIntoView(sel, chromedp.ByQuery),
		chromedp.WaitVisible(sel, chromedp.ByQuery),
		chromedp.WaitVisible(sel2, chromedp.ByQuery),
		chromedp.Sleep(1 * time.Second),
		chromedp.OuterHTML(sel, res, chromedp.NodeReady, chromedp.ByQuery),
	}
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel, sel2 string, res *[]byte) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.EmulateViewport(1920, 1280),
		chromedp.Reload(),
		chromedp.Sleep(1 * time.Second),
		chromedp.ScrollIntoView(sel, chromedp.ByQuery),
		chromedp.WaitVisible(sel, chromedp.ByQuery),
		chromedp.WaitVisible(sel2, chromedp.ByQuery),
		chromedp.Sleep(10 * time.Second),
		chromedp.Screenshot(sel, res, chromedp.NodeReady, chromedp.ByQuery),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
