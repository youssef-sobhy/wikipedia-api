package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/youssef1337/wikipedia-api/internal"
)

var _ = BeforeSuite(func() {
	// block all HTTP requests
	httpmock.Activate()
})

var _ = AfterSuite(func() {
	// unblock all HTTP requests
	httpmock.DeactivateAndReset()
})

var _ = Describe("WikipediaApi", func() {
	BeforeEach(func() {
		// remove any mocks
		httpmock.Reset()
	})

	Describe("/health", func() {
		Context("when the server is up", func() {
			It("should return 200 and status operational", func() {
				http.NewRequest("GET", "/api", nil)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				internal.Health(c)
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				defer w.Result().Body.Close()

				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(response["status"]).To(Equal("operational"))
			})
		})
	})

	Describe("/search", func() {
		Context("when the query parameter is missing", func() {
			It("should return 400 and a 'Query parameter is required.' message", func() {
				http.NewRequest("GET", "/api/search", nil)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				internal.Search(c)
				var response internal.ErrorResponse
				json.Unmarshal(w.Body.Bytes(), &response)

				defer w.Result().Body.Close()

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(response.Errors[0].Detail).To(Equal("Query parameter is required."))
			})
		})

		Context("when the query parameter is present", func() {
			Context("when the Wikipedia API returns the result we are looking for", func() {
				It("should return 200 and the short description", func() {
					httpmock.RegisterResponder(
						"GET",
						"https://en.wikipedia.org/w/api.php?action=query&prop=revisions&titles=Yoshua_Bengio&rvlimit=1&formatversion=2&format=json&rvprop=content",

						httpmock.NewStringResponder(
							200,
							`{
								"query": {
									"pages": [
										{
											"pageid": 47749536,
											"ns": 0,
											"title": "Yoshua Bengio",
											"revisions": [
												{
													"contentformat": "text/x-wiki",
													"contentmodel": "wikitext",
													"content": "{{Short description|Canadian computer scientist}}\n{{Use mdy dates|date=March 2019}}"
												}
											]
										}
									]
								}
							}`,
						),
					)

					req, _ := http.NewRequest("GET", "/api/search?query=Yoshua_Bengio", nil)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = req
					internal.Search(c)
					var response internal.SuccessResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					defer w.Result().Body.Close()

					Expect(w.Code).To(Equal(http.StatusOK))
					Expect(response.Data.ShortDescription).To(Equal("Canadian computer scientist"))
				})
			})

			Context("when the Wikipedia API returns that the page is missing", func() {
				It("should return 200 and a 'No wikipedia article found.' message", func() {
					httpmock.RegisterResponder(
						"GET",
						"https://en.wikipedia.org/w/api.php?action=query&prop=revisions&titles=Yoshua_Bengio~&rvlimit=1&formatversion=2&format=json&rvprop=content",

						httpmock.NewStringResponder(
							200,
							`{
								"query": {
									"pages": [
										{
											"ns": 0,
											"title": "Yoshua_Bengio~",
											"missing": true
										}
									]
								}
							}`,
						),
					)

					req, _ := http.NewRequest("GET", "/api/search?query=Yoshua_Bengio~", nil)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = req
					internal.Search(c)
					var response internal.MissingResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					defer w.Result().Body.Close()

					Expect(w.Code).To(Equal(http.StatusOK))
					Expect(response.Status).To(Equal("success"))
					Expect(response.Message).To(Equal("No wikipedia article found."))
					Expect(response.Missing).To(Equal(true))
				})
			})

			Context("when the Wikipedia API does not return a short description", func() {
				It("should return 200 and a 'No short description found for this article.' message", func() {
					httpmock.RegisterResponder(
						"GET",
						"https://en.wikipedia.org/w/api.php?action=query&prop=revisions&titles=Kim&rvlimit=1&formatversion=2&format=json&rvprop=content",

						httpmock.NewStringResponder(
							200,
							`{
								"query": {
									"pages": [
										{
											"pageid": 627030,
											"ns": 0,
											"title": "Kim",
											"revisions": [
												{
													"contentformat": "text/x-wiki",
													"contentmodel": "wikitext",
													"content": "{{wiktionary|Kim|kim}}"
												}
											]
										}
									]
								}
							}`,
						),
					)

					req, _ := http.NewRequest("GET", "/api/search?query=Kim", nil)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = req
					internal.Search(c)
					var response internal.NoDescriptionResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					defer w.Result().Body.Close()

					Expect(w.Code).To(Equal(http.StatusOK))
					Expect(response.Status).To(Equal("success"))
					Expect(response.Message).To(Equal("No short description found for this article."))
					Expect(response.Missing).To(Equal(false))
				})
			})

			Context("when the Wikipedia API returns an error", func() {
				It("should return 500 and a 'Wikipedia API error.' message", func() {
					httpmock.RegisterResponder(
						"GET",
						"https://en.wikipedia.org/w/api.php?action=query&prop=revisions&titles=Kim&rvlimit=1&formatversion=2&format=json&rvprop=content",

						httpmock.NewStringResponder(500, `{}`),
					)

					req, _ := http.NewRequest("GET", "/api/search?query=Kim", nil)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = req
					internal.Search(c)
					var response internal.ErrorResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					defer w.Result().Body.Close()

					Expect(w.Code).To(Equal(http.StatusInternalServerError))
					Expect(response.Status).To(Equal("error"))
					Expect(response.Errors[0].Detail).To(Equal("An error occurred while communicating with the wikipedia API with http code 500. Please find more information at https://en.wikipedia.org/w/api.php."))
				})
			})

		})
	})

})

func TestWikipediaApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WikipediaApi Suite")
}
