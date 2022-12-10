package internal

type SuccessResponse struct {
	Status string `json:"status" example:"success"`
	Data   Data   `json:"data"`
}

type CheckHealthResponse struct {
	Status string `json:"status" example:"operational"`
}

type MissingResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"No wikipedia article found."`
	Missing bool   `json:"missing" example:"true"`
}

type NoDescriptionResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"No short description found for this article."`
	Missing bool   `json:"missing" example:"false"`
}

type ErrorResponse struct {
	Status string      `json:"status" example:"error"`
	Errors []HTTPError `json:"errors"`
}

type InternalServerErrorResponse struct {
	Status string                `json:"status" example:"error"`
	Errors []InternalServerError `json:"errors"`
}

type WikipediaApiErrorResponse struct {
	Status string              `json:"status" example:"error"`
	Errors []WikipediaApiError `json:"errors"`
}

type Data struct {
	ShortDescription string `json:"short_description" example:"A short description of the person, place, or thing you searched for."`
}

type HTTPError struct {
	Code      int    `json:"code" example:"400"`
	RequestID string `json:"request_id" example:"f7a4c0c0-5b5e-4b4c-9c1f-1b5c1b5c1b5c"`
	Detail    string `json:"detail" example:"Query parameter is required"`
}

type WikipediaApiError struct {
	Code      int    `json:"code" example:"500"`
	RequestID string `json:"request_id" example:"f7a4c0c0-5b5e-4b4c-9c1f-1b5c1b5c1b5c"`
	Detail    string `json:"detail" example:"An error occurred while communicating with the wikipedia API with http code 500. Please find more information at https://en.wikipedia.org/w/api.php."`
}

type InternalServerError struct {
	Code      int    `json:"code" example:"500"`
	RequestID string `json:"request_id" example:"f7a4c0c0-5b5e-4b4c-9c1f-1b5c1b5c1b5c"`
	Detail    string `json:"detail" example:"An internal server error occurred. Please contact the developer at youssefsobhy22@gmail.com and provide the request ID."`
}

type WikipediaResponse struct {
	Query Query `json:"query"`
}

type Query struct {
	Pages []Page `json:"pages"`
}

type Page struct {
	PageID    int        `json:"pageid"`
	Ns        int        `json:"ns"`
	Title     string     `json:"title"`
	Revisions []Revision `json:"revisions"`
	Missing   bool       `json:"missing"`
}

type Revision struct {
	Content string `json:"content"`
}
