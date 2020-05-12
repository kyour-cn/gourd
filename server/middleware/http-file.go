package middleware

import (
	"github.com/kyour-cn/gourd/application"
	"github.com/kyour-cn/gourd/common"
	"net/http"
	"os"
	"path/filepath"
)

func FileMiddleware(w http.ResponseWriter, r *http.Request) {

	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//获取App配置
	var config application.Config

	_ = common.ReadConfig("app", &config)

	publicPath := config.Http.Path
	indexPath := config.Http.Index

	// prepend the path with the path to the static directory
	path = filepath.Join(publicPath, r.URL.Path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html

		http.ServeFile(w, r, filepath.Join(path, indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, path)

}

/*
func File(url string) http.Handler {

	//获取App配置
	var config application.Config

	_ = common.ReadConfig("app", &config)

	var dir string

	flag.StringVar(&dir, "dir", config.Http.Path, "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	return http.StripPrefix(url, http.FileServer(http.Dir(dir)))

}

*/
