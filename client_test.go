package bdd_talk_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	bddtalk "github.com/poy/bdd-talk"

	. "github.com/poy/eachers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {

	var (
		client    *bddtalk.Client
		responses chan []byte

		mockServer         *httptest.Server
		mockResponseStorer *mockResponseStorer
	)

	var handler = func(w http.ResponseWriter, req *http.Request) {
		defer GinkgoRecover()
		var resp []byte
		Expect(responses).To(Receive(&resp))

		_, err := w.Write(resp)
		Expect(err).ToNot(HaveOccurred())
	}

	BeforeEach(func() {
		responses = make(chan []byte, 100)
		mockServer = httptest.NewServer(http.HandlerFunc(handler))
		mockResponseStorer = newMockResponseStorer()

		client = bddtalk.NewClient(mockServer.URL, mockResponseStorer)
	})

	Describe("HitServer()", func() {

		Context("server is available", func() {
			JustBeforeEach(func() {
				responses <- nil
				close(mockResponseStorer.StoreOutput.ret0)
			})

			Context("ResponseServer does not return an error", func() {

				It("hits server", func() {
					Expect(client.HitServer()).To(Succeed())
				})

				It("stores the response with the correct index", func() {
					client.HitServer()

					Expect(mockResponseStorer.StoreInput.index).To(Receive(Equal(1)))
				})

				Describe("burst of hits", func() {

					BeforeEach(func() {
						for i := 0; i < 5-1; i++ {
							responses <- nil
						}
					})

					It("stores a burst of hits", func() {
						for i := 0; i < 5; i++ {
							client.HitServer()
						}

						Expect(mockResponseStorer.StoreInput.index).To(EqualEach(1, 2, 3, 4, 5))
					})
				})
			})

			Context("ResponseStorer returns an error", func() {
				var (
					expectedErr error
				)

				BeforeEach(func() {
					expectedErr = fmt.Errorf("some-error")
					mockResponseStorer.StoreOutput.ret0 <- expectedErr
				})

				It("returns an error", func() {
					Expect(client.HitServer()).To(MatchError(expectedErr))
				})
			})
		})

		Context("server is not available", func() {
			BeforeEach(func() {
				mockServer.Close()
			})

			It("returns error", func() {
				Expect(client.HitServer()).ToNot(Succeed())
			})

			It("does not call the ResponseStorer", func() {
				client.HitServer()

				Expect(mockResponseStorer.StoreCalled).ToNot(Receive())
			})
		})
	})

})
