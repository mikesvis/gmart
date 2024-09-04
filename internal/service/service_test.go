package service_test

import (
	"context"
	"net/http/httptest"
	"time"

	"github.com/mikesvis/gmart/internal/domain"
	"github.com/mikesvis/gmart/internal/randomstring"
	"github.com/mikesvis/gmart/internal/service/order"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	Context("When new user registers", func() {
		It("returns user ID", func() {
			userID, err := UserService.RegisterUser(context.Background(), randomstring.RandStringRunes(10), "password")
			Expect(userID).To(BeNumerically("~", 0, 9999))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns error for existing user", func() {
			userID, err := UserService.RegisterUser(context.Background(), ExistingUser.Login, ExistingUser.Password)
			Expect(userID).To(BeZero())
			Expect(err).Should(HaveOccurred())
		})

		It("returns error for empty login or password or both", func() {
			userID, err := UserService.RegisterUser(context.Background(), "", "")
			Expect(userID).To(BeZero())
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("When user logging in", func() {
		It("returns nil", func() {
			w := httptest.NewRecorder()
			err := UserService.Login(context.Background(), w, ExistingUser.ID)
			Expect(err).ShouldNot(HaveOccurred())
			header := w.Header().Get("Authorization")
			Expect(header).ShouldNot(BeEmpty())
		})
	})

	Context("When getting user ID by login and password", func() {
		It("returns user id for existing user", func() {
			userID, err := UserService.GetUserID(context.Background(), ExistingUser.Login, ExistingUser.Password)
			Expect(userID).Should(Equal(ExistingUser.ID))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns error for empty login or password or both", func() {
			userID, err := UserService.GetUserID(context.Background(), "", "")
			Expect(userID).To(BeZero())
			Expect(err).Should(HaveOccurred())
		})

		It("returns error for user ID is zero", func() {
			userID, err := UserService.GetUserID(context.Background(), "", "")
			Expect(userID).To(BeZero())
			Expect(err).Should(HaveOccurred())
		})
	})

})

var _ = Describe("Order", func() {
	Context("When fetching user orders", func() {
		It("returns empty if user has no orders", func() {
			orders, err := OrderService.GetOrdersByUser(context.Background(), EmptyOrdersExistingUserID)
			Expect(orders).To(BeNil())
			Expect(err).To(MatchError(order.ErrNoOrdersInSet))
		})

		It("lists user orders if any", func() {
			orders, err := OrderService.GetOrdersByUser(context.Background(), NotEmptyOrdersExistingUserID)
			Expect(orders).ShouldNot(BeEmpty())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

var _ = Describe("Accrual", func() {
	Context("When fetching user balance", func() {
		It("returns zero balance if user has no orders", func() {
			balance, err := AccrualService.GetUserBalance(context.Background(), EmptyOrdersExistingUserID)
			Expect(balance).To(BeEquivalentTo(&domain.UserBalance{Current: 0, Withdrawn: 0}))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns balance and withdrawns if user has orders", func() {
			balance, err := AccrualService.GetUserBalance(context.Background(), NotEmptyOrdersExistingUserID)
			Expect(balance).To(BeEquivalentTo(&domain.UserBalance{Current: 60000, Withdrawn: 10000}))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns error for user ID is zero", func() {
			balance, err := AccrualService.GetUserBalance(context.Background(), 0)
			Expect(balance).To(BeNil())
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("When fetching user withdrawals", func() {
		It("returns error of empty list if user has no withdrawals", func() {
			withdrawals, err := AccrualService.GetUserWithdrawals(context.Background(), EmptyOrdersExistingUserID)
			Expect(withdrawals).To(BeNil())
			Expect(err).Should(HaveOccurred())
		})

		It("returns withdrawals for existing user withdrawals", func() {
			withdrawals, err := AccrualService.GetUserWithdrawals(context.Background(), NotEmptyOrdersExistingUserID)
			Expect(withdrawals).NotTo(BeNil())
			for _, withdrawal := range withdrawals {
				Expect(withdrawal.OrderID).To(Equal(ExistingOrderNumberWithWithdrawal))
				Expect(withdrawal.Sum).To(Equal(uint64(10000)))
				Expect(withdrawal.ProcessedAt).Should(BeTemporally("~", time.Now(), 5*time.Second))
			}
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns error for user ID is zero", func() {
			withdrawals, err := AccrualService.GetUserWithdrawals(context.Background(), 0)
			Expect(withdrawals).To(BeNil())
			Expect(err).Should(HaveOccurred())
		})
	})
})
