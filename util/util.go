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
				Set("user-agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0").
				Set("accept-language", "en-US,en;q=0.5").
				Set("Cookie", "ccsid=651-9718854-2247275; __qca=P0-1757457643-1637915367666; __gads=ID=78739481e834c2ac-2212f59a58cf00fb:T=1638529674:S=ALNI_MYvFNydt4A5JUFw8vMLJmz1grBhMw; locale=en; _session_id2=f1a8fde7dff388586ee5b72a8797e4ae; blocking_sign_in_interstitial=true; u=6jlV6VcjB6ZJTb6rgkg-c4mWTLVi_Fq_5PZp1MPqXPQWHZe_; p=xBoYCVkGNQIpFQzfLOeDzIw2re7839FAYaaSMAzuq7sdWtG0; likely_has_account=true").
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
