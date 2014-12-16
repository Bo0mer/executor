package utils_test

import (
	"errors"

	. "github.com/Bo0mer/executor/utils"
	fakes "github.com/Bo0mer/executor/utils/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sniffer", func() {
	var s Sniffer
	var fakeWriter *fakes.FakeWriter
	var data []byte
	var n int
	var writeError error

	BeforeEach(func() {
		data = []byte("some_data_here")
		fakeWriter = new(fakes.FakeWriter)
		s = NewSniffer(fakeWriter)
	})

	JustBeforeEach(func() {
		n, writeError = s.Write(data)
	})

	Context("When writing succeeds", func() {
		BeforeEach(func() {
			fakeWriter.WriteReturns(len(data), nil)
		})

		It("should have returned correct result", func() {
			Expect(n).To(Equal(len(data)))
			Expect(writeError).ToNot(HaveOccurred())
		})

		It("should have wrote all the data to the Writer", func() {
			Expect(fakeWriter.WriteCallCount()).To(Equal(1))
			Expect(fakeWriter.WriteArgsForCall(0)).To(Equal(data))
		})

		It("should have the data sniffed", func() {
			Expect(s.SniffedData()).To(Equal(data))
		})
	})

	Context("when writing partially succeeds", func() {
		var err error

		BeforeEach(func() {
			err = errors.New("Wrote only 2 bytes")
			fakeWriter.WriteReturns(2, err)
		})

		It("should return the correct results", func() {
			Expect(n).To(Equal(2))
			Expect(writeError).To(HaveOccurred())
			Expect(writeError).To(Equal(err))
		})

		It("should have sniffed only the writed data", func() {
			Expect(s.SniffedData()).To(Equal(data[:2]))
		})
	})

	Context("when writing fails", func() {
		var err error

		BeforeEach(func() {
			err = errors.New("write failed")
			fakeWriter.WriteReturns(0, err)
		})

		It("should return the correct results", func() {
			Expect(n).To(Equal(0))
			Expect(writeError).To(HaveOccurred())
			Expect(writeError).To(Equal(err))
		})

		It("should have no sniffed data", func() {
			Expect(s.SniffedData()).To(BeEmpty())
		})
	})
})
