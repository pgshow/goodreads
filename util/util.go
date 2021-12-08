package util

import (
	"fmt"
	"github.com/avast/retry-go/v3"
	"github.com/parnurzeal/gorequest"
	"time"
)

// RetriableError is a custom error that contains a positive duration for the next retry
type RetriableError struct {
	Err        error
	RetryAfter time.Duration
}

// Error returns error message and a Retry-After duration
func (e *RetriableError) Error() string {
	return fmt.Sprintf("%s (retry after %v)", e.Err.Error(), e.RetryAfter)
}

var _ error = (*RetriableError)(nil)

func Fetch(url string) (body string, err error) {
	_ = retry.Do(
		func() error {
			resp, bodyTmp, errs := gorequest.New().Get(url).
				Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.57").
				Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6,zh-TW;q=0.5,ja;q=0.4").
				Set("Cookie", "ccsid=483-6782365-5243135; __qca=P0-1854119153-1637916332186; logged_out_browsing_page_count=2; u=5CtA35AROoKyKzsUz3f-6lQwMSedZK88edsPhHxpBSIYoMST; p=xBoYCVkGNQIpFQzfLOeDzIw2re7839FAYaaSMAzuq7sdWtG0; likely_has_account=true; _session_id2=3298b33fe951996dd9d9b5da8adc2872; locale=en").
				Timeout(45 * time.Second).
				End()

			body = bodyTmp

			// 自动重试只需要返回 error 即可
			if len(errs) > 0 {
				err = errs[0]
			}

			// 自定义错误重试部分，当页面返回404时的错误重试
			if err == nil {
				if resp.StatusCode != 200 {
					err = fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)

					if resp.StatusCode == 404 {
						fmt.Println("The page status is: 404")
						return retry.Unrecoverable(err)
					}

					return &RetriableError{
						Err:        err,
						RetryAfter: time.Duration(20) * time.Second,
					}
				}
			}

			return err
		},

		// 设置重试次数
		retry.Attempts(15),

		// 重试延迟
		retry.Delay(20),

		// 重试策略
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			fmt.Println("Scrape fails with: " + err.Error())
			if retriable, ok := err.(*RetriableError); ok {
				fmt.Printf("Retry after %v\n", retriable.RetryAfter)
				return retriable.RetryAfter
			}
			// apply a default exponential back off strategy
			return retry.BackOffDelay(n, err, config)
		}),
	)

	return
}
