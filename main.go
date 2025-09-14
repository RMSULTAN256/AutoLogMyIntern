package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func GetElements(ctx context.Context, sel string) bool {
	ctxShort, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return chromedp.Run(ctxShort, chromedp.WaitVisible(sel, chromedp.ByID)) == nil
}

func GetElementsQuery(ctx context.Context, sel string) bool {
	ctxShort, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return chromedp.Run(ctxShort, chromedp.WaitVisible(sel, chromedp.ByQuery)) == nil
}

func main() {
	//user := "dsadwwe"
	//pass := "wfasd"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("window-size", "1200,800"),
		chromedp.Flag("window-position", "200 500"),
	)

	allowctx, cancelAllow := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAllow()

	ctx, cancelCtx := chromedp.NewContext(allowctx)
	defer cancelCtx()

	var currenturl string

	const (
		UserQ      = `#nim`
		PassQ      = `#password`
		Submit     = `#sign_in_button`
		CheckInOut = `button[class="btn btn-info btn-sm btn-block"]`
		attendance = `#attendance_type`
		iframeCSS  = `iframe.cke_wysiwyg_frame[title^="Rich Text Editor"]`
		jsBody     = `document.querySelector(` + "`" + `iframe.cke_wysiwyg_frame[title^="Rich Text Editor"]` + "`" + `).contentDocument.body`
		text       = `Memperbaikki dan menambahkan deteksi dimana apakah akun sudah di dashboard walaupun belum otp dan memperbaikki beberapa logic code yang tertera`
	)

	wait := 2 * time.Second
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://myinternship.id/index.php?page=student_login"),
		chromedp.Sleep(wait),
		chromedp.Location(&currenturl),

		chromedp.ActionFunc(func(ctx context.Context) error {
			switch {
			case GetElements(ctx, UserQ):
				if err := chromedp.SetValue(UserQ, "sultan.4332201024@students.polibatam.ac.id", chromedp.ByQuery).Do(ctx); err != nil {
					return err
				}
				if err := chromedp.SetValue(PassQ, "Sultan1200", chromedp.ByQuery).Do(ctx); err != nil {
					return err
				}
				chromedp.Sleep(wait).Do(ctx)
				chromedp.Click(Submit, chromedp.ByQuery).Do(ctx)
				return nil

			default:
				chromedp.Sleep(wait).Do(ctx)
				log.Println("Url down")
				return nil
			}
		}),

		chromedp.WaitReady("body", chromedp.ByQuery),

		chromedp.ActionFunc(func(ctx context.Context) error {
			switch {
			case GetElementsQuery(ctx, CheckInOut):
				if err := chromedp.Click(CheckInOut, chromedp.ByQuery).Do(ctx); err != nil {
					log.Printf("There an error:%v", err)
					return err
				}

				switch {
				case GetElements(ctx, attendance):
					chromedp.WaitVisible(`#attendance_type`, chromedp.ByID).Do(ctx)
					chromedp.WaitEnabled(`#attendance_type`, chromedp.ByID).Do(ctx)
					chromedp.SetValue(`#attendance_type`, "Present", chromedp.ByID).Do(ctx)
					chromedp.Evaluate(`(function(){
					const el = document.getElementById('attendance_type');
					el.dispatchEvent(new Event('input', {bubbles:true}));
					el.dispatchEvent(new Event('change', {bubbles:true}));
					})()`, nil).Do(ctx)
					chromedp.Click(`#submit_btn`, chromedp.ByID).Do(ctx)
					chromedp.WaitReady(`body`, chromedp.ByQuery).Do(ctx)
					chromedp.Location(&currenturl).Do(ctx)
					return nil

				case GetElementsQuery(ctx, iframeCSS):
					chromedp.WaitReady("body", chromedp.ByQuery).Do(ctx)
					chromedp.WaitVisible(iframeCSS, chromedp.ByQuery).Do(ctx)
					chromedp.Focus(jsBody, chromedp.ByJSPath).Do(ctx)
					chromedp.Sleep(wait).Do(ctx)
					chromedp.SendKeys(jsBody, text, chromedp.ByJSPath).Do(ctx)
					chromedp.Sleep(wait)
					chromedp.Click(`#submit_btn`, chromedp.ByID)
					return nil
				}
			default:
				chromedp.Location(&currenturl).Do(ctx)
				chromedp.Stop().Do(ctx)
				log.Fatalf("Nothing")
				return nil
			}
			return nil
		}),
		chromedp.Location(&currenturl),
		chromedp.Sleep(wait),
		chromedp.Sleep(wait),
	); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Your Current url:", currenturl)
}
