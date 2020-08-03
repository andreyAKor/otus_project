package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreyAKor/otus_project/internal/cache"
	cacheMocks "github.com/andreyAKor/otus_project/internal/cache/mocks"
	clientMocks "github.com/andreyAKor/otus_project/internal/http/client/mocks"
	imageMocks "github.com/andreyAKor/otus_project/internal/image/mocks"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	emptyURI     = "/"
	normalURI    = "/2000/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg"
	normalSource = "www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg"
	imageContent = "some image content"
)

//nolint:funlen
func TestProcess(t *testing.T) {
	t.Run("request uri parsing fail", func(t *testing.T) {
		t.Run("empty uri", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := getEmptyCache(ctrl, emptyURI)

			srv, err := New(nil, nil, c, "", 0, 0)
			require.NoError(t, err)

			req := httptest.NewRequest("GET", emptyURI, nil)
			w := httptest.NewRecorder()

			_, err = srv.process(w, req)
			require.Error(t, err)

			rsp := w.Result()
			defer rsp.Body.Close()
			require.Equal(t, http.StatusBadRequest, rsp.StatusCode)
		})
	})
	t.Run("errors", func(t *testing.T) {
		t.Run("getting image from cache fail", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cacheError := errors.New("get some cache error")

			c := cacheMocks.NewMockCache(ctrl)
			c.EXPECT().
				Get(cache.Key(normalURI)).
				Return(nil, false, cacheError)

			srv, err := New(nil, nil, c, "", 0, 0)
			require.NoError(t, err)

			req := httptest.NewRequest("GET", normalURI, nil)
			w := httptest.NewRecorder()

			_, err = srv.process(w, req)
			require.Equal(t, errors.Unwrap(errors.Unwrap(err)), errors.Cause(cacheError))

			rsp := w.Result()
			defer rsp.Body.Close()
			require.Equal(t, http.StatusInternalServerError, rsp.StatusCode)
		})
		t.Run("client request fail", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			clientError := errors.New("get some client error")

			req := httptest.NewRequest("GET", normalURI, nil)
			w := httptest.NewRecorder()

			cl := clientMocks.NewMockClient(ctrl)
			cl.EXPECT().
				Request(normalSource, req).
				Return(nil, nil, clientError)

			c := getEmptyCache(ctrl, normalURI)

			srv, err := New(cl, nil, c, "", 0, 0)
			require.NoError(t, err)

			_, err = srv.process(w, req)
			require.Equal(t, errors.Unwrap(errors.Unwrap(errors.Unwrap(errors.Unwrap(err)))), errors.Cause(clientError))

			rsp := w.Result()
			defer rsp.Body.Close()
			require.Equal(t, http.StatusInternalServerError, rsp.StatusCode)
		})
		t.Run("bad gateway", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req := httptest.NewRequest("GET", normalURI, nil)
			w := httptest.NewRecorder()

			cl := clientMocks.NewMockClient(ctrl)
			cl.EXPECT().
				Request(normalSource, req).
				Return(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil, nil)

			c := getEmptyCache(ctrl, normalURI)

			srv, err := New(cl, nil, c, "", 0, 0)
			require.NoError(t, err)

			_, err = srv.process(w, req)
			require.Equal(t, errors.Unwrap(errors.Unwrap(err)), ErrBadGateway)

			rsp := w.Result()
			defer rsp.Body.Close()
			require.Equal(t, http.StatusBadGateway, rsp.StatusCode)
		})
		t.Run("image runner fail", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			imageError := errors.New("get some image error")
			content := []byte(imageContent)

			req := httptest.NewRequest("GET", normalURI, nil)
			w := httptest.NewRecorder()

			cl := getClientRequestOK(ctrl, normalSource, req)

			i := imageMocks.NewMockImage(ctrl)
			i.EXPECT().
				Resize(content, 2000, 200).
				Return(nil, imageError)

			c := getEmptyCache(ctrl, normalURI)

			srv, err := New(cl, i, c, "", 0, 0)
			require.NoError(t, err)

			_, err = srv.process(w, req)
			require.Equal(t, errors.Unwrap(errors.Unwrap(errors.Unwrap(errors.Unwrap(err)))), errors.Cause(imageError))

			rsp := w.Result()
			defer rsp.Body.Close()
			require.Equal(t, http.StatusInternalServerError, rsp.StatusCode)
		})
		t.Run("setting image to cache fail", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cacheError := errors.New("set some cache error")
			content := []byte(imageContent)

			req := httptest.NewRequest("GET", normalURI, nil)
			w := httptest.NewRecorder()

			cl := getClientRequestOK(ctrl, normalSource, req)
			i := getNormalImage(ctrl, 2000, 200)
			c := getEmptyCache(ctrl, normalURI)
			c.EXPECT().
				Set(cache.Key(normalURI), content).
				Return(cacheError)

			srv, err := New(cl, i, c, "", 0, 0)
			require.NoError(t, err)

			_, err = srv.process(w, req)
			require.Equal(t, errors.Unwrap(errors.Unwrap(errors.Unwrap(errors.Unwrap(err)))), errors.Cause(cacheError))

			rsp := w.Result()
			defer rsp.Body.Close()
			require.Equal(t, http.StatusInternalServerError, rsp.StatusCode)
		})
	})
	t.Run("normal", func(t *testing.T) {
		t.Run("from cache", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			content := []byte(imageContent)

			req := httptest.NewRequest("GET", normalURI, nil)
			w := httptest.NewRecorder()

			c := cacheMocks.NewMockCache(ctrl)
			c.EXPECT().
				Get(cache.Key(normalURI)).
				Return(content, true, nil)

			srv, err := New(nil, nil, c, "", 0, 0)
			require.NoError(t, err)

			res, err := srv.process(w, req)
			require.NoError(t, err)
			require.Equal(t, res, content)

			rsp := w.Result()
			defer rsp.Body.Close()
			require.Equal(t, http.StatusOK, rsp.StatusCode)
		})
		t.Run("from source", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			content := []byte(imageContent)

			req := httptest.NewRequest("GET", normalURI, nil)
			w := httptest.NewRecorder()

			cl := getClientRequestOK(ctrl, normalSource, req)
			i := getNormalImage(ctrl, 2000, 200)
			c := getEmptyCache(ctrl, normalURI)
			c.EXPECT().
				Set(cache.Key(normalURI), content).
				Return(nil)

			srv, err := New(cl, i, c, "", 0, 0)
			require.NoError(t, err)

			res, err := srv.process(w, req)
			require.NoError(t, err)
			require.Equal(t, res, content)

			rsp := w.Result()
			defer rsp.Body.Close()
			require.Equal(t, http.StatusOK, rsp.StatusCode)
		})
	})
}

func getEmptyCache(ctrl *gomock.Controller, uri string) *cacheMocks.MockCache {
	c := cacheMocks.NewMockCache(ctrl)
	c.EXPECT().
		Get(cache.Key(uri)).
		Return(nil, false, nil)

	return c
}

func getClientRequestOK(ctrl *gomock.Controller, src string, req *http.Request) *clientMocks.MockClient {
	content := []byte(imageContent)

	cl := clientMocks.NewMockClient(ctrl)
	cl.EXPECT().
		Request(src, req).
		Return(&http.Response{
			StatusCode: http.StatusOK,
		}, content, nil)

	return cl
}

func getNormalImage(ctrl *gomock.Controller, width, height int) *imageMocks.MockImage {
	content := []byte(imageContent)

	i := imageMocks.NewMockImage(ctrl)
	i.EXPECT().
		Resize(content, width, height).
		Return(content, nil)

	return i
}
