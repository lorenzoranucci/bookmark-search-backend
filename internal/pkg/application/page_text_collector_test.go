package application

import (
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageTextCollector_CollectText(t *testing.T) {
	type args struct {
		readCloser io.ReadCloser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				readCloser: ioutil.NopCloser(
					strings.NewReader(
						`<!doctype html>
<html>
<head>
    <title>Example Domain</title>

    <meta charset="utf-8" />
    <meta http-equiv="Content-type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style type="text/css">
    body {
        background-color: #f0f0f2;
        margin: 0;
        padding: 0;
        font-family: -apple-system, system-ui, BlinkMacSystemFont, "Segoe UI", "Open Sans", "Helvetica Neue", Helvetica, Arial, sans-serif;
        
    }
    div {
        width: 600px;
        margin: 5em auto;
        padding: 2em;
        background-color: #fdfdff;
        border-radius: 0.5em;
        box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);
    }
    a:link, a:visited {
        color: #38488f;
        text-decoration: none;
    }
    @media (max-width: 700px) {
        div {
            margin: 0 auto;
            width: auto;
        }
    }
    </style>    
</head>

<body>
<div>
    <h1>Example Domain</h1>
    <p>This domain is for use in illustrative examples in documents. You may use this
    domain in literature without prior coordination or asking for permission.</p>
    <p><a href="https://www.iana.org/domains/example">More information...</a></p>
</div>
</body>
</html>
`,
					),
				),
			},
			want: `Example Domain

    Example Domain
    This domain is for use in illustrative examples in documents. You may use this
    domain in literature without prior coordination or asking for permission.
    More information...



`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pbc := BookmarkContentCollector{}
			got, err := pbc.CollectText(tt.args.readCloser)
			if (err != nil) != tt.wantErr {
				t.Errorf("CollectText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotBytes, err := ioutil.ReadAll(got)
			assert.NoError(t, err)
			gotString := string(gotBytes)
			if !reflect.DeepEqual(gotString, tt.want) {
				t.Errorf("CollectText() got = %v, want %v", got, tt.want)
			}
		})
	}
}
