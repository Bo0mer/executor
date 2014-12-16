package utils_test

import (
	"io/ioutil"

	. "github.com/Bo0mer/executor/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sniffer", func() {

	Describe("SniffedData", func() {
		var s Sniffer
		var data []byte

		BeforeEach(func() {
			s = NewSniffer(ioutil.Discard)
			data = []byte("sniffed")
			s.Write(data)
		})

		It("should return the writer data", func() {
			Expect(s.SniffedData()).To(Equal(data))
		})
	})
})
