// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils_test

import (
	mockutils "github.com/gardener/gardener/pkg/mock/gardener/utils"
	. "github.com/gardener/gardener/pkg/utils"
	"github.com/gardener/gardener/test"
	"github.com/golang/mock/gomock"

	"errors"
	"time"

	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	errwrap "github.com/pkg/errors"
)

var (
	testErr = errors.New("severe")
)

var _ = Describe("utils", func() {

	var (
		ctrl         *gomock.Controller
		canceledCtx  context.Context
		zeroDuration = time.Duration(0)
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		canceledCtx = func() context.Context {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			return ctx
		}()
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Context("#RetryUntil", func() {
		closedChan := make(chan struct{})
		close(closedChan)

		It("should abort immediately on a severe error and return it", func() {
			ct := 0
			err := RetryUntil(context.Background(), zeroDuration, func() (ok, severe bool, err error) {
				if ct > 0 {
					Fail("Function called multiple times although should have already failed")
				}
				ct++
				return false, true, testErr
			})

			Expect(err).To(Equal(testErr))
		})

		It("should not error if the function exits cleanly", func() {
			err := RetryUntil(context.Background(), 0*time.Second, func() (ok, severe bool, err error) {
				return true, false, nil
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not timeout early and use the value of the delayed function", func() {
			err := RetryUntil(canceledCtx, zeroDuration, func() (ok, severe bool, err error) {
				return true, false, nil
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should exit with a timeout after the interval sleep due to a closed channel", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ct := 0
			err := RetryUntil(ctx, zeroDuration, func() (ok, severe bool, err error) {
				ct++
				cancel()
				return false, false, nil
			})

			Expect(err).To(HaveOccurred())
			Expect(IsTimedOut(err)).To(BeTrue())
		})

		It("should fail due to a timeout containing the last error", func() {
			err := RetryUntil(canceledCtx, zeroDuration, func() (ok, severe bool, err error) {
				return false, false, testErr
			})

			Expect(err).To(HaveOccurred())
			Expect(IsTimedOut(err)).To(BeTrue())
			Expect(LastErrorOfTimedOutWithError(err)).To(Equal(testErr))
			Expect(errwrap.Cause(err)).To(Equal(testErr))
		})

		It("should fail due to a timeout containing no last error", func() {
			err := RetryUntil(canceledCtx, zeroDuration, func() (ok, severe bool, err error) {
				return false, false, nil
			})

			Expect(err).To(HaveOccurred())
			Expect(IsTimedOut(err)).To(BeTrue())
			Expect(errwrap.Cause(err)).To(Equal(err))
		})

		It("should fail with a wait time greater or equal to the sleep time of a function", func() {
			leastWaitTime := 10 * time.Millisecond
			err := Retry(0*time.Second, 0*time.Second, func() (ok, severe bool, err error) {
				time.Sleep(leastWaitTime)
				return false, false, nil
			})

			Expect(err).To(HaveOccurred())
			Expect(IsTimedOut(err)).To(BeTrue())
			Expect(WaitTimeOfTimedOut(err)).To(BeNumerically(">=", leastWaitTime))
		})
	})

	Context("#Retry", func() {
		It("should fail due to a timeout", func() {
			err := Retry(0*time.Second, 0*time.Second, func() (ok, severe bool, err error) {
				return false, false, nil
			})

			Expect(err).To(HaveOccurred())
			Expect(IsTimedOut(err)).To(BeTrue())
		})

		It("should not fail due to a timeout", func() {
			err := Retry(0*time.Second, 10*time.Second, func() (ok, severe bool, err error) {
				return true, false, nil
			})

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("#Sleep", func() {
		It("should sleep for the given duration", func() {
			var (
				duration = 10 * time.Second
				newTimer = mockutils.NewMockNewTimer(ctrl)
				timer    = mockutils.NewMockTimer(ctrl)
				c        = make(chan time.Time, 1)
			)
			c <- time.Now()
			defer close(c)

			defer test.RevertableSet(&NewTimer, newTimer.Do)()

			gomock.InOrder(
				newTimer.EXPECT().Do(duration).Return(timer),
				timer.EXPECT().C().Return(c),
				timer.EXPECT().Stop(),
			)

			Expect(Sleep(context.Background(), duration)).To(Succeed())
		})

		It("should return the context's error if duration is 0 but the context is expired", func() {
			err := Sleep(canceledCtx, 0*time.Second)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(canceledCtx.Err()))
		})

		It("should error with the context's error if the channel does not fire", func() {
			var (
				duration = 10 * time.Second
				newTimer = mockutils.NewMockNewTimer(ctrl)
				timer    = mockutils.NewMockTimer(ctrl)
				c        = make(chan time.Time)
			)
			defer close(c)

			defer test.RevertableSet(&NewTimer, newTimer.Do)()

			gomock.InOrder(
				newTimer.EXPECT().Do(duration).Return(timer),
				timer.EXPECT().C().Return(c),
				timer.EXPECT().Stop(),
			)

			err := Sleep(canceledCtx, duration)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(canceledCtx.Err()))
		})
	})
})
