package jenkinsallure

import (
	"context"
	"log"
	"time"

	"math"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func CaptureAllureResult(url, jobName string) ([]byte, []byte, []byte, error) {
	log.Println(url)
	_ctx, _cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer _cancel()

	ctx, cancel := chromedp.NewContext(_ctx)
	defer cancel()

	// capture screenshot of an element
	var summaryResult []byte
	var behaviorsResult []byte
	var trendResult []byte
	summaryArea := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(1)`
	summaryInfo := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(1) > div.widget__body > div > div > div:nth-child(1) > div > div.splash__subtitle`
	if err := chromedp.Run(ctx, elementScreenshot(url, summaryArea, summaryInfo, &summaryResult)); err != nil {
		return nil, nil, nil, err
	}

	behaviorsArea := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(4)`
	behaviorsInfo := `#content > div > div.app__content > div > div:nth-child(1) > div:nth-child(4) > div.widget__body > div > h2 > span`
	if err := chromedp.Run(ctx, elementScreenshot(url, behaviorsArea, behaviorsInfo, &behaviorsResult)); err != nil {
		return nil, nil, nil, err
	}

	trendArea := `#content > div > div.app__content > div > div:nth-child(2) > div:nth-child(1)`
	trendInfo := `#content > div > div.app__content > div > div:nth-child(2) > div:nth-child(1) > div.widget__body > div > div`
	if err := chromedp.Run(ctx, elementScreenshot(url, trendArea, trendInfo, &trendResult)); err != nil {
		return nil, nil, nil, err
	}

	return summaryResult, behaviorsResult, trendResult, nil
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel, sel2 string, res *[]byte) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ScrollIntoView(sel, chromedp.ByQuery),
		chromedp.WaitVisible(sel, chromedp.ByQuery),
		chromedp.WaitReady(sel, chromedp.ByQuery),
		chromedp.WaitVisible(sel2, chromedp.ByQuery),
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