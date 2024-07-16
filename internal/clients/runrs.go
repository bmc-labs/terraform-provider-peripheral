// Manually implemented helpers

package runrs

import "net/http"

func (r *CreateResponse) GetError() *Error {
	var err *Error

	switch r.StatusCode() {
	case http.StatusCreated:
	case http.StatusBadRequest:
		err = r.JSON400
	case http.StatusInternalServerError:
		err = r.JSON500
	default:
		err = &Error{
			ErrType: "unknown",
			Msg:     string(r.Body),
		}
	}

	return err
}

func (r *ReadResponse) GetError() *Error {
	var err *Error

	switch r.StatusCode() {
	case http.StatusOK:
	case http.StatusNotFound:
		err = r.JSON404
	case http.StatusInternalServerError:
		err = r.JSON500
	default:
		err = &Error{
			ErrType: "unknown",
			Msg:     string(r.Body),
		}
	}

	return err
}

func (r *ListResponse) GetError() *Error {
	var err *Error

	switch r.StatusCode() {
	case http.StatusOK:
	case http.StatusNotFound:
		err = r.JSON404
	case http.StatusInternalServerError:
		err = r.JSON500
	default:
		err = &Error{
			ErrType: "unknown",
			Msg:     string(r.Body),
		}
	}

	return err
}

func (r *UpdateResponse) GetError() *Error {
	var err *Error

	switch r.StatusCode() {
	case http.StatusOK:
	case http.StatusNotFound:
		err = r.JSON404
	case http.StatusInternalServerError:
		err = r.JSON500
	default:
		err = &Error{
			ErrType: "unknown",
			Msg:     string(r.Body),
		}
	}

	return err
}

func (r *DeleteResponse) GetError() *Error {
	var err *Error

	switch r.StatusCode() {
	case http.StatusOK:
	case http.StatusNotFound:
		err = r.JSON404
	case http.StatusInternalServerError:
		err = r.JSON500
	default:
		err = &Error{
			ErrType: "unknown",
			Msg:     string(r.Body),
		}
	}

	return err
}
