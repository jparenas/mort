package object

import (
	"github.com/aldor007/mort/pkg/config"
	"github.com/aldor007/mort/pkg/transforms"
	"net/url"
	"path"
	"strconv"
)

func init() {
	RegisterParser("query", decodeQuery)
}

func decodeQuery(url *url.URL, bucketConfig config.Bucket, obj *FileObject) (string, error) {
	trans := bucketConfig.Transform

	var err error
	obj.Transforms, err = queryToTransform(url.Query())

	if obj.HasTransform() {
		parent := url.Path
		if trans.ParentBucket != "" {
			parent = "/" + path.Join(trans.ParentBucket, obj.Key)
		}

		obj.Key = hashKey(obj.Transforms.Hash(), parent)
		return parent, err
	}

	return "", err
}

func queryToTransform(query url.Values) (transforms.Transforms, error) {
	var trans transforms.Transforms
	if len(query) == 0 {
		return trans, nil
	}

	var err error
	opt := query.Get("operation")
	if opt == "" {
		var w, h int
		w, _ = queryToInt(query, "width")
		h, _ = queryToInt(query, "height")
		err = trans.Resize(w, h, false)
	} else {
		for qsKey, values := range query {
			if qsKey == "operation" {
				for _, o := range values {
					switch o {
					case "resize":
						var w, h int
						w, _ = queryToInt(query, "width")
						h, _ = queryToInt(query, "height")

						err = trans.Resize(w, h, false)
					case "crop":
						var w, h int
						w, err = queryToInt(query, "width")
						h, err = queryToInt(query, "height")

						err = trans.Crop(w, h, query.Get("gravity"), false)
					case "watermark":
						var opacity float64
						opacity, err = strconv.ParseFloat(query.Get("opacity"), 32)
						if err != nil {
							return trans, err
						}
						err = trans.Watermark(query.Get("image"), query.Get("position"), float32(opacity))
					case "blur":
						var sigma, minAmpl float64
						sigma, err = strconv.ParseFloat(query.Get("sigma"), 32)
						if err != nil {
							return trans, err
						}

						minAmpl, _ = strconv.ParseFloat(query.Get("minAmpl"), 32)
						err = trans.Blur(sigma, minAmpl)
					case "rotate":
						var a int
						a, err = queryToInt(query, "angle")
						if err != nil {
							return trans, err
						}
						err = trans.Rotate(a)
					}

				}
				break
			}
		}
	}

	var q int
	q, err = queryToInt(query, "quality")
	err = trans.Quality(q)
	if format, ok := query["format"]; ok {
		err = trans.Format(format[0])
	}

	if _, ok := query["grayscale"]; ok {
		trans.Grayscale()
	}

	return trans, err
}

func queryToInt(q url.Values, k string) (int, error) {
	r, err := strconv.ParseInt(q.Get(k), 10, 32)
	return int(r), err

}
