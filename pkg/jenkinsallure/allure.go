package jenkinsallure

import (
	"context"
	"time"

	"math"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func CaptureAllureResult(url, jobName string) ([]byte, []byte, error) {
	_ctx, _cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer _cancel()

	ctx, cancel := chromedp.NewContext(_ctx)
	defer cancel()

	// capture screenshot of an element
	var summaryResult []byte
	var behaviorsResult []byte
	summaryArea := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(1)`
	summaryInfo := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(1) > div.widget__body > div > div > div.widget__column.summary-widget__chart > div > svg`
	if err := chromedp.Run(ctx, elementScreenshot(url, summaryArea, summaryInfo, &summaryResult)); err != nil {
		return nil, nil, err
	}

	behaviorsArea := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(4)`
	behaviorsInfo := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(4) > div.widget__body > div > h2 > span`
	if err := chromedp.Run(ctx, elementScreenshot(url, behaviorsArea, behaviorsInfo, &behaviorsResult)); err != nil {
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
	summaryArea := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(1)`
	summaryInfo := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(1) > div.widget__body > div > div > div.widget__column.summary-widget__chart > div > svg`
	if err := chromedp.Run(ctx, elementHTML(url, summaryArea, summaryInfo, &summaryResult)); err != nil {
		return "", "", err
	}

	behaviorsArea := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(4)`
	behaviorsInfo := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(4) > div.widget__body > div > h2 > span`
	if err := chromedp.Run(ctx, elementHTML(url, behaviorsArea, behaviorsInfo, &behaviorsResult)); err != nil {
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
		chromedp.ScrollIntoView(sel, chromedp.ByQuery),
		chromedp.WaitVisible(sel, chromedp.ByQuery),
		chromedp.WaitVisible(sel2, chromedp.ByQuery),
		chromedp.Sleep(3 * time.Second),
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
