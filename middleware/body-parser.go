package middleware

// func BodyParser_legacy(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ln := r.ContentLength
// 		body := make([]byte, ln)
// 		r.Body.Read(body)
// 		obj, _ := rnjson.Unmarshal(string(body))
// 		f(w, r.WithContext(context.WithValue(r.Context(), "rnbody", obj)))
// 	}
// }

// func BodyParser(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		f(w, r.WithContext(context.WithValue(r.Context(), "rnparser", parser)))
// 	}
// }
